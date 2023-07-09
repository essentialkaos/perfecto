################################################################################

%{!?_without_check: %define _with_check 1}

################################################################################

Summary:            Test spec for perfecto
Name:               perfecto-spec
Version:            1.0.0
Release:            0%{?dist}
Group:              System Environment/Base
License:            MIT
URL:                https://domain.com

BuildRoot:          %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

Source0:            https://domain.com/%{name}-%{version}.tar.gz

# perfecto:ignore
Source1:            http://domain.com/%{name}-%{version}.tar.gz

################################################################################

%description
Test spec for perfecto app.

################################################################################

%package magic

Summary:            Test subpackage for perfecto
Group:              System Environment/Base

%description magic
Test subpackage for perfecto app.

################################################################################

%prep
%setup -qn %{name}-%{version}

%build
%{__make} %{?_smp_mflags}

%install
rm -rf %{buildroot}

%{make_install} PREFIX=%{buildroot}%{_prefix}

%clean
# perfecto:ignore 2
rm -rf %{buildroot}

%check
%if %{?_with_check:1}%{?_without_check:0}
%{make} check
%endif

%post
chmod 0755 %{_datadir}/%{name}
chown operator: %{_datadir}/%{name}
chown -h operator: %{_datadir}/%{name}
chown --no-dereference operator: %{_datadir}/%{name}

%preun
%{__chkconfig} --del %{name} &> /dev/null || :

%postun
%{__chkconfig} --del %{name} &> /dev/null || :

################################################################################

%files
%defattr(-,root,root,-)
%{_bindir}/%{name}
%{_datadir}/%{name}

%files magic
%defattr(-,root,root,-)
%{_bindir}/%{name}-magic

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog record
