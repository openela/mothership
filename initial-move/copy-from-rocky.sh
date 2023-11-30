#!/usr/bin/env bash
set -euo pipefail

COMMAND="${1:-}"
shift

if [[ -z "${COMMAND}" ]]; then
  echo "Usage: $0 <command>"
  exit 1
fi

# Latest branch
BRANCH_8="el-8.8"
BRANCH_9="el-9.2"

common_repos=(
  libguestfs
  cockpit
  abrt
  anaconda
  anaconda-user-help
  cloud-init
  crash
  dhcp
  dnf
  firefox
  fwupd
  gcc
  gdb
  gnome-session         # This may occur in r10 - does not apply to r8
  gnome-settings-daemon
  grub2
  initial-setup
  kernel
  kernel-rt
  libdnf
  libreoffice
  libreport
  lorax-templates-rhel
  nginx
  openscap
  osbuild               # still required, since osbuild is not smart enough to use the rhel/centos runner
  osbuild-composer      # unified patches, all should match (exception is maintenance mode)
  PackageKit
  python-pip
  redhat-rpm-config
  scap-security-guide   # unified patches
  shim
  shim-unsigned-x64
  shim-unsigned-aarch64
  subscription-manager
  subscription-manager-cockpit
  systemd
  thunderbird
  WALinuxAgent          ## temporary until 2.8.x.y
)

r8=(
  beignet
  ceph
  clufter
  compat-guile18
  efibootmgr
  fontforge
  gegl04
  gflags
  gomtree
  ibus-typing-booster
  jss
  kyotocabinet
  libucil
  mariadb
  mesa-libGLU
  nss-softokn
  perl-AnyEvent
  perl-Net-IDN-Encode
  perl-TimeDate
  python-cffi
  python-cups
  python-iso8601
  python-systemd
  rpm-ostree
  xorg-x11-drv-qxl
  dotnet
  dotnet3.0
  gnome-boxes
  libguestfs            # temporary
  oscap-anaconda-addon  # change "Red Hat" to "Rocky"
  pcs                   # Might be needed for 9
  plymouth
  python2
  python3
)

r9=(
  owasp-java-encoder
  mariadb
  nss
  flexiblas
  munge
  gnome-session
  mysql
  systemd
  ruby
)

r8_stream_1_0=(
  cobbler
  rhncfg
  rhnpush
  rhn-custom-info
  rhn-virtualization
  spacewalk-abrt
  spacewalk-backend
  spacewalk-client-cert
  spacewalk-koan
  spacewalk-oscap
  spacewalk-remote-utils
  spacewalk-usix
)

r8_stream_10_3=(
  mariadb
)

r8_stream_2=(
  disruptor
  jctools
  log4j
)

r8_stream_2_6=(
  ruby
)

r8_stream_2_7=(
  python-psycopg2
)

r8_stream_5_24=(
  perl
  perl-TimeDate
)

r8_stream_5_26=(
  perl
  perl-TimeDate
)

r8_stream_8_0=(
  mysql
)

# All packages that were declared in all repos above (including stream repos)
all_packages=(
  "${common_repos[@]}"
  "${r8[@]}"
  "${r9[@]}"
  "${r8_stream_1_0[@]}"
  "${r8_stream_10_3[@]}"
  "${r8_stream_2[@]}"
  "${r8_stream_2_6[@]}"
  "${r8_stream_2_7[@]}"
  "${r8_stream_5_24[@]}"
  "${r8_stream_5_26[@]}"
  "${r8_stream_8_0[@]}"
)

unique_packages=($(for pkg in "${all_packages[@]}"; do echo "${pkg}"; done | sort -u))

copy_patches() {
  PKG="${1}"

  # Clone into a temporary directory, delete the directory on exit
  TMPDIR=$(mktemp -d)
  trap "rm -rf ${TMPDIR}" EXIT

  # Clone the repo
  git clone "https://git.rockylinux.org/staging/patch/$PKG" "${TMPDIR}/${PKG}" || return
  pushd "${TMPDIR}/${PKG}"

  # List all branches that start with "r\d"
  # If no branches are found, then return
  BRANCHES=()
  NO_CHECKOUT=0
  if ! git branch -a | sed "s/.*remotes\/origin\///g" | grep -E -i "^r\d(-stream-.+|$)" >/dev/null 2>&1; then
    # Check if main branch exists
    if ! git branch -a | sed "s/.*remotes\/origin\///g" | grep -E -i "^main$" >/dev/null 2>&1; then
      echo "No branches found for ${PKG}"
      return
    else
      BRANCHES=($BRANCH_8 $BRANCH_9)
      NO_CHECKOUT=1
    fi
  else
    BRANCHES=($(git branch -a | sed "s/.*remotes\/origin\///g" | grep -E -i "^r\d(-stream-.+|$)" | sort -u))
  fi

  # Clone the target/destination repo
  gh repo clone "openela-main/$PKG" "${TMPDIR}/${PKG}-openela-main"

  # For every branch, create a new branch in the destination repo with the new name
  for branch in "${BRANCHES[@]}"; do
    # Replace "rX" with "${BRANCH}", where branch is decided by X ($BRANCH_X)
    if [[ "${NO_CHECKOUT}" == "0" ]]; then
      major_version=$(echo "${branch}" | sed "s/r\([0-9]\).*/\1/g")
      NEW_BRANCH="$(eval echo "\$BRANCH_${major_version}")"
    else
      NEW_BRANCH="${branch}"
    fi

    # Checkout branch, then pushd to the destination repo, create a new branch, then popd back to the source repo
    if [[ "${NO_CHECKOUT}" == "0" ]]; then
      git checkout "${branch}"
    fi
    pushd "${TMPDIR}/${PKG}-openela-main"
    # Check if the branch already exists, then delete it
    if git rev-parse --verify "${NEW_BRANCH}" >/dev/null 2>&1; then
      git checkout -b main || git checkout main
      git branch -D "${NEW_BRANCH}"
    fi
    git checkout --orphan "${NEW_BRANCH}"
    git rm -rf . || true
    rm -rf * || true
    mkdir PATCHES
    cp -r ${TMPDIR}/${PKG}/* PATCHES/
    pushd PATCHES
    mv ROCKY/CFG/*.cfg . || true
    mv ROCKY/_supporting/* . || true
    rm -rf ROCKY || true
    rm -rf ROCKY_REMOVE || true
    # In .cfg files do the following replaces (case sensitive):
    #   - "ROCKY/_supporting" -> "PATCHES"
    #   - "ROCKY/CFG" -> "PATCHES"
    sed -i.bak "s/ROCKY\/_supporting/PATCHES/g" *.cfg || true
    sed -i.bak "s/ROCKY\/CFG/PATCHES/g" *.cfg || true
    # In all files do the following replaces (case sensitive):
    #   - "rockylinux.org" -> "openela.org"
    #   - "Rocky Linux" -> "OpenELA"
    #   - "ROCKY" -> "OpenELA"
    #   - "Rocky" -> "OpenELA"
    #   - "rocky" -> "openela"
    #   - "openela.org>" -> "rockylinux.org>"
    sed -i.bak "s/rockylinux.org/openela.org/g" * || true
    sed -i.bak "s/Rocky Linux/OpenELA/g" * || true
    sed -i.bak "s/ROCKY/OpenELA/g" * || true
    sed -i.bak "s/Rocky/OpenELA/g" * || true
    sed -i.bak "s/rocky/openela/g" * || true
    sed -i.bak "s/openela.org>/rockylinux.org>/g" * || true
    # Replace occurrences in file names:
    #   - "rocky-linux" -> "openela"
    #   - "rocky" -> "openela"
    #   - "Rocky-Linux" -> "OpenELA"
    #   - "Rocky" -> "OpenELA"
    #   - "ROCKY" -> "OpenELA"
    rename "s/rocky-linux/openela/g" * || true
    rename "s/rocky/openela/g" * || true
    rename "s/Rocky-Linux/OpenELA/g" * || true
    rename "s/Rocky/OpenELA/g" * || true
    rename "s/ROCKY/OpenELA/g" * || true
    # Remove all .bak files
    rm -f *.bak
    popd
    git add .
    # Verify that there are changes to commit
    if [[ -z "$(git status --porcelain)" ]]; then
      echo "No changes for ${PKG} on branch ${branch}"
      popd
      continue
    fi
    git commit -m "Copy patches from Rocky Linux"
    git push -u origin "${NEW_BRANCH}" -f
    popd
  done
}

if [[ "${COMMAND}" == "list" ]]; then
  for pkg in "${unique_packages[@]}"; do
    echo "${pkg}"
  done
elif [[ "${COMMAND}" == "copy-patches" ]]; then
  ONE_PKG="${1:-}"

  if [[ -z "${ONE_PKG}" ]]; then
    for pkg in "${unique_packages[@]}"; do
      echo "Copying patches for ${pkg}"
      copy_patches "${pkg}"
    done
  else
    copy_patches "${ONE_PKG}"
  fi
elif [[ "${COMMAND}" == "create-repos" ]]; then
  for pkg in "${unique_packages[@]}"; do
    echo "Creating repo for ${pkg}"
    gh repo create --private --disable-issues --disable-wiki "openela-main/$pkg" || echo "Repo already exists"
  done
else
  echo "Invalid command: ${COMMAND}"
  exit 1
fi
