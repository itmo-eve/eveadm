#!/bin/sh

echo ./eveadm help
./eveadm help
echo ========================================

m=test

echo ./eveadm help $m
./eveadm help $m
echo ========================================
echo ./eveadm $m
./eveadm $m
echo ========================================  
echo ./eveadm help $m
./eveadm help $m
echo ========================================
echo ./eveadm $m
./eveadm $m
echo ========================================
echo ./eveadm $m ps x 
./eveadm $m ps x
echo ========================================
echo ./eveadm -v $m ps x 
./eveadm -v $m ps x
echo ========================================
echo ./eveadm -v $m ls
./eveadm -v $m ls
echo ========================================
echo ./eveadm -v $m ls qwerty
./eveadm -v $m ls qwerty
echo ========================================
echo time ./eveadm -v $m sleep 100
time ./eveadm -v $m sleep 100
echo ========================================
echo time ./eveadm -v $m sleep 100 -t 1
time ./eveadm -v $m sleep 100 -t 1
echo ========================================
echo ./eveadm -v $m date
./eveadm -v $m date
echo ========================================
echo ./eveadm -v $m date --env "LANG=zh_CN.UTF-8 TZ=Asia/Shanghai"
./eveadm -v $m date --env "LANG=zh_CN.UTF-8 TZ=Asia/Shanghai"
