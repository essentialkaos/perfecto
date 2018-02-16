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

%package magic

Summary:            Test subpackage for perfecto
Group:              System Environment/Base

%description magic
Test subpackge for perfecto app.

################################################################################

%package docs

Summary:            Test subpackage for perfecto
Group:              Documentation

%description docs
Test subpackge for perfecto app.

################################################################################

%prep
%setup -qn %{name}-%{version}

%build
%{__make} %{?_smp_mflags}

%install
%{__rm} -rf %{buildroot}

%{make_install} PREFIX=%{buildroot}%{_prefix}

%clean
# perfecto:absolve 2
rm -rf %{buildroot}

%post
%{__chkconfig} --add %{name} &>/dev/null || :

%preun
%{__chkconfig} --del %{name} &> /dev/null || :

%postun
%{__chkconfig} --del %{name} &> /dev/null || :

################################################################################

%files
%{_bindir}/%{name}

%files magic
%{_bindir}/%{name}-magic

%files docs
%defattr(-,root,root,-)
%{_docdir}/%{name}

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- Test changelog record
