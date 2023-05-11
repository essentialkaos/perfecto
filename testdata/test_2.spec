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
Test subpackage for perfecto app.

################################################################################

%prep
%setup -qn %{name}-%{version}

%build
export CFLAGS="$RPM_OPT_FLAGS"
export LDFLAGS="$RPM_LD_FLAGS"
export PREFIX="%{_prefix}"
export LIBDIR="%{_libdir}"
export DOCDIR="$RPM_DOC_DIR"

export CURRENT_OS="$RPM_OS"
export CURRENT_ARCH="$RPM_ARCH"
export PKG_NAME="$RPM_PACKAGE_NAME"
export PKG_VER="$RPM_PACKAGE_VERSION"
export PKG_REL="$RPM_PACKAGE_RELEASE"

%make_build

rm -f $RPM_BUILD_DIR/output.log

%install
rm -rf %{buildroot}

export PATH="$PATH:/usr/sbin/test"

install -pm file $RPM_BUILD_ROOT/usr/
install -pm file2 %{buildroot}/etc/

cp $RPM_SOURCE_DIR/app.conf %{buildroot}/root/

wget github.com/etcd-io/etcd/Documentation

%{__sed} -i -e 's|/usr/bin/traceroute|/bin/traceroute|' %{buildroot}%{_sysconfdir}/file

rm -f %{buildroot}/%{_usr}/file >/dev/null 2>&1 || :
rm -f %{buildroot}%{_usr}/file 2>&1 >/dev/null
rm -f %{buildroot}%{_usr}/file >/dev/null  2>/dev/null
rm -f %{buildroot}%{_usr}/file 2>/dev/null >/dev/null || exit 0

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
* Thu Jan 25 2018 Anton Novojilov <andy@essential-kaos.com> - 1.0.1
- Test changelog record #2

* Wed Jan 24 2018 Anton Novojilov <andy@essential-kaos.com> 1.0.0-0
- Test changelog record
#- Test changelog record
