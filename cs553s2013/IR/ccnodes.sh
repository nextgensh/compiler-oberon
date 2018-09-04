#! /bin/bash

name=$1
display=$2
file=$3

echo "
type $1 struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
}

func (n *$1) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *$1) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *$1) GetChildNos() (num int){
    return len(n.children)
}

func (n *$1) NodePrint() (buf string){
    return \"$2\"
}

func (n *$1) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *$1) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}

" >> $3

echo $1 >> list.txt
