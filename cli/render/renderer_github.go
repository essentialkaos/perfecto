package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/path"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// GithubRenderer renders report using github actions workflow commands
type GithubRenderer struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

// Report renders alerts from perfecto report
func (r *GithubRenderer) Report(file string, report *check.Report) {
	if report.Notices.Total() != 0 {
		r.renderActionAlerts("notice", file, report.Notices)
	}

	if report.Warnings.Total() != 0 {
		r.renderActionAlerts("warning", file, report.Warnings)
	}

	if report.Errors.Total() != 0 {
		r.renderActionAlerts("error", file, report.Errors)
	}

	if report.Criticals.Total() != 0 {
		r.renderActionAlerts("error", file, report.Criticals)
	}
}

// Perfect renders message about perfect spec
func (r *GithubRenderer) Perfect(file string, report *check.Report) {
	specName := strutil.Exclude(path.Base(file), ".spec")
	fmtc.Printf("{g}{*}%s.spec{!*} is perfect!{!}\n", specName)
}

// Skipped renders message about skipped check
func (r *GithubRenderer) Skipped(file string, report *check.Report) {
	specName := strutil.Exclude(path.Base(file), ".spec")
	fmtc.Printf("{s}{*}%s.spec{!*} check skipped due to non-applicable target{!}\n", specName)
}

// Error renders global error message
func (r *GithubRenderer) Error(file string, err error) {
	fmt.Printf("::error file=%s::%v\n", file, err)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// renderActionAlerts renders alert with control commands
func (r *GithubRenderer) renderActionAlerts(level, file string, alerts []check.Alert) {
	for _, alert := range alerts {
		title := "Global"

		if alert.ID != "" {
			title = alert.ID
		}

		if alert.Line.Index == -1 {
			fmt.Printf(
				"::%s file=%s,title=%s::%s\n",
				level, file, title, alert.Info,
			)
		} else {
			fmt.Printf(
				"::%s file=%s,line=%d,title=%s::%s\n",
				level, file, alert.Line.Index, title, alert.Info,
			)
		}
	}
}
