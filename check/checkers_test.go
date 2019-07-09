package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
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

	alerts := checkForUselessSpaces(s)

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

	alerts := checkForLineLength(s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Line is longer than 80 symbols")
	c.Assert(alerts[0].Line.Index, chk.Equals, 16)
	c.Assert(alerts[1].Line.Index, chk.Equals, 64)
}

func (sc *CheckSuite) TestCheckForDist(c *chk.C) {
	s, err := spec.Read("../testdata/test_1.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDist(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Release tag must contains %{?dist} as part of release")
	c.Assert(alerts[0].Line.Index, chk.Equals, 6)
}

func (sc *CheckSuite) TestCheckForNonMacroPaths(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForNonMacroPaths(s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Path \"/usr\" should be used as macro \"%{_usr}\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 41)
	c.Assert(alerts[1].Info, chk.Equals, "Path \"/etc\" should be used as macro \"%{_sysconfdir}\"")
	c.Assert(alerts[1].Line.Index, chk.Equals, 42)
}

func (sc *CheckSuite) TestCheckForBuildRoot(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForBuildRoot(s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Build root path must be used as macro %{buildroot}")
	c.Assert(alerts[0].Line.Index, chk.Equals, 41)
	c.Assert(alerts[1].Info, chk.Equals, "Slash after %{buildroot} macro is useless")
	c.Assert(alerts[1].Line.Index, chk.Equals, 48)
}

func (sc *CheckSuite) TestCheckForDevNull(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDevNull(s)

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

	alerts := checkChangelogHeaders(s)

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

	alerts := checkForMakeMacro(s)

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

	c.Assert(checkForHeaderTags(s), chk.HasLen, 0)

	s, err = spec.Read("../testdata/test_3.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForHeaderTags(s)

	c.Assert(alerts, chk.HasLen, 3)
	c.Assert(alerts[0].Info, chk.Equals, "Main package must contain URL tag")
	c.Assert(alerts[1].Info, chk.Equals, "Main package must contain Group tag")
	c.Assert(alerts[2].Info, chk.Equals, "Package magic must contain Group tag")
}

func (sc *CheckSuite) TestCheckForUnescapedPercent(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForUnescapedPercent(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Symbol % must be escaped by another % (i.e % → %%)")
	c.Assert(alerts[0].Line.Index, chk.Equals, 67)
}

func (sc *CheckSuite) TestCheckForMacroDefenitionPosition(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForMacroDefenitionPosition(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Move %define and %global to top of your spec")
	c.Assert(alerts[0].Line.Index, chk.Equals, 35)
}

func (sc *CheckSuite) TestCheckForSeparatorLength(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForSeparatorLength(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Separator must be 80 symbols long")
	c.Assert(alerts[0].Line.Index, chk.Equals, 63)
}

func (sc *CheckSuite) TestCheckForDefAttr(c *chk.C) {
	s, err := spec.Read("../testdata/test_5.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDefAttr(s)

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

	alerts := checkForUselessBinaryMacro(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Useless macro %{__rm} used for executing rm binary")
	c.Assert(alerts[0].Line.Index, chk.Equals, 47)
}

func (sc *CheckSuite) TestCheckForEmptySections(c *chk.C) {
	s, err := spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForEmptySections(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Section %check is empty")
	c.Assert(alerts[0].Line.Index, chk.Equals, 45)
}

func (sc *CheckSuite) TestCheckForIndentInFilesSection(c *chk.C) {
	s, err := spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForIndentInFilesSection(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Don't use indent in %files section")
	c.Assert(alerts[0].Line.Index, chk.Equals, 66)
}

func (sc *CheckSuite) TestCheckForSetupArguments(c *chk.C) {
	s, err := spec.Read("../testdata/test_4.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForSetupArguments(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Arguments \"-q -c -n\" can be simplified to \"-qcn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 33)

	s, err = spec.Read("../testdata/test_5.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForSetupArguments(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Arguments \"-c -n\" can be simplified to \"-cn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 41)

	s, err = spec.Read("../testdata/test_6.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForSetupArguments(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Arguments \"-q -n\" can be simplified to \"-qn\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 31)
}

func (sc *CheckSuite) TestCheckForEmptyLinesAtEnd(c *chk.C) {
	s, err := spec.Read("../testdata/test_8.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForEmptyLinesAtEnd(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Spec file should have empty line at the end")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)

	s, err = spec.Read("../testdata/test_9.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts = checkForEmptyLinesAtEnd(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Too much empty lines at the end of the spec")
	c.Assert(alerts[0].Line.Index, chk.Equals, -1)
}

func (sc *CheckSuite) TestCheckBashLoops(c *chk.C) {
	s, err := spec.Read("../testdata/test_10.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkBashLoops(s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Place 'do' keyword on the same line with for/while (for ... ; do)")
	c.Assert(alerts[0].Line.Index, chk.Equals, 37)
	c.Assert(alerts[1].Info, chk.Equals, "Place 'do' keyword on the same line with for/while (for ... ; do)")
	c.Assert(alerts[1].Line.Index, chk.Equals, 49)
}

func (sc *CheckSuite) TestWithEmptyData(c *chk.C) {
	s := &spec.Spec{}

	c.Assert(checkForUselessSpaces(s), chk.IsNil)
	c.Assert(checkForLineLength(s), chk.IsNil)
	c.Assert(checkForDist(s), chk.IsNil)
	c.Assert(checkForNonMacroPaths(s), chk.IsNil)
	c.Assert(checkForBuildRoot(s), chk.IsNil)
	c.Assert(checkForDevNull(s), chk.IsNil)
	c.Assert(checkChangelogHeaders(s), chk.IsNil)
	c.Assert(checkForMakeMacro(s), chk.IsNil)
	c.Assert(checkForHeaderTags(s), chk.IsNil)
	c.Assert(checkForUnescapedPercent(s), chk.IsNil)
	c.Assert(checkForMacroDefenitionPosition(s), chk.IsNil)
	c.Assert(checkForSeparatorLength(s), chk.IsNil)
	c.Assert(checkForDefAttr(s), chk.IsNil)
	c.Assert(checkForUselessBinaryMacro(s), chk.IsNil)
	c.Assert(checkForEmptySections(s), chk.IsNil)
	c.Assert(checkForIndentInFilesSection(s), chk.IsNil)
	c.Assert(checkForSetupArguments(s), chk.IsNil)
	c.Assert(checkForEmptyLinesAtEnd(s), chk.IsNil)
}

func (sc *CheckSuite) TestRPMLint(c *chk.C) {
	s, err := spec.Read("../testdata/test.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	r := Check(s, true, "")

	c.Assert(r, chk.NotNil)
	c.Assert(r.IsPerfect(), chk.Equals, true)

	s, err = spec.Read("../testdata/test_7.spec")

	c.Assert(s, chk.NotNil)
	c.Assert(err, chk.IsNil)

	r = Check(s, true, "")

	c.Assert(r, chk.NotNil)
	c.Assert(r.IsPerfect(), chk.Equals, false)

	rpmLintBin = "echo"
	s = &spec.Spec{File: ""}
	c.Assert(Lint(s, ""), chk.IsNil)

	s = &spec.Spec{File: "test.spec"}
	c.Assert(Lint(s, "test.conf"), chk.IsNil)
}

func (sc *CheckSuite) TestRPMLintParser(c *chk.C) {
	i, s1, s2 := extractAlertData("test.spec: W: no-buildroot-tag")
	c.Assert(i, chk.Equals, -1)
	c.Assert(s1, chk.Equals, "W")
	c.Assert(s2, chk.Equals, "no-buildroot-tag")

	i, s1, s2 = extractAlertData("test.spec: E: specfile-error error: line 356: Unknown tag: Release1")
	c.Assert(i, chk.Equals, 356)
	c.Assert(s1, chk.Equals, "E")
	c.Assert(s2, chk.Equals, "Unknown tag: Release1")

	i, s1, s2 = extractAlertData("test.spec:67: W: macro-in-%changelog %record")
	c.Assert(i, chk.Equals, 67)
	c.Assert(s1, chk.Equals, "W")
	c.Assert(s2, chk.Equals, "macro-in-%changelog %record")

	i, s1, s2 = extractAlertData("test.spec: E: specfile-error error: line A: Unknown tag: Release1")
	c.Assert(i, chk.Equals, -1)
	c.Assert(s1, chk.Equals, "")
	c.Assert(s2, chk.Equals, "")

	i, s1, s2 = extractAlertData("test.spec:A: W: macro-in-%changelog %record")
	c.Assert(i, chk.Equals, -1)
	c.Assert(s1, chk.Equals, "")
	c.Assert(s2, chk.Equals, "")
}

func (sc *CheckSuite) TestAux(c *chk.C) {
	// This test will fail if new checkers was added
	c.Assert(getCheckers(), chk.HasLen, 19)

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

	a := AlertSlice{Alert{}, Alert{}}
	a.Swap(0, 1)
	c.Assert(a.Len(), chk.Equals, 2)
	c.Assert(a.Less(0, 1), chk.Equals, false)

	al, _ := parseAlertLine("../testdata/test_7.spec: E: specfile-error warning: some error", &spec.Spec{})
	c.Assert(al.Level, chk.Equals, LEVEL_ERROR)
	c.Assert(al.Info, chk.Equals, "[rpmlint] some error")
}
