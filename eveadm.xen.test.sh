#!/bin/sh
if ! [ $(id -u) = 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi
echo 'home_dir=/tmp/xen_test
rm -rf "$home_dir"
mkdir "$home_dir"
echo "$home_dir"
name="testxen"'
home_dir=/tmp/xen_test
rm -rf "$home_dir"
mkdir "$home_dir"
echo "$home_dir"
name="testxen"
echo ========================================
echo "download image and create config"
echo ========================================
echo 'wget http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img -O "$home_dir"/cirros.qcow2'
wget http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img -O "$home_dir"/cirros.qcow2
echo 'cat << EOF > "$home_dir"/config.cfg
name = "$name"
bootloader = "pygrub"
extra = "console=hvc0 root=/dev/xvda1"
memory = 128
vcpus = 1
vif = [ '\''bridge=xenbr0'\'' ]
disk = [ '\''$home_dir/cirros.qcow2,qcow2,xvda,rw'\'' ]
EOF'
cat << EOF > "$home_dir"/config.cfg
name = "$name"
bootloader = "pygrub"
extra = "console=hvc0 root=/dev/xvda1"
memory = 128
vcpus = 1
vif = [ 'bridge=xenbr0' ]
disk = [ '$home_dir/cirros.qcow2,qcow2,xvda,rw' ]
EOF
echo 'brctl show|grep xenbr0||brctl addbr xenbr0'
brctl show|grep xenbr0||brctl addbr xenbr0
echo ========================================
echo "create vm"
echo ========================================
echo './eveadm xen create --xen-cfg-filename="$home_dir"/config.cfg --paused'
./eveadm xen create --xen-cfg-filename="$home_dir"/config.cfg --paused
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "domid vm"
echo ========================================
echo 'domuuid=$(./eveadm xen info --domname $name)'
domuuid=$(./eveadm xen info --domname $name)
echo 'echo "$domuuid"'
echo "$domuuid"
echo ========================================
echo "start vm"
echo ========================================
echo './eveadm xen start "$domuuid"'
./eveadm xen start "$domuuid"
echo ========================================
echo "vm info"
echo ========================================
echo './eveadm xen info "$domuuid"'
./eveadm xen info "$domuuid"
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "vm stop"
echo ========================================
echo './eveadm xen stop "$domuuid"'
./eveadm xen stop "$domuuid"
echo ========================================
echo "vm delete"
echo ========================================
echo './eveadm xen delete "$domuuid"'
./eveadm xen delete "$domuuid"
echo ========================================