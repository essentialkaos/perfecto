package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
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
	c.Assert(alerts[1].Line.Index, chk.Equals, 46)
}

func (sc *CheckSuite) TestCheckForDevNull(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForDevNull(s)

	c.Assert(alerts, chk.HasLen, 1)
	c.Assert(alerts[0].Info, chk.Equals, "Use \"&>/dev/null || :\" instead of \">/dev/null 2>&1 || :\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 46)
}

func (sc *CheckSuite) TestCheckChangelogHeaders(c *chk.C) {
	s, err := spec.Read("../testdata/test_2.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkChangelogHeaders(s)

	c.Assert(alerts, chk.HasLen, 2)
	c.Assert(alerts[0].Info, chk.Equals, "Changelog record header must contain release")
	c.Assert(alerts[0].Line.Index, chk.Equals, 69)
	c.Assert(alerts[1].Info, chk.Equals, "Misformatted changelog record header")
	c.Assert(alerts[1].Line.Index, chk.Equals, 72)
}

func (sc *CheckSuite) TestCheckForMakeMacro(c *chk.C) {
	s, err := spec.Read("../testdata/test_3.spec")

	c.Assert(err, chk.IsNil)
	c.Assert(s, chk.NotNil)

	alerts := checkForMakeMacro(s)

	c.Assert(alerts, chk.HasLen, 3)
	c.Assert(alerts[0].Info, chk.Equals, "Use %{__make} macro instead of \"make\"")
	c.Assert(alerts[0].Line.Index, chk.Equals, 34)
	c.Assert(alerts[1].Info, chk.Equals, "Don't forget to use %{?_smp_mflags} macro with make command")
	c.Assert(alerts[1].Line.Index, chk.Equals, 34)
	c.Assert(alerts[2].Info, chk.Equals, "Use %{make_install} macro instead of \"make install\"")
	c.Assert(alerts[2].Line.Index, chk.Equals, 39)
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
}

func (sc *CheckSuite) TestAux(c *chk.C) {
	c.Assert(getCheckers(), chk.HasLen, 13)
}
