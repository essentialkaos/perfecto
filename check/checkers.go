package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2018 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"pkg.re/essentialkaos/ek.v9/sliceutil"
	"pkg.re/essentialkaos/ek.v9/strutil"

	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Checker is spec check function
type Checker func(s *spec.Spec) []Alert

// ////////////////////////////////////////////////////////////////////////////////// //

type macro struct {
	Value string
	Name  string
}

var pathMacroSlice = []macro{
	{"/etc/init", "%{_initddir}"},
	{"/etc/rc.d/init.d", "%{_initddir}"},
	{"/etc", "%{_sysconfdir}"},
	{"/usr/bin", "%{_bindir}"},
	{"/usr/include", "%{_includedir}"},
	{"/usr/lib", "%{_libdir}"},
	{"/usr/lib64", "%{_libdir}"},
	{"/usr/libexec", "%{_libexecdir}"},
	{"/usr/sbin", "%{_sbindir}"},
	{"/usr/share/doc", "%{_defaultdocdir}"},
	{"/usr/share/doc", "%{_docdir}"},
	{"/usr/share/info", "%{_infodir}"},
	{"/usr/share/java", "%{_javadir}"},
	{"/usr/share/javadoc", "%{_javadocdir}"},
	{"/usr/share/man", "%{_mandir}"},
	{"/usr/share", "%{_datarootdir}"},
	{"/usr/src", "%{_usrsrc}"},
	{"/usr", "%{_usr}"},
	{"/var/lib", "%{_sharedstatedir}"},
	{"/var", "%{_var}"},
}

var emptyLine = spec.Line{-1, "", false}

// ////////////////////////////////////////////////////////////////////////////////// //

// getCheckers return slice with all supported checkers
func getCheckers() []Checker {
	return []Checker{
		checkForUselessSpaces,
		checkForLineLength,
		checkForDist,
		checkForNonMacroPaths,
		checkForBuildRoot,
		checkForDevNull,
		checkChangelogHeaders,
		checkForMakeMacro,
		checkForHeaderTags,
		checkForUnescapedPercent,
		checkForMacroDefenitionPosition,
		checkForSeparatorLength,
		checkForDefAttr,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkForUselessSpaces checks for useless spaces
func checkForUselessSpaces(s *spec.Spec) []Alert {
	var result []Alert

	for _, line := range s.Data {
		if contains(line, " ") {
			if strings.TrimSpace(line.Text) == "" {
				impLine := spec.Line{line.Index, strings.Replace(line.Text, " ", "▒", -1), line.Skip}
				result = append(result, Alert{LEVEL_NOTICE, "Line contains useless spaces", impLine})
			} else if strings.TrimRight(line.Text, " ") != line.Text {
				cleanLine := strings.TrimRight(line.Text, " ")
				spaces := len(line.Text) - len(cleanLine)
				impLine := spec.Line{line.Index, cleanLine + strings.Repeat("▒", spaces), line.Skip}
				result = append(result, Alert{LEVEL_NOTICE, "Line contains spaces at the end of line", impLine})
			}
		}
	}

	return result
}

// checkForLineLength checks changelog and description lines for 80 symbols limit
func checkForLineLength(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{
		spec.SECTION_DESCRIPTION,
		spec.SECTION_CHANGELOG,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			// Ignore changelog headers
			if section.Name == spec.SECTION_CHANGELOG && prefix(line, "* ") {
				continue
			}

			if strutil.Len(line.Text) > 80 {
				result = append(result, Alert{LEVEL_WARNING, "Line is longer than 80 symbols", line})
			}
		}
	}

	return result
}

// checkForDist checks for dist macro in release tag
func checkForDist(s *spec.Spec) []Alert {
	var result []Alert

	for _, header := range s.GetHeaders() {
		for _, line := range header.Data {
			if strings.HasPrefix(line.Text, "Release:") {
				if !contains(line, "%{?dist}") {
					result = append(result, Alert{LEVEL_ERROR, "Release tag must contains %{?dist} as part of release", line})
				}
			}
		}
	}

	return result
}

// checkForNonMacroPaths checks if standart path not used as macro
func checkForNonMacroPaths(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{
		spec.SECTION_BUILD,
		spec.SECTION_CHECK,
		spec.SECTION_CLEAN,
		spec.SECTION_FILES,
		spec.SECTION_INSTALL,
		spec.SECTION_PACKAGE,
		spec.SECTION_POST,
		spec.SECTION_POSTTRANS,
		spec.SECTION_POSTUN,
		spec.SECTION_PRE,
		spec.SECTION_PREP,
		spec.SECTION_PRETRANS,
		spec.SECTION_PREUN,
		spec.SECTION_SETUP,
		spec.SECTION_TRIGGERIN,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
		spec.SECTION_VERIFYSCRIPT,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			// Ignore comments and env vars exports
			if contains(line, "export") || prefix(line, "#") {
				continue
			}

			// Ignore sed replacements
			if contains(line, "sed") {
				continue
			}

			text := line.Text

			for _, macro := range pathMacroSlice {
				if strings.Contains(text, macro.Value) {
					text = strings.Replace(text, macro.Value, "", -1)
					result = append(result, Alert{LEVEL_WARNING, fmt.Sprintf("Path \"%s\" should be used as macro \"%s\"", macro.Value, macro.Name), line})
				}
			}
		}
	}

	return result
}

// checkForBuildRoot checks for build root path used as $RPM_BUILD_ROOT
func checkForBuildRoot(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{
		spec.SECTION_INSTALL,
		spec.SECTION_CLEAN,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if contains(line, "$RPM_BUILD_ROOT") {
				result = append(result, Alert{LEVEL_ERROR, "Build root path must be used as macro %{buildroot}", line})
			}

			if contains(line, "%{buildroot}/%{_") {
				result = append(result, Alert{LEVEL_WARNING, "Slash after %{buildroot} macro is useless", line})
			}
		}
	}

	return result
}

// checkForDevNull checks for devnull redirect format
func checkForDevNull(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{
		spec.SECTION_BUILD,
		spec.SECTION_CHECK,
		spec.SECTION_INSTALL,
		spec.SECTION_POST,
		spec.SECTION_POSTTRANS,
		spec.SECTION_POSTUN,
		spec.SECTION_PRE,
		spec.SECTION_PREP,
		spec.SECTION_PRETRANS,
		spec.SECTION_PREUN,
		spec.SECTION_SETUP,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
		spec.SECTION_VERIFYSCRIPT,
		spec.SECTION_VERIFYSCRIPT,
	}

	devNull := strings.Replace(">/dev/null 2>&1 || :", " ", "", -1)

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if strings.Contains(strings.Replace(line.Text, " ", "", -1), devNull) {
				result = append(result, Alert{LEVEL_NOTICE, "Use \"&>/dev/null || :\" instead of \">/dev/null 2>&1 || :\"", line})
			}
		}
	}

	return result
}

// checkChangelogHeaders checks changelog for misformatted records
func checkChangelogHeaders(s *spec.Spec) []Alert {
	var result []Alert

	for _, section := range s.GetSections(spec.SECTION_CHANGELOG) {
		for _, line := range section.Data {
			// Ignore changelog records text
			if !prefix(line, "* ") {
				continue
			}

			if !contains(line, " - ") {
				result = append(result, Alert{LEVEL_WARNING, "Misformatted changelog record header", line})
			} else {
				separator := strings.Index(line.Text, " - ")
				if !strings.Contains(strutil.Substr(line.Text, separator+3, 999999), "-") {
					result = append(result, Alert{LEVEL_WARNING, "Changelog record header must contain release", line})
				}
			}
		}
	}

	return result
}

// checkForMakeMacro checks if make used not as macro
func checkForMakeMacro(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{
		spec.SECTION_BUILD,
		spec.SECTION_INSTALL,
		spec.SECTION_CHECK,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if !contains(line, "make") {
				continue
			}

			if prefix(line, "make") {
				result = append(result, Alert{LEVEL_WARNING, "Use %{__make} macro instead of \"make\"", line})
			}

			if section.Name == spec.SECTION_INSTALL && containsField(line, "install") && contains(line, "DESTDIR") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					result = append(result, Alert{LEVEL_WARNING, "Use %{make_install} macro instead of \"make install\"", line})
				}
			}

			if section.Name == spec.SECTION_BUILD && !contains(line, "%{?_smp_mflags}") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					if line.Text == "make" || line.Text == "%{__make}" || containsField(line, "all") {
						result = append(result, Alert{LEVEL_WARNING, "Don't forget to use %{?_smp_mflags} macro with make command", line})
					}
				}
			}
		}
	}

	return result
}

// checkForHeaderTags check headers for required tags
func checkForHeaderTags(s *spec.Spec) []Alert {
	var result []Alert

	for _, header := range s.GetHeaders() {
		if header.Package == "" {
			if !containsTag(header.Data, "URL:") {
				result = append(result, Alert{LEVEL_ERROR, "Main package must contain URL tag", emptyLine})
			}
		}

		if !containsTag(header.Data, "Group:") {
			if header.Package == "" {
				result = append(result, Alert{LEVEL_WARNING, "Main package must contain Group tag", emptyLine})
			} else {
				result = append(result, Alert{LEVEL_WARNING, fmt.Sprintf("Package %s must contain Group tag", header.Package), emptyLine})
			}
		}
	}

	return result
}

// checkForUnescapedPercent check changelog and descriptions for unescaped percent symbol
func checkForUnescapedPercent(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{spec.SECTION_CHANGELOG}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if contains(line, "%") && !contains(line, "%%") {
				result = append(result, Alert{LEVEL_ERROR, "Symbol % must be escaped by another % (i.e % → %%)", line})
			}
		}
	}

	return result
}

// checkForMacroDefenitionPosition check for macro defined after description
func checkForMacroDefenitionPosition(s *spec.Spec) []Alert {
	var result []Alert

	var underDescription bool

	for _, line := range s.Data {
		if !underDescription && prefix(line, "%description") {
			underDescription = true
		}

		if prefix(line, "%files") {
			break
		}

		if underDescription {
			if contains(line, "%global ") || contains(line, "%define ") {
				result = append(result, Alert{LEVEL_WARNING, "Move %define and %global to top of your spec", line})
			}
		}
	}

	return result
}

// checkForSeparatorLength check for separator length
func checkForSeparatorLength(s *spec.Spec) []Alert {
	var result []Alert

	for _, line := range s.Data {
		if contains(line, "#") && strings.Trim(line.Text, "#") == "" && strings.Count(line.Text, "#") != 80 {
			result = append(result, Alert{LEVEL_NOTICE, "Separator must be 80 symbols long", line})
		}
	}

	return result
}

// checkForDefAttr check spec for %defattr macro in %files sections
func checkForDefAttr(s *spec.Spec) []Alert {
	var result []Alert

	for _, section := range s.GetSections(spec.SECTION_FILES) {
		hasDefAttr := false

		for _, line := range section.Data {
			if prefix(line, "%defattr") {
				hasDefAttr = true
			}
		}

		if hasDefAttr {
			continue
		}

		packageName := section.GetPackageName()

		switch packageName {
		case "":
			result = append(result, Alert{LEVEL_ERROR, "%files section must contains %defattr macro", emptyLine})
		default:
			result = append(result, Alert{LEVEL_ERROR, "%files section for package " + packageName + " must contains %defattr macro", emptyLine})
		}
	}

	return result
}

// ////////////////////////////////////////////////////////////////////////////////// //

// prefix is strings.HasPrefix wrapper
func prefix(line spec.Line, value string) bool {
	return strings.HasPrefix(strings.TrimLeft(line.Text, " "), value)
}

// contains is strings.Contains wrapper
func contains(line spec.Line, value string) bool {
	return strings.Contains(line.Text, value)
}

// containsField return true if line contains given field
func containsField(line spec.Line, value string) bool {
	return sliceutil.Contains(strutil.Fields(line.Text), value)
}

// containsTag check if data contains given tag
func containsTag(data []spec.Line, tag string) bool {
	for _, line := range data {
		if prefix(line, tag) {
			return true
		}
	}

	return false
}
