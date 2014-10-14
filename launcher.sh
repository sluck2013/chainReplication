#!/bin/bash
if [ "$#" -ne 3 ]; then
    echo Usage : $0 '<configFileName> <min server id> <max server id>'
    exit
fi

pkill server.exe
for i in `seq $2 $3`; do
   ./server.exe $1 $i &
done
echo servers launched!
