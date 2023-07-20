package check

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/essentialkaos/ek/v12/cache"
	"github.com/essentialkaos/ek/v12/req"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/perfecto/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Checker is spec check function
type Checker func(id string, s *spec.Spec) []Alert

type macro struct {
	Value string
	Name  string
}

// ////////////////////////////////////////////////////////////////////////////////// //

var httpCheckCache *cache.Cache

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

var binariesAsMacro = []string{
	"7zip", "bzip2", "bzr", "cat", "chgrp", "chmod", "chown", "cp", "cpio",
	"file", "git", "grep", "gzip", "hg", "id", "install", "ld", "lrzip", "lzip",
	"mkdir", "mv", "nm", "objcopy", "objdump", "patch", "quilt",
	"rm", "rsh", "sed", "semodule", "ssh", "strip", "tar", "unzip", "xz",
}

var emptyLine = spec.Line{-1, "", false}

var macroRegExp = regexp.MustCompile(`\%\{?\??([a-zA-Z0-9_\?\:]+)\}?`)

// ////////////////////////////////////////////////////////////////////////////////// //

// getCheckers return slice with all supported checkers
func getCheckers() map[string]Checker {
	return map[string]Checker{
		"PF1":  checkForUselessSpaces,
		"PF2":  checkForLineLength,
		"PF3":  checkForDist,
		"PF4":  checkForNonMacroPaths,
		"PF5":  checkForVariables,
		"PF6":  checkForDevNull,
		"PF7":  checkChangelogHeaders,
		"PF8":  checkForMakeMacro,
		"PF9":  checkForHeaderTags,
		"PF10": checkForUnescapedPercent,
		"PF11": checkForMacroDefinitionPosition,
		"PF12": checkForSeparatorLength,
		"PF13": checkForDefAttr,
		"PF14": checkForUselessBinaryMacro,
		"PF15": checkForEmptySections,
		"PF16": checkForIndentInFilesSection,
		"PF17": checkForSetupOptions,
		"PF18": checkForEmptyLinesAtEnd,
		"PF19": checkBashLoops,
		"PF20": checkURLForHTTPS,
		"PF21": checkForCheckMacro,
		"PF22": checkIfClause,
		"PF23": checkForUselessSlash,
		"PF24": checkForEmptyIf,
		"PF25": checkForDotInSummary,
		"PF26": checkForChownAndChmod,
		"PF27": checkForUnclosedCondition,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// checkForUselessSpaces checks for useless spaces
func checkForUselessSpaces(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, line := range s.Data {
		if contains(line, " ") {
			if strings.TrimSpace(line.Text) == "" {
				impLine := spec.Line{line.Index, strings.Replace(line.Text, " ", "░", -1), line.Ignore}
				result = append(result, NewAlert(id, LEVEL_NOTICE, "Line contains useless spaces", impLine))
			} else if strings.TrimRight(line.Text, " ") != line.Text {
				cleanLine := strings.TrimRight(line.Text, " ")
				spaces := len(line.Text) - len(cleanLine)
				impLine := spec.Line{line.Index, cleanLine + strings.Repeat("░", spaces), line.Ignore}
				result = append(result, NewAlert(id, LEVEL_NOTICE, "Line contains spaces at the end of line", impLine))
			}
		}
	}

	return result
}

// checkForLineLength checks changelog and description lines for 80 symbols limit
func checkForLineLength(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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

			if strings.IndexRune(strutil.Substr(line.Text, 2, 999), ' ') == -1 {
				continue
			}

			if strutil.Len(line.Text) > 80 {
				result = append(result, NewAlert(id, LEVEL_WARNING, "Line is longer than 80 symbols", line))
			}
		}
	}

	return result
}

// checkForDist checks for dist macro in release tag
func checkForDist(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, header := range s.GetHeaders() {
		for _, line := range header.Data {
			if isComment(line) {
				continue
			}

			if prefix(line, "Release:") {
				if !containsMacro(line, "autorelease") && !containsMacro(line, "dist") {
					result = append(result, NewAlert(id, LEVEL_ERROR, "Release tag must contains %{?dist} as part of release", line))
				}
			}
		}
	}

	return result
}

// checkForNonMacroPaths checks if standard path not used as macro
func checkForNonMacroPaths(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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
			if isComment(line) {
				continue
			}

			// Ignore comments and env vars exports
			if contains(line, "export") {
				continue
			}

			// Ignore sed replacements
			if contains(line, "sed") {
				continue
			}

			text := line.Text

			for _, macro := range pathMacroSlice {
				re := regexp.MustCompile(macro.Value + `(\/|$|%)`)
				if re.MatchString(text) {
					result = append(result, NewAlert(id, LEVEL_WARNING, fmt.Sprintf("Path \"%s\" should be used as macro \"%s\"", macro.Value, macro.Name), line))
				}
			}
		}
	}

	return result
}

// checkForVariables checks for using variables instead of macros
func checkForVariables(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{
		spec.SECTION_BUILD,
		spec.SECTION_INSTALL,
		spec.SECTION_CLEAN,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			switch {
			case contains(line, "$RPM_BUILD_ROOT"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Build root path must be used as macro %{buildroot}", line))
			case contains(line, "$RPM_OPT_FLAGS"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Optimization flags must be used as macro %{optflags}", line))
			case contains(line, "$RPM_LD_FLAGS"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Linking flags must be used as macro %{build_ldflags}", line))
			case contains(line, "$RPM_DOC_DIR"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Linking flags must be used as macro %{_docdir}", line))
			case contains(line, "$RPM_SOURCE_DIR"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Path to source directory must be used as macro %{_sourcedir}", line))
			case contains(line, "$RPM_BUILD_DIR"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Path to build directory must be used as macro %{_builddir}", line))
			case contains(line, "$RPM_ARCH"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Arch value must be used as macro %{_arch}", line))
			case contains(line, "$RPM_OS"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "OS value must be used as macro %{_os}", line))
			case contains(line, "$RPM_PACKAGE_NAME"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Package name value must be used as macro %{name}", line))
			case contains(line, "$RPM_PACKAGE_VERSION"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Package version value must be used as macro %{version}", line))
			case contains(line, "$RPM_PACKAGE_RELEASE"):
				result = append(result, NewAlert(id, LEVEL_ERROR, "Package release value must be used as macro %{release}", line))
			}
		}
	}

	return result
}

// checkForDevNull checks for devnull redirect format
func checkForDevNull(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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

	variations := []string{
		">/dev/null 2>&1",
		"2>&1 >/dev/null",
		">/dev/null 2>/dev/null",
		"2>/dev/null >/dev/null",
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			for _, v := range variations {
				if strings.Contains(strutil.Exclude(line.Text, " "), strutil.Exclude(v, " ")) {
					result = append(result, NewAlert(id, LEVEL_NOTICE, fmt.Sprintf("Use \"&>/dev/null || :\" instead of \"%s || :\"", v), line))
				}
			}

			if contains(line, "|| exit 0") {
				result = append(result, NewAlert(id, LEVEL_NOTICE, "Use \" || :\" instead of \" || exit 0\"", line))
			}
		}
	}

	return result
}

// checkChangelogHeaders checks changelog for misformatted records
func checkChangelogHeaders(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, section := range s.GetSections(spec.SECTION_CHANGELOG) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			// Ignore changelog records text
			if !prefix(line, "* ") {
				continue
			}

			if !contains(line, " - ") {
				result = append(result, NewAlert(id, LEVEL_WARNING, "Misformatted changelog record header", line))
			} else {
				separator := strings.Index(line.Text, " - ")
				if !strings.Contains(strutil.Substr(line.Text, separator+3, 999999), "-") {
					result = append(result, NewAlert(id, LEVEL_WARNING, "Changelog record header must contain release", line))
				}
			}
		}
	}

	return result
}

// checkForMakeMacro checks if make used not as macro
func checkForMakeMacro(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{
		spec.SECTION_BUILD,
		spec.SECTION_INSTALL,
		spec.SECTION_CHECK,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			if !contains(line, "make") {
				continue
			}

			if prefix(line, "make") {
				result = append(result, NewAlert(id, LEVEL_WARNING, "Use %{__make} macro instead of \"make\"", line))
			}

			if section.Name == spec.SECTION_INSTALL && containsField(line, "install") && contains(line, "DESTDIR") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					result = append(result, NewAlert(id, LEVEL_WARNING, "Use %{make_install} macro instead of \"make install\"", line))
				}
			}

			if section.Name == spec.SECTION_BUILD && !contains(line, "%{?_smp_mflags}") {
				if prefix(line, "make") || prefix(line, "%{__make}") {
					if line.Text == "make" || line.Text == "%{__make}" || containsField(line, "all") {
						result = append(result, NewAlert(id, LEVEL_WARNING, "Don't forget to use %{?_smp_mflags} macro with make command", line))
					}
				}
			}
		}
	}

	return result
}

// checkForHeaderTags checks headers for required tags
func checkForHeaderTags(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, header := range s.GetHeaders() {
		if header.Package == "" {
			if !containsTag(header.Data, "URL:") {
				result = append(result, NewAlert(id, LEVEL_ERROR, "Main package must contain URL tag", emptyLine))
			}
		}

		if !containsTag(header.Data, "Group:") {
			if header.Package == "" {
				result = append(result, NewAlert(id, LEVEL_WARNING, "Main package must contain Group tag", emptyLine))
			} else {
				result = append(result, NewAlert(id, LEVEL_WARNING, fmt.Sprintf("Package %s must contain Group tag", header.Package), emptyLine))
			}
		}
	}

	return result
}

// codebeat:disable[BLOCK_NESTING]

// checkForUnescapedPercent checks changelog and descriptions for unescaped percent symbol
func checkForUnescapedPercent(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{spec.SECTION_CHANGELOG}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if containsMacro(line, "autochangelog") {
				continue
			}

			for _, word := range strings.Fields(line.Text) {
				if strings.HasPrefix(word, "%") && !strings.HasPrefix(word, "%%") {
					result = append(result, NewAlert(id, LEVEL_ERROR, "Symbol % must be escaped by another % (i.e % → %%)", line))
				}
			}
		}
	}

	return result
}

// codebeat:enable[BLOCK_NESTING]

// checkForMacroDefinitionPosition checks for macro defined after description
func checkForMacroDefinitionPosition(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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
				result = append(result, NewAlert(id, LEVEL_WARNING, "Move %define and %global to top of your spec", line))
			}
		}
	}

	return result
}

// checkForSeparatorLength checks for separator length
func checkForSeparatorLength(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, line := range s.Data {
		if contains(line, "#") && strings.Trim(line.Text, "#") == "" && strings.Count(line.Text, "#") != 80 {
			result = append(result, NewAlert(id, LEVEL_NOTICE, "Separator must be 80 symbols long", line))
		}
	}

	return result
}

// checkForDefAttr checks spec for %defattr macro in %files sections
func checkForDefAttr(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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
			result = append(result, NewAlert(id, LEVEL_ERROR, "%files section must contains %defattr macro", emptyLine))
		default:
			result = append(result, NewAlert(id, LEVEL_ERROR, "%files section for package "+packageName+" must contains %defattr macro", emptyLine))
		}
	}

	return result
}

// checkForUselessBinaryMacro checks spec for useless binary macro
func checkForUselessBinaryMacro(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, line := range s.Data {
		for _, binary := range binariesAsMacro {
			if contains(line, "%{__"+binary+"}") {
				result = append(result, NewAlert(id, LEVEL_NOTICE, fmt.Sprintf("Useless macro %%{__%s} used for executing %s binary", binary, binary), line))
			}
		}
	}

	return result
}

// checkForEmptySections checks spec for empty sections
func checkForEmptySections(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{
		spec.SECTION_CHECK,
		spec.SECTION_POST,
		spec.SECTION_POSTTRANS,
		spec.SECTION_POSTUN,
		spec.SECTION_PRE,
		spec.SECTION_PRETRANS,
		spec.SECTION_PREUN,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
		spec.SECTION_VERIFYSCRIPT,
		spec.SECTION_VERIFYSCRIPT,
	}

	for _, section := range s.GetSections(sections...) {
		if len(section.Args) == 0 && isEmptyData(section.Data) {
			result = append(result, NewAlert(id, LEVEL_ERROR, fmt.Sprintf("Section %%%s is empty", section.Name), s.GetLine(section.Start)))
		}
	}

	return result
}

// checkForIndentInFilesSection checks spec for prefixes in %files section
func checkForIndentInFilesSection(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, section := range s.GetSections(spec.SECTION_FILES) {
		for _, line := range section.Data {
			if strings.HasPrefix(line.Text, " ") || strings.HasPrefix(line.Text, "\t") {
				result = append(result, NewAlert(id, LEVEL_NOTICE, "Don't use indent in %files section", line))
			}
		}
	}

	return result
}

// checkForSetupOptions checks setup arguments
func checkForSetupOptions(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, section := range s.GetSections(spec.SECTION_SETUP) {
		switch {
		case containsArgs(section, "-q", "-c", "-n"):
			result = append(result, NewAlert(id, LEVEL_NOTICE, "Options \"-q -c -n\" can be simplified to \"-qcn\"", s.GetLine(section.Start)))
		case containsArgs(section, "-q", "-n"):
			result = append(result, NewAlert(id, LEVEL_NOTICE, "Options \"-q -n\" can be simplified to \"-qn\"", s.GetLine(section.Start)))
		case containsArgs(section, "-c", "-n"):
			result = append(result, NewAlert(id, LEVEL_NOTICE, "Options \"-c -n\" can be simplified to \"-cn\"", s.GetLine(section.Start)))
		}
	}

	return result
}

// checkForEmptyLinesAtEnd checks spec for empty lines at the end
func checkForEmptyLinesAtEnd(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	totalLines := len(s.Data)
	lastLine := s.Data[totalLines-1]

	if lastLine.Text != "" {
		return []Alert{NewAlert(id, LEVEL_NOTICE, "Spec file should have empty line at the end", emptyLine)}
	}

	emptyLines := 0

	for i := totalLines - 1; i > 0; i-- {
		if s.Data[i].Text == "" {
			emptyLines++
		} else {
			if emptyLines > 1 {
				return []Alert{NewAlert(id, LEVEL_NOTICE, "Too much empty lines at the end of the spec", emptyLine)}
			}

			break
		}
	}

	return nil
}

// checkBashLoops checks bash loops format
func checkBashLoops(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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
		spec.SECTION_TRIGGERIN,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
		spec.SECTION_VERIFYSCRIPT,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if !prefix(line, "for") && !prefix(line, "while") {
				continue
			}

			nextLine := s.GetLine(line.Index + 1)
			nextLineText := strings.TrimLeft(nextLine.Text, "\t ")

			if !suffix(nextLine, ";do") && nextLineText == "do" {
				result = append(result, NewAlert(id, LEVEL_NOTICE, "Place 'do' keyword on the same line with for/while (for ... ; do)", line))
			}
		}
	}

	return result
}

// checkURLForHTTPS checks if source domain supports HTTPS
func checkURLForHTTPS(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	if httpCheckCache == nil {
		httpCheckCache = cache.New(time.Hour, 0)
	}

	var result []Alert

	urls := s.GetSources()

	for _, header := range s.GetHeaders() {
		for _, line := range header.Data {
			if prefix(line, "URL:") {
				urls = append(urls, line)
			}
		}
	}

	for _, line := range urls {
		lineText := strings.TrimLeft(line.Text, "\t ")
		url := strutil.ReadField(lineText, 1, true, " ")

		if !strings.HasPrefix(url, "http://") {
			continue
		}

		domain := extractDomainFromURL(url)

		if domain == "" {
			continue
		}

		if isHostSupportsHTTPS(domain) {
			result = append(result, NewAlert(
				id, LEVEL_WARNING,
				fmt.Sprintf("Domain %s supports HTTPS. Replace http by https in URL.", domain),
				line,
			))
		}
	}

	return result
}

// checkForCheckMacro checks check section for a macro which allows skipping the check
func checkForCheckMacro(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	if !s.HasSection(spec.SECTION_CHECK) {
		return nil
	}

	for _, section := range s.GetSections(spec.SECTION_CHECK) {
		if section.IsEmpty() {
			return nil
		}

		for _, line := range section.Data {
			if contains(line, "?_without_check") && contains(line, "?_with_check") {
				return nil
			}
		}
	}

	return []Alert{
		NewAlert(id, LEVEL_WARNING, "Use %{_without_check} and %{_with_check} macros for controlling tests execution", emptyLine),
	}
}

// checkIfClause checks if clause for using single equals symbol instead of two
func checkIfClause(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, line := range s.Data {
		if !prefix(line, "%if ") {
			continue
		}

		if contains(line, " = ") {
			result = append(result, NewAlert(id, LEVEL_ERROR, "Use two equals symbols for comparison in %if clause", line))
		}
	}

	return result
}

// checkForUselessSlash checks for useless slash after %{buildroot} macro
func checkForUselessSlash(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			for _, macro := range pathMacroSlice {
				if contains(line, "%{buildroot}/"+macro.Name) {
					desc := fmt.Sprintf("Slash between %%{buildroot} and %s macros is useless", macro.Name)
					result = append(result, NewAlert(id, LEVEL_WARNING, desc, line))
				}
			}
		}
	}

	return result
}

// checkForEmptyIf checks for possible empty if clauses
func checkForEmptyIf(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

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

	var clauseOpen, macroOpen, hasContent bool
	var clauseLine spec.Line

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			if prefix(line, "if ") && !macroOpen {
				clauseOpen = true
				clauseLine = line
				continue
			}

			if prefix(line, "%else") {
				hasContent = true
				continue
			}

			if prefix(line, "%if") {
				if !macroOpen {
					macroOpen = true
				} else {
					hasContent = true
				}
			}

			if prefix(line, "%endif") && macroOpen {
				macroOpen = false
				continue
			}

			if !macroOpen && clauseOpen && !prefix(line, "fi") {
				hasContent = true
			}

			if prefix(line, "fi") {
				if clauseOpen && !hasContent {
					desc := fmt.Sprintf("Evaluated if clause can be empty. Change the order of clauses (i.e. %%if → if instead of if → %%if).")
					result = append(result, NewAlert(id, LEVEL_WARNING, desc, clauseLine))
				}

				clauseOpen, macroOpen, hasContent = false, false, false
			}
		}

		clauseOpen, macroOpen, hasContent = false, false, false
	}

	return result
}

// checkForDotInSummary checks for dot on the end of summary
func checkForDotInSummary(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	for _, header := range s.GetHeaders() {
		for _, line := range header.Data {
			if isComment(line) {
				continue
			}

			if prefix(line, "Summary:") && suffix(line, ".") {
				result = append(result, NewAlert(id, LEVEL_WARNING, "The summary contains useless dot at the end", line))
			}
		}
	}

	return result
}

// checkForChownAndChmod checks scriptlets for chown and chmod commands
func checkForChownAndChmod(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{
		spec.SECTION_POST,
		spec.SECTION_POSTTRANS,
		spec.SECTION_POSTUN,
		spec.SECTION_PRE,
		spec.SECTION_PREP,
		spec.SECTION_PRETRANS,
		spec.SECTION_PREUN,
		spec.SECTION_TRIGGERIN,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
	}

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			if prefix(line, "chmod ") {
				result = append(result, NewAlert(id, LEVEL_ERROR, "Do not change file or directory mode in scriptlets", line))
			}

			if prefix(line, "chown ") {
				if !contains(line, " -h ") && !contains(line, " --no-dereference ") {
					result = append(result, NewAlert(id, LEVEL_ERROR, "Do not change file or directory owner without --no-dereference option", line))
				}
			}
		}
	}

	return result
}

// checkForChownAndChmod checks scriptlets for unclosed conditions
func checkForUnclosedCondition(id string, s *spec.Spec) []Alert {
	if len(s.Data) == 0 {
		return nil
	}

	var result []Alert

	sections := []string{
		spec.SECTION_POST,
		spec.SECTION_POSTTRANS,
		spec.SECTION_POSTUN,
		spec.SECTION_PRE,
		spec.SECTION_PREP,
		spec.SECTION_PRETRANS,
		spec.SECTION_PREUN,
		spec.SECTION_TRIGGERIN,
		spec.SECTION_TRIGGERPOSTUN,
		spec.SECTION_TRIGGERUN,
	}

	var conditions []spec.Line

	for _, section := range s.GetSections(sections...) {
		for _, line := range section.Data {
			if isComment(line) {
				continue
			}

			if prefix(line, "if ") && contains(line, ";") && contains(line, "then") {
				if !contains(line, " fi") {
					conditions = append(conditions, line)
				}
			}

			if prefix(line, "fi") && len(conditions) != 0 {
				conditions = conditions[:len(conditions)-1]
			}
		}
	}

	if len(conditions) != 0 {
		for _, line := range conditions {
			result = append(result, NewAlert(id, LEVEL_CRITICAL, "Scriptlet contains unclosed IF condition", line))
		}
	}

	return result
}

// ////////////////////////////////////////////////////////////////////////////////// //

// prefix is strings.HasPrefix wrapper
func prefix(line spec.Line, value string) bool {
	return strings.HasPrefix(strings.TrimLeft(line.Text, "\t "), value)
}

// suffix is strings.HasSuffix wrapper
func suffix(line spec.Line, value string) bool {
	return strings.HasSuffix(strings.TrimLeft(line.Text, "\t "), value)
}

// contains is strings.Contains wrapper
func contains(line spec.Line, value string) bool {
	return strings.Contains(line.Text, value)
}

// contains checks if line contain given macro
func containsMacro(line spec.Line, macro string) bool {
	for _, found := range macroRegExp.FindAllStringSubmatch(line.Text, -1) {
		if found[1] == macro {
			return true
		}
	}

	return false
}

// containsField return true if line contains given field
func containsField(line spec.Line, value string) bool {
	return sliceutil.Contains(strutil.Fields(line.Text), value)
}

// isComment return true if current line is commented
func isComment(line spec.Line) bool {
	return prefix(line, "#")
}

// isEmptyData check if data is empty or contains only spaces
func isEmptyData(data []spec.Line) bool {
	for _, line := range data {
		if strings.Replace(line.Text, " ", "", -1) != "" {
			return false
		}
	}

	return true
}

// containsArgs return true if section contains given args
func containsArgs(section *spec.Section, args ...string) bool {
	for _, arg := range args {
		if !sliceutil.Contains(section.Args, arg) {
			return false
		}
	}

	return true
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

// extractDomainFromURL extracts domain name from source URL
func extractDomainFromURL(url string) string {
	url = strutil.Exclude(url, "http://")
	return strutil.ReadField(url, 0, false, "/")
}

// isHostSupportsHTTPS return true if domain supports HTTPS protocol
func isHostSupportsHTTPS(domain string) bool {
	if httpCheckCache.Has(domain) {
		return httpCheckCache.Get(domain).(bool)
	}

	_, err := req.Request{
		URL:         "https://" + domain,
		AutoDiscard: true,
	}.Head()

	supported := err == nil

	httpCheckCache.Set(domain, supported)

	return supported
}
