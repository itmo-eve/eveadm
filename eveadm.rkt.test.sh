#!/bin/sh
if ! [ $(id -u) = 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi
curdir=$PWD
home_dir=/tmp/pods_test
rm -rf "$home_dir"
mkdir "$home_dir"
echo "$home_dir"
echo ========================================
echo "create image"
echo ========================================
IMAGE_HASH=$(./eveadm rkt create -i coreos.com/etcd:v2.0.0 --dir="$home_dir")
echo "$IMAGE_HASH"
echo ========================================
echo "list image"
echo ========================================
./eveadm rkt list -i --dir="$home_dir"
echo ========================================
echo "info image"
echo ========================================
./eveadm rkt info -i --dir="$home_dir" "$IMAGE_HASH"
echo ========================================
echo "create container"
echo ========================================
systemd-run "$curdir"/eveadm rkt create --dir="$home_dir" "$IMAGE_HASH" --no-overlay=true --stage1-path=""
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "start container"
echo ========================================
./eveadm rkt start --dir="$home_dir" --stage1-type=common "$CONTAINER_UUID"
echo ========================================
echo "list container"
echo ========================================
until [ "$CONTAINERS" ]
do
CONTAINERS=$(./eveadm rkt list --dir="$home_dir" --no-legend=true)
done
echo "$CONTAINERS"
CONTAINER_UUID=$(echo "$CONTAINERS"|cut -f1)
echo ========================================
echo "info container"
echo ========================================
./eveadm rkt info --dir="$home_dir" "$CONTAINER_UUID"
echo ========================================
echo "stop container"
echo ========================================
./eveadm rkt stop --dir="$home_dir" "$CONTAINER_UUID"
sleep 5
echo ========================================
echo "delete container"
echo ========================================
./eveadm rkt delete --dir="$home_dir" "$CONTAINER_UUID"
echo ========================================
echo "delete image"
echo ========================================
./eveadm rkt delete -i --dir="$home_dir" "$IMAGE_HASH"