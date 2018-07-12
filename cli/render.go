package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/strutil"

	"github.com/essentialkaos/perfecto/check"
	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var levelsPrefixes = map[uint8]string{
	check.LEVEL_NOTICE:   "(N)",
	check.LEVEL_WARNING:  "(W)",
	check.LEVEL_ERROR:    "(E)",
	check.LEVEL_CRITICAL: "(C)",
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

var fallbackLevel = map[uint8]string{
	check.LEVEL_NOTICE:   "N",
	check.LEVEL_WARNING:  "W",
	check.LEVEL_ERROR:    "E",
	check.LEVEL_CRITICAL: "C",
}

// ////////////////////////////////////////////////////////////////////////////////// //

// renderError render error for given format
func renderError(format, file string, err error) {
	filename := strings.Replace(path.Base(file), ".spec", "", -1)

	switch format {
	case FORMAT_TINY:
		fmtc.Printf("%24s: {r}✖ {!} %v\n", filename, err)
	case FORMAT_JSON:
		fmt.Printf("{\"error\":\"%v\"}\n", err)
	case FORMAT_XML:
		fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
		fmt.Println("<alerts>")
		fmt.Printf("  <error>%v</error>\n", err)
		fmt.Println("</alerts>")
	case "", FORMAT_SUMMARY, FORMAT_SHORT:
		printError(err.Error())
	}
}

// renderPerfect render message about perfect spec
func renderPerfect(format, file string) {
	switch format {
	case FORMAT_TINY:
		fmtc.Printf("%24s: {g}✔ {!}\n", file)
	case FORMAT_JSON:
		fmt.Println("{}")
	case FORMAT_XML:
		fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
		fmt.Println("<alerts>\n</alerts>")
	case "", FORMAT_SUMMARY, FORMAT_SHORT:
		fmtc.Println("{g}This spec is perfect!{!}")
	}
}

// renderFullReport render all alerts from report
func renderFullReport(r *check.Report) {
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

	renderSummary(r)

	fmtc.NewLine()
}

// renderFullReport render all alerts from report in short format (used in rpmbuilder)
func renderShortReport(r *check.Report) {
	if len(r.Notices) != 0 {
		renderShortAlerts(check.LEVEL_NOTICE, r.Notices)
	}

	if len(r.Warnings) != 0 {
		renderShortAlerts(check.LEVEL_WARNING, r.Warnings)
	}

	if len(r.Errors) != 0 {
		renderShortAlerts(check.LEVEL_ERROR, r.Errors)
	}

	if len(r.Criticals) != 0 {
		renderShortAlerts(check.LEVEL_CRITICAL, r.Criticals)
	}

	fmtc.NewLine()

	renderSummary(r)
}

// renderTinyReport render tiny report (useful for mass check)
func renderTinyReport(s *spec.Spec, r *check.Report) {
	fmtc.Printf("%24s: ", s.GetFileName())

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

	for _, level := range levels {
		alerts := categories[level]

		if len(alerts) == 0 {
			continue
		}

		for _, alert := range alerts {
			if options.GetB(OPT_NO_COLOR) {
				if alert.Line.Skip {
					fmtc.Printf("X ")
				} else {
					fmtc.Printf(fallbackLevel[level] + " ")
				}
			} else {
				if alert.Line.Skip {
					fmtc.Printf("{s-}%s{!}", "•")
				} else {
					fmtc.Printf(fgColor[level]+"%s{!}", "•")
				}
			}
		}
	}

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

// renderAlerts render all alerts from slice
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

// renderShortAlerts render all alerts from slice
func renderShortAlerts(level uint8, alerts []check.Alert) {
	for _, alert := range alerts {
		renderShortAlert(alert)
	}
}

// renderAlert render alert
func renderAlert(alert check.Alert) {
	fg := fgColor[alert.Level]
	hl := hlColor[alert.Level]

	fmtc.Printf(fg + "│ {!}")

	if alert.Line.Index != -1 {
		if alert.Line.Skip {
			fmtc.Printf(hl+"[%d]{!} {s}[A]{!} "+fg+"%s{!}\n", alert.Line.Index, alert.Info)
		} else {
			fmtc.Printf(hl+"[%d]{!} "+fg+"%s{!}\n", alert.Line.Index, alert.Info)
		}

		if alert.Line.Text != "" {
			text := strutil.Ellipsis(alert.Line.Text, 86)
			fmtc.Printf(fg+"│ {s-}%s{!}\n", text)
		}
	} else {
		fmtc.Printf(hl+"[global]{!} "+fg+"%s{!}\n", alert.Info)
	}
}

// renderShortAlert render short alert
func renderShortAlert(alert check.Alert) {
	fg := fgColor[alert.Level]
	hl := hlColor[alert.Level]

	if fmtc.DisableColors {
		fmtc.Printf(levelsPrefixes[alert.Level] + " ")
	}

	if alert.Line.Index != -1 {
		if alert.Line.Skip {
			fmtc.Printf(hl+"[%d]{!} {s}[A]{!} "+fg+"%s{!}\n", alert.Line.Index, alert.Info)
		} else {
			fmtc.Printf(hl+"[%d]{!} "+fg+"%s{!}\n", alert.Line.Index, alert.Info)
		}
	} else {
		fmtc.Printf(hl+"[global]{!} "+fg+"%s{!}\n", alert.Info)
	}
}

// renderSummary print number for each alert type
func renderSummary(r *check.Report) {
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

	fmtc.Printf("{*}Summary:{!} ")

	for _, level := range levels {
		alerts := categories[level]

		if len(alerts) == 0 {
			continue
		}

		actual, absolved := splitAlertsCount(alerts)

		if absolved != 0 {
			result = append(result,
				fgColor[level]+headers[level]+": "+strconv.Itoa(actual)+"{s-}/"+strconv.Itoa(absolved)+"{!}",
			)
		} else {
			result = append(result,
				fgColor[level]+headers[level]+": "+strconv.Itoa(actual)+"{!}",
			)
		}
	}

	fmtc.Println(strings.Join(result, "{s-} • {!}"))
}
