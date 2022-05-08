package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
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
func (r *JSONRenderer) Report(file string, report *check.Report) error {
	data, err := json.MarshalIndent(report, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}

// Perfect renders message about perfect spec
func (r *JSONRenderer) Perfect(file string) {
	fmt.Println("{}")
}

// Error renders global error message
func (r *JSONRenderer) Error(file string, err error) {
	fmt.Printf("{\"error\":\"%v\"}\n", err)
}

// ////////////////////////////////////////////////////////////////////////////////// //
