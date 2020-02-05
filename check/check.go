package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"sort"

	"pkg.re/essentialkaos/ek.v11/sliceutil"
	"pkg.re/essentialkaos/ek.v11/sortutil"

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
	Notices   []Alert `json:"notices,omitempty"`
	Warnings  []Alert `json:"warnings,omitempty"`
	Errors    []Alert `json:"errors,omitempty"`
	Criticals []Alert `json:"criticals,omitempty"`
}

// Alert contain basic alert info
type Alert struct {
	ID      string    `json:"id"`
	Level   uint8     `json:"level"`
	Info    string    `json:"info"`
	Line    spec.Line `json:"line"`
	Absolve bool      `json:"absolve"`
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

// NewAlert creates new alert
func NewAlert(id string, level uint8, info string, line spec.Line) Alert {
	return Alert{id, level, info, line, false}
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

// IDs returns slice with all mentioned checks ID's
func (r *Report) IDs() []string {
	ids := make(map[string]bool)

	for _, a := range r.Notices {
		ids[a.ID] = true
	}

	for _, a := range r.Warnings {
		ids[a.ID] = true
	}

	for _, a := range r.Errors {
		ids[a.ID] = true
	}

	for _, a := range r.Criticals {
		ids[a.ID] = true
	}

	if len(ids) == 0 {
		return nil
	}

	var result []string

	for id := range ids {
		if id == "" {
			continue
		}

		result = append(result, id)
	}

	sortutil.StringsNatural(result)

	return result
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Check check spec
func Check(s *spec.Spec, lint bool, linterConfig string, absolved []string) *Report {
	report := &Report{}
	checkers := getCheckers()

	if lint {
		alerts := Lint(s, linterConfig)
		appendLinterAlerts(report, alerts)
	}

	for id, checker := range checkers {
		alerts := checker(id, s)

		if len(alerts) == 0 {
			continue
		}

		absolve := sliceutil.Contains(absolved, id)

		for _, alert := range alerts {
			if absolve || alert.Line.Skip {
				alert.Absolve = true
			}

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
func appendLinterAlerts(r *Report, alerts []Alert) {
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
