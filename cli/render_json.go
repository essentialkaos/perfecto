package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2021 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"encoding/json"
	"fmt"

	"github.com/essentialkaos/perfecto/check"
)

// renderJSONReport render report in JSON format
func renderJSONReport(r *check.Report) {
	data, err := json.MarshalIndent(r, "", "  ")

	if err != nil {
		printErrorAndExit(err.Error())
	}

	fmt.Println(string(data))
}
