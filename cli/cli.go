package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strings"

	"pkg.re/essentialkaos/ek.v9/env"
	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/strutil"
	"pkg.re/essentialkaos/ek.v9/usage"
	"pkg.re/essentialkaos/ek.v9/usage/update"

	"github.com/essentialkaos/perfecto/check"
	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "Perfecto"
	VER  = "0.0.1"
	DESC = "Tool for checking perfectly written RPM specs"
)

// Options
const (
	OPT_SUMMARY  = "s:summary"
	OPT_NO_LINT  = "nl:no-lint"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var optMap = options.Map{
	OPT_SUMMARY:  {Type: options.BOOL},
	OPT_NO_LINT:  {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

var bgColor = map[uint8]string{
	check.LEVEL_NOTICE:   "{*@c}",
	check.LEVEL_WARNING:  "{*@y}",
	check.LEVEL_ERROR:    "{*@r}",
	check.LEVEL_CRITICAL: "{*@r}",
}

var fgColor = map[uint8]string{
	check.LEVEL_NOTICE:   "{c}",
	check.LEVEL_WARNING:  "{y}",
	check.LEVEL_ERROR:    "{r}",
	check.LEVEL_CRITICAL: "{r}",
}

var hlColor = map[uint8]string{
	check.LEVEL_NOTICE:   "{c*-}",
	check.LEVEL_WARNING:  "{y*-}",
	check.LEVEL_ERROR:    "{r*-}",
	check.LEVEL_CRITICAL: "{r*-}",
}

var headers = map[uint8]string{
	check.LEVEL_NOTICE:   "Notice",
	check.LEVEL_WARNING:  "Warning",
	check.LEVEL_ERROR:    "Error",
	check.LEVEL_CRITICAL: "Critical",
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

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) == 0 {
		showUsage()
		return
	}

	process(args[0])
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
func process(file string) {
	s, err := spec.Read(file)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if !options.GetB(OPT_NO_LINT) && !isLinterInstalled() {
		printErrorAndExit("Can't run linter: rpmlint not installed. Install rpmlint or use option '--no-lint'.")
	}

	report := check.Check(s, !options.GetB(OPT_NO_LINT))

	if report.IsPerfect() {
		fmtc.Println("{g}Your spec is perfect. Great job!{!}")
		os.Exit(0)
	}

	if options.GetB(OPT_SUMMARY) {
		renderResume(report)
	} else {
		renderReport(report)
	}

	os.Exit(1)
}

// renderReport render all alerts from report
func renderReport(r *check.Report) {
	fmtc.NewLine()

	if len(r.Notices) != 0 {
		renderHeader(check.LEVEL_NOTICE, len(r.Notices))
		renderAlerts(check.LEVEL_NOTICE, r.Notices)
	}

	if len(r.Warnings) != 0 {
		renderHeader(check.LEVEL_WARNING, len(r.Warnings))
		renderAlerts(check.LEVEL_WARNING, r.Warnings)
	}

	if len(r.Errors) != 0 {
		renderHeader(check.LEVEL_ERROR, len(r.Errors))
		renderAlerts(check.LEVEL_ERROR, r.Errors)
	}

	if len(r.Criticals) != 0 {
		renderHeader(check.LEVEL_CRITICAL, len(r.Criticals))
		renderAlerts(check.LEVEL_CRITICAL, r.Criticals)
	}

	fmtc.Printf("{s}%s{!}\n\n", strings.Repeat("-", 88))

	renderResume(r)

	fmtc.NewLine()
}

// renderHeader render header
func renderHeader(level uint8, count int) {
	header := headers[level] + fmt.Sprintf(" (%d)", count)

	fg := fgColor[level]
	bg := bgColor[level]

	fmtc.Printf(bg+" ••• %-83s{!}\n", header)
	fmtc.Printf(fg + "│{!}\n")
}

// renderAlerts render all alerts fron slice
func renderAlerts(level uint8, alerts []check.Alert) {
	totalAlerts := len(alerts)

	for index, alert := range alerts {
		renderAlert(alert)

		if index+1 < totalAlerts {
			fmtc.Printf(fgColor[level] + "│{!}\n")
		}
	}

	fmtc.NewLine()
}

// renderAlert render alert
func renderAlert(alert check.Alert) {
	fg := fgColor[alert.Level]
	hl := hlColor[alert.Level]

	if alert.Line.Index != -1 {
		fmtc.Printf(fg+"│ "+hl+"[%d]{!} "+fg+"%s{!}\n", alert.Line.Index, alert.Info)

		if alert.Line.Text != "" {
			text := strutil.Ellipsis(alert.Line.Text, 86)
			fmtc.Printf(fg+"│ {s-}%s{!}\n", text)
		}
	} else {
		fmtc.Printf(fg+"│ "+hl+"[global]{!} "+fg+"%s{!}\n", alert.Info)
	}
}

// renderResume print number for each alert type
func renderResume(r *check.Report) {
	categories := map[uint8][]check.Alert{
		check.LEVEL_NOTICE:   r.Notices,
		check.LEVEL_WARNING:  r.Warnings,
		check.LEVEL_ERROR:    r.Errors,
		check.LEVEL_CRITICAL: r.Criticals,
	}

	levels := []uint8{
		check.LEVEL_NOTICE,
		check.LEVEL_WARNING,
		check.LEVEL_ERROR,
		check.LEVEL_CRITICAL,
	}

	var result []string

	fmtc.Printf("{*}Resume:{!} ")

	for _, level := range levels {
		alerts := categories[level]

		if len(alerts) == 0 {
			continue
		}

		result = append(result,
			fgColor[level]+headers[level]+": "+fmtc.Sprintf("%d", len(alerts))+"{!}",
		)
	}
	fmtc.Println(strings.Join(result, "{s-} • {!}"))
}

// isLinterInstalled checks if rpmlint is installed
func isLinterInstalled() bool {
	return env.Which("rpmlint") != ""
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

// showUsage show usage info
func showUsage() {
	info := usage.NewInfo("spec-file")

	info.AddOption(OPT_SUMMARY, "Print only summary info")
	info.AddOption(OPT_NO_LINT, "Disable rpmlint checks")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.Render()
}

// showAbout show info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2006,
		Owner:         "ESSENTIAL KAOS",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/perfecto", update.GitHubChecker},
	}

	about.Render()
}
