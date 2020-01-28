#!/bin/sh
if ! [ $(id -u) = 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi
home_dir=/tmp/xen_test
rm -rf "$home_dir"
mkdir "$home_dir"
echo "$home_dir"
name="testxen"
echo ========================================
echo "download image and create config"
echo ========================================
wget http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img -O "$home_dir"/cirros.qcow2
cat << EOF > "$home_dir"/config.cfg
name = "$name"
kernel = "/boot/vmlinuz"
extra = "root=/dev/xvda1"
memory = 128
vcpus = 1
vif = [ '' ]
disk = [ '$home_dir/cirros.qcow2,xvda,rw' ]
EOF
echo ========================================
echo "create vm"
echo ========================================
./eveadm xen create --xen-cfg-filename="$home_dir"/config.cfg --paused
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "domid vm"
echo ========================================
domuuid=$(./eveadm xen info --domname $name)
echo "$domuuid"
echo ========================================
echo "start vm"
echo ========================================
./eveadm xen start "$domuuid"
echo ========================================
echo "vm info"
echo ========================================
./eveadm xen info "$domuuid"
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "vm stop"
echo ========================================
./eveadm xen stop "$domuuid"
echo ========================================
echo "vm delete"
echo ========================================
./eveadm xen delete "$domuuid"
echo ========================================