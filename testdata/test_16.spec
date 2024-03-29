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
%{__chkconfig} --add %{name} &>/dev/null || :

%preun
if [[ $? -ne 1 ]] ; then
  if [[ $? -gt 2 ]] ;then
    %{__chkconfig} --del %{name} &> /dev/null || :
    fi
fi

%postun
if [[ $? -ne 1 ]] ; then
  if [[ $? -gt 2 ]] ;then
    if [[ $? -lt 5 ]] ; then
    %{__chkconfig} --del %{name} &> /dev/null || :
fi

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
