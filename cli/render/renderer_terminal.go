package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/path"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// TerminalRenderer renders report to terminal
type TerminalRenderer struct {
	Format string

	levelsPrefixes map[uint8]string
	bgColor        map[uint8]string
	fgColor        map[uint8]string
	hlColor        map[uint8]string
	headers        map[uint8]string
	fallbackLevel  map[uint8]string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Report renders alerts from perfecto report
func (r *TerminalRenderer) Report(file string, report *check.Report) {
	r.initUI()

	switch r.Format {
	case "summary":
		r.renderSummary(report)
	case "short":
		r.renderShortReport(report)
	case "tiny":
		r.renderTinyReport(file, report)
	default:
		r.renderFull(report)
	}
}

// Perfect renders message about perfect spec
func (r *TerminalRenderer) Perfect(file string, report *check.Report) {
	r.initUI()

	specName := strutil.Exclude(path.Base(file), ".spec")

	switch r.Format {
	case "tiny":
		fmtc.Printf("%24s{s}:{!} {g}✔ {!}\n", specName)
	case "summary":
		r.renderSummary(report)
	default:
		fmtc.Printf("{g}{*}%s.spec{!*} is perfect!{!}\n", specName)
	}
}

// Skipped renders message about skipped check
func (r *TerminalRenderer) Skipped(file string, report *check.Report) {
	specName := strutil.Exclude(path.Base(file), ".spec")

	switch r.Format {
	case "tiny":
		fmtc.Printf("%24s{s}:{!} {s}—{!}\n", specName)
	default:
		fmtc.Printf("{s}{*}%s.spec{!*} check skipped due to non-applicable target{!}\n", specName)
	}
}

// Error renders global error message
func (r *TerminalRenderer) Error(file string, err error) {
	specName := strutil.Exclude(path.Base(file), ".spec")

	switch r.Format {
	case "tiny":
		fmtc.Printf("%24s{s}:{!} {r}✖ (%v){!}\n", specName, err)
	default:
		fmtc.Fprintf(os.Stderr, "{r}%v{!}\n", err)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// initUI initialize UI
func (r *TerminalRenderer) initUI() {
	r.levelsPrefixes = map[uint8]string{
		check.LEVEL_NOTICE:   "<N>",
		check.LEVEL_WARNING:  "<W>",
		check.LEVEL_ERROR:    "<E>",
		check.LEVEL_CRITICAL: "<C>",
	}

	r.bgColor = map[uint8]string{
		check.LEVEL_NOTICE:   "{*@c}",
		check.LEVEL_WARNING:  "{*@y}",
		check.LEVEL_ERROR:    "{*@r}",
		check.LEVEL_CRITICAL: "{*@r}",
	}

	r.fgColor = map[uint8]string{
		check.LEVEL_NOTICE:   "{c}",
		check.LEVEL_WARNING:  "{y}",
		check.LEVEL_ERROR:    "{r}",
		check.LEVEL_CRITICAL: "{r}",
	}

	r.hlColor = map[uint8]string{
		check.LEVEL_NOTICE:   "{c*-}",
		check.LEVEL_WARNING:  "{y*-}",
		check.LEVEL_ERROR:    "{r*-}",
		check.LEVEL_CRITICAL: "{r*-}",
	}

	r.headers = map[uint8]string{
		check.LEVEL_NOTICE:   "Notice",
		check.LEVEL_WARNING:  "Warning",
		check.LEVEL_ERROR:    "Error",
		check.LEVEL_CRITICAL: "Critical",
	}

	r.fallbackLevel = map[uint8]string{
		check.LEVEL_NOTICE:   "N",
		check.LEVEL_WARNING:  "W",
		check.LEVEL_ERROR:    "E",
		check.LEVEL_CRITICAL: "C",
	}

	// Set color to orange for errors if 256 colors are supported
	if fmtc.Is256ColorsSupported() {
		r.bgColor[check.LEVEL_ERROR] = "{*@}{#208}"
		r.fgColor[check.LEVEL_ERROR] = "{#208}"
		r.hlColor[check.LEVEL_ERROR] = "{*}{#214}"
	}
}

// renderFull prints full report
func (r *TerminalRenderer) renderFull(report *check.Report) {
	fmtc.NewLine()

	if len(report.Notices) != 0 {
		r.renderHeader(check.LEVEL_NOTICE, len(report.Notices))
		r.renderAlerts(check.LEVEL_NOTICE, report.Notices)
	}

	if len(report.Warnings) != 0 {
		r.renderHeader(check.LEVEL_WARNING, len(report.Warnings))
		r.renderAlerts(check.LEVEL_WARNING, report.Warnings)
	}

	if len(report.Errors) != 0 {
		r.renderHeader(check.LEVEL_ERROR, len(report.Errors))
		r.renderAlerts(check.LEVEL_ERROR, report.Errors)
	}

	if len(report.Criticals) != 0 {
		r.renderHeader(check.LEVEL_CRITICAL, len(report.Criticals))
		r.renderAlerts(check.LEVEL_CRITICAL, report.Criticals)
	}

	r.renderLinks(report)

	fmtutil.Separator(true)

	fmtc.NewLine()

	r.renderSummary(report)

	fmtc.NewLine()
}

// renderFullReport prints all alerts from report in short format (used in rpmbuilder)
func (r *TerminalRenderer) renderShortReport(report *check.Report) {
	if report.Notices.Total() != 0 {
		r.renderShortAlerts(check.LEVEL_NOTICE, report.Notices)
	}

	if report.Warnings.Total() != 0 {
		r.renderShortAlerts(check.LEVEL_WARNING, report.Warnings)
	}

	if report.Errors.Total() != 0 {
		r.renderShortAlerts(check.LEVEL_ERROR, report.Errors)
	}

	if report.Criticals.Total() != 0 {
		r.renderShortAlerts(check.LEVEL_CRITICAL, report.Criticals)
	}
}

// renderTinyReport prints tiny report (useful for mass check)
func (r *TerminalRenderer) renderTinyReport(file string, report *check.Report) {
	specName := strutil.Exclude(path.Base(file), ".spec")

	fmtc.Printf("%24s{s}:{!} ", specName)

	categories := map[uint8][]check.Alert{
		check.LEVEL_NOTICE:   report.Notices,
		check.LEVEL_WARNING:  report.Warnings,
		check.LEVEL_ERROR:    report.Errors,
		check.LEVEL_CRITICAL: report.Criticals,
	}

	levels := []uint8{
		check.LEVEL_NOTICE,
		check.LEVEL_WARNING,
		check.LEVEL_ERROR,
		check.LEVEL_CRITICAL,
	}

	for _, level := range levels {
		alerts := categories[level]

		if len(alerts) == 0 {
			continue
		}

		for _, alert := range alerts {
			if fmtc.DisableColors {
				if alert.IsIgnored {
					fmtc.Printf("X ")
				} else {
					fmtc.Printf(r.fallbackLevel[level] + " ")
				}
			} else {
				if alert.IsIgnored {
					fmtc.Printf("{s-}%s{!}", "•")
				} else {
					fmtc.Printf(r.fgColor[level]+"%s{!}", "•")
				}
			}
		}
	}

	fmtc.NewLine()
}

// renderHeader prints level header
func (r *TerminalRenderer) renderHeader(level uint8, count int) {
	header := r.headers[level] + fmt.Sprintf(" (%d)", count)

	fg := r.fgColor[level]
	bg := r.bgColor[level]

	fmtc.Printf(bg+" ••• %-83s{!}\n", header)
	fmtc.Printf(fg + "│{!}\n")
}

// renderAlerts prints all alerts from given slice
func (r *TerminalRenderer) renderAlerts(level uint8, alerts []check.Alert) {
	totalAlerts := len(alerts)

	for index, alert := range alerts {
		r.renderAlert(alert)

		if index+1 < totalAlerts {
			fmtc.Printf(r.fgColor[level] + "│{!}\n")
		}
	}

	fmtc.NewLine()
}

// renderAlert prints detailed info about given alert
func (r *TerminalRenderer) renderAlert(alert check.Alert) {
	fg := r.fgColor[alert.Level]
	hl := r.hlColor[alert.Level]
	lc := fg

	if alert.IsIgnored {
		fg = "{s}"
		hl = "{s*}"
	}

	fmtc.Printf(lc + "│ {!}")

	if alert.Line.Index != -1 {
		fmtc.Printf(hl+"[%d]{!} ", alert.Line.Index)
	} else {
		fmtc.Printf(hl + "[global]{!} ")
	}

	if alert.IsIgnored {
		fmtc.Printf("{s}[I]{!} ")
	}

	if alert.ID != "" {
		fmtc.Printf(fg+"(%s) %s{!}\n", alert.ID, alert.Info)
	} else {
		fmtc.Printf(fg+"(rpmlint) %s{!}\n", alert.Info)
	}

	if alert.Line.Text != "" {
		text := strutil.Ellipsis(alert.Line.Text, 86)
		if alert.IsIgnored {
			fmtc.Printf(lc+"│ {s-}%s{!}\n", text)
		} else {
			fmtc.Printf(lc+"│ {s}%s{!}\n", text)
		}
	}
}

// renderLinks prints links to mentioned failed checks
func (r *TerminalRenderer) renderLinks(report *check.Report) {
	ids := report.IDs()

	if len(ids) == 0 {
		return
	}

	fmtutil.Separator(true)

	fmtc.Println("\n{*}Links:{!}\n")

	for _, id := range report.IDs() {
		fmtc.Printf(" {s}•{!} %s%s\n", "https://kaos.sh/perfecto/w/", id)
	}

	fmtc.NewLine()
}

// renderSummary prints report statistics
func (r *TerminalRenderer) renderSummary(report *check.Report) {
	categories := map[uint8]check.Alerts{
		check.LEVEL_NOTICE:   report.Notices,
		check.LEVEL_WARNING:  report.Warnings,
		check.LEVEL_ERROR:    report.Errors,
		check.LEVEL_CRITICAL: report.Criticals,
	}

	levels := []uint8{
		check.LEVEL_NOTICE,
		check.LEVEL_WARNING,
		check.LEVEL_ERROR,
		check.LEVEL_CRITICAL,
	}

	var result []string

	for _, level := range levels {
		alerts := categories[level]

		if len(alerts) == 0 {
			result = append(result,
				r.fgColor[level]+"{*}"+r.headers[level]+":{!*} 0{!}",
			)
			continue
		}

		total, ignored := alerts.Total(), alerts.Ignored()
		actual := total - ignored

		if ignored != 0 {
			result = append(result,
				r.fgColor[level]+"{*}"+r.headers[level]+":{!*} "+strconv.Itoa(actual)+"{s-}/"+strconv.Itoa(ignored)+"{!}",
			)
		} else {
			result = append(result,
				r.fgColor[level]+"{*}"+r.headers[level]+":{!*} "+strconv.Itoa(actual),
			)
		}
	}

	fmtc.Println(strings.Join(result, "{s-} • {!}"))
}

// renderShortAlerts prints all alerts from slice in short format
func (r *TerminalRenderer) renderShortAlerts(level uint8, alerts []check.Alert) {
	for _, alert := range alerts {
		r.renderShortAlert(alert)
	}
}

// renderShortAlert render short alert
func (r *TerminalRenderer) renderShortAlert(alert check.Alert) {
	fg := r.fgColor[alert.Level]
	hl := r.hlColor[alert.Level]

	if fmtc.DisableColors {
		fmtc.Printf(r.levelsPrefixes[alert.Level] + " ")
	}

	if alert.Line.Index != -1 {
		fmtc.Printf(hl+"[%d]{!} ", alert.Line.Index)
	} else {
		fmtc.Printf(hl + "[global]{!} ")
	}

	if alert.IsIgnored {
		fmtc.Printf("{s}[I]{!} ")
	}

	if alert.ID != "" {
		fmtc.Printf(fg+"(%s) %s{!}\n", alert.ID, alert.Info)
	} else {
		fmtc.Printf(fg+"(rpmlint) %s{!}\n", alert.Info)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //
