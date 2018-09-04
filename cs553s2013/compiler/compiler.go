// Author : Shravan Aras < shravanaras@cs.arizona.edu >
// Date : 5/2/13

package compiler

import (
	"github.com/proebsting/cs553s2013/lexer"
	"cs553s2013/mylexer"
	//"fmt"
)

// Note - If the `size` variable >= 0 then it is an array

type Sym_table struct {
	ntype string
	size int		// Makes sense only in case of arrays.
	Line int	// book keeping for reporting errors
	// The following fields are valid only if the entry is a procedure
	proc_ret string 	// Return type of the procedure
	arg_type []string	// types of all the arguments
	arg_val []int	// Indicates if this is a value param  
	formal_val int // is set to 1 if the variable is a value parameter
	visible int	// Set after this procedure is declared

	argument_names[] string /* Names of the arguments passed */	

	/* These fields are used for IR generation */
	offset int /* Used to store the starting variable offset from $fp */
	argument_param bool	/* Indicates if this is a argument parameter */
}

// These constants impose semantic meaning to the AST nodes
const (
	Operator = 1	// a + b type operator
	Procedure = 2
	Indexed = 3	// Anything of the type a[..]
	Soperator = 4	// +b | -b type operator
	Integer = 5
	Real = 6	// I never knew we don't have to code REAL :(
	Customtype = 7		// Given to an ident which has a new type
	Type = 8	// This ident has been used as a type
	Nil = 9
	Boolean = 10
	String = 11
	Explist = 12	// A special node which is used nest all the function parameters
	Statseq = 13
	Ifn = 14
	Elseifn = 15
	Elsen = 16
	Whilen = 17
	Identifier = 18		// Normal identifiers
	Forn = 20			// A node which indicates the FOR loop
	Forstep = 21	// This node contains the integer value of the for stepcounter
	Declseq = 22	// Just used to group all the procedure decl / bodies together.
	Procedureret = 23	// Special node for return of procedure
	Moduleename = 24	// Name given after the end of the module
	Module = 25
	Assoperator = 26	// Assignment operator. Makes life simple to have this seperate.
	Roperator = 27	// Indicates relational operators
	Boperator = 28	// Indicates boolean operators
)

// This structure represents nodes inside the AST 
type Node struct {
	children [] Node	
	Ntype int
	Tok lexer.Token
	void string
}

// Helper functions to manage the the tree
func (n *Node) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *Node) GetVoid() (val string){
	return n.void;
}

func (n *Node) getChild(index int) (child *Node){
	return &n.children[index]
}

func (n *Node) GetChild(index int) (child *Node){
	return n.getChild(index)
}

func (n *Node) getChildNos() (num int){
	return len(n.children)
}

func (n *Node) GetChildNos() (num int){
	return n.getChildNos()
}

// The main function which initiates the semantic / syntax checking phase
func Checker(input string) (ok bool, error string, ast *Node, sym_table map[string]map[string]Sym_table) {

	// reused code from the driver function
	ch := make(chan lexer.Token, 100)
	go mylexer.Lexer(input, ch)
	v := <- ch
	lex := lexType{ch, v}
	ret, msg, r, t := module(&lex)
	//fmt.Println(r)
	//fmt.Println(t)

	if !ret {
		//fmt.Println(msg)
		return false, msg, nil, t
	}

	ret, msg = Ast_module(&r, t)

	if !ret {
		//fmt.Println(msg)
		return false, msg, nil, t
	}

	return true, "ok", &r, t
}

func Enum2String (val int) (s string) {
	switch val {
		case String :
			return "STRING"
		case Integer :
			return "INTEGER"
	}

	return ""
}

func (n Sym_table) GetNtype() (s string){
	return n.ntype;
}

func (n Sym_table) GetSize() (s int){
	return n.size;
}

func (n Sym_table) GetVisibile() (visible int) {
	return n.visible;			
}

func (n Sym_table) SetVisibile(visible int) {
	n.visible = visible			
}

func (n Sym_table) GetFormalVal() (visible int) {
	return n.formal_val;			
}

func (n Sym_table) GetOffset() (visible int) {
	return n.offset;			
}

func (n Sym_table) GetArgVal() (arg_val []int) {
	return n.arg_val		
}

func (n Sym_table) GetArgName() (argument_names []string) {
	return n.argument_names		
}

func (n Sym_table) GetArgParam() (param bool) {
	return n.argument_param
}

func CreateSymEntry(offset1 int, ntype1 string, size1 int) (entry Sym_table) {
	return Sym_table{offset : offset1, ntype : ntype1, size : size1}
}
