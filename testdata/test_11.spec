################################################################################

Summary:            Test spec for perfecto
Name:               perfecto-spec
Version:            1.0.0
Release:            0%{?dist}
Group:              System Environment/Base
License:            MIT
URL:                https://domain.com

BuildRoot:          %{_tmppath}/%{name}-%{version}-%{release}-root-%(%{__id_u} -n)

Source0:            http://kaos.st/%{name}-%{version}-1.tar.gz

# perfecto:ignore
Source1:            http://kaos.st/%{name}-%{version}-2.tar.gz
Source2:            http://kaos.st/%{name}-%{version}-3.tar.gz
Source3:            %{name}-%{version}-3.tar.gz
Source4:            http://

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
rm -rf %{buildroot}

%check
%{make} check

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

%files magic
%defattr(-,root,root,-)
%{_bindir}/%{name}-magic

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog record
