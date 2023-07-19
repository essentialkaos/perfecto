package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type SpecSuite struct{}

var _ = Suite(&SpecSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *SpecSuite) TestFileCheck(c *C) {
	tmpDir := c.MkDir()
	tmpFile1 := tmpDir + "test1.spec"
	tmpFile2 := tmpDir + "test2.spec"

	os.WriteFile(tmpFile1, []byte(""), 0644)
	os.WriteFile(tmpFile2, []byte("TEST"), 0200)
}

func (s *SpecSuite) TestParsing(c *C) {
	spec, err := Read("../testdata/test1.spec")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = Read("../testdata/broken.spec")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = readFile("../testdata/_unknown_")

	c.Assert(err, NotNil)
	c.Assert(spec, IsNil)

	spec, err = Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.GetLine(-1), DeepEquals, Line{-1, "", false})
	c.Assert(spec.GetLine(99), DeepEquals, Line{-1, "", false})
	c.Assert(spec.GetLine(46), DeepEquals, Line{46, "%{__make} %{?_smp_mflags}", false})
}

func (s *SpecSuite) TestSections(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.HasSection(SECTION_BUILD), Equals, true)
	c.Assert(spec.HasSection(SECTION_PRETRANS), Equals, false)

	sections := spec.GetSections()
	c.Assert(sections, HasLen, 15)
	sections = spec.GetSections(SECTION_BUILD)
	c.Assert(sections, HasLen, 1)
	c.Assert(sections[0].Data, HasLen, 2)
	c.Assert(sections[0].Start, Equals, 45)
	c.Assert(sections[0].End, Equals, 47)
	c.Assert(sections[0].IsEmpty(), Equals, false)
	sections = spec.GetSections(SECTION_SETUP)
	c.Assert(sections[0].Name, Equals, "setup")
	c.Assert(sections[0].Args, DeepEquals, []string{"-qn", "%{name}-%{version}"})
	c.Assert(sections[0].Data, HasLen, 1)
	sections = spec.GetSections(SECTION_FILES)
	c.Assert(sections, HasLen, 2)
	c.Assert(sections[1].GetPackageName(), Equals, "magic")

	c.Assert(spec.Targets, DeepEquals, []string{"ubuntu", "el8", "el9", "@rhel"})

	spec, err = Read("../testdata/test_12.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	c.Assert(spec.HasSection(SECTION_CHECK), Equals, true)

	sections = spec.GetSections(SECTION_CHECK)

	c.Assert(sections, HasLen, 1)
	c.Assert(sections[0].IsEmpty(), Equals, true)
}

func (s *SpecSuite) TestHeaders(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	headers := spec.GetHeaders()
	c.Assert(headers, HasLen, 2)
	c.Assert(headers[0].Package, Equals, "")
	c.Assert(headers[0].IsSubpackage, Equals, false)
	c.Assert(headers[0].Data, HasLen, 16)
	c.Assert(headers[1].Package, Equals, "magic")
	c.Assert(headers[1].IsSubpackage, Equals, true)
	c.Assert(headers[1].Data, HasLen, 4)

	pkgName, subPkg := parsePackageName("%package magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, true)
	pkgName, subPkg = parsePackageName("%package -n magic")
	c.Assert(pkgName, Equals, "magic")
	c.Assert(subPkg, Equals, false)
}

func (s *SpecSuite) TestSourceExtractor(c *C) {
	spec, err := Read("../testdata/test.spec")

	c.Assert(err, IsNil)
	c.Assert(spec, NotNil)

	sources := spec.GetSources()

	c.Assert(sources, HasLen, 1)
}

func (s *SpecSuite) TestIgnoreDirective(c *C) {
	c.Assert(extractIgnoreCount("# perfecto:ignore"), Equals, 1)
	c.Assert(extractIgnoreCount("# perfecto:ignore ABC"), Equals, 0)
	c.Assert(extractIgnoreCount("# perfecto:ignore 1"), Equals, 1)
	c.Assert(extractIgnoreCount("# perfecto:ignore 10"), Equals, 10)
	c.Assert(extractIgnoreCount("# perfecto:absolve 10"), Equals, 10)
}

func (s *SpecSuite) TestTargetDirective(c *C) {
	c.Assert(extractTargets("# perfecto:target"), IsNil)
	c.Assert(extractTargets("# perfecto:target el7"), DeepEquals, []string{"el7"})
	c.Assert(extractTargets("# perfecto:target EL8"), DeepEquals, []string{"el8"})
	c.Assert(extractTargets("# perfecto:target EL8,EL9"), DeepEquals, []string{"el8", "el9"})
	c.Assert(extractTargets("# perfecto:target EL7 EL8"), DeepEquals, []string{"el7", "el8"})
}

func (s *SpecSuite) TestSectionPackageParsing(c *C) {
	section := Section{"test", []string{}, []Line{}, 0, 0}
	c.Assert(section.GetPackageName(), Equals, "")
	section = Section{"test", []string{"test1"}, []Line{}, 0, 0}
	c.Assert(section.GetPackageName(), Equals, "test1")
	section = Section{"test", []string{"-n", "test2"}, []Line{}, 0, 0}
	c.Assert(section.GetPackageName(), Equals, "test2")
}
