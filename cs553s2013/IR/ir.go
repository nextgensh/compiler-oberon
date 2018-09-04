
// Author : Shravan Aras < shravanaras@cs.arizona.edu >
// Date : 28/2/13

package IR

import (
	//"github.com/proebsting/cs553s2013/lexer"
	//"cs553s2013/mylexer"
	"cs553s2013/compiler"
	"fmt"
)

// This queue maintains the order of trees in the forest
var seq_queue [] Node

// A generic interface which outlines the way nodes will look
type Node interface {
	AddChild(child Node)
	GetChild(index int) (child Node)
	GetChildNos() (num int)
	NodePrint() string
	AddValue(val string) 
	GetValue() (val string)
	SetSun(sun int)	// Set the setti ullman number for this node
	GetSun() (sun int)
	SetReg(reg int)
	GetReg() (reg int) 
	EmitSpim() (code string)
}

// The main function which initiates the semantic / syntax checking phase
func IRGen(input string) (ok bool, error string) {
	ok, err, r, t := compiler.Checker(input)
	if !ok {
		return ok, err
	}

	Ast_module(r, t)

	//start_print()
	code := start_eval()
	//fmt.Print(code)			

	return true, code
}

// Code which adds tree to the forest ue
func Add_tree (n Node){
	seq_queue = append(seq_queue, n)		
}

func PrintCode (n Node) (code string){

	if n.GetChildNos() == 0 {
		return n.EmitSpim()
	}

	if n.GetChildNos () == 1 {
		code = PrintCode(n.GetChild(0))
	} else {
		code = PrintCode(n.GetChild(0))
		code = code + "\n"
		code = code + PrintCode(n.GetChild(1))
	}

	code = code + "\n"
	code = code + n.EmitSpim()
	code = code + "\n"

	return code
}

// Temp code only for milestone 1 which traverses the queue to print things
func start_print() {
	for a:=0; a < len(seq_queue); a++{
		print(seq_queue[a])	
		fmt.Println("")
	}
}

// Code which traverses the queue to generate SUN for each expr.
func start_eval() (code string) {
	for a:=0; a < len(seq_queue); a++ {
		CleanReg()
		sun_allocation(seq_queue[a])	// Assign SUN
		sun_assign(seq_queue[a])	// Generate SPIM code	
		code = code + PrintCode(seq_queue[a])
	}

	return code
}

func sun_assign (n Node) (code string){

	if n.GetChildNos() == 0 {
		return n.EmitSpim()
	}

	if n.GetChildNos () == 1 {
		code  = sun_assign(n.GetChild(0))
	} else if n.GetChild(1).GetSun() > n.GetChild(0).GetSun() {
		code = sun_assign(n.GetChild(1))
		code = code + "\n"
		code = code + sun_assign(n.GetChild(0))
	} else {
		code = sun_assign(n.GetChild(0))
		code = code + "\n"
		code = code + sun_assign(n.GetChild(1))
	}

	code = code + "\n"
	code = code + n.EmitSpim()
	code = code + "\n"

	return code
}

func sun_allocation (n Node) {
	if n.GetChildNos() == 0 {
		n.SetSun(1)
		return
	}

	// Call this function on the left and right 
	// sub tree
	sun_allocation(n.GetChild(0))
	if n.GetChildNos() > 1 {
		sun_allocation(n.GetChild(1))
	}

	// For uniary operations we just need to copy the allocation number up
	if n.GetChildNos() == 1 {
		n.SetSun(n.GetChild(0).GetSun())
		return
	}

	if n.GetChild(0).GetSun() == n.GetChild(1).GetSun() {
		n.SetSun(n.GetChild(0).GetSun() + 1)
		return
	}

	if n.GetChild(0).GetSun() > n.GetChild(1).GetSun() {
		n.SetSun(n.GetChild(0).GetSun())
	} else {
		n.SetSun(n.GetChild(1).GetSun())
	}

	return
}

func print(n Node){

	fmt.Print(n.NodePrint()+" ")

	if n.GetChildNos() <= 0{
		return 
	}

	fmt.Print("[")
	print(n.GetChild(0))
	if n.GetChildNos() > 1 {
		fmt.Print(",")
		print(n.GetChild(1))
		if len(n.GetValue()) > 0 {
			fmt.Print(",")
			fmt.Print(n.GetValue())
		}
	}

	if len(n.GetValue()) > 0 {
		fmt.Print(",")
		fmt.Print(n.GetValue())
	}

	fmt.Print("]")
}
