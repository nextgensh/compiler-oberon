/* The nodes which make up the control statement all have 3 nodes, 
 * rather than the expected 2 nodes. The 3rd is the label node which
 * signals were things will jump to. These things are populated 
 * by the optimizecontrol function */

package IR

import (
	"strconv"
	"fmt"
	"cs553s2013/compiler"
)

// STORE node can be used for assignment
type Store struct {
	children [] Node	
	val string
	sun int 
	reg int
}

func (n *Store) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Store) GetChild(index int) (child Node){
	return n.children[index]
}

func (n *Store) GetChildNos() (num int){
	return len(n.children)
}

func (n *Store) NodePrint() (buf string){
	return "STORE"
}

func (n *Store) AddValue(value string){
	n.val = value
}

func (n *Store) SetSun(val int) {
	n.sun = val
}

func (n *Store) GetValue() (val string){
	return n.val
}

func (n *Store) GetSun() (sun int) {
	return n.sun
}

// ADD node can be used for +
type Add struct {
	children [] Node
	sun int 
	reg int
}

func (n *Add) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Add) GetChild(index int) (child Node){
	return n.children[index]
}

func (n *Add) GetChildNos() (num int){
	return len(n.children)
}

func (n *Add) NodePrint() (buf string){
	return "ADD"
}

func (n *Add) AddValue(value string){
	// Do notihing
}

func (n *Add) GetValue() (val string){
	return ""
}
func (n *Add) SetSun(val int) {
	n.sun = val
}
func (n *Add) GetSun() (sun int) {
	return n.sun
}

// CONST - used to handle constant integer values
type Const struct {
	val int
	sun int 
	reg int
}

func (n *Const) AddChild(child Node) {
	//n.children = append(n.children, child)
}

func (n *Const) GetChild(index int) (child Node){
	// This is the leaf
	return nil
}

func (n *Const) GetChildNos() (num int){
	return 0
}

func (n *Const) NodePrint() (buf string){
	return "CONST "+strconv.Itoa(n.val)
}

func (n *Const)  AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Const) GetValue() (val string){
	return strconv.Itoa(n.val)
}
func (n *Const) SetSun(val int) {
	n.sun = val
}
func (n *Const) GetSun() (sun int) {
	return n.sun
}

// ADDRGP node used to get an address of a variable
type Addrgp struct {
	val string // For now this just stores the name of variable
	sun int 
	reg int
}

func (n *Addrgp) AddChild(child Node) {
	// Do nothing. We cannot add a child to a leaf
}

func (n *Addrgp) GetChild(index int) (child Node){
	// This is the leaf
	return nil
}

func (n *Addrgp) GetChildNos() (num int){
	return 0
}

func (n *Addrgp) NodePrint() (buf string){
	return "ADDRGP _"+n.val
}

func (n *Addrgp) AddValue(value string){
	n.val = value
}

func (n *Addrgp) GetValue() (val string){
	return n.val
}
func (n *Addrgp) SetSun(val int) {
	n.sun = val
}
func (n *Addrgp) GetSun() (sun int) {
	return n.sun
}

// LTCONTROL node, takes care of all the less that control stuff
type Ltcontrol struct {
	children [] Node
	false_jmp int 	// Label number to which the control must jump if false	
	sun int 
	reg int
}

func (n *Ltcontrol) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Ltcontrol) GetChild(index int) (child Node){
	return n.children[index]
}

func (n *Ltcontrol) GetChildNos() (num int){
	return len(n.children)
}

func (n *Ltcontrol) NodePrint() (buf string){
	return "GTE"
}

func (n *Ltcontrol) AddValue(value string){
	n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Ltcontrol) GetValue() (val string){
	return strconv.Itoa(n.false_jmp)
}
func (n *Ltcontrol) SetSun(val int) {
	n.sun = val
}
func (n *Ltcontrol) GetSun() (sun int) {
	return n.sun
}

// JMP node, takes care of uncond jumps
type Jmp struct {
	jmp int 	// Label number to which the control must jump if false	
	sun int 
	reg int
}

func (n *Jmp) AddChild(child Node) {
	// Do nothing
}

func (n *Jmp) GetChild(index int) (child Node){
	// There shall be no children
	return child
}

func (n *Jmp) GetChildNos() (num int){
	return 0 
}

func (n *Jmp) NodePrint() (buf string){
	return "JMP "+strconv.Itoa(n.jmp)
}

func (n *Jmp) AddValue(value string){
	n.jmp, _ = strconv.Atoi(value)
}

func (n *Jmp) GetValue() (val string){
	return strconv.Itoa(n.jmp)
}
func (n *Jmp) SetSun(val int) {
	n.sun = val
}
func (n *Jmp) GetSun() (sun int) {
	return n.sun
}

// Label  node, takes care of all the label stuff
type Label struct {
	label int 	// Label number to which the control must jump if false	
	sun int 
	reg int
}

func (n *Label) AddChild(child Node) {
	// Do nothing
}

func (n *Label) GetChild(index int) (child Node){
	// There shall be no children
	return child
}

func (n *Label) GetChildNos() (num int){
	return 0 
}

func (n *Label) NodePrint() (buf string){
	return strconv.Itoa(n.label)+":"
}

func (n *Label) AddValue(value string){
	n.label, _ = strconv.Atoi(value)
}

func (n *Label) GetValue() (val string){
	return strconv.Itoa(n.label)
}
func (n *Label) SetSun(val int) {
	n.sun = val
}
func (n *Label) GetSun() (sun int) {
	return n.sun
}

// INDIR node can be used for indir
type Indir struct {
	children [] Node
	val string
	sun int 
	reg int
}

func (n *Indir) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Indir) GetChild(index int) (child Node){
	return n.children[index]
}

func (n *Indir) GetChildNos() (num int){
	return len(n.children)
}

func (n *Indir) NodePrint() (buf string){
	return "INDIR"
}

func (n *Indir) AddValue(value string){
	n.val = value
}

func (n *Indir) GetValue() (val string){
	return n.val
}
func (n *Indir) SetSun(val int) {
	n.sun = val
}
func (n *Indir) GetSun() (sun int) {
	return n.sun
}

// SUB node can be used for -
type Sub struct {
	children [] Node
	sun int 
	reg int
}

func (n *Sub) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Sub) GetChild(index int) (child Node){
	return n.children[index]
}

func (n *Sub) GetChildNos() (num int){
	return len(n.children)
}

func (n *Sub) NodePrint() (buf string){
	return "SUB"
}

func (n *Sub) AddValue(value string){
	// Do notihing
}

func (n *Sub) GetValue() (val string){
	return ""
}
func (n *Sub) SetSun(val int) {
	n.sun = val
}
func (n *Sub) GetSun() (sun int) {
	return n.sun
}


// Node for Mul

type Mul struct {
	    children [] Node    
		sun int 
	reg int
}

func (n *Mul) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Mul) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Mul) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Mul) NodePrint() (buf string){
	    return "MUL"
}

func (n *Mul) AddValue(value string){
	    // Do notihing
}

func (n *Mul) GetValue() (val string){
	    return ""
}
func (n *Mul) SetSun(val int) {
	n.sun = val
}
func (n *Mul) GetSun() (sun int) {
	return n.sun
}



// Node for Div

type Div struct {
	    children [] Node    
	sun int 
	reg int
}

func (n *Div) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Div) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Div) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Div) NodePrint() (buf string){
	    return "DIV"
}

func (n *Div) AddValue(value string){
	    // Do notihing
}

func (n *Div) GetValue() (val string){
	    return ""
}
func (n *Div) SetSun(val int) {
	n.sun = val
}
func (n *Div) GetSun() (sun int) {
	return n.sun
}



// Node for Mod

type Mod struct {
	    children [] Node    
	sun int 
	reg int
}

func (n *Mod) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Mod) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Mod) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Mod) NodePrint() (buf string){
	    return "MOD"
}

func (n *Mod) AddValue(value string){
	    // Do notihing
}

func (n *Mod) GetValue() (val string){
	    return ""
}
func (n *Mod) SetSun(val int) {
	n.sun = val
}
func (n *Mod) GetSun() (sun int) {
	return n.sun
}



// Node for BAnd

type BAnd struct {
	children [] Node    
	sun int 
	reg int
}

func (n *BAnd) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *BAnd) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *BAnd) GetChildNos() (num int){
	    return len(n.children)
}

func (n *BAnd) NodePrint() (buf string){
	    return "BAND"
}

func (n *BAnd) AddValue(value string){
	    // Do notihing
}

func (n *BAnd) GetValue() (val string){
	    return ""
}
func (n *BAnd) SetSun(val int) {
	n.sun = val
}
func (n *BAnd) GetSun() (sun int) {
	return n.sun
}



// Node for BOr

type BOr struct {
	    children [] Node    
	sun int 
	reg int
}

func (n *BOr) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *BOr) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *BOr) GetChildNos() (num int){
	    return len(n.children)
}

func (n *BOr) NodePrint() (buf string){
	    return "BOR"
}

func (n *BOr) AddValue(value string){
	    // Do notihing
}

func (n *BOr) GetValue() (val string){
	    return ""
}
func (n *BOr) SetSun(val int) {
	n.sun = val
}
func (n *BOr) GetSun() (sun int) {
	return n.sun
}



// Node for UAdd

type UAdd struct {
	    children [] Node    
	sun int 
	reg int
}

func (n *UAdd) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *UAdd) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *UAdd) GetChildNos() (num int){
	    return len(n.children)
}

func (n *UAdd) NodePrint() (buf string){
	    return "UADD"
}

func (n *UAdd) AddValue(value string){
	    // Do notihing
}

func (n *UAdd) GetValue() (val string){
	return ""
}

func (n *UAdd) SetSun(val int) {
	n.sun = val
}
func (n *UAdd) GetSun() (sun int) {
	return n.sun
}



// Node for USub

type USub struct {
	    children [] Node    
	sun int 
	reg int
}

func (n *USub) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *USub) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *USub) GetChildNos() (num int){
	    return len(n.children)
}

func (n *USub) NodePrint() (buf string){
	    return "USUB"
}

func (n *USub) AddValue(value string){
	    // Do notihing
}

func (n *USub) GetValue() (val string){
	    return ""
}
func (n *USub) SetSun(val int) {
	n.sun = val
}
func (n *USub) GetSun() (sun int) {
	return n.sun
}


type Gtcontrol struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
	sun int 
	reg int
}

func (n *Gtcontrol) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *Gtcontrol) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *Gtcontrol) GetChildNos() (num int){
    return len(n.children)
}

func (n *Gtcontrol) NodePrint() (buf string){
    return "LTE"
}

func (n *Gtcontrol) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Gtcontrol) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}
func (n *Gtcontrol) SetSun(val int) {
	n.sun = val
}
func (n *Gtcontrol) GetSun() (sun int) {
	return n.sun
}



type Ltecontrol struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
	sun int 
	reg int
}

func (n *Ltecontrol) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *Ltecontrol) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *Ltecontrol) GetChildNos() (num int){
    return len(n.children)
}

func (n *Ltecontrol) NodePrint() (buf string){
    return "GT"
}

func (n *Ltecontrol) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Ltecontrol) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}
func (n *Ltecontrol) SetSun(val int) {
	n.sun = val
}
func (n *Ltecontrol) GetSun() (sun int) {
	return n.sun
}



type Gtecontrol struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
	sun int 
	reg int
}

func (n *Gtecontrol) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *Gtecontrol) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *Gtecontrol) GetChildNos() (num int){
    return len(n.children)
}

func (n *Gtecontrol) NodePrint() (buf string){
    return "LT"
}

func (n *Gtecontrol) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Gtecontrol) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}
func (n *Gtecontrol) SetSun(val int) {
	n.sun = val
}
func (n *Gtecontrol) GetSun() (sun int) {
	return n.sun
}



type Eqcontrol struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
	sun int 
	reg int
}

func (n *Eqcontrol) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *Eqcontrol) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *Eqcontrol) GetChildNos() (num int){
    return len(n.children)
}

func (n *Eqcontrol) NodePrint() (buf string){
    return "NE"
}

func (n *Eqcontrol) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Eqcontrol) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}
func (n *Eqcontrol) SetSun(val int) {
	n.sun = val
}
func (n *Eqcontrol) GetSun() (sun int) {
	return n.sun
}



type Necontrol struct {
    children [] Node
    false_jmp int   // Label number to which the control must jump if false 
	sun int 
	reg int
}

func (n *Necontrol) AddChild(child Node) {
    n.children = append(n.children, child)
}

func (n *Necontrol) GetChild(index int) (child Node){
    return n.children[index]
}

func (n *Necontrol) GetChildNos() (num int){
    return len(n.children)
}

func (n *Necontrol) NodePrint() (buf string){
    return "EQ"
}

func (n *Necontrol) AddValue(value string){
    n.false_jmp, _ = strconv.Atoi(value)
}

func (n *Necontrol) GetValue() (val string){
    return strconv.Itoa(n.false_jmp)
}
func (n *Necontrol) SetSun(val int) {
	n.sun = val
}

func (n *Necontrol) GetSun() (sun int) {
	return n.sun
}

	
	func (n *Mul) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Div) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Mod) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *BAnd) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *BOr) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *UAdd) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *USub) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Gtcontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Ltecontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Gtecontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Eqcontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Necontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Store) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Add) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Const) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Addrgp) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Ltcontrol) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Jmp) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Label) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Indir) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Sub) SetReg(reg int) {
		n.reg = reg
	}
	

	
	func (n *Mul) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Div) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Mod) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *BAnd) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *BOr) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *UAdd) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *USub) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Gtcontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Ltecontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Gtecontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Eqcontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Necontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Store) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Add) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Const) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Addrgp) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Ltcontrol) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Jmp) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Label) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Indir) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Sub) GetReg() (reg int) {
		return n.reg
	}
	

	
	func (n *Mul) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
		n.reg = reg1

		ReleaseReg(reg2)

		return fmt.Sprintf("mul $t%d, $t%d, $t%d", reg1, reg1, reg2)
	}
	

	
	func (n *Div) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
			
		n.reg = reg1

		ReleaseReg(reg2)

		code = fmt.Sprintf("div $t%d, $t%d, $t%d", n.reg, reg1, reg2)

		return code
	}
	

	
	func (n *Mod) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
			
		n.reg = reg1

		ReleaseReg(reg2)

		code = fmt.Sprintf("rem $t%d, $t%d, $t%d",n.reg, reg1, reg2)
		

		return code
	}
	

	
	func (n *BAnd) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
			
		n.reg = reg1

		ReleaseReg(reg2)

		code = fmt.Sprintf("and $t%d, $t%d, $t%d", reg1, reg1, reg2)

		return code
	}
	

	
	func (n *BOr) EmitSpim() (code string) {

		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
			
		n.reg = reg1

		ReleaseReg(reg2)

		code = fmt.Sprintf("or $t%d, $t%d, $t%d", reg1, reg1, reg2)

		return code
	}
	

	
	func (n *UAdd) EmitSpim() (code string) {
		n.reg = n.GetChild(0).GetReg()	

		code = "nop"
		
		return code
	}
	

	
	func (n *USub) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()

		n.reg = reg1

		code = fmt.Sprintf("negu $t%d, $t%d", reg1, reg1)

		return code
	}
	

	
	func (n *Gtcontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()


		n.reg = reg1

		code = fmt.Sprintf("sgt $t%d, $t%d, $t%d\n", reg1, reg1, reg2)

		ReleaseReg(reg2)

		return code
	}
	

	
	func (n *Ltecontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		n.reg = reg1

		code = fmt.Sprintf("sle  $t%d, $t%d, $t%d", reg1, reg1, reg2)

		ReleaseReg(reg2)

		return code

	}
	

	
	func (n *Gtecontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		n.reg = reg1

		code = fmt.Sprintf("sge $t%d, $t%d, $t%d", reg1, reg1, reg2)
	
		ReleaseReg(reg2)

		return code

	}
	

	
	func (n *Eqcontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		n.reg = reg1

		code = fmt.Sprintf("subu $a0, $t%d, $t%d \n", reg1, reg2)
		code = code + fmt.Sprintf("li $a1, 1 \n")
		code = code + fmt.Sprintf("li $t%d, 0 \n", n.reg)
		code = code + fmt.Sprintf("movz $t%d, $a1, $a0 \n", reg1)

		ReleaseReg(reg2)

		return code

	}
	

	
	func (n *Necontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		code = fmt.Sprintf("sub $a0, $t%d, $t%d \n", reg1, reg2)
		code = code + fmt.Sprintf("li $a1, 1 \n")
		code = code + fmt.Sprintf("li $t%d,00 \n", n.reg)
		code = code + fmt.Sprintf("movn $t%d, $a1, $a0 \n", reg1)

		ReleaseReg(reg2)

		return code

	}
	

	
	func (n *Store) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		ReleaseReg(reg1)
		ReleaseReg(reg2)

		instruction := "sw"
		size := 4

		sym_ele, _ := Ast_SymCheckIdent(n.val)
		
		/* Decide if I should use a lb or a lw */
		temp_type, _ := Ast_SymGetType(sym_ele.GetNtype())
		temp_type, _ = Ast_SymGetType(temp_type)

		if compiler.Ast_IsArray(temp_type) {
			array_of := compiler.Ast_ArrayOf(temp_type)
			size = IR_GetByteSize(array_of)
		} else {
			size = IR_GetByteSize(temp_type)
		}

		if size == 1 {
			instruction = "sb"
		}


		return fmt.Sprintf("%s $t%d, 0($t%d)",instruction, reg2, reg1)
	}
	

	
	func (n *Add) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
		n.reg = reg1

		ReleaseReg(reg2)

		return fmt.Sprintf("add $t%d, $t%d, $t%d", reg1, reg1, reg2)
	}
	

	
	func (n *Const) EmitSpim() (code string) {
		reg := GetNextReg()
		n.reg = reg
		return fmt.Sprintf("li $t%d, %d", reg, n.val)
	}
	

	
	func (n *Addrgp) EmitSpim() (code string) {
		reg := GetNextReg()
		n.reg = reg

		sym_ele, _ := Ast_SymCheckIdent(n.val)
		
		if sym_ele.GetOffset() >= 0 {
			if sym_ele.GetFormalVal() == 1 {
				/* This will be a value parameter so copy the address of stack */
				code = fmt.Sprintf("addi $t%d, $fp, -%d", n.reg, sym_ele.GetOffset())
			} else {
				/* This will be a reference parameter so just used the value inside the stack. */ 
				code = fmt.Sprintf("lw $t%d, -%d($fp)", n.reg, sym_ele.GetOffset())
			}
		} else {
			/* This value will need to be picked from the data section */
			code = fmt.Sprintf("la $t%d, _%s\n", n.reg, n.val)
		}

		return code 
	}
	

	
	func (n *Ltcontrol) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()

		n.reg = reg1

		code = fmt.Sprintf("slt $t%d, $t%d, $t%d \n", reg1, reg1, reg2)	

		ReleaseReg(reg2)

		return code
	}
	

	
	func (n *Jmp) EmitSpim() (code string) {
		return fmt.Sprintf("b %s\n", "_L"+n.GetValue())
	}
	

	
	func (n *Label) EmitSpim() (code string) {
		return "_L"+n.GetValue()+":\n"
	}
	

	
	func (n *Indir) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		n.reg = reg1

		instruction := "lw"
		size := 4

		sym_ele, _ := Ast_SymCheckIdent(n.val)
		
		/* Decide if I should use a lb or a lw */
		temp_type, _ := Ast_SymGetType(sym_ele.GetNtype())
		temp_type, _ = Ast_SymGetType(temp_type)

		if compiler.Ast_IsArray(temp_type) {
			array_of := compiler.Ast_ArrayOf(temp_type)
			size = IR_GetByteSize(array_of)
		} else {
			size = IR_GetByteSize(temp_type)
		}

		if size == 1 {
			instruction = "lb"
		}

		return fmt.Sprintf("%s $t%d, 0($t%d)",instruction, reg1, reg1)
	}
	

	
	func (n *Sub) EmitSpim() (code string) {
		reg1 := n.GetChild(0).GetReg()
		reg2 := n.GetChild(1).GetReg()
		n.reg = reg1

		ReleaseReg(reg2)

		return fmt.Sprintf("subu $t%d, $t%d, $t%d", reg1, reg1, reg2)

		
	}
	


// Node for Pstart

type Pstart struct {
	val string
	sun int
	reg int
}

func (n *Pstart) AddChild(child Node) {
		/* Do nothing */
}

func (n *Pstart) GetChild(index int) (child Node){
		return nil
}

func (n *Pstart) GetChildNos() (num int){
		return 0
}

func (n *Pstart) NodePrint() (buf string){
	    return "_Pstart"
}

func (n *Pstart) AddValue(value string){
		n.val = value
}

func (n *Pstart) GetValue() (val string){
	    return n.val
}

func (n *Pstart) EmitSpim() (code string) {
	code = "_"+n.val+":\n"
	code = code + "sw $ra, 0($fp)\n"
	
	st_ele := global_st[n.val]
	st_stack = append(st_stack, st_ele)

	return code
}

	func (n *Pstart) GetReg() (reg int) {
		return n.reg
	}

	func (n *Pstart) SetReg(reg int) {
		n.reg = reg
	}

	func (n *Pstart) GetSun() (sun int) {
		return n.sun
	}

	func (n *Pstart) SetSun(val int) {
		n.sun = val
	}


// Node for Pend

type Pend struct {
	val string
	reg int
	sun int
}

func (n *Pend) AddChild(child Node) {
		/* Do nothing */
}

func (n *Pend) GetChild(index int) (child Node){
		return nil
}

func (n *Pend) GetChildNos() (num int){
		return 0
}

func (n *Pend) NodePrint() (buf string){
	    return "_Pend"
}

func (n *Pend) AddValue(value string){
		n.val = value
}

func (n *Pend) GetValue() (val string){
	    return n.val
}

func (n *Pend) EmitSpim() (code string) {
	code = "lw $ra, 0($fp)"
	code = code + "\n"
	code = code + "jr $ra"
	code = code + "\n"

	st_stack = st_stack[:len(st_stack)-1]

	return code
}

	func (n *Pend) GetReg() (reg int) {
		return n.reg
	}

	func (n *Pend) SetReg(reg int) {
		n.reg = reg
	}

func (n *Pend) GetSun() (sun int) {
	return n.sun
}
func (n *Pend) SetSun(val int) {
	n.sun = val
}


// Node for PStartActivation

type PStartActivation struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *PStartActivation) AddChild(child Node) {
	/* Do nothing */
}

func (n *PStartActivation) GetChild(index int) (child Node){
		return nil
}

func (n *PStartActivation) GetChildNos() (num int){
	    return 0 
}

func (n *PStartActivation) NodePrint() (buf string){
	    return "PSA"
}

func (n *PStartActivation) AddValue(value string){
		n.val, _ = strconv.Atoi(value) 
}

func (n *PStartActivation) GetValue() (val string){
	    return ""
}

func (n *PStartActivation) GetReg() (reg int) {
	return n.reg
}

func (n *PStartActivation) SetReg(reg int) {
	n.reg = reg
}

func (n *PStartActivation) GetSun() (sun int) {
	return n.sun
}
func (n *PStartActivation) SetSun(val int) {
	n.sun = val
}
func (n *PStartActivation) EmitSpim() (code string) {
	code = "sw $fp, -4($sp) \n" 	/* Save the $fp */
	code = code + "sw $sp, -8($sp) \n"	/* Save the $sp */
	code = code + "move $fp, $sp \n"	/* $fp = $sp */
	if n.val < 12 {
		/* This is wrong, adjust it back to 12 */
		n.val = 12
	}
	code = code + fmt.Sprintf("addi $sp, $sp, -%d \n", n.val)	/* Adjust the $sp */ 

	return code
}



// Node for PEndActivation

type PEndActivation struct {
	    children [] Node    
		reg int
		sun int
}

func (n *PEndActivation) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *PEndActivation) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *PEndActivation) GetChildNos() (num int){
	    return len(n.children)
}

func (n *PEndActivation) NodePrint() (buf string){
	    return "PEA"
}

func (n *PEndActivation) AddValue(value string){
	    // Do notihing
}

func (n *PEndActivation) GetValue() (val string){
	    return ""
}

func (n *PEndActivation) GetReg() (reg int) {
	return n.reg
}

func (n *PEndActivation) SetReg(reg int) {
	n.reg = reg
}

func (n *PEndActivation) GetSun() (sun int) {
	return n.sun
}
func (n *PEndActivation) SetSun(val int) {
	n.sun = val
}
func (n *PEndActivation) EmitSpim() (code string) {
	code = "lw $sp, -8($fp) \n"	/* Restore the $sp */
	code = code + "lw $fp, -4($fp) \n"	/* Restore the $fp */

	return code
}



// Node for IF

type IF struct {
	    children [] Node    
		jmp_false int
		reg int
		sun int
}

func (n *IF) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *IF) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *IF) GetChildNos() (num int){
	    return len(n.children)
}

func (n *IF) NodePrint() (buf string){
	    return "IF"
}

func (n *IF) AddValue(value string){
	n.jmp_false, _ = strconv.Atoi(value)
}

func (n *IF) GetValue() (val string){
	return strconv.Itoa(n.jmp_false)
}

func (n *IF) GetReg() (reg int) {
	return n.reg
}

func (n *IF) SetReg(reg int) {
	n.reg = reg
}

func (n *IF) GetSun() (sun int) {
	return n.sun
}
func (n *IF) SetSun(val int) {
	n.sun = val
}
func (n *IF) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	code = fmt.Sprintf("beqz $t%d, _L%s \n", reg1, n.GetValue())

	ReleaseReg(reg1)

	return code
}



// Node for PushW

type PushW struct {
	    children [] Node    
		offset int
		reg int
		sun int
}

func (n *PushW) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *PushW) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *PushW) GetChildNos() (num int){
	    return len(n.children)
}

func (n *PushW) NodePrint() (buf string){
	    return "PushW"
}

func (n *PushW) AddValue(value string){
	n.offset, _ = strconv.Atoi(value)
}

func (n *PushW) GetValue() (val string){
	    return ""
}

func (n *PushW) GetReg() (reg int) {
	return n.reg
}

func (n *PushW) SetReg(reg int) {
	n.reg = reg
}

func (n *PushW) GetSun() (sun int) {
	return n.sun
}
func (n *PushW) SetSun(val int) {
	n.sun = val
}
func (n *PushW) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	code = fmt.Sprintf("sw $t%d, -%d($sp) \n", reg1, n.offset)

	ReleaseReg(reg1)

	return code
}



// Node for Call

type Call struct {
	    children [] Node    
		proc_name string
		reg int
		sun int
}

func (n *Call) AddChild(child Node) {
	/* Do nothing */
}

func (n *Call) GetChild(index int) (child Node){
	return nil
}

func (n *Call) GetChildNos() (num int){
	    return 0
}

func (n *Call) NodePrint() (buf string){
	    return "Call "+n.proc_name
}

func (n *Call) AddValue(value string){
	n.proc_name = value
}

func (n *Call) GetValue() (val string){
	    return ""
}

func (n *Call) GetReg() (reg int) {
	return n.reg
}

func (n *Call) SetReg(reg int) {
	n.reg = reg
}

func (n *Call) GetSun() (sun int) {
	return n.sun
}
func (n *Call) SetSun(val int) {
	n.sun = val
}
func (n *Call) EmitSpim() (code string) {
	code = fmt.Sprintf("jal _%s \n", n.proc_name)

	return code
}



// Node for Not

type Not struct {
	    children [] Node    
		reg int
		sun int
}

func (n *Not) AddChild(child Node) {
		n.children = append(n.children, child)
}

func (n *Not) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Not) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Not) NodePrint() (buf string){
	    return "NOT"
}

func (n *Not) AddValue(value string){
	    // Do notihing
}

func (n *Not) GetValue() (val string){
	    return ""
}

func (n *Not) GetReg() (reg int) {
	return n.reg
}

func (n *Not) SetReg(reg int) {
	n.reg = reg
}

func (n *Not) GetSun() (sun int) {
	return n.sun
}
func (n *Not) SetSun(val int) {
	n.sun = val
}
func (n *Not) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1

	code = fmt.Sprintf("not $t%d, $t%d \n", reg1, reg1)
	code = code + fmt.Sprintf("andi $t%d, $t%d, 1 \n", reg1, reg1)

	return code
}

// Node for PushB

type PushB struct {
	    children [] Node    
		offset int
		reg int
		sun int
}

func (n *PushB) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *PushB) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *PushB) GetChildNos() (num int){
	    return len(n.children)
}

func (n *PushB) NodePrint() (buf string){
	    return "PushB"
}

func (n *PushB) AddValue(value string){
	n.offset, _ = strconv.Atoi(value)
}

func (n *PushB) GetValue() (val string){
	    return ""
}

func (n *PushB) GetReg() (reg int) {
	return n.reg
}

func (n *PushB) SetReg(reg int) {
	n.reg = reg
}

func (n *PushB) GetSun() (sun int) {
	return n.sun
}
func (n *PushB) SetSun(val int) {
	n.sun = val
}
func (n *PushB) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	code = fmt.Sprintf("sb $t%d, -%d($sp) \n", reg1, n.offset)

	ReleaseReg(reg1)

	return code
}


// Node for ModuleStart

type ModuleStart struct {
	    children [] Node    
		reg int
		sun int
}

func (n *ModuleStart) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *ModuleStart) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *ModuleStart) GetChildNos() (num int){
	    return len(n.children)
}

func (n *ModuleStart) NodePrint() (buf string){
	    return "ModuleStart"
}

func (n *ModuleStart) AddValue(value string){
	    // Do notihing
}

func (n *ModuleStart) GetValue() (val string){
	    return ""
}

func (n *ModuleStart) GetReg() (reg int) {
	return n.reg
}

func (n *ModuleStart) SetReg(reg int) {
	n.reg = reg
}

func (n *ModuleStart) GetSun() (sun int) {
	return n.sun
}
func (n *ModuleStart) SetSun(val int) {
	n.sun = val
}
func (n *ModuleStart) EmitSpim() (code string) {
	code = "main : \n"
	
	return code
}



// Node for DataSection

type DataSection struct {
	    children [] Node    
		val string
		reg int
		sun int
}

func (n *DataSection) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *DataSection) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *DataSection) GetChildNos() (num int){
	    return len(n.children)
}

func (n *DataSection) NodePrint() (buf string){
	    return "DataSection"
}

func (n *DataSection) AddValue(value string){
	n.val = value
}

func (n *DataSection) GetValue() (val string){
	return n.val
}

func (n *DataSection) GetReg() (reg int) {
	return n.reg
}

func (n *DataSection) SetReg(reg int) {
	n.reg = reg
}

func (n *DataSection) GetSun() (sun int) {
	return n.sun
}
func (n *DataSection) SetSun(val int) {
	n.sun = val
}
func (n *DataSection) EmitSpim() (code string) {
	return n.val
}



// Node for TopTemplate

type TopTemplate struct {
	    children [] Node    
		reg int
		sun int
}

func (n *TopTemplate) AddChild(child Node) {
	/* Do nothing */
}

func (n *TopTemplate) GetChild(index int) (child Node){
	return nil
}

func (n *TopTemplate) GetChildNos() (num int){
	return 0
}

func (n *TopTemplate) NodePrint() (buf string){
	    return "TopTemplate"
}

func (n *TopTemplate) AddValue(value string){
	    // Do notihing
}

func (n *TopTemplate) GetValue() (val string){
	    return ""
}

func (n *TopTemplate) GetReg() (reg int) {
	return n.reg
}

func (n *TopTemplate) SetReg(reg int) {
	n.reg = reg
}

func (n *TopTemplate) GetSun() (sun int) {
	return n.sun
}
func (n *TopTemplate) SetSun(val int) {
	n.sun = val
}
func (n *TopTemplate) EmitSpim() (code string) {
	code = ".text \n"
	code = code + ".globl main \n"

	return code
}



// Node for EndTemplate

type EndTemplate struct {
	    children [] Node    
		reg int
		sun int
}

func (n *EndTemplate) AddChild(child Node) {
	/* Do Nothing */
}

func (n *EndTemplate) GetChild(index int) (child Node){
	    return nil
}

func (n *EndTemplate) GetChildNos() (num int){
	    return 0
}

func (n *EndTemplate) NodePrint() (buf string){
	    return "EndTemplate"
}

func (n *EndTemplate) AddValue(value string){
	    // Do notihing
}

func (n *EndTemplate) GetValue() (val string){
	    return ""
}

func (n *EndTemplate) GetReg() (reg int) {
	return n.reg
}

func (n *EndTemplate) SetReg(reg int) {
	n.reg = reg
}

func (n *EndTemplate) GetSun() (sun int) {
	return n.sun
}
func (n *EndTemplate) SetSun(val int) {
	n.sun = val
}
func (n *EndTemplate) EmitSpim() (code string) {
	code = "li $v0, 10 \n"
	code = code + "syscall \n"
	code = code + ` 
___LEN:
sw $ra, 0($fp)
lw $s2, -12($fp)
lw $v1, 0($s2)
lw $ra, 0($fp)
jr $ra 
`

code = code + `
___ORD:
sw $ra, 0($fp)
lb $a1, -12($fp)
move $v1, $a1 

# Test and see if they are negative values
# If yes then add an extra 255 to it
bgez $v1, _LO1
	addi $v1, 255
_LO1:

lw $ra, 0($fp)
jr $ra
`

code = code + `
___CHR:
sw $ra, 0($fp)
lb $a1, -12($fp)
move $v1, $a1 
lw $ra, 0($fp)
jr $ra
`
code = code + `
___WRITE:
sw $ra, 0($fp)
addi $t0, $fp, -20
li $t1, 0
sw $t1, 0($t0)
lw $t0, -12($fp)
sw $t0, -12($sp) 

sw $fp, -4($sp) 
sw $sp, -8($sp) 
move $fp, $sp 
addi $sp, $sp, -16 
jal ___LEN 
lw $sp, -8($fp) 
lw $fp, -4($fp) 
addi $t0, $fp, -16
move $t1, $v1 

sw $t1, 0($t0)
_LI0:
addi $t0, $fp, -20
lw $t0, 0($t0)

addi $t1, $fp, -16
lw $t1, 0($t1)

slt $t0, $t0, $t1 


beqz $t0, _LI1 

li $t0, 0
sw $t0, -12($sp) 

sw $fp, -4($sp) 
sw $sp, -8($sp) 
move $fp, $sp 
addi $sp, $sp, -16 
jal ___CHR 
lw $sp, -8($fp) 
lw $fp, -4($fp) 
addi $t0, $fp, -20
lw $t0, 0($t0)

li $t1, 1
mul $t0, $t0, $t1

lw $t1, -12($fp)
add $t1, $t1, $t0

li $t0, 4
add $t1, $t1, $t0

lb $t1, 0($t1)

move $t0, $v1 

sub $s0, $t1, $t0 
li $s1, 0 
li $t0, 1 
movn $t1, $s1, $s0 


beqz $t0, _LI4 

addi $t0, $fp, -20
lw $t0, 0($t0)

li $t1, 1
mul $t0, $t0, $t1

lw $t1, -12($fp)
add $t1, $t1, $t0

li $t0, 4
add $t1, $t1, $t0

lb $t1, 0($t1)

addi $t0, $fp, -24
sb $t1, 0($t0)

lb $a0, 0($t0)
beqz $a0, _LI2
li $v0, 11
syscall

b _LI5
_LI4:
addi $t0, $fp, -20
addi $t1, $fp, -16
lw $t1, 0($t1)

sw $t1, 0($t0)
_LI5:
addi $t0, $fp, -20
lw $t0, 0($t0)

li $t1, 1
add $t0, $t0, $t1

addi $t1, $fp, -20
sw $t0, 0($t1)
b _LI0
b _LI2
_LI1:
_LI2:
lw $ra, 0($fp)
jr $ra

`
code = code + `

___COPY:
sw $ra, 0($fp)
addi $t0, $fp, -28
li $t1, 0
sw $t1, 0($t0)
lw $t0, -12($fp)
sw $t0, -12($sp) 

sw $fp, -4($sp) 
sw $sp, -8($sp) 
move $fp, $sp 
addi $sp, $sp, -16 
jal ___LEN 
lw $sp, -8($fp) 
lw $fp, -4($fp) 
addi $t0, $fp, -20
move $t1, $v1 

sw $t1, 0($t0)
lw $t0, -16($fp)
sw $t0, -12($sp) 

sw $fp, -4($sp) 
sw $sp, -8($sp) 
move $fp, $sp 
addi $sp, $sp, -16 
jal ___LEN 
lw $sp, -8($fp) 
lw $fp, -4($fp) 
addi $t0, $fp, -24
move $t1, $v1 

sw $t1, 0($t0)
addi $t0, $fp, -20
lw $t0, 0($t0)

addi $t1, $fp, -24
lw $t1, 0($t1)

sgt $t0, $t0, $t1

beqz $t0, _LC1 

li $v0, 10
syscall

b _LC2
_LC1:
_LC2:
_LC3:
addi $t0, $fp, -28
lw $t0, 0($t0)

addi $t1, $fp, -20
lw $t1, 0($t1)

slt $t0, $t0, $t1 


beqz $t0, _LC4 

addi $t0, $fp, -28
lw $t0, 0($t0)

li $t1, 1
mul $t0, $t0, $t1

lw $t1, -16($fp)
add $t1, $t1, $t0

li $t0, 4
add $t1, $t1, $t0

addi $t0, $fp, -28
lw $t0, 0($t0)

li $t2, 1
mul $t0, $t0, $t2

lw $t2, -12($fp)
add $t2, $t2, $t0

li $t0, 4
add $t2, $t2, $t0

lb $t2, 0($t2)

sb $t2, 0($t1)
addi $t0, $fp, -28
lw $t0, 0($t0)

li $t1, 1
add $t0, $t0, $t1

addi $t1, $fp, -28
sw $t0, 0($t1)
b _LC3
b _LC5
_LC4:
_LC5:
_LC6:
addi $t0, $fp, -28
lw $t0, 0($t0)

addi $t1, $fp, -24
lw $t1, 0($t1)

slt $t0, $t0, $t1 


beqz $t0, _LC7 

li $t0, 0
sw $t0, -12($sp) 

sw $fp, -4($sp) 
sw $sp, -8($sp) 
move $fp, $sp 
addi $sp, $sp, -16 
jal ___CHR 
lw $sp, -8($fp) 
lw $fp, -4($fp) 
addi $t0, $fp, -28
lw $t0, 0($t0)

li $t1, 1
mul $t0, $t0, $t1

lw $t1, -16($fp)
add $t1, $t1, $t0

li $t0, 4
add $t1, $t1, $t0

move $t0, $v1 

sb $t0, 0($t1)
addi $t0, $fp, -28
lw $t0, 0($t0)

li $t1, 1
add $t0, $t0, $t1

addi $t1, $fp, -28
sw $t0, 0($t1)
b _LC6
b _LC8
_LC7:
_LC8:
lw $ra, 0($fp)
jr $ra


`

	return code
}



// Node for CharConstant

type CharConstant struct {
	    children [] Node    
		val int
		reg int
		sun int
}

func (n *CharConstant) AddChild(child Node) {
	/* Do Nothing */
}

func (n *CharConstant) GetChild(index int) (child Node){
	    return nil
}

func (n *CharConstant) GetChildNos() (num int){
	    return 0
}

func (n *CharConstant) NodePrint() (buf string){
	    return "CharConstant"
}

func (n *CharConstant) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *CharConstant) GetValue() (val string){
	    return strconv.Itoa(n.val)
}

func (n *CharConstant) GetReg() (reg int) {
	return n.reg
}

func (n *CharConstant) SetReg(reg int) {
	n.reg = reg
}

func (n *CharConstant) GetSun() (sun int) {
	return n.sun
}
func (n *CharConstant) SetSun(val int) {
	n.sun = val
}
func (n *CharConstant) EmitSpim() (code string) {
	reg := GetNextReg()
	n.reg = reg

	return fmt.Sprintf("li $t%d, %d", reg, n.val)
}



// Node for StringConstant

type StringConstant struct {
	    children [] Node    
		val string
		reg int
		sun int
}

func (n *StringConstant) AddChild(child Node) {
	/* Do Nothing */
}

func (n *StringConstant) GetChild(index int) (child Node){
	    return nil
}

func (n *StringConstant) GetChildNos() (num int){
	    return 0
}

func (n *StringConstant) NodePrint() (buf string){
	    return "StringConstant"
}

func (n *StringConstant) AddValue(value string){
	n.val = value
}

func (n *StringConstant) GetValue() (val string){
	    return n.val
}

func (n *StringConstant) GetReg() (reg int) {
	return n.reg
}

func (n *StringConstant) SetReg(reg int) {
	n.reg = reg
}

func (n *StringConstant) GetSun() (sun int) {
	return n.sun
}
func (n *StringConstant) SetSun(val int) {
	n.sun = val
}
func (n *StringConstant) EmitSpim() (code string) {
	return "";
}



// Node for Return

type Return struct {
	    children [] Node   
		val int
		reg int
		sun int
}

func (n *Return) AddChild(child Node) {
	/* Do nothing */
}

func (n *Return) GetChild(index int) (child Node){
	    return nil
}

func (n *Return) GetChildNos() (num int){
	    return 0
}

func (n *Return) NodePrint() (buf string){
	    return "Return"
}

func (n *Return) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Return) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Return) GetReg() (reg int) {
	return n.reg
}

func (n *Return) SetReg(reg int) {
	n.reg = reg
}

func (n *Return) GetSun() (sun int) {
	return n.sun
}
func (n *Return) SetSun(val int) {
	n.sun = val
}
func (n *Return) EmitSpim() (code string) {
	n.reg = GetNextReg() 

	return fmt.Sprintf("move $t%d, $s%d \n", n.reg, n.val)
}



// Node for Dummy

type Dummy struct {
	    children [] Node    
		val int
		reg int
		sun int
}

func (n *Dummy) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Dummy) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Dummy) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Dummy) NodePrint() (buf string){
	    return "Dummy"
}

func (n *Dummy) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Dummy) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Dummy) GetReg() (reg int) {
	return n.reg
}

func (n *Dummy) SetReg(reg int) {
	n.reg = reg
}

func (n *Dummy) GetSun() (sun int) {
	return n.sun
}
func (n *Dummy) SetSun(val int) {
	n.sun = val
}
func (n *Dummy) EmitSpim() (code string) {
	return fmt.Sprintf("move $s%d, $v1\n", n.val)
}



// Node for Addi

type Addi struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *Addi) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Addi) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Addi) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Addi) NodePrint() (buf string){
	    return "Addi"
}

func (n *Addi) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Addi) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Addi) GetReg() (reg int) {
	return n.reg
}

func (n *Addi) SetReg(reg int) {
	n.reg = reg
}

func (n *Addi) GetSun() (sun int) {
	return n.sun
}
func (n *Addi) SetSun(val int) {
	n.sun = val
}
func (n *Addi) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1;

	return fmt.Sprintf("addi $t%d, $t%d, %d \n", n.reg, n.reg, n.val);
}




// Node for BAndi

type BAndi struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *BAndi) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *BAndi) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *BAndi) GetChildNos() (num int){
	    return len(n.children)
}

func (n *BAndi) NodePrint() (buf string){
	    return "BAndi"
}

func (n *BAndi) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *BAndi) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *BAndi) GetReg() (reg int) {
	return n.reg
}

func (n *BAndi) SetReg(reg int) {
	n.reg = reg
}

func (n *BAndi) GetSun() (sun int) {
	return n.sun
}
func (n *BAndi) SetSun(val int) {
	n.sun = val
}
func (n *BAndi) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1;

	return fmt.Sprintf("andi $t%d, $t%d, %d \n", n.reg, n.reg, n.val);
}



// Node for BOri

type BOri struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *BOri) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *BOri) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *BOri) GetChildNos() (num int){
	    return len(n.children)
}

func (n *BOri) NodePrint() (buf string){
	    return "BOri"
}

func (n *BOri) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *BOri) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *BOri) GetReg() (reg int) {
	return n.reg
}

func (n *BOri) SetReg(reg int) {
	n.reg = reg
}

func (n *BOri) GetSun() (sun int) {
	return n.sun
}
func (n *BOri) SetSun(val int) {
	n.sun = val
}
func (n *BOri) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1;

	return fmt.Sprintf("ori $t%d, $t%d, %d \n", n.reg, n.reg, n.val);
}


// Node for OptimizedLoadW

type OptimizedLoadW struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *OptimizedLoadW) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedLoadW) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedLoadW) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedLoadW) NodePrint() (buf string){
	    return "INDIR"
}

func (n *OptimizedLoadW) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedLoadW) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedLoadW) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedLoadW) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedLoadW) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedLoadW) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedLoadW) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	n.reg = reg1

	return fmt.Sprintf("lw $t%d, %d($t%d)\n", n.reg, n.val, n.reg)
}



// Node for OptimizedLoadB

type OptimizedLoadB struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *OptimizedLoadB) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedLoadB) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedLoadB) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedLoadB) NodePrint() (buf string){
	    return "INDIR"
}

func (n *OptimizedLoadB) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedLoadB) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedLoadB) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedLoadB) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedLoadB) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedLoadB) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedLoadB) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	n.reg = reg1

	return fmt.Sprintf("lb $t%d, %d($t%d)\n", n.reg, n.val, n.reg)
}



// Node for OptimizedStoreB

type OptimizedStoreB struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *OptimizedStoreB) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedStoreB) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedStoreB) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedStoreB) NodePrint() (buf string){
	    return "STORE"
}

func (n *OptimizedStoreB) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedStoreB) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedStoreB) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedStoreB) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedStoreB) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedStoreB) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedStoreB) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return fmt.Sprintf("sb $t%d, %d($t%d)", reg2, n.val, reg1)
}



// Node for OptimizedStoreW

type OptimizedStoreW struct {
	    children [] Node    
		reg int
		sun int
		val int
}

func (n *OptimizedStoreW) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedStoreW) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedStoreW) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedStoreW) NodePrint() (buf string){
	    return "STORE"
}

func (n *OptimizedStoreW) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedStoreW) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedStoreW) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedStoreW) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedStoreW) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedStoreW) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedStoreW) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return fmt.Sprintf("sw $t%d, %d($t%d)", reg2, n.val, reg1)
}



// Node for Beq

type Beq struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Beq) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Beq) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Beq) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Beq) NodePrint() (buf string){
	    return "Beq"
}

func (n *Beq) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Beq) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Beq) GetReg() (reg int) {
	return n.reg
}

func (n *Beq) SetReg(reg int) {
	n.reg = reg
}

func (n *Beq) GetSun() (sun int) {
	return n.sun
}
func (n *Beq) SetSun(val int) {
	n.sun = val
}
func (n *Beq) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("beq $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Bge

type Bge struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Bge) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Bge) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Bge) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Bge) NodePrint() (buf string){
	    return "Bge"
}

func (n *Bge) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Bge) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Bge) GetReg() (reg int) {
	return n.reg
}

func (n *Bge) SetReg(reg int) {
	n.reg = reg
}

func (n *Bge) GetSun() (sun int) {
	return n.sun
}
func (n *Bge) SetSun(val int) {
	n.sun = val
}
func (n *Bge) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("bge $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Bgt

type Bgt struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Bgt) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Bgt) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Bgt) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Bgt) NodePrint() (buf string){
	    return "Bgt"
}

func (n *Bgt) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Bgt) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Bgt) GetReg() (reg int) {
	return n.reg
}

func (n *Bgt) SetReg(reg int) {
	n.reg = reg
}

func (n *Bgt) GetSun() (sun int) {
	return n.sun
}
func (n *Bgt) SetSun(val int) {
	n.sun = val
}
func (n *Bgt) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("bgt $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Ble

type Ble struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Ble) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Ble) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Ble) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Ble) NodePrint() (buf string){
	    return "Ble"
}

func (n *Ble) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Ble) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Ble) GetReg() (reg int) {
	return n.reg
}

func (n *Ble) SetReg(reg int) {
	n.reg = reg
}

func (n *Ble) GetSun() (sun int) {
	return n.sun
}
func (n *Ble) SetSun(val int) {
	n.sun = val
}
func (n *Ble) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("ble $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Blt

type Blt struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Blt) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Blt) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Blt) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Blt) NodePrint() (buf string){
	    return "Blt"
}

func (n *Blt) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Blt) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Blt) GetReg() (reg int) {
	return n.reg
}

func (n *Blt) SetReg(reg int) {
	n.reg = reg
}

func (n *Blt) GetSun() (sun int) {
	return n.sun
}
func (n *Blt) SetSun(val int) {
	n.sun = val
}
func (n *Blt) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("blt $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Bne

type Bne struct {
		jump int
	    children [] Node    
		reg int
		sun int
}

func (n *Bne) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Bne) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Bne) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Bne) NodePrint() (buf string){
	    return "Bne"
}

func (n *Bne) AddValue(value string){
	n.jump, _ = strconv.Atoi(value)
}

func (n *Bne) GetValue() (val string){
	return strconv.Itoa(n.jump)
}

func (n *Bne) GetReg() (reg int) {
	return n.reg
}

func (n *Bne) SetReg(reg int) {
	n.reg = reg
}

func (n *Bne) GetSun() (sun int) {
	return n.sun
}
func (n *Bne) SetSun(val int) {
	n.sun = val
}
func (n *Bne) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()
	reg2 := n.GetChild(1).GetReg()

	code = fmt.Sprintf("bne $t%d, $t%d, _L%d \n", reg1, reg2, n.jump)

	ReleaseReg(reg1)
	ReleaseReg(reg2)

	return code
}



// Node for Muli

type Muli struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *Muli) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Muli) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Muli) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Muli) NodePrint() (buf string){
	    return "Muli"
}

func (n *Muli) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Muli) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Muli) GetReg() (reg int) {
	return n.reg
}

func (n *Muli) SetReg(reg int) {
	n.reg = reg
}

func (n *Muli) GetSun() (sun int) {
	return n.sun
}
func (n *Muli) SetSun(val int) {
	n.sun = val
}
func (n *Muli) EmitSpim() (code string) {

	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1

	code = fmt.Sprintf("mul $t%d, $t%d, %d \n", reg1, reg1, n.val)

	return code
}



// Node for Subi

type Subi struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *Subi) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Subi) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Subi) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Subi) NodePrint() (buf string){
	    return "Subi"
}

func (n *Subi) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Subi) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Subi) GetReg() (reg int) {
	return n.reg
}

func (n *Subi) SetReg(reg int) {
	n.reg = reg
}

func (n *Subi) GetSun() (sun int) {
	return n.sun
}
func (n *Subi) SetSun(val int) {
	n.sun = val
}
func (n *Subi) EmitSpim() (code string) {

	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1

	code = fmt.Sprintf("sub $t%d, $t%d, %d \n", reg1, reg1, n.val)

	return code
}



// Node for TrueAssign

type TrueAssign struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *TrueAssign) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *TrueAssign) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *TrueAssign) GetChildNos() (num int){
	    return len(n.children)
}

func (n *TrueAssign) NodePrint() (buf string){
	    return "TrueAssign"
}

func (n *TrueAssign) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *TrueAssign) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *TrueAssign) GetReg() (reg int) {
	return n.reg
}

func (n *TrueAssign) SetReg(reg int) {
	n.reg = reg
}

func (n *TrueAssign) GetSun() (sun int) {
	return n.sun
}
func (n *TrueAssign) SetSun(val int) {
	n.sun = val
}
func (n *TrueAssign) EmitSpim() (code string) {

	code = fmt.Sprintf("li $v0, 1\n")

	return code
}



// Node for FalseAssign

type FalseAssign struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *FalseAssign) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *FalseAssign) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *FalseAssign) GetChildNos() (num int){
	    return len(n.children)
}

func (n *FalseAssign) NodePrint() (buf string){
	    return "FalseAssign"
}

func (n *FalseAssign) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *FalseAssign) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *FalseAssign) GetReg() (reg int) {
	return n.reg
}

func (n *FalseAssign) SetReg(reg int) {
	n.reg = reg
}

func (n *FalseAssign) GetSun() (sun int) {
	return n.sun
}
func (n *FalseAssign) SetSun(val int) {
	n.sun = val
}
func (n *FalseAssign) EmitSpim() (code string) {

	code = fmt.Sprintf("li $v0, 0\n")

	return code
}



// Node for ControlMove

type ControlMove struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *ControlMove) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *ControlMove) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *ControlMove) GetChildNos() (num int){
	    return len(n.children)
}

func (n *ControlMove) NodePrint() (buf string){
	    return "ControlMove"
}

func (n *ControlMove) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *ControlMove) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *ControlMove) GetReg() (reg int) {
	return n.reg
}

func (n *ControlMove) SetReg(reg int) {
	n.reg = reg
}

func (n *ControlMove) GetSun() (sun int) {
	return n.sun
}
func (n *ControlMove) SetSun(val int) {
	n.sun = val
}
func (n *ControlMove) EmitSpim() (code string) {

	reg := GetNextReg()
	n.reg = reg
	code = fmt.Sprintf("move $t%d, $v0\n", n.reg)

	return code
}



// Node for IFn

type IFn struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *IFn) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *IFn) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *IFn) GetChildNos() (num int){
	    return len(n.children)
}

func (n *IFn) NodePrint() (buf string){
	    return "IFn"
}

func (n *IFn) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *IFn) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *IFn) GetReg() (reg int) {
	return n.reg
}

func (n *IFn) SetReg(reg int) {
	n.reg = reg
}

func (n *IFn) GetSun() (sun int) {
	return n.sun
}
func (n *IFn) SetSun(val int) {
	n.sun = val
}
func (n *IFn) EmitSpim() (code string) {

	reg1 := n.GetChild(0).GetReg()

	code = fmt.Sprintf("bnez $t%d, _L%s \n", reg1, n.GetValue())

	ReleaseReg(reg1)

	return code
}



// Node for Divi

type Divi struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *Divi) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Divi) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Divi) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Divi) NodePrint() (buf string){
	    return "Divi"
}

func (n *Divi) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Divi) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Divi) GetReg() (reg int) {
	return n.reg
}

func (n *Divi) SetReg(reg int) {
	n.reg = reg
}

func (n *Divi) GetSun() (sun int) {
	return n.sun
}
func (n *Divi) SetSun(val int) {
	n.sun = val
}
func (n *Divi) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1

	code = fmt.Sprintf("div $t%d, $t%d, %d \n", n.reg, n.reg, n.val)

	return code
}



// Node for Modi

type Modi struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *Modi) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *Modi) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *Modi) GetChildNos() (num int){
	    return len(n.children)
}

func (n *Modi) NodePrint() (buf string){
	    return "Modi"
}

func (n *Modi) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *Modi) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *Modi) GetReg() (reg int) {
	return n.reg
}

func (n *Modi) SetReg(reg int) {
	n.reg = reg
}

func (n *Modi) GetSun() (sun int) {
	return n.sun
}
func (n *Modi) SetSun(val int) {
	n.sun = val
}
func (n *Modi) EmitSpim() (code string) {
	reg1 := n.GetChild(0).GetReg()

	n.reg = reg1

	code = fmt.Sprintf("rem $t%d, $t%d, %d \n", n.reg, n.reg, n.val)

	return code
}



// Node for OptimizedStoreFB

type OptimizedStoreFB struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *OptimizedStoreFB) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedStoreFB) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedStoreFB) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedStoreFB) NodePrint() (buf string){
	    return "OptimizedStoreFB"
}

func (n *OptimizedStoreFB) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedStoreFB) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedStoreFB) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedStoreFB) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedStoreFB) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedStoreFB) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedStoreFB) EmitSpim() (code string) {

	reg1 := n.GetChild(0).GetReg()

	ReleaseReg(reg1);	

	code = fmt.Sprintf("sb $t%d, %d($fp)", reg1, n.val);

	return code
}



// Node for OptimizedStoreFW

type OptimizedStoreFW struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *OptimizedStoreFW) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedStoreFW) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedStoreFW) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedStoreFW) NodePrint() (buf string){
	    return "OptimizedStoreFW"
}

func (n *OptimizedStoreFW) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedStoreFW) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedStoreFW) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedStoreFW) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedStoreFW) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedStoreFW) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedStoreFW) EmitSpim() (code string) {

	reg1 := n.GetChild(0).GetReg()

	ReleaseReg(reg1);	

	code = fmt.Sprintf("sw $t%d, %d($fp)", reg1, n.val);

	return code

}



// Node for OptimizedLoadFB

type OptimizedLoadFB struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *OptimizedLoadFB) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedLoadFB) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedLoadFB) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedLoadFB) NodePrint() (buf string){
	    return "OptimizedLoadFB"
}

func (n *OptimizedLoadFB) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedLoadFB) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedLoadFB) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedLoadFB) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedLoadFB) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedLoadFB) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedLoadFB) EmitSpim() (code string) {

	n.reg = GetNextReg();

	code = fmt.Sprintf("lb $t%d, %d($fp)", n.reg, n.val);

	return code
}



// Node for OptimizedLoadFW

type OptimizedLoadFW struct {
		val int
	    children [] Node    
		reg int
		sun int
}

func (n *OptimizedLoadFW) AddChild(child Node) {
	    n.children = append(n.children, child)
}

func (n *OptimizedLoadFW) GetChild(index int) (child Node){
	    return n.children[index]
}

func (n *OptimizedLoadFW) GetChildNos() (num int){
	    return len(n.children)
}

func (n *OptimizedLoadFW) NodePrint() (buf string){
	    return "OptimizedLoadFW"
}

func (n *OptimizedLoadFW) AddValue(value string){
	n.val, _ = strconv.Atoi(value)
}

func (n *OptimizedLoadFW) GetValue() (val string){
	return strconv.Itoa(n.val)
}

func (n *OptimizedLoadFW) GetReg() (reg int) {
	return n.reg
}

func (n *OptimizedLoadFW) SetReg(reg int) {
	n.reg = reg
}

func (n *OptimizedLoadFW) GetSun() (sun int) {
	return n.sun
}
func (n *OptimizedLoadFW) SetSun(val int) {
	n.sun = val
}
func (n *OptimizedLoadFW) EmitSpim() (code string) {

	n.reg = GetNextReg()

	code = fmt.Sprintf("lw $t%d, %d($fp)", n.reg, n.val);

	return code

}

