package render

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/perfecto/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Renderer is interface for perfecto data
type Renderer interface {

	// Report renders alerts from perfecto report
	Report(file string, report *check.Report) error

	// Perfect renders message about perfect spec
	Perfect(file string)

	// Error renders global error message
	Error(file string, err error)
}

// ////////////////////////////////////////////////////////////////////////////////// //
