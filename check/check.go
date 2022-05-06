package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"sort"

	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/sortutil"

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

// Report contains info about all alerts
type Report struct {
	Notices   Alerts `json:"notices,omitempty"`
	Warnings  Alerts `json:"warnings,omitempty"`
	Errors    Alerts `json:"errors,omitempty"`
	Criticals Alerts `json:"criticals,omitempty"`
}

// Alert contains basic alert info
type Alert struct {
	ID      string    `json:"id"`
	Level   uint8     `json:"level"`
	Info    string    `json:"info"`
	Line    spec.Line `json:"line"`
	Absolve bool      `json:"absolve"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Alerts is slice with alerts
type Alerts []Alert

func (s Alerts) Len() int      { return len(s) }
func (s Alerts) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Alerts) Less(i, j int) bool {
	return s[i].Line.Index < s[j].Line.Index
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewAlert creates new alert
func NewAlert(id string, level uint8, info string, line spec.Line) Alert {
	return Alert{id, level, info, line, false}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsPerfect returns true if report doesn't have any alerts
func (r *Report) IsPerfect() bool {
	return r.Total() == 0
}

// Total returns total number of alerts (including absolved)
func (r *Report) Total() int {
	return r.Notices.Total() +
		r.Warnings.Total() +
		r.Errors.Total() +
		r.Criticals.Total()
}

// Absolved returns number of absolved (skipped) alerts
func (r *Report) Absolved() int {
	return r.Notices.Absolved() +
		r.Warnings.Absolved() +
		r.Errors.Absolved() +
		r.Criticals.Absolved()
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

// HasAlerts returns true if alerts contains at least one non-absolved alert
func (a Alerts) HasAlerts() bool {
	for _, alert := range a {
		if alert.Absolve {
			continue
		}

		return true
	}

	return false
}

// Absolved returns total number of alerts
func (a Alerts) Total() int {
	return len(a)
}

// Absolved returns number of absolved (skipped) alerts
func (a Alerts) Absolved() int {
	var counter int

	for _, alert := range a {
		if alert.Absolve {
			counter++
		}
	}

	return counter
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
			case LEVEL_CRITICAL:
				report.Criticals = append(report.Criticals, alert)
			}
		}
	}

	sort.Sort(Alerts(report.Notices))
	sort.Sort(Alerts(report.Warnings))
	sort.Sort(Alerts(report.Errors))
	sort.Sort(Alerts(report.Criticals))

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
