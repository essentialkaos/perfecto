package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"sort"

	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Alert levels
const (
	LEVEL_NOTICE uint8 = iota
	LEVEL_WARNING
	LEVEL_ERROR
	LEVEL_CRITICAL
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Report contain alerts
type Report struct {
	Notices   []Alert
	Warnings  []Alert
	Errors    []Alert
	Criticals []Alert
}

// Alert contain basic alert info
type Alert struct {
	Level uint8
	Info  string
	Line  spec.Line
}

// ////////////////////////////////////////////////////////////////////////////////// //

// AlertSlice is alerts slice
type AlertSlice []Alert

func (s AlertSlice) Len() int      { return len(s) }
func (s AlertSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s AlertSlice) Less(i, j int) bool {
	return s[i].Line.Index < s[j].Line.Index
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsPerfect return true if report doesn't have any alerts
func (r *Report) IsPerfect() bool {
	switch {
	case len(r.Notices) != 0:
		return false
	case len(r.Warnings) != 0:
		return false
	case len(r.Errors) != 0:
		return false
	case len(r.Criticals) != 0:
		return false
	}

	return true
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Check check spec
func Check(s *spec.Spec, lint bool, linterConfig string) *Report {
	report := &Report{}
	checkers := getCheckers()

	if lint {
		appendLinterAlerts(s, report, linterConfig)
	}

	for _, checker := range checkers {
		alerts := checker(s)

		if len(alerts) == 0 {
			continue
		}

		for _, alert := range alerts {
			switch alert.Level {
			case LEVEL_NOTICE:
				report.Notices = append(report.Notices, alert)
			case LEVEL_WARNING:
				report.Warnings = append(report.Warnings, alert)
			case LEVEL_ERROR:
				report.Errors = append(report.Errors, alert)
			}
		}
	}

	sort.Sort(AlertSlice(report.Notices))
	sort.Sort(AlertSlice(report.Warnings))
	sort.Sort(AlertSlice(report.Errors))
	sort.Sort(AlertSlice(report.Criticals))

	return report
}

// ////////////////////////////////////////////////////////////////////////////////// //

// appendLinterAlerts append rpmlint alerts to report
func appendLinterAlerts(s *spec.Spec, r *Report, linterConfig string) {
	alerts := Lint(s, linterConfig)

	if len(alerts) == 0 {
		return
	}

	for _, alert := range alerts {
		switch alert.Level {
		case LEVEL_ERROR:
			r.Errors = append(r.Errors, alert)
		case LEVEL_CRITICAL:
			r.Criticals = append(r.Criticals, alert)
		}
	}
}
