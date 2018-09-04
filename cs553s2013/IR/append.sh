#! /bin/bash

line=`wc -l list.txt | cut -d " " -f 1`

i=1
while [ $i -le $line ]
do

	name=`head -n $i list.txt | tail -n 1`

	echo "
	
	func (n *$name) EmitSpim() (code string) {
		panic("function not implemented")
	}
	" >> $1;

	(( i++ ))

done	
