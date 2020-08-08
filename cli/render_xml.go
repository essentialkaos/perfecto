package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"unicode/utf8"

	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// renderXMLReport render report in XML format
func renderXMLReport(r *check.Report) {
	fmt.Println(`<?xml version="1.0" encoding="UTF-8"?>`)
	fmt.Println("<alerts>")

	if len(r.Notices) != 0 {
		renderAlertsAsXML("notices", r.Notices)
	}

	if len(r.Warnings) != 0 {
		renderAlertsAsXML("warnings", r.Warnings)
	}

	if len(r.Errors) != 0 {
		renderAlertsAsXML("errors", r.Errors)
	}

	if len(r.Criticals) != 0 {
		renderAlertsAsXML("criticals", r.Criticals)
	}

	fmt.Println("</alerts>")
}

// renderAlertsAsXML render alerts category as xml node
func renderAlertsAsXML(category string, alerts []check.Alert) {
	fmt.Printf("  <%s>\n", category)

	for _, alert := range alerts {
		fmt.Printf("    <alert id=\"%s\" absolve=\"%t\">\n", alert.ID, alert.Absolve)
		fmt.Printf("      <info>%s</info>\n", escapeStringForXML(alert.Info))

		if alert.Line.Index != -1 {
			fmt.Printf(
				"      <line index=\"%d\" skip=\"%t\">%s</line>\n",
				alert.Line.Index, alert.Line.Skip,
				escapeStringForXML(alert.Line.Text),
			)
		}

		fmt.Println("    </alert>")
	}

	fmt.Printf("  </%s>\n", category)
}

// escapeStringForXML return properly escaped XML equivalent
// of the plain text data
func escapeStringForXML(s string) string {
	var result, esc string
	var last int

	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width
		switch r {
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
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
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

// Decide whether the given rune is in the XML Character Range
func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}
