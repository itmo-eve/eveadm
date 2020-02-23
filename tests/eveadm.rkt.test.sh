#!/bin/sh
if [ $(id -u) != 0 ]; then
   echo "The script need to be run as root." >&2
   exit 1
fi

curdir=$(realpath $0)
echo curdir=$curdir
curdir=$(dirname $curdir)
echo curdir=$curdir
EVEADM=$curdir/../eveadm
echo EVEADM=$EVEADM
home_dir=/tmp/pods_test
echo "$home_dir"
echo rm -rf "$home_dir"
rm -rf "$home_dir"
echo mkdir "$home_dir"
mkdir "$home_dir"

echo ========================================
echo "create image"
echo ========================================
echo IMAGE_HASH=$($EVEADM rkt create -i coreos.com/etcd:v2.0.0 --dir="$home_dir")
IMAGE_HASH=$($EVEADM rkt create -i coreos.com/etcd:v2.0.0 --dir="$home_dir")
echo 'echo $IMAGE_HASH'
echo "$IMAGE_HASH"
echo ========================================
echo "list image"
echo ========================================
echo $EVEADM rkt list -i --dir="$home_dir"
$EVEADM rkt list -i --dir="$home_dir"
echo ========================================
echo "info image"
echo ========================================
echo $EVEADM rkt info -i --dir="$home_dir" "$IMAGE_HASH"
$EVEADM rkt info -i --dir="$home_dir" "$IMAGE_HASH"
echo ========================================
echo "create container"
echo ========================================
echo systemd-run $EVEADM rkt create --dir="$home_dir" "$IMAGE_HASH" --no-overlay=true --stage1-path="$curdir/stage1-xen.aci"
systemd-run $EVEADM rkt create --dir="$home_dir" "$IMAGE_HASH" --no-overlay=true --stage1-path="$curdir/stage1-xen.aci"
echo ========================================
echo "sleep 5"
echo ========================================
sleep 5
echo ========================================
echo "list container"
echo ========================================
echo 'until [ "$CONTAINERS" ]
do'
echo CONTAINERS=$($EVEADM rkt list --dir="$home_dir" --no-legend=true)
echo done
until [ "$CONTAINERS" ]
do
CONTAINERS=$($EVEADM rkt list --dir="$home_dir" --no-legend=true)
done
echo 'echo "$CONTAINERS"'
echo "$CONTAINERS"
echo 'CONTAINER_UUID=$(echo "$CONTAINERS"|cut -f1)'
CONTAINER_UUID=$(echo "$CONTAINERS"|cut -f1)
echo 'echo "$CONTAINER_UUID"'
echo "$CONTAINER_UUID"
echo ========================================
echo "start container"
echo ========================================
echo $EVEADM rkt start --dir="$home_dir" --stage1-type=common "$CONTAINER_UUID"
$EVEADM rkt start --dir="$home_dir" --stage1-type=common "$CONTAINER_UUID"
echo 'sleep 5'
sleep 5
echo ========================================
echo "info container"
echo ========================================
echo $EVEADM rkt info --dir="$home_dir" "$CONTAINER_UUID"
$EVEADM rkt info --dir="$home_dir" "$CONTAINER_UUID"
echo ========================================
echo "stop container"
echo ========================================
echo $EVEADM rkt stop --dir="$home_dir" "$CONTAINER_UUID"
$EVEADM rkt stop --dir="$home_dir" "$CONTAINER_UUID"
echo 'sleep 5'
sleep 5
echo ========================================
echo "delete container"
echo ========================================
echo $EVEADM rkt delete --dir="$home_dir" "$CONTAINER_UUID"
$EVEADM rkt delete --dir="$home_dir" "$CONTAINER_UUID"
echo ========================================
echo "delete image"
echo ========================================
echo $EVEADM rkt delete -i --dir="$home_dir" "$IMAGE_HASH"
$EVEADM rkt delete -i --dir="$home_dir" "$IMAGE_HASH"
echo ========================================
