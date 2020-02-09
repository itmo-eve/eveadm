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
adam_dir="$tmp_dir"/adam
eve_dir="$tmp_dir"/eve
apt update
apt upgrade -y
snap install --classic go
apt-get install -y git make docker.io qemu-system-x86 qemu-utils openssl jq
touch ~/.rnd
cd "$tmp_dir" || exit
git clone https://github.com/itmo-eve/eve.git
git clone https://github.com/giggsoff/adam.git
echo ========================================
echo "Prepare and run ADAM"
echo ========================================
IP=`hostname -I|cut -d' ' -f1`
cd $adam_dir || exit
make build-docker
mkdir -p run/adam
mkdir -p run/config
cp "$config_file" run/
cd run/adam||exit
onboarduuid=$(uuidgen)
openssl genrsa -out rootCA.key 4096
openssl req -x509 -new -nodes -key rootCA.key -sha256 -subj "/C=RU/ST=SPB/O=MyOrg, Inc./CN=test" -days 1024 -out rootCA.crt
openssl ecparam -name prime256v1 -genkey -out server-key.pem
openssl req -new -sha256 -key server-key.pem -subj "/C=RU/ST=SPB/O=MyOrg, Inc./CN=mydomain.com" -reqexts SAN -config <(cat /etc/ssl/openssl.cnf <(printf "\n[SAN]\nsubjectAltName=DNS:mydomain.com,IP:$IP")) -out server.csr
openssl x509 -req -extfile <(printf "subjectAltName=DNS:mydomain.com,IP:$IP") -days 365 -in server.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out server.pem
openssl ecparam -name prime256v1 -genkey -out onboard.key
openssl req -new -sha256 -key onboard.key -subj "/C=RU/ST=SPB/O=MyOrg, Inc./CN=$onboarduuid" -out onboard.pem.csr
openssl x509 -req -in onboard.pem.csr -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out onboard.pem -days 500 -sha256
cp rootCA.crt ../config/root-certificate.pem
cp onboard.pem ../config/onboard.cert.pem
cp onboard.key ../config/onboard.key.pem
echo "$IP" mydomain.com>../config/hosts
echo mydomain.com:$unused_port>../config/server
sudo chmod 644 ../config/*.pem
cd "$adam_dir" || exit
nohup docker run -v "$adam_dir"/run:/adam/run -p $unused_port:8080 lfedge/adam server --conf-dir ./run/config/adam >$tmp_dir/adam.log 2>&1 &
echo $! >../adam.pid
echo ========================================
echo "ADAM pid:"
cat ../adam.pid
echo ========================================
max_retry=5
counter=0
until docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":$unused_port onboard add --path /adam/run/config/onboard.cert.pem; do
  [[ counter -eq $max_retry ]] && echo "Failed to add onboard!" && exit 1
  sleep 5
  echo "Trying again. Try #$counter"
  ((counter++))
done
echo '*' >run/adam/onboard/$onboarduuid/onboard-serials.txt
echo ========================================
echo "Prepare and run EVE"
echo ========================================
cd $eve_dir||exit
sed -i "s/SandyBridge/host/g" Makefile
make CONF_DIR=../adam/run/config/ live
nohup make ACCEL=true SSH_PORT=$ssh_port CONF_DIR=../adam/run/config/ run >$tmp_dir/eve.log 2>&1 &
echo $! >../eve.pid
echo ========================================
echo "EVE pid:"
cat ../eve.pid
echo ========================================
echo "Try to modify EVE config"
echo ========================================
cd $adam_dir||exit
UUID="$(docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":$unused_port device list)"
max_retry=20
counter=0
until [ "$UUID" ]
do
  [[ counter -eq $max_retry ]] && echo "Failed to add onboard!" && exit 1
  echo "Trying again. Try #$counter"
  sleep 30
  UUID="$(docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":$unused_port device list)"
  ((counter++))
done
UUID=$(echo "$UUID" | xargs)
echo ========================================
echo "EVE device UUID:"
echo $UUID
echo ========================================
sed -i "s/DEVICE_UUID/$UUID/g" run/cfg.json
docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":$unused_port device config set --uuid "$UUID" --config-path ./run/cfg.json
echo ========================================
echo "EVE config successfull"
read -rsn1 -p"Press any key to cleanup";echo
kill `cat $tmp_dir/eve.pid`
kill `cat $tmp_dir/adam.pid`
sleep 5
if [ -n "$use_custom_dir" ]
then
rm -rf $use_custom_dir/*
else
rm -rf $tmp_dir
fi
