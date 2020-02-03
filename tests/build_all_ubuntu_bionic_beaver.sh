#!/bin/bash
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
apt update
apt upgrade -y
snap install --classic go
apt-get install -y git build-essential python-dev gettext uuid-dev libncurses5-dev libyajl-dev libaio-dev pkg-config libglib2.0-dev libssl-dev libpixman-1-dev bridge-utils wget libfdt-dev bin86 bcc liblzma-dev iasl libc6-dev-i386 libglib2.0-dev libpixman-1-dev libcap-dev libattr1-dev automake libacl1-dev libsystemd-dev busybox-static jq libcap-ng-dev libelf-dev
cd /root || exit
git clone -b stable-4.13 git://xenbits.xen.org/xen.git
git clone git://git.qemu.org/qemu.git
git clone https://github.com/rkt/rkt.git
git clone https://github.com/rkt/stage1-xen.git
echo "XEN build"
cd xen || exit
./configure --prefix=/usr --with-system-qemu=/usr/lib/xen/bin/qemu-system-i386 --disable-stubdom --disable-qemu-traditional --disable-rombios
make -j4 PYTHON_PREFIX_ARG=
make install PYTHON_PREFIX_ARG=
update-rc.d xencommons defaults
echo 'GRUB_CMDLINE_XEN_DEFAULT="dom0_mem=1024M,max:1024M"'>>/etc/default/grub
sed -i 's/GRUB_DEFAULT=[0-9]*/GRUB_DEFAULT=2/g' /etc/default/grub
sed -i 's/GRUB_TIMEOUT_STYLE=.*/GRUB_TIMEOUT_STYLE=menu/g' /etc/default/grub
sed -i 's/GRUB_TIMEOUT=[0-9]*/GRUB_TIMEOUT=10/g' /etc/default/grub
update-grub
brctl addbr xenbr0
echo "QEMU build"
cd /root/qemu || exit
export DIR=/root/xen
./configure --enable-xen --target-list=i386-softmmu \
                --extra-cflags="-I$DIR/tools/include \
                -I$DIR/tools/libs/toollog/include \
                -I$DIR/tools/libs/evtchn/include \
                -I$DIR/tools/libs/gnttab/include \
                -I$DIR/tools/libs/foreignmemory/include \
                -I$DIR/tools/libs/devicemodel/include \
                -I$DIR/tools/libxc/include \
                -I$DIR/tools/xenstore/include \
                -I$DIR/tools/xenstore/compat/include" \
                --extra-ldflags="-L$DIR/tools/libxc \
                -L$DIR/tools/xenstore \
                -L$DIR/tools/libs/evtchn \
                -L$DIR/tools/libs/gnttab \
                -L$DIR/tools/libs/foreignmemory \
                -L$DIR/tools/libs/call \
                -L$DIR/tools/libs/devicemodel \
                -Wl,-rpath-link=$DIR/tools/libs/toollog \
                -Wl,-rpath-link=$DIR/tools/libs/evtchn \
                -Wl,-rpath-link=$DIR/tools/libs/gnttab \
                -Wl,-rpath-link=$DIR/tools/libs/call \
                -Wl,-rpath-link=$DIR/tools/libs/foreignmemory \
                -Wl,-rpath-link=$DIR/tools/libs/call \
                -Wl,-rpath-link=$DIR/tools/libs/devicemodel" \
                --disable-kvm --enable-virtfs
make -j4
make install
cp i386-softmmu/qemu-system-i386 /usr/lib/xen/bin/
echo "RKT build"
cd /root/rkt || exit
./autogen.sh
./configure --enable-sdjournal=no --disable-tpm --with-stage1-flavors=host
make
cp build-rkt*/target/bin/rkt /usr/sbin
echo "Stage1-xen build"
cd /root/stage1-xen || exit
wget https://raw.githubusercontent.com/lf-edge/eve/master/pkg/rkt-stage1/0001-Go-12-upgrade.patch
wget https://raw.githubusercontent.com/lf-edge/eve/master/pkg/rkt-stage1/0003-rkt-seed-xl.patch
wget https://raw.githubusercontent.com/lf-edge/eve/master/pkg/rkt-stage1/0004-Adding-STAGE1_XL_OPTS-and-fixing-one-nit.patch
wget https://raw.githubusercontent.com/lf-edge/eve/master/pkg/rkt-stage1/0006-Enable-local-networking-in-containers.patch
wget https://raw.githubusercontent.com/lf-edge/eve/master/pkg/rkt-stage1/0007-Set-up-environment-for-the-stage2-container.patch
wget https://gist.githubusercontent.com/giggsoff/9a6e57265279bce158ebbc2f60c29577/raw/9448cf52259b79d4ce53b21599be36357136e0bd/0008-Fix-dir-for-stage1.patch
for p in *.patch ; do patch -p1 < "$p" ; done
cp /bin/busybox /bin/busybox.static
export GOPATH=/root/go
mkdir -p /root/go/src
bash build.sh
cp stage1-xen.aci /root/
cp stage1-xen.aci /usr/sbin/
echo "Build done"
read -rsn1 -p"Press any key to reboot";echo
reboot now