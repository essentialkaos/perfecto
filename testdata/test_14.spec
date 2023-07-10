################################################################################

%{!?_without_check: %define _with_check 1}

################################################################################

Summary:            Test spec for perfecto.
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

%pre
if [[ $1 -ge 1 ]] ; then
%if 0%{?rhel} >= 7
  %{__systemctl} daemon-reload &>/dev/null || :
%endif
%{__service} %{name} restart &>/dev/null || :
fi

%post
if [[ $1 -eq 1 ]] ; then
  %if 0%{?rhel} >= 7
    %{__systemctl} enable %{name}.service &>/dev/null || :
  %else
    %{__chkconfig} --add %{name}
  %endif
fi

%preun
if [[ $1 -eq 0 ]] ; then
  %if 0%{?rhel} >= 7
    %{__systemctl} --no-reload disable %{name}.service &>/dev/null || :
    %{__systemctl} stop %{name}.service &>/dev/null || :
  %else
    %{__service} %{name} stop &>/dev/null || :
    %{__chkconfig} --del %{name}
  %endif
fi

if [[ $1 -eq 0 ]] ; then
  echo 1
fi

%postun
if [[ $1 -ge 1 ]] ; then
%if 0%{?rhel} >= 7
  %{__systemctl} daemon-reload &>/dev/null || :
%endif
fi

%postun magic
%if 0%{?rhel} >= 7
%if %{?_with_check:1}%{?_without_check:0}
  %{__systemctl} daemon-reload &>/dev/null || :
%endif
%endif

################################################################################

%files
%defattr(-,root,root,-)
%{_bindir}/%{name}

%files magic
%defattr(-,root,root,-)
%{_bindir}/%{name}-magic

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog record
