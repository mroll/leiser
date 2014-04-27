#!/bin/bash

filename="$1"
table="$2"

cat $filename |
while read line ; do 
	mysql -uleiser -p123botp@ss leiserDB << EOF
	insert into "$table" (body) values ("$line");
EOF
done
