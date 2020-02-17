#!/bin/bash
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi
eve_repo=https://github.com/itmo-eve/eve.git
adam_repo=https://github.com/giggsoff/adam.git
memory_to_use=4096
while [ -n "$1" ]
do
case "$1" in
-m) memory_to_use="$2"
echo "Use with memory $memory_to_use"
shift ;;
--) shift
break ;;
*) echo "$1 is not an option";;
esac
shift
done
tmp_dir=$(mktemp -d -t eveadam-"$(date +%Y-%m-%d-%H-%M-%S)"-XXXXXXXXXX)
unused_port=$(comm -23 <(seq 49152 49252 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1)
ssh_port=$(comm -23 <(seq 49252 49352 | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf | head -n 1)
config_file="$PWD"/cfg.json
echo ========================================
echo "Temp directory for test: $tmp_dir"
echo ========================================
subnet1_prefix=""
subnet2_prefix=""
for ((i = 0; i <= 254; i++)); do
  ip a | grep -F -q "192.168.$i"
  if [[ $? -ne 0 ]]; then
    if [ -z $subnet1_prefix ]; then
      subnet1_prefix="192\.168\.$i"
      continue
    fi
    if [ -z $subnet2_prefix ]; then
      subnet2_prefix="192\.168\.$i"
      continue
    fi
    break
  fi
done
adam_dir="$tmp_dir"/adam
eve_dir="$tmp_dir"/eve
apt update
apt upgrade -y
snap install --classic go
apt-get install -y git make docker.io qemu-system-x86 qemu-utils openssl jq
touch ~/.rnd
cd "$tmp_dir" || exit
git clone $eve_repo
git clone $adam_repo
echo ========================================
echo "Generate keypair for ssh (no overwrite if exists)"
echo ========================================
ssh-keygen -t rsa -f /root/.ssh/id_rsa -q -N "" <<< n
echo
echo ========================================
echo "Prepare and run ADAM"
echo ========================================
IP=$(hostname -I | cut -d' ' -f1)
cd $adam_dir || exit
#make build-docker
mkdir -p run/adam
mkdir -p run/config
cp "$config_file" run/
cd run/adam || exit
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
echo "$IP" mydomain.com >../config/hosts
echo mydomain.com:$unused_port >../config/server
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
cd $eve_dir || exit
sed -i "s/eth0,net=192\.168\.1\.0\/24,dhcpstart=192\.168\.1\.10/eth0,net=$subnet1_prefix\.0\/24,dhcpstart=$subnet1_prefix\.10/g" Makefile
sed -i "s/eth1,net=192\.168\.2\.0\/24,dhcpstart=192\.168\.2\.10/eth1,net=$subnet2_prefix\.0\/24,dhcpstart=$subnet2_prefix\.10/g" Makefile
sed -i "s/SandyBridge/host/g" Makefile
sed -i "s/-m 4096/-m $memory_to_use/g" Makefile
make CONF_DIR=../adam/run/config/ live
nohup make ACCEL=true SSH_PORT=$ssh_port CONF_DIR=../adam/run/config/ run >$tmp_dir/eve.log 2>&1 &
echo $! >../eve.pid
echo ========================================
echo "EVE pid:"
cat ../eve.pid
echo ========================================
echo "Try to modify EVE config"
echo ========================================
cd $adam_dir || exit
UUID="$(docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":"$unused_port" device list)"
max_retry=30
counter=0
until [ "$UUID" ]; do
  [[ counter -eq $max_retry ]] && echo "Failed to add onboard!" && exit 1
  echo "Trying again. Try #$counter"
  sleep 30
  UUID="$(docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":"$unused_port" device list)"
  ((counter++))
done
UUID=$(echo "$UUID" | xargs)
echo ========================================
echo "EVE device UUID:"
echo $UUID
echo ========================================
sed -i "s/DEVICE_UUID/$UUID/g" run/cfg.json
sed -i -e "s/SSH_KEY/$(sed 's:/:\\/:g' /root/.ssh/id_rsa.pub)/" run/cfg.json
docker run -v "$adam_dir"/run:/adam/run lfedge/adam admin --server https://"$IP":$unused_port device config set --uuid "$UUID" --config-path ./run/cfg.json
echo ========================================
echo "EVE config successfull"
echo "You can connect to node via ssh"
echo "sudo ssh -p $ssh_port 127.0.0.1"
while true; do
    read -p "Do you want to cleanup? (y/n)" yn
    case $yn in
        [Yy]* )
          kill `cat $tmp_dir/eve.pid`
          kill $(cat $tmp_dir/adam.pid)
          sleep 5
          rm -rf $tmp_dir ; break;;
        [Nn]* ) exit;;
        * ) echo "Please answer y or n.";;
    esac
done
