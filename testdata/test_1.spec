################################################################################

Summary:            Test spec for perfecto
Name:               perfecto-spec
Version:            1.0.0
Release:            0
Group:              System Environment/Base
License:            MIT 
URL:                https://domain.com
 
Source0:            https://domain.com/%{name}-%{version}.tar.gz

################################################################################

%description
Nam libero tempore, cum soluta nobis est eligendi option cumque nihil impedit quo minus
id quod maxime placeat facere possimus, omnis voluptas assumenda est.

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

################################################################################

%changelog
* Wed Jan 24 2018 Anton Novojilov <andy@essentialkaos.com> - 1.0.0-0
- At vero eos et accusamus et iusto odio dignissimos ducimus qui blanditiis praesentium
- https://github.com/organization/project/issues?q=milestone:versionAlpha+label:ReadyForReview
