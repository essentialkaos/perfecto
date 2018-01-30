package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type SpecSuite struct{}

var _ = Suite(&SpecSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *SpecSuite) TestParsing(c *C) {
	spec, err := Read("../testdata/test1.spec")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = Read("../testdata/broken.spec")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.GetFileName(), Equals, "test")

	c.Assert(spec.GetLine(-1), DeepEquals, Line{-1, "", false})
	c.Assert(spec.GetLine(99), DeepEquals, Line{-1, "", false})
	c.Assert(spec.GetLine(34), DeepEquals, Line{34, "%{__make} %{?_smp_mflags}", false})
}

func (s *SpecSuite) TestSections(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.HasSection(SECTION_BUILD), Equals, true)
	c.Assert(spec.HasSection(SECTION_CHECK), Equals, false)

	sections := spec.GetSections()
	c.Assert(sections, HasLen, 13)
	sections = spec.GetSections(SECTION_BUILD)
	c.Assert(sections, HasLen, 1)
	c.Assert(sections[0].Data, HasLen, 1)
	sections = spec.GetSections(SECTION_SETUP)
	c.Assert(sections[0].Name, Equals, "setup")
	c.Assert(sections[0].Args, DeepEquals, []string{"-qn", "%{name}-%{version}"})
	c.Assert(sections[0].Data, HasLen, 0)
}

func (s *SpecSuite) TestHeaders(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	headers := spec.GetHeaders()
	c.Assert(headers, HasLen, 2)
	c.Assert(headers[0].Package, Equals, "")
	c.Assert(headers[0].Subpackage, Equals, false)
	c.Assert(headers[0].Data, HasLen, 11)
	c.Assert(headers[1].Package, Equals, "magic")
	c.Assert(headers[1].Subpackage, Equals, true)
	c.Assert(headers[1].Data, HasLen, 4)

	pkgName, subPkg := parsePackageName("%package magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, true)
	pkgName, subPkg = parsePackageName("%package -n magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, false)
}

func (s *SpecSuite) TestSkipTag(c *C) {
	c.Assert(isSkipTag("# perfecto:absolve 3"), Equals, true)
	c.Assert(isSkipTag("# abcd 1"), Equals, false)

	c.Assert(extractSkipCount("# perfecto:absolve"), Equals, 1)
	c.Assert(extractSkipCount("# perfecto:absolve ABC"), Equals, 0)
	c.Assert(extractSkipCount("# perfecto:absolve 1"), Equals, 1)
	c.Assert(extractSkipCount("# perfecto:absolve 10"), Equals, 10)
}
