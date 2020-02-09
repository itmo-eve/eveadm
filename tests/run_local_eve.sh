#!/bin/bash
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi
use_custom_dir=$1
if [ -n "$use_custom_dir" ]
then
tmp_dir=$use_custom_dir
else
tmp_dir=$(mktemp -d -t eveadam-"$(date +%Y-%m-%d-%H-%M-%S)"-XXXXXXXXXX)
fi
unused_port=`comm -23 <(seq 49152 49252 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1`
ssh_port=`comm -23 <(seq 49252 49352 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1`
config_file="$PWD"/cfg.json
echo ========================================
echo "Temp directory for test: $tmp_dir"
echo ========================================
eve_dir="$tmp_dir"/eve
apt update
apt upgrade -y
snap install --classic go
apt-get install -y git make docker.io qemu-system-x86 qemu-utils openssl jq
touch ~/.rnd
cd "$tmp_dir" || exit
git clone https://github.com/itmo-eve/eve.git
echo ========================================
echo "Prepare and run EVE"
echo ========================================
cd $eve_dir||exit
sed -i "s/SandyBridge/host/g" Makefile
make live
nohup make ACCEL=true SSH_PORT=$ssh_port run >$tmp_dir/eve.log 2>&1 &
echo $! >../eve.pid
echo ========================================
echo "EVE pid:"
cat ../eve.pid
echo ========================================
echo "EVE config successfull"
echo ========================================
echo "Please use Onboarding Key"
echo "5d0767ee-0547-4569-b530-387e526f8cb9"
echo "and Serial Number"
echo "31415926"
echo "in zedcloud.zededa.net"
read -rsn1 -p"Press any key to cleanup";echo
kill `cat $tmp_dir/eve.pid`
sleep 5
if [ -n "$use_custom_dir" ]
then
rm -rf $use_custom_dir/*
else
rm -rf $tmp_dir
fi
