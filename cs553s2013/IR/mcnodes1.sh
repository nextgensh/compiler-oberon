#! /bin/bash

name=$1
display=$2
file=$3

echo "

// Node for $1

type $1 struct {
	val string
}

func (n *$1) AddChild(child Node) {
		/* Do nothing */
}

func (n *$1) GetChild(index int) (child Node){
		return nil
}

func (n *$1) GetChildNos() (num int){
		return 0
}

func (n *$1) NodePrint() (buf string){
	    return \"_$2\"
}

func (n *$1) AddValue(value string){
		n.val = value
}

func (n *$1) GetValue() (val string){
	    return n.val
}

func (n *$1) EmitSpim() (code string) {
	panic(\"function not implemented\")
}
" >> $3

echo $1 >> list.txt
