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
	{"/var", "%{_var}"},
	{"/usr", "%{_usr}"},
	{"/usr/src", "%{_usrsrc}"},
	{"/usr/share/doc", "%{_docdir}"},
	{"/etc", "%{_sysconfdir}"},
	{"/usr/bin", "%{_bindir}"},
	{"/usr/lib", "%{_libdir}"},
	{"/usr/lib64", "%{_libdir}"},
	{"/usr/libexec", "%{_libexecdir}"},
	{"/usr/sbin", "%{_sbindir}"},
	{"/var/lib", "%{_sharedstatedir}"},
	{"/usr/share", "%{_datarootdir}"},
	{"/usr/include", "%{_includedir}"},
	{"/usr/share/info", "%{_infodir}"},
	{"/usr/share/man", "%{_mandir}"},
	{"/etc/rc.d/init.d", ""},
	{"/etc/init", "%{_initddir}"},
	{"/usr/share/java", "%{_javadir}"},
	{"/usr/share/javadoc", "%{_javadocdir}"},
	{"/usr/share/doc", "%{_defaultdocdir}"},
}

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
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkForUselessSpaces checks for useless spaces
func checkForUselessSpaces(s *spec.Spec) []Alert {
	var result []Alert

	for _, line := range s.Data {
		if contains(line, " ") {
			if strings.TrimSpace(line.Text) == "" {
				result = append(result, Alert{LEVEL_NOTICE, "Line contains useless spaces", spec.Line{line.Index, ""}})
			} else if strings.TrimRight(line.Text, " ") != line.Text {
				cleanLine := strings.TrimRight(line.Text, " ")
				spaces := len(line.Text) - len(cleanLine)
				impLine := spec.Line{line.Index, cleanLine + strings.Repeat("▒", spaces)}
				result = append(result, Alert{LEVEL_NOTICE, "Line contains spaces at the end of line", impLine})
			}
		}
	}

	return result
}

// checkForLineLength checks changelog and description lines for 80 symbols limit
func checkForLineLength(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{"description", "changelog"}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			// Ignore changelog headers
			if section.Name == "changelog" && prefix(line, "* ") {
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
		"prep", "setup", "build", "install", "check",
		"files", "package", "verifyscript", "pretrans",
		"pre", "post", "preun", "postun", "posttrans",
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			// Ignore env vars exports
			if contains(line, "export") {
				continue
			}

			for _, macro := range pathMacroSlice {
				if contains(line, macro.Value) {
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

	for _, section := range s.GetSections("install") {
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
		"prep", "setup", "build", "install", "check",
		"verifyscript", "pretrans", "pre", "post",
		"preun", "postun", "posttrans",
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

	for _, section := range s.GetSections("changelog") {
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

	sections := []string{"build", "install", "check"}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if !contains(line, "make") {
				continue
			}

			if prefix(line, "make") {
				result = append(result, Alert{LEVEL_WARNING, "Use %{__make} macro instead of \"make\"", line})
			}

			if section.Name == "install" && contains(line, " install") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					result = append(result, Alert{LEVEL_WARNING, "Use %{make_install} macro instead of \"make install\"", line})
				}
			}

			if section.Name == "build" && !contains(line, "%{?_smp_mflags}") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					result = append(result, Alert{LEVEL_WARNING, "Don't forget to use %{?_smp_mflags} macro with make command", line})
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
				result = append(result, Alert{LEVEL_ERROR, "Main package must contain URL tag", spec.Line{-1, ""}})
			}
		}

		if !containsTag(header.Data, "Group:") {
			if header.Package == "" {
				result = append(result, Alert{LEVEL_WARNING, "Main package must contain Group tag", spec.Line{-1, ""}})
			} else {
				result = append(result, Alert{LEVEL_WARNING, fmt.Sprintf("Package %s must contain Group tag", header.Package), spec.Line{-1, ""}})
			}
		}
	}

	return result
}

// checkForUnescapedPercent check changelog and descriptions for unescaped percent symbol
func checkForUnescapedPercent(s *spec.Spec) []Alert {
	var result []Alert

	sections := []string{"changelog"}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if prefix(line, "%") {
				continue
			}

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

// ////////////////////////////////////////////////////////////////////////////////// //

// prefix is strings.HasPrefix wrapper
func prefix(line spec.Line, value string) bool {
	return strings.HasPrefix(strings.TrimLeft(line.Text, " "), value)
}

// contains is strings.Contains wrapper
func contains(line spec.Line, value string) bool {
	return strings.Contains(line.Text, value)
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
