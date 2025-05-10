package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"fmt"

	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// JSONRenderer renders report in JSON format
type JSONRenderer struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

// Report renders alerts from perfecto report
func (r *JSONRenderer) Report(file string, report *check.Report) {
	encodeReport(report)
}

// Perfect renders message about perfect spec
func (r *JSONRenderer) Perfect(file string, report *check.Report) {
	encodeReport(report)
}

// Skipped renders message about skipped check
func (r *JSONRenderer) Skipped(file string, report *check.Report) {
	encodeReport(report)
}

// Error renders global error message
func (r *JSONRenderer) Error(file string, err error) {
	fmt.Printf("{\"error\":\"%v\"}\n", err)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func encodeReport(report *check.Report) {
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
}
