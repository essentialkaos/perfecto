package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
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
