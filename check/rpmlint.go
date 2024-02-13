package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// RPMLINT_CHECK_ID is ID for all rpmlint checks
const RPMLINT_CHECK_ID = "LNT0"

// ////////////////////////////////////////////////////////////////////////////////// //

var rpmLintBin = "rpmlint"

// ////////////////////////////////////////////////////////////////////////////////// //

// Lint run rpmlint and return alerts from it
func Lint(s *spec.Spec, linterConfig string) []Alert {
	if env.Which(rpmLintBin) == "" {
		return nil // RPMLint not installed
	}

	cmd := exec.Command(rpmLintBin)

	if linterConfig != "" {
		cmd.Args = append(cmd.Args, "-f", linterConfig)
	}

	cmd.Args = append(cmd.Args, s.File)

	output, _ := cmd.Output()

	if len(output) < 2 {
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

	if strings.Contains(desc, "specfile-error warning") {
		level = "W"
		desc = strutil.Exclude(desc, "specfile-error warning: ")
	}

	switch level {
	case "W":
		return NewAlert(RPMLINT_CHECK_ID, LEVEL_ERROR, desc, s.GetLine(line)), true
	case "E":
		return NewAlert(RPMLINT_CHECK_ID, LEVEL_CRITICAL, desc, s.GetLine(line)), true
	}

	return Alert{}, false
}

// extractAlertData extract data from rpmlint alert
func extractAlertData(text string) (int, string, string) {
	if strings.Count(text, ":") < 2 {
		return -1, "", ""
	}

	lineSlice := strings.Split(text, ":")

	// Alert with error type and line number in text of alert
	if strings.Contains(text, "specfile-error error: line ") && len(lineSlice) > 4 {
		line, err := strconv.Atoi(strings.Trim(lineSlice[3], "line "))

		if err != nil {
			return -1, "", ""
		}

		return line, "E", strings.TrimSpace(strings.Join(lineSlice[4:], ":"))
	}

	// Alert with type and without line number
	if strings.HasPrefix(lineSlice[1], " ") {
		level := strings.TrimSpace(lineSlice[1])
		return -1, level, strings.TrimSpace(strings.Join(lineSlice[2:], ":"))
	}

	// Alert with type and line number
	level := strings.TrimSpace(lineSlice[2])
	line, err := strconv.Atoi(lineSlice[1])

	if err != nil {
		return -1, "", ""
	}

	return line, level, strings.TrimSpace(strings.Join(lineSlice[3:], ":"))
}
