package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/mathutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/perfecto/check"
	"github.com/essentialkaos/perfecto/spec"

	"github.com/essentialkaos/perfecto/cli/render"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "Perfecto"
	VER  = "4.0.0"
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
	FORMAT_FULL    = "full"
	FORMAT_SUMMARY = "summary"
	FORMAT_SHORT   = "short"
	FORMAT_TINY    = "tiny"
	FORMAT_GITHUB  = "github"
	FORMAT_JSON    = "json"
	FORMAT_XML     = "xml"
)

// Levels
const (
	LEVEL_NOTICE   = "notice"
	LEVEL_WARNING  = "warning"
	LEVEL_ERROR    = "error"
	LEVEL_CRITICAL = "critical"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap is map with all supported options
var optMap = options.Map{
	OPT_ABSOLVE:     {Mergeble: true},
	OPT_FORMAT:      {},
	OPT_LINT_CONFIG: {},
	OPT_ERROR_LEVEL: {},
	OPT_QUIET:       {Type: options.BOOL},
	OPT_NO_LINT:     {Type: options.BOOL},
	OPT_NO_COLOR:    {Type: options.BOOL},
	OPT_HELP:        {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:         {Type: options.BOOL, Alias: "ver"},

	OPT_COMPLETION: {},
}

// formats is slice with all supported formats
var formats = []string{
	FORMAT_FULL,
	FORMAT_SUMMARY,
	FORMAT_SHORT,
	FORMAT_TINY,
	FORMAT_GITHUB,
	FORMAT_JSON,
	FORMAT_XML,
	"",
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

	strutil.EllipsisSuffix = "…"
}

// process start spec file processing
func process(files options.Arguments) {
	var exitCode int

	format := getFormat(files)

	if !sliceutil.Contains(formats, format) {
		printErrorAndExit("Output format %q is not supported", format)
	}

	for _, file := range files {
		ec := checkSpec(file.Clean().String(), format)
		exitCode = mathutil.Max(ec, exitCode)
	}

	os.Exit(exitCode)
}

// codebeat:disable[ABC]

// checkSpec check spec file
func checkSpec(file, format string) int {
	rnd := getRenderer(format)
	s, err := spec.Read(file)

	if err != nil && !options.GetB(OPT_QUIET) {
		rnd.Error(file, err)
		return 1
	}

	report := check.Check(
		s, !options.GetB(OPT_NO_LINT),
		options.GetS(OPT_LINT_CONFIG),
		strings.Split(options.GetS(OPT_ABSOLVE), ","),
	)

	if report.IsPerfect() {
		if !options.GetB(OPT_QUIET) {
			rnd.Perfect(s.GetFileName())
		}

		return 0
	}

	err = rnd.Report(s.GetFileName(), report)

	if err != nil {
		printError(err.Error())
		return 1
	}

	return getExitCode(report)
}

// getFormat returns output format
func getFormat(files options.Arguments) string {
	format := options.GetS(OPT_FORMAT)

	if len(files) > 1 {
		switch format {
		case FORMAT_JSON, FORMAT_XML:
			printErrorAndExit("Can't check multiple files with %q output format", format)
		case "":
			format = FORMAT_TINY
		}
	} else {
		if format == "" && os.Getenv("GITHUB_ACTIONS") == "true" {
			format = FORMAT_GITHUB
		}
	}

	return format
}

// getRenderer returns renderer for given format
func getRenderer(format string) render.Renderer {
	switch format {
	case FORMAT_SUMMARY:
		return &render.TerminalRenderer{Format: FORMAT_SUMMARY}
	case FORMAT_TINY:
		return &render.TerminalRenderer{Format: FORMAT_TINY}
	case FORMAT_SHORT:
		return &render.TerminalRenderer{Format: FORMAT_SHORT}
	case FORMAT_GITHUB:
		return &render.GithubRenderer{}
	case FORMAT_JSON:
		return &render.JSONRenderer{}
	case FORMAT_XML:
		return &render.XMLRenderer{}
	default:
		return &render.TerminalRenderer{Format: FORMAT_FULL}
	}
}

// codebeat:enable[ABC]

// getExitCode return exit code based on report data
func getExitCode(r *check.Report) int {
	var maxLevel int
	var nonZero bool

	switch {
	case r.Criticals.HasAlerts():
		maxLevel = 4
	case r.Errors.HasAlerts():
		maxLevel = 3
	case r.Warnings.HasAlerts():
		maxLevel = 2
	case r.Notices.HasAlerts():
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
	info.AddOption(OPT_FORMAT, "Output format {s-}(summary|tiny|short|github|json|xml){!}", "format")
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

	info.AddExample(
		"--format tiny app.spec",
		"Check spec and print tiny report",
	)

	info.AddExample(
		"--format summary app.spec",
		"Check spec and print summary",
	)

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
