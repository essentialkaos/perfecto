################################################################################

# rpmbuilder:relative-pack true

################################################################################

%define  debug_package %{nil}

################################################################################

Summary:         Tool for checking perfectly written RPM specs
Name:            perfecto
Version:         1.0.1
Release:         0%{?dist}
Group:           Development/Tools
License:         EKOL
URL:             https://github.com/essentialkaos/perfecto

Source0:         https://source.kaos.io/%{name}/%{name}-%{version}.tar.bz2

BuildRoot:       %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

BuildRequires:   golang >= 1.9

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
* Thu Feb 08 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.1-0
- Added more checks

* Wed Jan 31 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Added support of tag 'perfecto:absolve' for "absolving" some alerts
- Improved spec parsing
- Improved checks

* Sun May 21 2017 Anton Novojilov <andy@essentialkaos.com> - 0.0.1-0
- First public release
