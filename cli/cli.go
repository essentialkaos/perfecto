package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/mathutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/sliceutil"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	"github.com/essentialkaos/perfecto/check"
	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "Perfecto"
	VER  = "3.7.0"
	DESC = "Tool for checking perfectly written RPM specs"
)

// Options
const (
	OPT_FORMAT      = "f:format"
	OPT_LINT_CONFIG = "c:lint-config"
	OPT_ERROR_LEVEL = "e:error-level"
	OPT_ABSOLVE     = "A:absolve"
	OPT_QUIET       = "q:quiet"
	OPT_NO_LINT     = "nl:no-lint"
	OPT_NO_COLOR    = "nc:no-color"
	OPT_HELP        = "h:help"
	OPT_VER         = "v:version"

	OPT_COMPLETION = "completion"
)

// Supported formats
const (
	FORMAT_SUMMARY = "summary"
	FORMAT_SHORT   = "short"
	FORMAT_TINY    = "tiny"

	FORMAT_JSON = "json"
	FORMAT_XML  = "xml"
)

// Levels
const (
	LEVEL_NOTICE   = "notice"
	LEVEL_WARNING  = "warning"
	LEVEL_ERROR    = "error"
	LEVEL_CRITICAL = "critical"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// options map
var optMap = options.Map{
	OPT_ABSOLVE:     {Mergeble: true},
	OPT_FORMAT:      {Type: options.STRING},
	OPT_LINT_CONFIG: {Type: options.STRING},
	OPT_ERROR_LEVEL: {Type: options.STRING},
	OPT_QUIET:       {Type: options.BOOL},
	OPT_NO_LINT:     {Type: options.BOOL},
	OPT_NO_COLOR:    {Type: options.BOOL},
	OPT_HELP:        {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:         {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main function of cli
func Init() {
	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.Has(OPT_COMPLETION) {
		genCompletion()
	}

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	process(args)
}

// configureUI configure UI on start
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	// Set color to orange for errors if 256 colors are supported
	if fmtc.Is256ColorsSupported() {
		bgColor[check.LEVEL_ERROR] = "{*@}{#208}"
		fgColor[check.LEVEL_ERROR] = "{#208}"
		hlColor[check.LEVEL_ERROR] = "{*}{#214}"
	}

	strutil.EllipsisSuffix = "…"
}

// process start spec file processing
func process(files []string) {
	var exitCode int

	format := options.GetS(OPT_FORMAT)

	if !sliceutil.Contains([]string{FORMAT_TINY, FORMAT_SHORT, FORMAT_SUMMARY, FORMAT_JSON, FORMAT_XML, ""}, format) {
		printErrorAndExit("Output format \"%s\" is not supported", format)
	}

	if len(files) > 1 {
		format = FORMAT_TINY
	}

	for _, file := range files {
		ec := checkSpec(file, format)
		exitCode = mathutil.Max(ec, exitCode)
	}

	os.Exit(exitCode)
}

// codebeat:disable[ABC]

// checkSpec check spec file
func checkSpec(file, format string) int {
	s, err := spec.Read(file)

	if err != nil && !options.GetB(OPT_QUIET) {
		renderError(format, file, err)
		return 1
	}

	report := check.Check(
		s, !options.GetB(OPT_NO_LINT),
		options.GetS(OPT_LINT_CONFIG),
		strings.Split(options.GetS(OPT_ABSOLVE), ","),
	)

	if report.IsPerfect() {
		if !options.GetB(OPT_QUIET) {
			renderPerfect(format, s.GetFileName())
		}

		return 0
	}

	if !options.GetB(OPT_QUIET) {
		switch format {
		case FORMAT_SUMMARY:
			renderSummary(report)
		case FORMAT_TINY:
			renderTinyReport(s, report)
		case FORMAT_SHORT:
			renderShortReport(report)
		case FORMAT_JSON:
			renderJSONReport(report)
		case FORMAT_XML:
			renderXMLReport(report)
		case "":
			renderFullReport(report)
		}
	}

	return getExitCode(report)
}

// codebeat:enable[ABC]

// getExitCode return exit code based on report data
func getExitCode(r *check.Report) int {
	var maxLevel int
	var nonZero bool

	switch {
	case countAlerts(r.Criticals) != 0:
		maxLevel = 4
	case countAlerts(r.Errors) != 0:
		maxLevel = 3
	case countAlerts(r.Warnings) != 0:
		maxLevel = 2
	case countAlerts(r.Notices) != 0:
		maxLevel = 1
	}

	switch options.GetS(OPT_ERROR_LEVEL) {
	case LEVEL_NOTICE:
		nonZero = maxLevel >= 1
	case LEVEL_WARNING:
		nonZero = maxLevel >= 2
	case LEVEL_ERROR:
		nonZero = maxLevel >= 3
	case LEVEL_CRITICAL:
		nonZero = maxLevel == 4
	default:
		nonZero = maxLevel != 0
	}

	if nonZero {
		return 1
	}

	return 0
}

// countAlerts return number of actual alerts
func countAlerts(alerts []check.Alert) int {
	var counter int

	for _, alert := range alerts {
		if !alert.Absolve {
			counter++
		}
	}

	return counter
}

// splitAlertsCount count actual and absolved alerts
func splitAlertsCount(alerts []check.Alert) (int, int) {
	actual := countAlerts(alerts)
	absolved := len(alerts) - actual

	return actual, absolved
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "file…")

	info.AddOption(OPT_ABSOLVE, "Disable some checks by their ID", "id…")
	info.AddOption(OPT_FORMAT, "Output format {s-}(summary|tiny|short|json|xml){!}", "format")
	info.AddOption(OPT_LINT_CONFIG, "Path to RPMLint configuration file", "file")
	info.AddOption(OPT_ERROR_LEVEL, "Return non-zero exit code if alert level greater than given {s-}(notice|warning|error|critical){!}", "level")
	info.AddOption(OPT_QUIET, "Suppress all normal output")
	info.AddOption(OPT_NO_LINT, "Disable RPMLint checks")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample("app.spec", "Check spec and print extended report")

	info.AddExample(
		"--no-lint app.spec",
		"Check spec without rpmlint and print extended report",
	)

	info.AddExample("--format tiny app.spec", "Check spec and print tiny report")
	info.AddExample("--format summary app.spec", "Check spec and print summary")

	info.AddExample(
		"--format json app.spec 1> report.json",
		"Check spec, generate report in JSON format and save as report.json",
	)

	return info
}

// genCompletion generates completion for different shells
func genCompletion() {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "perfecto", "spec"))
	case "fish":
		fmt.Printf(fish.Generate(info, "perfecto"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "perfecto", "*.spec"))
	default:
		os.Exit(1)
	}

	os.Exit(0)
}

// showAbout shows info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/perfecto", update.GitHubChecker},
	}

	about.Render()
}
