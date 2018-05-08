################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Tool for checking perfectly written RPM specs
Name:            perfecto
Version:         2.0.1
Release:         0%{?dist}
Group:           Development/Tools
License:         EKOL
URL:             https://github.com/essentialkaos/perfecto

Source0:         https://source.kaos.st/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.10

Requires:        rpmlint

Provides:        %{name} = %{version}-%{release}

################################################################################

%description
Tool for checking perfectly written RPM specs.

################################################################################

%prep
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

################################################################################

%files
%defattr(-,root,root,-)
%doc LICENSE.EN LICENSE.RU
%{_bindir}/%{name}

################################################################################

%changelog
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
