################################################################################

Summary:            Test spec for perfecto
Name:               perfecto-spec
Version:            1.0.0
Release:            0%{?dist}
Group:              System Environment/Base
License:            MIT
URL:                https://domain.com

Source0:            https://domain.com/%{name}-%{version}.tar.gz

################################################################################

%description
Test spec for perfecto app.

################################################################################

%if 1
%package magic

Summary:            Test subpackage for perfecto
Group:              System Environment/Base

%description magic
Test subpackage for perfecto app.
%endif

################################################################################

%prep
%setup -q -c -n %{name}-%{version}

%define _system /usr/system

%build
%{__make} %{?_smp_mflags}

%install
rm -rf %{buildroot}

%{make_install} PREFIX=%{buildroot}%{_prefix}

%clean
rm -rf %{buildroot}

%post
%{__chkconfig} --add %{name} &>/dev/null || :

%preun
%{__chkconfig} --del %{name} &> /dev/null || :

%postun
%{__chkconfig} --del %{name} &> /dev/null || :

################################################################################

%files
%defattr(-,root,root,-)
%{_bindir}/%{name}

###############################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog %record
