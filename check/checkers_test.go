package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2020 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	"github.com/essentialkaos/perfecto/spec"

	chk "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { chk.TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type CheckSuite struct{}

var _ = chk.Suite(&CheckSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (sc *CheckSuite) TestCheckForUselessSpaces(c *chk.C) {
	s, err := spec.Read("../testdata/test_1.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUselessSpaces("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Line contains spaces at the end of line")
	c.Assert(alerts[0].Line.Text, chk.Equals, "License:            MIT▒")
	c.Assert(alerts[1].Info, chk.Equals, "Line contains useless spaces")
	c.Assert(alerts[1].Line.Index, chk.Equals, 10)
}

func (sc *CheckSuite) TestCheckForLineLength(c *chk.C) {
	s, err := spec.Read("../testdata/test_1.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForLineLength("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Line is longer than 80 symbols")
	c.Assert(alerts[0].Line.Index, chk.Equals, 16)
	c.Assert(alerts[1].Line.Index, chk.Equals, 64)
}

func (sc *CheckSuite) TestCheckForDist(c *chk.C) {
	s, err := spec.Read("../testdata/test_1.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDist("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Release tag must contains %{?dist} as part of release")
	c.Assert(alerts[0].Line.Index, chk.Equals, 6)
}

func (sc *CheckSuite) TestCheckForNonMacroPaths(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForNonMacroPaths("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Path \"/usr\" should be used as macro \"%{_usr}\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 41)
	c.Assert(alerts[1].Info, chk.Equals, "Path \"/etc\" should be used as macro \"%{_sysconfdir}\"")
	c.Assert(alerts[1].Line.Index, chk.Equals, 42)
}

func (sc *CheckSuite) TestCheckForVariables(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForVariables("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Optimization flags must be used as macro %{optflags}")
	c.Assert(alerts[0].Line.Index, chk.Equals, 34)
	c.Assert(alerts[1].Info, chk.Equals, "Build root path must be used as macro %{buildroot}")
	c.Assert(alerts[1].Line.Index, chk.Equals, 41)
}

func (sc *CheckSuite) TestCheckForDevNull(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDevNull("", s)

	c.Assert(alerts, chk.HasLen, 5)
	c.Assert(alerts[0].Info, chk.Equals, "Use \"&>/dev/null || :\" instead of \">/dev/null 2>&1 || :\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 48)
	c.Assert(alerts[1].Info, chk.Equals, "Use \"&>/dev/null || :\" instead of \"2>&1 >/dev/null || :\"")
	c.Assert(alerts[1].Line.Index, chk.Equals, 49)
	c.Assert(alerts[2].Info, chk.Equals, "Use \"&>/dev/null || :\" instead of \">/dev/null 2>/dev/null || :\"")
	c.Assert(alerts[2].Line.Index, chk.Equals, 50)
	c.Assert(alerts[3].Info, chk.Equals, "Use \"&>/dev/null || :\" instead of \"2>/dev/null >/dev/null || :\"")
	c.Assert(alerts[3].Line.Index, chk.Equals, 51)
	c.Assert(alerts[4].Info, chk.Equals, "Use \" || :\" instead of \" || exit 0\"")
	c.Assert(alerts[4].Line.Index, chk.Equals, 51)
}

func (sc *CheckSuite) TestCheckChangelogHeaders(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkChangelogHeaders("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Changelog record header must contain release")
	c.Assert(alerts[0].Line.Index, chk.Equals, 74)
	c.Assert(alerts[1].Info, chk.Equals, "Misformatted changelog record header")
	c.Assert(alerts[1].Line.Index, chk.Equals, 77)
}

func (sc *CheckSuite) TestCheckForMakeMacro(c *chk.C) {
	s, err := spec.Read("../testdata/test_3.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForMakeMacro("", s)

	c.Assert(alerts, chk.HasLen, 3)
	c.Assert(alerts[0].Info, chk.Equals, "Use %{__make} macro instead of \"make\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 35)
	c.Assert(alerts[1].Info, chk.Equals, "Don't forget to use %{?_smp_mflags} macro with make command")
	c.Assert(alerts[1].Line.Index, chk.Equals, 35)
	c.Assert(alerts[2].Info, chk.Equals, "Use %{make_install} macro instead of \"make install\"")
	c.Assert(alerts[2].Line.Index, chk.Equals, 40)
}

func (sc *CheckSuite) TestCheckForHeaderTags(c *chk.C) {
	s, err := spec.Read("../testdata/test.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	c.Assert(checkForHeaderTags("", s), chk.HasLen, 0)

	s, err = spec.Read("../testdata/test_3.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForHeaderTags("", s)

	c.Assert(alerts, chk.HasLen, 3)
	c.Assert(alerts[0].Info, chk.Equals, "Main package must contain URL tag")
	c.Assert(alerts[1].Info, chk.Equals, "Main package must contain Group tag")
	c.Assert(alerts[2].Info, chk.Equals, "Package magic must contain Group tag")
}

func (sc *CheckSuite) TestCheckForUnescapedPercent(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUnescapedPercent("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Symbol % must be escaped by another % (i.e % → %%)")
	c.Assert(alerts[0].Line.Index, chk.Equals, 67)
}

func (sc *CheckSuite) TestCheckForMacroDefenitionPosition(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForMacroDefenitionPosition("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Move %define and %global to top of your spec")
	c.Assert(alerts[0].Line.Index, chk.Equals, 35)
}

func (sc *CheckSuite) TestCheckForSeparatorLength(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForSeparatorLength("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Separator must be 80 symbols long")
	c.Assert(alerts[0].Line.Index, chk.Equals, 63)
}

func (sc *CheckSuite) TestCheckForDefAttr(c *chk.C) {
	s, err := spec.Read("../testdata/test_5.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDefAttr("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "%files section must contains %defattr macro")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)
	c.Assert(alerts[1].Info, chk.Equals, "%files section for package magic must contains %defattr macro")
	c.Assert(alerts[1].Line.Index, chk.Equals, -1)
}

func (sc *CheckSuite) TestCheckForUselessBinaryMacro(c *chk.C) {
	s, err := spec.Read("../testdata/test_5.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUselessBinaryMacro("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Useless macro %{__rm} used for executing rm binary")
	c.Assert(alerts[0].Line.Index, chk.Equals, 47)
}

func (sc *CheckSuite) TestCheckForEmptySections(c *chk.C) {
	s, err := spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForEmptySections("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Section %check is empty")
	c.Assert(alerts[0].Line.Index, chk.Equals, 45)
}

func (sc *CheckSuite) TestCheckForIndentInFilesSection(c *chk.C) {
	s, err := spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForIndentInFilesSection("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Don't use indent in %files section")
	c.Assert(alerts[0].Line.Index, chk.Equals, 66)
}

func (sc *CheckSuite) TestCheckForSetupArguments(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForSetupOptions("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Options \"-q -c -n\" can be simplified to \"-qcn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 33)

	s, err = spec.Read("../testdata/test_5.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForSetupOptions("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Options \"-c -n\" can be simplified to \"-cn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 41)

	s, err = spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForSetupOptions("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Options \"-q -n\" can be simplified to \"-qn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 31)
}

func (sc *CheckSuite) TestCheckForEmptyLinesAtEnd(c *chk.C) {
	s, err := spec.Read("../testdata/test_8.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForEmptyLinesAtEnd("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Spec file should have empty line at the end")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)

	s, err = spec.Read("../testdata/test_9.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForEmptyLinesAtEnd("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Too much empty lines at the end of the spec")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)
}

func (sc *CheckSuite) TestCheckBashLoops(c *chk.C) {
	s, err := spec.Read("../testdata/test_10.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkBashLoops("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Place 'do' keyword on the same line with for/while (for ... ; do)")
	c.Assert(alerts[0].Line.Index, chk.Equals, 37)
	c.Assert(alerts[1].Info, chk.Equals, "Place 'do' keyword on the same line with for/while (for ... ; do)")
	c.Assert(alerts[1].Line.Index, chk.Equals, 49)
}

func (sc *CheckSuite) TestCheckURLForHTTPS(c *chk.C) {
	s, err := spec.Read("../testdata/test_11.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkURLForHTTPS("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Domain kaos.st supports HTTPS. Replace http by https in URL.")
	c.Assert(alerts[0].Line.Index, chk.Equals, 13)
}

func (sc *CheckSuite) TestCheckForCheckMacro(c *chk.C) {
	s, err := spec.Read("../testdata/test_11.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForCheckMacro("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Use %{_without_check} and %{_with_check} macroses for controlling tests execution")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)

	s, err = spec.Read("../testdata/test_12.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForCheckMacro("", s)

	c.Assert(alerts, chk.HasLen, 0)
}

func (sc *CheckSuite) TestCheckIfClause(c *chk.C) {
	s, err := spec.Read("../testdata/test_13.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkIfClause("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Use two equals symbols for comparison in %if clause")
	c.Assert(alerts[0].Line.Index, chk.Equals, 55)
}

func (sc *CheckSuite) TestCheckForUselessSlash(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUselessSlash("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Slash between %{buildroot} and %{_usr} macroses is useless")
	c.Assert(alerts[0].Line.Index, chk.Equals, 48)
}

func (sc *CheckSuite) TestCheckForEmptyIf(c *chk.C) {
	s, err := spec.Read("../testdata/test_14.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForEmptyIf("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Evaluated if clause can be empty. Change the order of clauses (i.e. %if → if instead of if → %if).")
	c.Assert(alerts[0].Line.Index, chk.Equals, 92)
}

func (sc *CheckSuite) TestCheckForDotInSummary(c *chk.C) {
	s, err := spec.Read("../testdata/test_14.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDotInSummary("", s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "The summary contains useless dot at the end")
	c.Assert(alerts[0].Line.Index, chk.Equals, 7)
}

func (sc *CheckSuite) TestCheckForChownAndChmod(c *chk.C) {
	s, err := spec.Read("../testdata/test_15.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForChownAndChmod("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Do not change file or directory mode in scriptlets")
	c.Assert(alerts[0].Line.Index, chk.Equals, 60)
	c.Assert(alerts[1].Info, chk.Equals, "Do not change file or directory owner without --no-dereference option")
	c.Assert(alerts[1].Line.Index, chk.Equals, 61)
}

func (sc *CheckSuite) TestCheckForUnclosedCondition(c *chk.C) {
	s, err := spec.Read("../testdata/test_16.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUnclosedCondition("", s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Scriptlet contains unclosed IF condition")
	c.Assert(alerts[0].Line.Index, chk.Equals, 70)
	c.Assert(alerts[1].Info, chk.Equals, "Scriptlet contains unclosed IF condition")
	c.Assert(alerts[1].Line.Index, chk.Equals, 71)
}

func (sc *CheckSuite) TestWithEmptyData(c *chk.C) {
	s := &spec.Spec{}

	for _, checker := range getCheckers() {
		c.Assert(checker("", s), chk.IsNil)
	}
}

func (sc *CheckSuite) TestRPMLint(c *chk.C) {
	s, err := spec.Read("../testdata/test.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	r := Check(s, true, "", nil)

	c.Assert(r, chk.NotNil)
	c.Assert(r.IsPerfect(), chk.Equals, true)
	c.Assert(r.IDs(), chk.HasLen, 0)

	s, err = spec.Read("../testdata/test_7.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	r = Check(s, true, "", nil)

	c.Assert(r, chk.NotNil)
	c.Assert(r.IsPerfect(), chk.Equals, false)

	s, err = spec.Read("../testdata/test_11.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	r = Check(s, true, "", []string{"PF20", "PF21"})

	c.Assert(r, chk.NotNil)
	c.Assert(r.Warnings, chk.HasLen, 3)
	c.Assert(r.Warnings[0].Absolve, chk.Equals, true)

	rpmLintBin = "echo"
	s = &spec.Spec{File: ""}
	c.Assert(Lint(s, ""), chk.IsNil)

	s = &spec.Spec{File: "test.spec"}
	c.Assert(Lint(s, "test.conf"), chk.IsNil)
}

func (sc *CheckSuite) TestRPMLintParser(c *chk.C) {
	report := &Report{}
	alerts := []Alert{}

	s, err := spec.Read("../testdata/test_7.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	a, ok := parseAlertLine("test.spec: W: no-buildroot-tag", s)

	c.Assert(ok, chk.Equals, true)
	c.Assert(a.ID, chk.Equals, "")
	c.Assert(a.Level, chk.Equals, LEVEL_ERROR)
	c.Assert(a.Info, chk.Equals, "no-buildroot-tag")
	c.Assert(a.Line.Index, chk.Equals, -1)
	alerts = append(alerts, a)

	a, ok = parseAlertLine("test.spec: E: specfile-error error: line 10: Unknown tag: Release1", s)

	c.Assert(ok, chk.Equals, true)
	c.Assert(a.ID, chk.Equals, "")
	c.Assert(a.Level, chk.Equals, LEVEL_CRITICAL)
	c.Assert(a.Info, chk.Equals, "Unknown tag: Release1")
	c.Assert(a.Line.Index, chk.Equals, 10)
	alerts = append(alerts, a)

	a, ok = parseAlertLine("test.spec:67: W: macro-in-%changelog %record", s)

	c.Assert(ok, chk.Equals, true)
	c.Assert(a.ID, chk.Equals, "")
	c.Assert(a.Level, chk.Equals, LEVEL_ERROR)
	c.Assert(a.Info, chk.Equals, "macro-in-%changelog %record")
	c.Assert(a.Line.Index, chk.Equals, 67)
	alerts = append(alerts, a)

	a, ok = parseAlertLine("test.spec: E: specfile-error error: line A: Unknown tag: Release1", s)

	c.Assert(ok, chk.Equals, false)

	a, ok = parseAlertLine("test.spec:A: W: macro-in-%changelog %record", s)

	c.Assert(ok, chk.Equals, false)

	appendLinterAlerts(report, alerts)

	c.Assert(report.Errors, chk.HasLen, 2)
	c.Assert(report.Criticals, chk.HasLen, 1)
}

func (sc *CheckSuite) TestAux(c *chk.C) {
	// This test will fail if new checkers was added
	c.Assert(getCheckers(), chk.HasLen, 27)

	r := &Report{}
	c.Assert(r.IsPerfect(), chk.Equals, true)
	r = &Report{Notices: []Alert{Alert{}}}
	c.Assert(r.IsPerfect(), chk.Equals, false)
	r = &Report{Warnings: []Alert{Alert{}}}
	c.Assert(r.IsPerfect(), chk.Equals, false)
	r = &Report{Errors: []Alert{Alert{}}}
	c.Assert(r.IsPerfect(), chk.Equals, false)
	r = &Report{Criticals: []Alert{Alert{}}}
	c.Assert(r.IsPerfect(), chk.Equals, false)

	r = &Report{
		Notices:   []Alert{Alert{}},
		Warnings:  []Alert{Alert{ID: "PF0"}},
		Errors:    []Alert{Alert{ID: "PF0"}},
		Criticals: []Alert{Alert{ID: "PF0"}},
	}

	c.Assert(r.IDs(), chk.HasLen, 1)

	a := AlertSlice{Alert{}, Alert{}}
	a.Swap(0, 1)
	c.Assert(a.Len(), chk.Equals, 2)
	c.Assert(a.Less(0, 1), chk.Equals, false)

	al, _ := parseAlertLine("../testdata/test_7.spec: E: specfile-error warning: some error", &spec.Spec{})
	c.Assert(al.Level, chk.Equals, LEVEL_ERROR)
	c.Assert(al.Info, chk.Equals, "some error")
}
