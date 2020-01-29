################################################################################

# rpmbuilder:relative-pack true

################################################################################

%global crc_check pushd ../SOURCES ; sha512sum -c %{SOURCE100} ; popd

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Tool for checking perfectly written RPM specs
Name:            perfecto
Version:         3.3.0
Release:         0%{?dist}
Group:           Development/Tools
License:         EKOL
URL:             https://kaos.sh/perfecto

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

Source100:       checksum.sha512

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.13

Requires:        rpmlint

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Tool for checking perfectly written RPM specs.

################################################################################

%prep
%{crc_check}

%setup -q

%build
export GOPATH=$(pwd)
go build src/github.com/essentialkaos/%{name}/%{name}.go

%install
rm -rf %{buildroot}

install -dm 755 %{buildroot}%{_bindir}
install -pm 755 %{name} %{buildroot}%{_bindir}/

%clean
rm -rf %{buildroot}

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
%doc LICENSE.EN LICENSE.RU
%{_bindir}/%{name}

################################################################################

%changelog
* Wed Jan 29 2020 Anton Novojilov <andy@essentialkaos.com> - 3.3.0-0
- Added check for possible empty evaluated if clauses (PF24)

* Sat Dec 21 2019 Anton Novojilov <andy@essentialkaos.com> - 3.2.0-0
- Added printing links to wiki articles about failed checks
- Improved check for using variables instead of macroses (PF5)

* Fri Dec 13 2019 Anton Novojilov <andy@essentialkaos.com> - 3.1.0-0
- Added URL check to PF20

* Sat Oct 26 2019 Anton Novojilov <andy@essentialkaos.com> - 3.0.0-0
- Improved all renderers
- Added check ID to output
- Added check for checking check scriptlet for using _without_check
  and _with_check macroses
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
