package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"unicode/utf8"

	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// JSONRenderer renders report in XML format
type XMLRenderer struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

// Report renders alerts from perfecto report
func (r *XMLRenderer) Report(file string, report *check.Report) error {
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println("<alerts>")

	if len(report.Notices) != 0 {
		r.renderAlertsAsXML("notices", report.Notices)
	}

	if len(report.Warnings) != 0 {
		r.renderAlertsAsXML("warnings", report.Warnings)
	}

	if len(report.Errors) != 0 {
		r.renderAlertsAsXML("errors", report.Errors)
	}

	if len(report.Criticals) != 0 {
		r.renderAlertsAsXML("criticals", report.Criticals)
	}

	fmt.Println("</alerts>")

	return nil
}

// Perfect renders message about perfect spec
func (r *XMLRenderer) Perfect(file string) {
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println("<alerts>\n</alerts>")
}

// Error renders global error message
func (r *XMLRenderer) Error(file string, err error) {
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println("<alerts>")
	fmt.Printf("  <error>%v</error>\n", err)
	fmt.Println("</alerts>")
}

// ////////////////////////////////////////////////////////////////////////////////// //

// renderAlertsAsXML renders alerts category as XML node
func (r *XMLRenderer) renderAlertsAsXML(category string, alerts []check.Alert) {
	fmt.Printf("  <%s>\n", category)

	for _, alert := range alerts {
		fmt.Printf("    <alert id=\"%s\" absolve=\"%t\">\n", alert.ID, alert.Absolve)
		fmt.Printf("      <info>%s</info>\n", r.escapeStringForXML(alert.Info))

		if alert.Line.Index != -1 {
			fmt.Printf(
				"      <line index=\"%d\" skip=\"%t\">%s</line>\n",
				alert.Line.Index, alert.Line.Skip,
				r.escapeStringForXML(alert.Line.Text),
			)
		}

		fmt.Println("    </alert>")
	}

	fmt.Printf("  </%s>\n", category)
}

// escapeStringForXML returns properly escaped XML equivalent
// of the plain text data
func (r *XMLRenderer) escapeStringForXML(s string) string {
	var result, esc string
	var last int

	for i := 0; i < len(s); {
		rn, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch rn {
		case '"':
			esc = "&#34;"
		case '\'':
			esc = "&#39;"
		case '&':
			esc = "&amp;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '\t':
			esc = "&#x9;"
		case '\n':
			esc = "&#xA;"
		case '\r':
			esc = "&#xD;"
		default:
			if !r.isInCharacterRange(rn) || (rn == 0xFFFD && width == 1) {
				esc = "\uFFFD"
				break
			}
			continue
		}

		result += s[last : i-width]
		result += esc

		last = i
	}

	result += s[last:]

	return result
}

// isInCharacterRange decides whether the given rune is in the XML Character Range
func (r *XMLRenderer) isInCharacterRange(rn rune) (inrange bool) {
	return rn == 0x09 ||
		rn == 0x0A ||
		rn == 0x0D ||
		rn >= 0x20 && rn <= 0xDF77 ||
		rn >= 0xE000 && rn <= 0xFFFD ||
		rn >= 0x10000 && rn <= 0x10FFFF
}
