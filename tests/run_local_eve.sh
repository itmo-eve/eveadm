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
ssh_port=`comm -23 <(seq 49252 49352 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1`
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
cd conf||exit
onboarduuid=$(uuidgen)
sn=$(head -200 /dev/urandom | cksum | cut -f1 -d " "|fold -w 8|head -n 1)
openssl ecparam -name prime256v1 -genkey -noout -out onboard.key.pem
openssl req -x509 -new -nodes -key onboard.key.pem -sha256 -subj "/C=RU/ST=SPB/O=MyOrg, Inc./CN=$onboarduuid" -days 1024 -out onboard.cert.pem
cd $eve_dir||exit
sed -i "s/SandyBridge/host/g" Makefile
sed -i "s/31415926/$sn/g" Makefile
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
echo "$onboarduuid"
echo "and Serial Number"
echo "$sn"
echo "in zedcloud.zededa.net"
read -p "Do you want to cleanup? (y/n)" -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
  kill `cat $tmp_dir/eve.pid`
  sleep 5
  if [ -n "$use_custom_dir" ]
  then
    rm -rf $use_custom_dir/*
  else
    rm -rf $tmp_dir
  fi
fi
