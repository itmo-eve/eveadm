#!/bin/sh

echo ./eveadm help
./eveadm help
echo ========================================

for m in test xen 
do
	echo ./eveadm help $m
	./eveadm help $m
	echo ========================================
	echo ./eveadm $m
	./eveadm $m
	echo ========================================  
	for a in create delete info list start stop update 
	do
		echo ./eveadm help $m $a
	        ./eveadm help $m $a
		echo ========================================
		echo ./eveadm $m $a
	        ./eveadm $m $a
		echo ========================================
		echo ./eveadm $m $a ps x 
                ./eveadm $m $a ps x
		echo ========================================
                echo ./eveadm -v $m $a ps x 
                ./eveadm -v $m $a ps x
                echo ========================================
		echo ./eveadm -v $m $a ls
                ./eveadm -v $m $a ls
		echo ========================================
		echo ./eveadm -v $m $a ls qwerty
                ./eveadm -v $m $a ls qwerty
		echo ========================================
		echo time ./eveadm -v $m $a sleep 100
                time ./eveadm -v $m $a sleep 100
		echo ========================================
		echo time ./eveadm -v $m $a sleep 100 -t 1
                time ./eveadm -v $m $a sleep 100 -t 1
		echo ========================================
		if [ "$m" = test ]
		then 
			echo ./eveadm -v $m $a date
        	        ./eveadm -v $m $a date
			echo ========================================
			echo ./eveadm -v $m $a date --env "LANG=zh_CN.UTF-8 TZ=Asia/Shanghai"
			./eveadm -v $m $a date --env "LANG=zh_CN.UTF-8 TZ=Asia/Shanghai"
		fi
	done
done
