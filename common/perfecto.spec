################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define debug_package  %{nil}

################################################################################

Summary:        Tool for checking perfectly written RPM specs
Name:           perfecto
Version:        6.3.1
Release:        0%{?dist}
Group:          Development/Tools
License:        Apache License, Version 2.0
URL:            https://kaos.sh/perfecto

Source0:        https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:      checksum.sha512

BuildRoot:      %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:  make golang >= 1.23

Provides:       %{name} = %{version}-%{release}

################################################################################

%description
Tool for checking perfectly written RPM specs.

################################################################################

%prep
%crc_check
%autosetup

if [[ ! -d "%{name}/vendor" ]] ; then
  echo -e "----\nThis package requires vendored dependencies\n----"
  exit 1
elif [[ -f "%{name}/%{name}" ]] ; then
  echo -e "----\nSources must not contain precompiled binaries\n----"
  exit 1
fi

%build
pushd %{name}
  %{make_build} all
  cp LICENSE ..
popd

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name}/%{name} %{buildroot}%{_bindir}/

install -pDm 644 %{name}/common/perfecto.toml %{buildroot}%{_sysconfdir}/xdg/rpmlint/perfecto.toml

%post
if [[ -d %{_sysconfdir}/bash_completion.d ]] ; then
  %{name} --completion=bash 1> %{_sysconfdir}/bash_completion.d/%{name} 2>/dev/null
fi

if [[ -d %{_datarootdir}/fish/vendor_completions.d ]] ; then
  %{name} --completion=fish 1> %{_datarootdir}/fish/vendor_completions.d/%{name}.fish 2>/dev/null
fi

if [[ -d %{_datadir}/zsh/site-functions ]] ; then
  %{name} --completion=zsh 1> %{_datadir}/zsh/site-functions/_%{name} 2>/dev/null
fi

%postun
if [[ $1 == 0 ]] ; then
  if [[ -f %{_sysconfdir}/bash_completion.d/%{name} ]] ; then
    rm -f %{_sysconfdir}/bash_completion.d/%{name} &>/dev/null || :
  fi

  if [[ -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish ]] ; then
    rm -f %{_datarootdir}/fish/vendor_completions.d/%{name}.fish &>/dev/null || :
  fi

  if [[ -f %{_datadir}/zsh/site-functions/_%{name} ]] ; then
    rm -f %{_datadir}/zsh/site-functions/_%{name} &>/dev/null || :
  fi
fi

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE
%{_bindir}/%{name}
%{_sysconfdir}/xdg/rpmlint/perfecto.toml

################################################################################

%changelog
* Sat May 10 2025 Anton Novojilov <andy@essentialkaos.com> - 6.3.1-0
- Code refactoring
- Dependencies update

* Wed Sep 04 2024 Anton Novojilov <andy@essentialkaos.com> - 6.3.0-0
- ek package updated to v13
- Code refactoring
- Dependencies update
- rpmlint removed from dependencies

* Wed Jul 03 2024 Anton Novojilov <andy@essentialkaos.com> - 6.2.1-0
- Level changed for PF28 to notice

* Sun Jun 23 2024 Anton Novojilov <andy@essentialkaos.com> - 6.2.0-0
- Added check PF28 for checking summary tag length
- Code refactoring
- Dependencies update

* Thu Mar 28 2024 Anton Novojilov <andy@essentialkaos.com> - 6.1.1-0
- Improved support information gathering
- Code refactoring
- Dependencies update

* Tue Dec 19 2023 Anton Novojilov <andy@essentialkaos.com> - 6.1.0-0
- Added '-P'/'--pager' option to use pager for long output
- Improved verbose version info generation
- Code refactoring
- Dependencies update

* Thu Jul 20 2023 Anton Novojilov <andy@essentialkaos.com> - 6.0.0-0
- Added 'target' directive
- Improved XML render
- Improved JSON render
- Improved github render
- Improved terminal render
- Added extra info to report for XML and JSON renders
- Code refactoring

* Mon Jul 10 2023 Anton Novojilov <andy@essentialkaos.com> - 5.0.0-0
- -A/--absolve option renamed to -I/--ignore
- 'absolve' directive renamed to 'ignore'
- Apply 'ignore' directive to rpmlint alerts

* Sun Jul 09 2023 Anton Novojilov <andy@essentialkaos.com> - 4.1.4-0
- Fixed bug with printing help content if no specs provided
- Fixed bug in terminal render with printing spec name with problems

* Sun Jul 09 2023 Anton Novojilov <andy@essentialkaos.com> - 4.1.3-0
- Added custom configuration for rpmlint ≥ 2

* Sat Jul 08 2023 Anton Novojilov <andy@essentialkaos.com> - 4.1.2-0
- Fixed using colored output on CI

* Thu Jun 29 2023 Anton Novojilov <andy@essentialkaos.com> - 4.1.1-0
- Do not disable colors on CI
- Dependencies update

* Wed Nov 30 2022 Anton Novojilov <andy@essentialkaos.com> - 4.1.0-0
- Added verbose version info output
- Dependencies update
- Fixed build using sources from source.kaos.st

* Sun Sep 18 2022 Anton Novojilov <andy@essentialkaos.com> - 4.0.1-0
- Improve PF5 check

* Fri May 06 2022 Anton Novojilov <andy@essentialkaos.com> - 4.0.0-0
- Added autochangelog and autorelease macro support
- Added renderer for github actions
- Code refactoring
- UI improvements

* Wed Mar 30 2022 Anton Novojilov <andy@essentialkaos.com> - 3.7.2-0
- Removed pkg.re usage
- Added module info
- Added Dependabot configuration

* Mon Aug 16 2021 Anton Novojilov <andy@essentialkaos.com> - 3.7.1-0
- Fixed compatibility with the latest version of ek package

* Sat Apr 03 2021 Anton Novojilov <andy@essentialkaos.com> - 3.7.0-0
- Do not run RPMLint if it isn't installed

* Fri Apr 02 2021 Anton Novojilov <andy@essentialkaos.com> - 3.6.3-0
- Ignore links longer than 80 symbols in PF2

* Fri Feb 26 2021 Anton Novojilov <andy@essentialkaos.com> - 3.6.2-0
- Fixed bash completion
- Fixed zsh completion

* Sat Jun 13 2020 Anton Novojilov <andy@essentialkaos.com> - 3.6.1-0
- Fixed false-positive PF27 alert for Lua code

* Sat Jun 13 2020 Anton Novojilov <andy@essentialkaos.com> - 3.6.0-0
- Added check PF27 for checking unclosed if conditions in scriptlets
- Minor improvements

* Fri May 15 2020 Anton Novojilov <andy@essentialkaos.com> - 3.5.0-0
- Added check for chown and chmod usage in scriptlets
- ek package updated to the latest version

* Sat Feb 08 2020 Anton Novojilov <andy@essentialkaos.com> - 3.4.0-0
- Added check for useless dot at the end of the package summary

* Wed Feb 05 2020 Anton Novojilov <andy@essentialkaos.com> - 3.3.1-0
- Fixed bug with extracting sections
- Fixed false-positive alerts for PF24

* Wed Jan 29 2020 Anton Novojilov <andy@essentialkaos.com> - 3.3.0-0
- Added check for possible empty evaluated if clauses (PF24)

* Sat Dec 21 2019 Anton Novojilov <andy@essentialkaos.com> - 3.2.0-0
- Added printing links to wiki articles about failed checks
- Improved check for using variables instead of macros (PF5)

* Fri Dec 13 2019 Anton Novojilov <andy@essentialkaos.com> - 3.1.0-0
- Added URL check to PF20

* Sat Oct 26 2019 Anton Novojilov <andy@essentialkaos.com> - 3.0.0-0
- Improved all renderers
- Added check ID to output
- Added check for checking check scriptlet for using _without_check
  and _with_check macros
- Added check for single equals symbol in %%if clause

* Sat Oct 26 2019 Anton Novojilov <andy@essentialkaos.com> - 2.5.1-0
- Fixed bug with counting absolved alerts

* Fri Oct 25 2019 Anton Novojilov <andy@essentialkaos.com> - 2.5.0-0
- Added option for disabling some checks for entire spec
- Added check for HTTPS support on a source domain
- Added cache for HTTPS checks
- ek package updated to the latest version

* Tue Jul 09 2019 Anton Novojilov <andy@essentialkaos.com> - 2.4.0-0
- Added check for bash loops syntax
- Added completions generators

* Fri Jul 05 2019 Anton Novojilov <andy@essentialkaos.com> - 2.3.1-0
- Fixed bug with checking default paths without macro

* Mon Jun 10 2019 Anton Novojilov <andy@essentialkaos.com> - 2.3.0-0
- Added new checker for checking the number of empty lines at the end of
  the spec
- Improved spec parser
- Minor code refactoring

* Sun Nov 04 2018 Anton Novojilov <andy@essentialkaos.com> - 2.2.0-0
- Added quiet mode (option -q/--quiet)
- Improved RPMLint output parser

* Thu Jul 12 2018 Anton Novojilov <andy@essentialkaos.com> - 2.1.0-0
- Added mass check feature

* Fri Jun 15 2018 Anton Novojilov <andy@essentialkaos.com> - 2.0.2-0
- Improved check for unescaped percent symbol

* Wed May 09 2018 Anton Novojilov <andy@essentialkaos.com> - 2.0.1-0
- Fixed bug with wrong exit code when '--error-level error' is used

* Sun May 06 2018 Anton Novojilov <andy@essentialkaos.com> - 2.0.0-0
- Added short format support
- Added JSON format support
- Added XML format support
- Fixed bug with help content generation
- Code refactoring

* Thu Apr 05 2018 Anton Novojilov <andy@essentialkaos.com> - 1.3.0-0
- Improved check for output redirect to /dev/null

* Wed Mar 28 2018 Anton Novojilov <andy@essentialkaos.com> - 1.2.0-0
- Fixed bug with extracting sections names
- Fixed bug with extracting section data
- Added check for empty sections
- Added check for indent in files section
- Added check for setup section arguments

* Sat Feb 17 2018 Anton Novojilov <andy@essentialkaos.com> - 1.1.0-0
- Added check for useless binary macro usage
- Improved spec processing

* Wed Feb 14 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.2-0
- Fixed bug with selecting proper exit code if max alert level wasn't defined

* Thu Feb 08 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Added check for defattr macro in files section

* Wed Jan 31 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added support of tag 'perfecto:absolve' for "absolving" some alerts
- Improved spec parsing
- Improved checks

* Sun May 21 2017 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- First public release
