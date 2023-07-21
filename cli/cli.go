package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/mathutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/perfecto/check"
	"github.com/essentialkaos/perfecto/spec"

	"github.com/essentialkaos/perfecto/cli/render"
	"github.com/essentialkaos/perfecto/cli/support"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "Perfecto"
	VER  = "6.0.0"
	DESC = "Tool for checking perfectly written RPM specs"
)

// Options
const (
	OPT_FORMAT      = "f:format"
	OPT_LINT_CONFIG = "c:lint-config"
	OPT_ERROR_LEVEL = "e:error-level"
	OPT_IGNORE      = "I:ignore"
	OPT_QUIET       = "q:quiet"
	OPT_NO_LINT     = "nl:no-lint"
	OPT_NO_COLOR    = "nc:no-color"
	OPT_HELP        = "h:help"
	OPT_VER         = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
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
	OPT_IGNORE:      {Mergeble: true, Alias: "A:absolve"},
	OPT_FORMAT:      {},
	OPT_LINT_CONFIG: {},
	OPT_ERROR_LEVEL: {},
	OPT_QUIET:       {Type: options.BOOL},
	OPT_NO_LINT:     {Type: options.BOOL},
	OPT_NO_COLOR:    {Type: options.BOOL},
	OPT_HELP:        {Type: options.BOOL},
	OPT_VER:         {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
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

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	process(args)
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	term := os.Getenv("TERM")

	fmtc.DisableColors = true
	strutil.EllipsisSuffix = "…"

	if term != "" {
		switch {
		case strings.Contains(term, "xterm"),
			strings.Contains(term, "color"),
			term == "screen":
			fmtc.DisableColors = false
		}
	}

	// Check for output redirect using pipes
	if fsutil.IsCharacterDevice("/dev/stdin") &&
		!fsutil.IsCharacterDevice("/dev/stdout") &&
		os.Getenv("FAKETTY") == "" {
		fmtc.DisableColors = true
	}

	if os.Getenv("CI") != "" {
		fmtc.DisableColors = false
	}

	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configure UI on start
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// process start spec file processing
func process(files options.Arguments) {
	var exitCode int

	format := getFormat(files)

	if !sliceutil.Contains(formats, format) {
		printErrorAndExit("Output format %q is not supported", format)
	}

	rndr := getRenderer(format, files)

	for _, file := range files {
		ec := checkSpec(file.Clean().String(), rndr)
		exitCode = mathutil.Max(ec, exitCode)
	}

	os.Exit(exitCode)
}

// codebeat:disable[ABC]

// checkSpec check spec file
func checkSpec(file string, rndr render.Renderer) int {
	var ignoreChecks []string

	s, err := spec.Read(file)

	if err != nil && !options.GetB(OPT_QUIET) {
		rndr.Error(file, err)
		return 1
	}

	if options.Has(OPT_IGNORE) {
		ignoreChecks = strings.Split(options.GetS(OPT_IGNORE), ",")
	}

	report := check.Check(
		s, !options.GetB(OPT_NO_LINT),
		options.GetS(OPT_LINT_CONFIG),
		ignoreChecks,
	)

	switch {
	case report.IsSkipped:
		if !options.GetB(OPT_QUIET) {
			rndr.Skipped(file, report)
		}
		return 0
	case report.IsPerfect:
		if !options.GetB(OPT_QUIET) {
			rndr.Perfect(file, report)
		}
		return 0
	}

	rndr.Report(file, report)

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
func getRenderer(format string, files options.Arguments) render.Renderer {
	maxFilenameSize := getMaxFilenameSize(files)

	switch format {
	case FORMAT_SUMMARY:
		return &render.TerminalRenderer{Format: FORMAT_SUMMARY, FilenameSize: maxFilenameSize}
	case FORMAT_TINY:
		return &render.TerminalRenderer{Format: FORMAT_TINY, FilenameSize: maxFilenameSize}
	case FORMAT_SHORT:
		return &render.TerminalRenderer{Format: FORMAT_SHORT, FilenameSize: maxFilenameSize}
	case FORMAT_GITHUB:
		return &render.GithubRenderer{}
	case FORMAT_JSON:
		return &render.JSONRenderer{}
	case FORMAT_XML:
		return &render.XMLRenderer{}
	default:
		return &render.TerminalRenderer{Format: FORMAT_FULL, FilenameSize: maxFilenameSize}
	}
}

// getMaxFilenameSize returns maximum filename size without extension
func getMaxFilenameSize(files options.Arguments) int {
	var result int

	for _, file := range files {
		filenameSize := strutil.Exclude(file.Base().Clean().String(), ".spec")
		result = mathutil.Max(result, len(filenameSize))
	}

	return result
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

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
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

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "spec…")

	info.AddOption(OPT_IGNORE, "Disable one or more checks by their ID", "id…")
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
		"--ignore PF2,PF12 app.spec",
		"Check spec without PF2 and PF12 checks",
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

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/perfecto", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
