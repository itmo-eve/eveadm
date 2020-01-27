#!/bin/sh
if ! [ $(id -u) = 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi
curdir=$PWD
home_dir=/tmp/pods
mkdir "$home_dir"
echo "$home_dir"
echo ========================================
IMAGE_HASH=$(./eveadm rkt create image --image-url=coreos.com/etcd:v2.0.0 --dir="$home_dir")
echo "$IMAGE_HASH"
echo ========================================
./eveadm rkt list image --dir="$home_dir"
echo ========================================
./eveadm rkt info image --dir="$home_dir" --image-hash="$IMAGE_HASH"
echo ========================================
systemd-run "$curdir"/eveadm rkt create --dir="$home_dir" --image-hash="$IMAGE_HASH" --no-overlay=true --stage1-path=""
echo ========================================
sleep 5
echo ========================================
CONTAINERS=$(./eveadm rkt list --dir="$home_dir" --no-legend=true)
echo "$CONTAINERS"
CONTAINER_UUID=$(echo "$CONTAINERS"|cut -f1)
echo ========================================
./eveadm rkt info --dir="$home_dir" --container-uuid="$CONTAINER_UUID"
echo ========================================
./eveadm rkt stop --dir="$home_dir" --container-uuid="$CONTAINER_UUID"
sleep 5
echo ========================================
./eveadm rkt delete --dir="$home_dir" --container-uuid="$CONTAINER_UUID"
echo ========================================
./eveadm rkt delete image --dir="$home_dir" --image-hash="$IMAGE_HASH"