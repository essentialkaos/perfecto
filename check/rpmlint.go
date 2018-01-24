package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Lint run rpmlint and return alerts from it
func Lint(s *spec.Spec) []Alert {
	cmd := exec.Command("rpmlint", s.File)

	output, _ := cmd.Output()

	if len(output) == 0 {
		return nil
	}

	return parseRPMLintOutput(string(output), s)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// parseRPMLintOutput parse rpmlint output
func parseRPMLintOutput(output string, s *spec.Spec) []Alert {
	var result []Alert

	for _, line := range strings.Split(output, "\n") {
		alert, parsed := parseAlertLine(line, s)

		if !parsed {
			continue
		}

		result = append(result, alert)
	}

	return result
}

// parseAlertLine parse rpmlint output line
func parseAlertLine(text string, s *spec.Spec) (Alert, bool) {
	line, level, desc := extractAlertData(text)

	desc = strings.Replace(desc, "specfile-error error: ", "", -1)

	if strings.Contains(desc, "specfile-error warning") {
		level = "W"
		desc = strings.Replace(desc, "specfile-error warning: ", "", -1)
	}

	desc = "[rpmlint] " + desc

	switch level {
	case "W":
		return Alert{LEVEL_ERROR, desc, s.GetLine(line)}, true
	case "E":
		return Alert{LEVEL_CRITICAL, desc, s.GetLine(line)}, true
	}

	return Alert{}, false
}

// extractAlertData extract data from rpmlint alert
func extractAlertData(text string) (int, string, string) {
	lineSlice := strings.Split(text, ":")

	if len(lineSlice) < 3 {
		return -1, "", ""
	}

	if strings.HasPrefix(lineSlice[1], " ") {
		level := strings.TrimSpace(lineSlice[1])
		return -1, level, strings.TrimSpace(strings.Join(lineSlice[2:], ":"))
	}

	level := strings.TrimSpace(lineSlice[2])
	line, err := strconv.Atoi(lineSlice[1])

	if err != nil {
		return -1, "", ""
	}

	return line, level, strings.TrimSpace(strings.Join(lineSlice[3:], ":"))
}
