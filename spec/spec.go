package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Sections
const (
	SECTION_BUILD         = "build"
	SECTION_CHANGELOG     = "changelog"
	SECTION_CHECK         = "check"
	SECTION_CLEAN         = "clean"
	SECTION_DESCRIPTION   = "description"
	SECTION_FILES         = "files"
	SECTION_INSTALL       = "install"
	SECTION_PACKAGE       = "package"
	SECTION_POST          = "post"
	SECTION_POSTTRANS     = "posttrans"
	SECTION_POSTUN        = "postun"
	SECTION_PRE           = "pre"
	SECTION_PREP          = "prep"
	SECTION_PRETRANS      = "pretrans"
	SECTION_PREUN         = "preun"
	SECTION_SETUP         = "setup"
	SECTION_TRIGGERIN     = "triggerin"
	SECTION_TRIGGERPOSTUN = "triggerpostun"
	SECTION_TRIGGERUN     = "triggerun"
	SECTION_VERIFYSCRIPT  = "verifyscript"
)

const (
	DIRECTIVE_IGNORE = "perfecto:ignore"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Spec spec contains data from spec file
type Spec struct {
	File    string   `json:"file"`
	Data    []Line   `json:"data"`
	Targets []string `json:"target"`
}

// Line contains line data and index
type Line struct {
	Index  int    `json:"index"`
	Text   string `json:"text"`
	Ignore bool   `json:"ignore"`
}

// Header header contains header info and data
type Header struct {
	Package      string `json:"package"`
	Data         []Line `json:"data"`
	IsSubpackage bool   `json:"is_subpackage"`
}

// Section contains section info and data
type Section struct {
	Name  string   `json:"name"`
	Args  []string `json:"args"`
	Data  []Line   `json:"data"`
	Start int      `json:"start"`
	End   int      `json:"end"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// sections is slice with rpm spec sections
var sections = []string{
	"prep",
	"setup",
	"build",
	"install",
	"check",
	"clean",
	"files",
	"changelog",
	"package",
	"description",
	"verifyscript",
	"pretrans",
	"pre",
	"post",
	"preun",
	"postun",
	"posttrans",
	"triggerin",
	"triggerun",
	"triggerpostun",
}

// tags is slice with header tags
var tags = []string{
	"BuildArch",
	"BuildConflicts",
	"BuildPreReq",
	"BuildRequires",
	"BuildRoot",
	"Conflicts",
	"ExcludeArch",
	"ExclusiveArch",
	"Group",
	"License",
	"Name",
	"Obsoletes",
	"Patch",
	"PreReq",
	"Provides",
	"Release",
	"Requires",
	"Requires(posttrans)",
	"Requires(post)",
	"Requires(postun)",
	"Requires(pre)",
	"Requires(pretrans)",
	"Requires(preun)",
	"Source",
	"Summary",
	"URL",
	"Version",
}

// regexpCache is regexp cache
var regexpCache = make(map[string]*regexp.Regexp)

// sectionRegex is section check regexp
var sectionRegex = regexp.MustCompile(`^%(prep|setup|build|install|check|clean|files|changelog|package|description|verifyscript|pretrans|pre|post|preun|postun|posttrans|triggerin|triggerun|triggerpostun)( |$)`)

// ////////////////////////////////////////////////////////////////////////////////// //

// Read read and parse rpm spec file
func Read(file string) (*Spec, error) {
	err := checkFile(file)

	if err != nil {
		return nil, err
	}

	return readFile(file)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// HasSection check if section is present in spec file
func (s *Spec) HasSection(section string) bool {
	return hasSection(s, section)
}

// GetSections return slice with sections
func (s *Spec) GetSections(names ...string) []*Section {
	return extractSections(s, names)
}

// GetHeaders return slice with headers
func (s *Spec) GetHeaders() []*Header {
	return extractHeaders(s)
}

// GetSources returns slice with sources
func (s *Spec) GetSources() []Line {
	return extractSources(s)
}

// GetLine return spec line by index
func (s *Spec) GetLine(index int) Line {
	if index < 0 {
		return Line{-1, "", false}
	}

	for _, line := range s.Data {
		if line.Index == index {
			return line
		}
	}

	return Line{-1, "", false}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// GetPackageName return package name if section is package specific
func (s *Section) GetPackageName() string {
	if len(s.Args) == 0 {
		return ""
	}

	if s.Args[0] == "-n" && len(s.Args) > 1 {
		return s.Args[1]
	}

	return s.Args[0]
}

func (s *Section) IsEmpty() bool {
	for _, line := range s.Data {
		if strings.Trim(line.Text, " \t") != "" {
			return false
		}
	}

	return true
}

// ////////////////////////////////////////////////////////////////////////////////// //

// readFile reads and parses spec file
func readFile(file string) (*Spec, error) {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer fd.Close()

	line, ignore := 1, 0
	spec := &Spec{File: file}
	r := bufio.NewReader(fd)

LOOP:
	for {
		text, err := r.ReadString('\n')

		if err != nil {
			spec.Data = append(spec.Data, Line{line, text, ignore != 0})
			break LOOP
		}

		text = strings.TrimRight(text, "\r\n")

		switch {
		case isIgnoreDirective(text):
			ignore = extractIgnoreCount(text)
			ignore++
		}

		spec.Data = append(spec.Data, Line{line, text, ignore != 0})

		if ignore != 0 {
			ignore--
		}

		line++
	}

	if !isSpec(spec) {
		return nil, fmt.Errorf("File %s is not a spec file", file)
	}

	return spec, nil
}

// checkFile checks file for errors
func checkFile(file string) error {
	if !fsutil.IsExist(file) {
		return fmt.Errorf("File %s doesn't exist", file)
	}

	if !fsutil.IsRegular(file) {
		return fmt.Errorf("%s isn't a regular file", file)
	}

	if !fsutil.IsReadable(file) {
		return fmt.Errorf("File %s isn't readable", file)
	}

	if !fsutil.IsNonEmpty(file) {
		return fmt.Errorf("File %s is empty", file)
	}

	return nil
}

// hasSection returns true if spec contains given section
func hasSection(s *Spec, section string) bool {
	for _, line := range s.Data {
		if getSectionRegexp(section).MatchString(line.Text) {
			return true
		}
	}

	return false
}

// extractSections extracts data for given sections
func extractSections(s *Spec, names []string) []*Section {
	var result []*Section
	var section *Section
	var start int

	for index, line := range s.Data {
		if isSectionHeader(line.Text) {
			if section != nil {
				if start+1 <= index-1 {
					section.Data = s.Data[start+1 : index]
					section.Start, section.End = start+1, index
				}
				result = append(result, section)
				section = nil
			}

			if !isSectionMatch(strutil.ReadField(line.Text, 0, true, " "), names) {
				continue
			}

			name, args := parseSectionName(line.Text)

			section = &Section{
				Name: name,
				Args: args,
			}

			start = index
		}
	}

	if section != nil {
		section.Data = s.Data[start+1:]
		section.Start, section.End = start+1, len(s.Data)
		result = append(result, section)
	}

	return result
}

// extractHeaders extracts packages' headers
func extractHeaders(s *Spec) []*Header {
	var result []*Header
	var header *Header
	var start int

	for index, line := range s.Data {
		if header == nil {
			if len(result) == 0 && isHeaderTag(line.Text) {
				header = &Header{}
				start = index
				continue
			} else if strings.HasPrefix(line.Text, "%package") {
				name, sub := parsePackageName(line.Text)
				header = &Header{Package: name, IsSubpackage: sub}
				start = index
				continue
			}
		}

		if isSectionHeader(line.Text) {
			if header != nil {
				header.Data = s.Data[start : index-1]
				result = append(result, header)
				header = nil
			}
		}
	}

	return result
}

// extractSources extracts sources from spec file
func extractSources(s *Spec) []Line {
	var result []Line

	for _, line := range s.Data {
		if line.Ignore {
			continue
		}

		if isSectionHeader(line.Text) {
			break
		}

		lineText := strings.TrimLeft(line.Text, "\t ")

		if strings.HasPrefix(lineText, "Source") {
			result = append(result, line)
		}
	}

	return result
}

// isSectionHeader returns if given string is package header
func isSectionHeader(text string) bool {
	return sectionRegex.MatchString(text)
}

// isHeaderTag returns if given string is header tag
func isHeaderTag(text string) bool {
	for _, tagName := range tags {
		if strings.HasPrefix(text, tagName) {
			return true
		}
	}

	return false
}

// parseSectionName parses section name
func parseSectionName(text string) (string, []string) {
	if !strings.Contains(text, " ") {
		return strings.TrimLeft(text, "%"), nil
	}

	sectionNameSlice := strutil.Fields(text)

	return strings.TrimLeft(sectionNameSlice[0], "%"), sectionNameSlice[1:]
}

// parsePackageName parses package name
func parsePackageName(text string) (string, bool) {
	if strutil.ReadField(text, 1, true) == "-n" {
		return strutil.ReadField(text, 2, true), false
	}

	return strutil.ReadField(text, 1, true), true
}

// isSectionMatch returns true if data contains name of any given sections
func isSectionMatch(data string, names []string) bool {
	if len(names) == 0 {
		return true
	}

	for _, name := range names {
		if getSectionRegexp(name).MatchString(data) {
			return true
		}
	}

	return false
}

// isSpec checks that given file contains spec data
func isSpec(spec *Spec) bool {
	var count int

	markers := []string{"%install", "%files", "%changelog"}

	for _, line := range spec.Data {
		for _, marker := range markers {
			if strings.HasPrefix(line.Text, marker) {
				count++
			}
		}
	}

	if count < 3 {
		return false
	}

	count = 0
	markers = []string{"Name:", "Version:", "Summary:"}

	for _, line := range spec.Data {
		for _, marker := range markers {
			if strings.HasPrefix(line.Text, marker) {
				count++
			}
		}
	}

	return count >= 3
}

// getSectionRegexp creates new regex struct and cache it
func getSectionRegexp(section string) *regexp.Regexp {
	_, exist := regexpCache[section]

	if exist {
		return regexpCache[section]
	}

	regexpCache[section] = regexp.MustCompile(`^%` + section + `($| )`)

	return regexpCache[section]
}

// isIgnoreDirective returns true if text contains ignore directive
func isIgnoreDirective(text string) bool {
	return strings.Contains(text, DIRECTIVE_IGNORE) ||
		strings.Contains(text, "perfecto:absolve")
}

// extractIgnoreCount returns number of lines to skip
func extractIgnoreCount(text string) int {
	count := strutil.ReadField(text, 2, true)

	if count == "" {
		return 1
	}

	countInt, err := strconv.Atoi(count)

	if err != nil || countInt <= 0 {
		return 0
	}

	return countInt
}
