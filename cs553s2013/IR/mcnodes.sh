#! /bin/bash

name=$1
display=$2
file=$3
instr=$4

echo "

// Node for $1

type $1 struct {
		val int
	    children [] Node    
		reg int
		sun int
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
	n.val, _ = strconv.Atoi(value)
}

func (n *$1) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *$1) GetReg() (reg int) {
	return n.reg
}

func (n *$1) SetReg(reg int) {
	n.reg = reg
}

func (n *$1) GetSun() (sun int) {
	return n.sun
}
func (n *$1) SetSun(val int) {
	n.sun = val
}
func (n *$1) EmitSpim() (code string) {

	n.reg = reg1

	reg1 := n.GetChild(0).GetReg()

	code = fmt.Sprintf(\"mul \$t%d, \$t%d \$t%d \n\", reg1, reg1, n.val)

	return code
}
" >> $3

echo $1 >> list.txt
