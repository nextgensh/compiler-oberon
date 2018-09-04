// Author : Shravan Aras < shravanaras@cs.arizona.edu >

// Note - The whole maner in which label allocation works needs to 
//			reviewed again, it does not feel correct.
//			The label fails for a simple IF :P [Fixed]
// Note - In the for loop, the constant step is added using a 
// a Add IR instruction, this should be changed to use the Addi
// instruction. Same must be done for other things also [fixed]

package IR

import (
	//"github.com/proebsting/cs553s2013/lexer"
	//"cs553s2013/mylexer"
	"cs553s2013/compiler"
	"fmt"
	"strings"
	"strconv"
)

/* An external map to define visibility in the IR phase */
var visible_scope map[string]int

var len_reg int

// This always points to the next free register
var regs []bool
var s_reg_count int

var st_stack []map[string]compiler.Sym_table
var global_st map[string]map[string]compiler.Sym_table

var label_count int

/* Store the content which is part of data section */
var data_section string

/* Stores the next literal count */
var literal_nos int

func AddDataSection (val string) {
	data_section = data_section + val + "\n"
}

/* Function which adds this literal to the symbol table 
   & returns the generated literal name */
func AddLiteral (literal string) (genname string) {
	/* We must check if this literal contains a \n, and if it does
	 * then escape it for spim */
	find_nl := strings.Index(literal, "\n")
	main_scope := global_st["module"]
	literal_name := "_"+strconv.Itoa(literal_nos)
	literal_type_name := "_t"+strconv.Itoa(literal_nos)
		/* Add this literal to the data section */
	size := len(literal)
	size = size + 4
	size = AlignAddr(size)
	AddDataSection("_"+literal_name+": .space "+strconv.Itoa(4))
	count := 0
	if find_nl < 0 || len(literal) <= 1{
		AddDataSection("__"+literal_name+": .ascii \""+literal+"\"")
	} else {
		temp_str := strings.Split(literal, "\n")
		for count=0; count < len(temp_str); count ++ {
			temp_str[count] = strings.TrimRight(temp_str[count], "\n")
			
			literal_nos = literal_nos + 1
			temp_literal_name := "_"+strconv.Itoa(literal_nos)
			AddDataSection("__"+temp_literal_name+": .ascii \""+temp_str[count]+"\"")
			literal_nos = literal_nos + 1
			temp_literal_name = "_"+strconv.Itoa(literal_nos)
			AddDataSection("__"+temp_literal_name+": .byte 92")
			literal_nos = literal_nos + 1
			temp_literal_name = "_"+strconv.Itoa(literal_nos)
			AddDataSection("__"+temp_literal_name+": .byte 110")
		}
	}

	/* Insert the type first */
	main_scope[literal_type_name] = compiler.CreateSymEntry(-1, "CHAR", len(literal)+count-1)
	/* Then connect the newly defined literal to this type */
	main_scope[literal_name] = compiler.CreateSymEntry(-1, literal_type_name, 0)
	/* Add this back to the original scope */
	global_st["module"] = main_scope


	AddDataSection(".align 2")
	/* Increment the literal counter */
	literal_nos = literal_nos + 1

	return literal_name
}

// These functions do not have bounds check on next_reg
// true means it is free and false means the register is currently in use
func GetNextReg() (a int) {
	for a=0; a < len(regs); a++ {
		if regs[a] {
			regs[a] = false
			return a
		}
	}

	// If we are here it means we could not find a single register :( bad !
	return 0
}

func ReleaseReg(a int) {
	regs[a] = true
}

func CleanReg() {
	for a:=0; a < len(regs); a++ {
		regs[a] = true
	}
}

/* Function which aligns the address */
func AlignAddr (addr int) (alignedaddr int) {
	for ; (addr % 4) != 0; {
		addr = addr + 1
	}

	return addr
}

/* Function to init the array lengths */
func InitArrayLengths (st_ele map[string]compiler.Sym_table) {
	for k, e := range st_ele {
		if e.GetSize() >= 0 {
			continue
		}
		
		temptype, _ := Ast_SymGetType(e.GetNtype())
		temptype, _ = Ast_SymGetType(temptype)

		/*  If this is a array and a local variable then we init its length here. */
		if compiler.Ast_IsArray(temptype) && !e.GetArgParam() {
			sym_type, _ := Ast_SymCheckIdent(e.GetNtype())	
			
			const_size := &Const{}
			const_size.AddValue(strconv.Itoa(sym_type.GetSize()))
			
			addrgp := &Addrgp{val : k}

			store := &Store{}
			store.AddChild(addrgp)
			store.AddChild(const_size)
			Add_tree(store)
		}
		
	}

}

func Decrement_s () {
	if s_reg_count > 2 {
		s_reg_count --
	}
}

func GetNextS() (val int){
	temp := s_reg_count
	s_reg_count = ((s_reg_count +1 ) % 8)
	
	return temp
}

func Ast_module (n *compiler.Node, st map[string]map[string]compiler.Sym_table) (ok bool, err string){
	// Load the module and predefined sym table onto the stack

	/* Initilialize the literal number to 0 */
	literal_nos = 0
	
	visible_scope = make(map[string]int)

	s_reg_count = 2

	toptemplate := &TopTemplate{}
	Add_tree(toptemplate)

	label_count = 0
	regs = make([]bool, 10)	// For now i only consider the 10 t registers
	CleanReg()

	
    st_ele := st["module"]
	/* We must set the visible bit of all procedure back to 0 except for 
     * the inbuilt procedures as the visibility order will be redefined again 
     */
	for k, _ := range st_ele {
		count := strings.Count(k, "_")
		if count == 1 {
			visible_scope[k] = 0
		}
	}

	/* This ordering is important */
	st_stack = st_stack[0:0]
    st_ele = st["module"]
    st_stack = append(st_stack, st_ele)


	AddDataSection(".data")
	for k, e := range st_ele {
		
		if e.GetSize() >= 0 {
			continue
		}

		temptype, _ := Ast_SymGetType(e.GetNtype())
		temptype, _ = Ast_SymGetType(temptype)


		if compiler.Ast_IsArray(temptype) {
			/* This is an array */
			size := IR_GetByteSize(compiler.Ast_ArrayOf(temptype))
			size = size * st_ele[e.GetNtype()].GetSize()
			size = AlignAddr(size)
			size = size + 4	/* Extra space to store the array length */
			AddDataSection("_"+k+": .space "+strconv.Itoa(size))
		} else {
			size := IR_GetByteSize(temptype)
			size = AlignAddr(size)	
			AddDataSection("_"+k+": .space "+strconv.Itoa(size))
		}
	}

	global_st = st
	
	pdecseq := n.GetChild(0);
		    
    if pdecseq.Ntype == compiler.Declseq {
        for a:=0; a < pdecseq.GetChildNos(); a++ {
            decseq := pdecseq.GetChild(a)
            ok, err := Ast_Procedure(decseq, st) 
            if !ok {
                return false, err 
            }   
        }   
    }

	module := &ModuleStart{}
	Add_tree(module)

	/* We must init the length field of all our predefined arrays */
	InitArrayLengths(st_ele)
	
	if n.GetChildNos() >=2 {
		decseq := n.GetChild(1)
		if decseq.Ntype == compiler.Statseq {
			ok, err = Ast_statementsequence(decseq)
		}
	}

	endtemplate := &EndTemplate{}
	Add_tree(endtemplate)

	ds := &DataSection{}
	ds.AddValue(data_section)
	Add_tree(ds)


	return ok, err
}

func Ast_Procedure (n *compiler.Node, st map[string]map[string]compiler.Sym_table) (ok bool, err string) {
	// Load the scope into the current scope stack
	name := n.GetVoid()

	st_ele := st[name]
	st_stack = append(st_stack, st_ele)

	visible_scope[name] = 1

	pstart := &Pstart{}
	pstart.AddValue(name)
	Add_tree(pstart)

	/* Initialize the array lengths inside the procedure */
	InitArrayLengths(st_ele)

	if n.GetChildNos() >= 1 {
		if n.GetChild(0).Ntype != compiler.Statseq {
		} else {
			ok, err := Ast_statementsequence(n.GetChild(0))
			if !ok {
				return false, err
			}
		}
	}

	// Make sure there was no return statement
	st_stack = st_stack[:len(st_stack)-1]

	pend := &Pend{}
	Add_tree(pend)

	return true, "ok"
}	

func Ast_statementsequence (n *compiler.Node) (ok bool, err string){
	// We iterate over all the children 
	
	for a:=0; a < n.GetChildNos(); a++ {
		c := n.GetChild(a)
		ntype := c.Ntype

		switch ntype {
			case compiler.Assoperator :
				ok, err = Ast_assignorcall(c)
			case compiler.Ifn :
				label_count ++
				ok, err = Ast_ifstatement(c, -1, -1, 0)
			case compiler.Whilen :
				ok, err = Ast_whilestatement(c)
			case compiler.Forn :
				ok, err = Ast_forstatement(c)
			case compiler.Procedure :
				 ok, err, _ = Ast_expression(c, false)
		}
		if !ok {
			return ok, err
		}
	}

	return true, "ok"
}

func Ast_forstatement (n *compiler.Node) (ok bool, err string) {
	
	cexpr := n.GetChild(0)

	/* assignment node */
	Ast_assignorcall(cexpr)

	start_label := GetLabel()
	label1 := &Label{}
	label1.AddValue(strconv.Itoa(start_label))
	Add_tree(label1)

	toexpr := n.GetChild(1)

	_, _, irnode2 := Ast_expression(toexpr, false)

	end_label := GetLabel()

	
	/* This is the control variable. */
	addrgp_i := &Addrgp{val : n.GetChild(0).GetChild(0).Tok.Value()}
	indir_i := &Indir{}
	indir_i.AddChild(addrgp_i)

	/* The actual less than or equal part goes here. */
	lte := &Ltecontrol{}
	lte.AddChild(indir_i)
	lte.AddChild(irnode2)
	
	/*
	ifn := &IF{}
	ifn.AddValue(strconv.Itoa(end_label))
	ifn.AddChild(lte)
	Add_tree(ifn)
	*/
	OptimizeControl(lte, end_label, false)

	statorby := n.GetChild(2)

	conststep := &Const{}

	if statorby.Ntype == compiler.Forstep {
		/* This is just the integer literal so we add this to the end */
		
		trim_str := strings.TrimRight(statorby.Tok.Value(), "H")
		if trim_str == statorby.Tok.Value() {
			conststep.AddValue(statorby.Tok.Value())
		} else {
			dec_val, _ := strconv.ParseInt(trim_str, 16, 64)
			conststep.AddValue(strconv.Itoa(int(dec_val)))
		}
		statorby = n.GetChild(3)
	} else {
		/* The default is to increment it just by 1. */
		conststep.AddValue(strconv.Itoa(1))
	}

	ok, err = Ast_statementsequence(statorby)
	if !ok {
		return false, err	
	}
	

	/* A new instance */
	addrgp_i = &Addrgp{val : n.GetChild(0).GetChild(0).Tok.Value()}
	indir_i = &Indir{}
	indir_i.AddChild(addrgp_i)

	/* Optimized this to use a addi */
	addstep := &Addi{}
	addstep.AddChild(indir_i)
	addstep.AddValue(conststep.GetValue())

	addrgp_i = &Addrgp{val : n.GetChild(0).GetChild(0).Tok.Value()}
	store := &Store{}
	store.AddChild(addrgp_i)
	store.AddChild(addstep)
	Add_tree(store)


	jmp := &Jmp{}
	jmp.AddValue(strconv.Itoa(start_label))
	Add_tree(jmp)


	
	label2 := &Label{}
	label2.AddValue(strconv.Itoa(end_label))
	Add_tree(label2)

	return true, "ok"
}

func Ast_assignorcall(n *compiler.Node) (ret bool, err string){

	if n.Ntype == compiler.Assoperator {

		lc := n.GetChild(0)
		lr := n.GetChild(1)

		ret, err, irnode1 := Ast_expression(lc, true)
		if !ret {
			return false, err 
		}
	
		ret, err, irnode2 := Ast_expression(lr, false)
		if !ret {
			return false, err	
		}

		else_label := GetLabel()
		end_label := GetLabel()

		logical := CanBoolOptimize(irnode2.NodePrint())
		if !logical {
			logical = false
		} else {
			logical = OptimizeControl(irnode2, else_label, false)
		}
		

		// Make a STORE node
		store := &Store{}
		store.AddChild(irnode1)
		if !logical {
			store.AddChild(irnode2)
		} else {
			/* We are dealing with some logical expression. */
			truecode := &TrueAssign{}
			Add_tree(truecode)
			jmp := &Jmp{}
			jmp.AddValue(strconv.Itoa(end_label))
			Add_tree(jmp)
			l1 := &Label{}
			l1.AddValue(strconv.Itoa(else_label))
			Add_tree(l1)
			falsecode := &FalseAssign{}
			Add_tree(falsecode)
			l2 := &Label{}
			l2.AddValue(strconv.Itoa(end_label))
			Add_tree(l2)
			cm := &ControlMove{}
			store.AddChild(cm)
		}
		/* Iterate till you reach addrgp and then pull the variable name from it */
		temp := irnode1
		for ; temp.NodePrint() == "ADD" || temp.NodePrint() == "Addi" ; {
			temp = temp.GetChild(0)	
		}

		/* Are we dealing with an array here ? */
		sym_ele, _ := Ast_SymCheckIdent(temp.GetValue())
		temp_type, _ := Ast_SymGetType(sym_ele.GetNtype())
		temp_type, _ = Ast_SymGetType(temp_type)
		if compiler.Ast_IsArray(temp_type) { 
			/* Check to see if this is an array, if it is can we use a 
			 * optimized store instruction for it. */
			ret, index, addrgp := OptimizeArray(irnode1)
			if ret {
				/* Wow this can be optimized, lets bring out the complex addressing mode store. */
				/* Lets dig a little deeper, can we actually get rid of the addrgp node */

				constant := ReturnConstant(addrgp); 
				if constant >= 0 {
					/* Wola we have an opportunity to get rid of the addrgp
					 * and reference directly from $fp */
					 index = index + (-constant);
				}

				array_of := compiler.Ast_ArrayOf(temp_type)
				size := IR_GetByteSize(array_of)
				if size == 1 {
					if constant >= 0 {
						/* 1 bytes array */
						optimizedstore := &OptimizedStoreFB{}
						optimizedstore.AddValue(strconv.Itoa(index))

						if !logical {
							optimizedstore.AddChild(irnode2)
						} else {
							cm := &ControlMove{}
							optimizedstore.AddChild(cm)
						}

						Add_tree(optimizedstore)
						return true, "ok"
					} else {
						/* 1 bytes array */
						optimizedstore := &OptimizedStoreB{}
						optimizedstore.AddValue(strconv.Itoa(index))
						optimizedstore.AddChild(addrgp)

						if !logical {
							optimizedstore.AddChild(irnode2)
						} else {
							cm := &ControlMove{}
							optimizedstore.AddChild(cm)
						}

						Add_tree(optimizedstore)
						return true, "ok"
					}

				} else {
					if constant >= 0 {
						/* Word aligned array */
						optimizedstore := &OptimizedStoreFW{}
						optimizedstore.AddValue(strconv.Itoa(index))

						if !logical {
							optimizedstore.AddChild(irnode2)
						} else {
							cm := &ControlMove{}
							optimizedstore.AddChild(cm)
						}

						Add_tree(optimizedstore)
						return true, "ok"
					} else {
						/* Word aligned array */
						optimizedstore := &OptimizedStoreW{}
						optimizedstore.AddValue(strconv.Itoa(index))
						optimizedstore.AddChild(addrgp)
						if !logical {
							optimizedstore.AddChild(irnode2)
						} else {
							cm := &ControlMove{}
							optimizedstore.AddChild(cm)
						}

						Add_tree(optimizedstore)
						return true, "ok"
					}
				}
			} else {
				store.AddValue(temp.GetValue())
				Add_tree(store)
			}
		} else {
			store.AddValue(temp.GetValue())
			Add_tree(store)
		}
	}

	return true, "ok"
}

func Ast_expression (n *compiler.Node, left bool) (ok bool, err string,
													ir Node){
	
	switch n.Ntype {

		case compiler.Procedure :

			expcount := 0

			proc_name := n.GetChild(0)
			sys_module := st_stack[0]
			proc_header, found := sys_module["_"+proc_name.Tok.Value()]
			my_proc_name := "_"+proc_name.Tok.Value()

			visibility := visible_scope[my_proc_name]
            /* If this is not found lets try again after applying the prebuilt conversion */
            if !found || visibility != 1 {
				my_proc_name = compiler.ConvertPrebuild("_"+proc_name.Tok.Value())
                proc_header, _ = sys_module[my_proc_name]
            } 

			
			if n.GetChildNos() <= 1 {
				expcount = 0
			 } else {
				expcount = n.GetChild(1).GetChildNos()
			 }

			/* Assert - We have the correct number of arguments */
			var explist *compiler.Node
			if (expcount > 0) {
				explist = n.GetChild(1)
			}
			for i:=0; i < expcount; i++ {		
				isvaltype := proc_header.GetArgVal()[i]
				argname := proc_header.GetArgName()[i]
				/* Get the symtable entry for this */
				/* temp load into the procedure scope */
				st_ele := global_st[my_proc_name]
				st_stack = append(st_stack, st_ele)
					argsym, _ := Ast_SymCheckIdent(argname)
					argtype, _ := Ast_SymGetType(argsym.GetNtype())
					argtype, _ = Ast_SymGetType(argtype)
				/* Return to original scope */
				st_stack = st_stack[:len(st_stack)-1]
				valtype := true
				if isvaltype != 1 || compiler.Ast_IsArray(argtype) {
					valtype = false
				}
				ok, err, ir = Ast_expression(explist.GetChild(i), !valtype)

				/* Check if this node can be optimized. */
				else_label := GetLabel()
				end_label := GetLabel()

				logical := CanBoolOptimize(ir.NodePrint())
				if !logical {
					logical = false
				} else {
					logical = OptimizeControl(ir, else_label, false)
				}

				if logical {
					/* We are dealing with some logical expression which can be 
					 * optimized. */
					truecode := &TrueAssign{}
					Add_tree(truecode)
					jmp := &Jmp{}
					jmp.AddValue(strconv.Itoa(end_label))
					Add_tree(jmp)
					l1 := &Label{}
					l1.AddValue(strconv.Itoa(else_label))
					Add_tree(l1)
					falsecode := &FalseAssign{}
					Add_tree(falsecode)
					l2 := &Label{}
					l2.AddValue(strconv.Itoa(end_label))
					Add_tree(l2)
					//store.AddChild(cm)
				}


				/* I might want to check the size here and change how I push this accordingly */
				if compiler.Ast_IsArray(argtype) {
					/* Check to see if a string literal was returned by Ast_expression */
					temp_len := 0
					literal_name := ""
					if ir.NodePrint() == "CharConstant" || ir.NodePrint() == "StringConstant" {
						if ir.NodePrint() == "StringConstant" {
							temp := strings.Split(ir.GetValue(),"\n")
							temp_len = len(ir.GetValue())+len(temp)-1
							literal_name = AddLiteral(ir.GetValue())
							literal_addrgp := &Addrgp{val : literal_name}
							ir = literal_addrgp
						} else {
							/* otherwise this is char const and we must convert it to char */
							intval, _ := strconv.Atoi(ir.GetValue());
							irval := fmt.Sprintf("%c", intval)
							temp_len = 1
							literal_name = AddLiteral(irval)
							literal_addrgp := &Addrgp{val : literal_name}
							ir = literal_addrgp
						}

						/* Also make a IR tree to initialize the length field of this literal */
						const_size := &Const{}
						const_size.AddValue(strconv.Itoa(temp_len))
						
						addrgp := &Addrgp{val : literal_name}

						store := &Store{}
						store.AddChild(addrgp)
						store.AddChild(const_size)
						Add_tree(store)
					}
				}
				if IR_GetByteSize(argtype) == 1 && valtype {
					pushb := &PushB{}
					pushb.AddValue(strconv.Itoa(argsym.GetOffset()))
					if !logical {
						pushb.AddChild(ir)
					} else {
						cm := &ControlMove{}
						pushb.AddChild(cm)
					}
					Add_tree(pushb)
				} else {
					pushw := &PushW{}
					pushw.AddValue(strconv.Itoa(argsym.GetOffset()))
					if !logical {
						pushw.AddChild(ir)
					} else {
						cm := &ControlMove{}
						pushw.AddChild(cm)
					}
					Add_tree(pushw)
				}
				if !ok {
					return false, err, ir
				}
				//Add_tree(ir)		
			}		

			/* Call the activation record to do the stack adjustment */
			pstart := &PStartActivation{}
			pstart.AddValue(strconv.Itoa(proc_header.GetOffset()))
			Add_tree(pstart)


			/* The actual jal spim instruction goes here. */
			call := &Call{}
			call.AddValue(my_proc_name)
			Add_tree(call)

			pend := &PEndActivation{}
			Add_tree(pend)

			dummy := &Dummy{}
			dummy.AddValue(strconv.Itoa(GetNextS()))
			Add_tree(dummy)

			ir := &Return{}
			ir.AddValue(dummy.GetValue())

			return true, "ok", ir

		case compiler.Operator : 
			ok, err, ir = Ast_BinaryOp(n)
			if !ok {
				return false, err, ir
			}
			return true, "ok", ir

		case compiler.Integer :
			// Return a CONST IR node
			trim_str := strings.TrimRight(n.Tok.Value(), "H")
			constant := &Const{}
			if trim_str == n.Tok.Value() {
				constant.AddValue(n.Tok.Value())
			} else {
				dec_val, _ := strconv.ParseInt(trim_str, 16, 64)
				constant.AddValue(strconv.Itoa(int(dec_val)))
			}
			return true, "ok", constant
			/* If we see a hex number we must take care of it here. */ 

		case compiler.Boolean :
			constant := &Const{}
			/* If this is true, we make this equal to 1 else it is 0*/
			if n.Tok.Value () == "TRUE" {
				constant.AddValue(strconv.Itoa(1))				
			} else {
				constant.AddValue(strconv.Itoa(0))
			}
			return true, "ok", constant

		case compiler.Identifier :
			addrgp := &Addrgp{val : n.Tok.Value()}
			if !left {
				indir := &Indir{}
				indir.AddValue(n.Tok.Value())
				indir.AddChild(addrgp)
				return true, "ok", indir
			}
			return true, "ok", addrgp

		case compiler.Boperator :
			if n.Tok.Value() == "&" {
				
				l1 := n.GetChild(0)
				l2 := n.GetChild(1)

				ret, err, irnode1 := Ast_expression(l1, false)
				if !ret {
					return false, err, ir
				}
				
				ret, err, irnode2 := Ast_expression(l2, false)
				if !ret {
					return false, err, ir
				}
				
				/*
				ret, cnode, onode := ContainsConst(irnode1, irnode2)
				if(ret) {		
					bandi := &BAndi{}
					bandi.AddValue(cnode.GetValue())
					bandi.AddChild(onode)

					return true, "ok", bandi
				} else {
				*/
					band := &BAnd{}
					band.AddChild(irnode1)
					band.AddChild(irnode2)

					return true, "ok", band
				//}
			} else if n.Tok.Value() == "OR" {

				l1 := n.GetChild(0)
				l2 := n.GetChild(1)


				ret, err, irnode1 := Ast_expression(l1, false)
				if !ret {
					return false, err, ir
				}

				ret, err, irnode2 := Ast_expression(l2, false)
				if !ret {
					return false, err, ir
				}

				//ret, cnode, onode := ContainsConst(irnode1, irnode2)
				/*if(ret) {		
					bori := &BOri{}
					bori.AddValue(cnode.GetValue())
					bori.AddChild(onode)

					return true, "ok", bori
				} else {*/
					bor := &BOr{}
					bor.AddChild(irnode1)
					bor.AddChild(irnode2)

					return true, "ok", bor
				//}
			} 

		case compiler.Roperator :
			ret, err, irnode1 := Ast_expression(n.GetChild(0), false)
			if !ret {
				return false, err, ir
			}
			
			ret, err, irnode2 := Ast_expression(n.GetChild(1), false)
			if !ret {
				return false, err, ir
			}

			if n.Tok.Value() == "<"{
				control := &Ltcontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			} else if n.Tok.Value() == ">"{
				control := &Gtcontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			}	else if n.Tok.Value() == "<="{
				control := &Ltecontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			}	else if n.Tok.Value() == ">="{
				control := &Gtecontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			} else if n.Tok.Value() == "="{
				control := &Eqcontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			} else if n.Tok.Value() == "#"{
				control := &Necontrol{}
				control.AddChild(irnode1)
				control.AddChild(irnode2)

				return true, "ok", control
			} 


		case compiler.Soperator :

			ret, err, irnode1 := Ast_expression(n.GetChild(0), false)

			if !ret {
				return false, err, ir
			}

			if n.Tok.Value() == "-" {
				opnode := &USub{}
				opnode.AddChild(irnode1)

				return true, "ok", opnode
			} else if n.Tok.Value() == "+" {
				opnode := &UAdd{}
				opnode.AddChild(irnode1)

				return true, "ok", opnode
			} else if n.Tok.Value() == "~" {
				opnode := &Not{}
				opnode.AddChild(irnode1)

				return true, "ok", opnode
			}

		case compiler.Indexed :
			lc := n.GetChild(0)
            lr := n.GetChild(1)

            sym_ele, ret := Ast_SymCheckIdent(lc.Tok.Value())

            arr_type, _ := Ast_SymGetType(sym_ele.GetNtype())
            ntype := Ast_ArrayOf(arr_type)

			byte_size := IR_GetByteSize(ntype)

            ret, err, irnode1 := Ast_expression(lr, false)

			/* Do the multiplication with the constant size */
			const_fact := &Const{}
			const_fact.AddValue(strconv.Itoa(byte_size))

			mul := &Muli{}
			mul.AddChild(irnode1)
			mul.AddValue(strconv.Itoa(byte_size))

			/* Get the base address */
			addrgp := &Addrgp{val : n.Tok.Value()}
			/* Add the base address. */
			addadd := &Add{}
			addadd.AddChild(addrgp)
			addadd.AddChild(mul)

			/* Always offset array by constant of 4, because length comes first */
			const4 := &Const{}
			const4.AddValue(strconv.Itoa(4))

			finaladd := &Addi{}
			finaladd.AddChild(addadd)
			finaladd.AddValue(strconv.Itoa(4))

			if !left {
				/* Check to see if this is an array form we can
				 * optimize */
				ret, index, _ := OptimizeArray(finaladd)
				if ret {
					/*This array has a constant index and can be optimized. */
					constant := ReturnConstant(addrgp); 
					if constant >= 0 {
						/* Wola we have an opportunity to get rid of the addrgp
						 * and reference directly from $fp */
						 index = index + (-constant);
					}

					if byte_size == 1 {
						if constant >=0 {
							optimizedload := &OptimizedLoadFB{}
							optimizedload.AddValue(strconv.Itoa(index))
							return true, "ok", optimizedload
						} else {
							optimizedload := &OptimizedLoadB{}
							optimizedload.AddValue(strconv.Itoa(index))
							optimizedload.AddChild(addrgp)	
							return true, "ok", optimizedload
						}
					} else {
						if constant >=0 {
							optimizedload := &OptimizedLoadFW{}
							optimizedload.AddValue(strconv.Itoa(index))
							return true, "ok", optimizedload
						} else {
							optimizedload := &OptimizedLoadW{}
							optimizedload.AddValue(strconv.Itoa(index))
							optimizedload.AddChild(addrgp)	
							return true, "ok", optimizedload
						}
					}
				} else {
					indir := &Indir{}
					indir.AddValue(n.Tok.Value())
					indir.AddChild(finaladd)
					return true, "ok", indir
				}
			}

			return true, "ok", finaladd

            if !ret {
                return false, err, ir
            }

            // Assert - the ident is an array and the index is a valid integer expression 

		case compiler.String :
			/* Handle the case of empty string */
			new_string, stat := strconv.Unquote(n.Tok.Value())
			var dec_str int64
			dec_str = 0
			if stat != nil {
				/* Now it is 0AX */
				trim_str := strings.TrimRight(n.Tok.Value(), "X")
				dec_str, _ = strconv.ParseInt(trim_str, 16, 64)
				new_string = fmt.Sprintf("%c", dec_str)
			} else {	
				if len(new_string) == 1 {
					dec_str = int64(new_string[0])
				}
			}

			if len(new_string) == 1 || dec_str > 127{
				const_char := &CharConstant{}
				const_char.AddValue(strconv.Itoa(int(dec_str)))

				return true, "ok", const_char
			} else {
				
				const_string := &StringConstant{}
				const_string.AddValue(new_string)

				return true, "ok", const_string
			}
	}

	return false, "fail", ir
}

func Ast_BinaryOp (n *compiler.Node) (ok bool, err string, ir Node) {

	lc := n.GetChild(0)
	lr := n.GetChild(1)

	ok, err, irnode1 := Ast_expression(lc, false)
	if !ok {
		return false, err, ir
	}
	ok, err, irnode2 := Ast_expression(lr, false)
	if !ok {
		return false, err, ir
	}

	// Depending on what binary op this is, generate the IR node for this
	if n.Tok.Value() == "+" {
		/* Check if any of this is a constant node and do things. */
		ret, cnode, onode := ContainsConst(irnode1, irnode2)
		if(ret){
			addi := &Addi{}
			addi.AddValue(cnode.GetValue())
			addi.AddChild(onode)

			return true, "ok", addi
		} else {
			add := &Add{}
			add.AddChild(irnode1)
			add.AddChild(irnode2)
			
			return true, "ok", add
		}
	} else if n.Tok.Value() == "-" {

		ret, cnode, onode := ContainsConst(irnode1, irnode2)
		if(ret){
			subi := &Subi{}
			subi.AddValue(cnode.GetValue())
			subi.AddChild(onode)

			/* If the left hand side was a constant then we need to 
			 * add a negation to this. */
			 if IsConst(irnode1.NodePrint()) {
				neg := &USub{}
				neg.AddChild(subi)

				return true, "ok", neg
			 }

			return true, "ok", subi
		} else {
			sub := &Sub{}
			sub.AddChild(irnode1)
			sub.AddChild(irnode2)
			
			return true, "ok", sub
		}

	} else if n.Tok.Value() == "*" {

		ret, cnode, onode := ContainsConst(irnode1, irnode2)
		if(ret){
			muli := &Muli{}
			muli.AddValue(cnode.GetValue())
			muli.AddChild(onode)

			return true, "ok", muli
		} else {
			mul := &Mul{}
			mul.AddChild(irnode1)
			mul.AddChild(irnode2)
			
			return true, "ok", mul
		}

	} else if n.Tok.Value() == "DIV" {
		ret, _, _ := ContainsConst(irnode1, irnode2)
	
		if ret && IsConst(irnode2.NodePrint()) && irnode2.GetValue() != "0" {
			divi := &Divi{}
			divi.AddValue(irnode2.GetValue())
			divi.AddChild(irnode1)

			return true, "ok", divi
		} else {
	
			sub := &Div{}
			sub.AddChild(irnode1)
			sub.AddChild(irnode2)

			return true, "ok", sub
		}
	} else if n.Tok.Value() == "MOD" {
		ret, _, _ := ContainsConst(irnode1, irnode2)

		if ret && IsConst(irnode2.NodePrint()) && irnode2.GetValue() != "0" {
			modi := &Modi{}
			modi.AddValue(irnode2.GetValue())
			modi.AddChild(irnode1)

			return true, "ok", modi
		} else {
			sub := &Mod{}
			sub.AddChild(irnode1)
			sub.AddChild(irnode2)

			return true, "ok", sub
		}
	}

	return false, "ok", ir
}

func Ast_whilestatement (n *compiler.Node) (v bool, err string) {

	start_label := GetLabel()
	label1 := &Label{label : start_label}
	Add_tree(label1)

	v, err = Ast_ifstatement(n, start_label, -1, 0)

	return v, err
}

	//This method is overloaded to the max

	//Arg0 : the if node
	//Arg1 : Is this being called from the while loop
	//Arg2 : The end label to which everyone jumps to in case of if-elseif struct
	//Arg3 : The recursion depth the call currently is

func Ast_ifstatement(n *compiler.Node, while_start_label, end_label int, depth int) (v bool, err string){
	
	expc := n.GetChild(0)
	statc := n.GetChild(1)

	else_label := GetLabel()

	if end_label < 0 {	
		end_label = GetLabel() 
	}

/*

	if !iswhile && depth <= 0{
		label_count ++
	}
*/
	
	rete, err, irnode1 := Ast_expression(expc, false)

	if !rete {
		return false, err
	}

/*
	ifnode := &IF{}
	ifnode.AddChild(irnode1)
	ifnode.AddValue(strconv.Itoa(else_label))
	Add_tree(ifnode)
*/
	/* Optimize this short circuit code */
	OptimizeControl(irnode1, else_label, false)

	rets, err := Ast_statementsequence(statc)

	if !rets {
		return false, err
	}
	
	if while_start_label >=0 {
		jmp := &Jmp{jmp : while_start_label}
		Add_tree(jmp)
	}

	jmp_to_end := &Jmp{}
	jmp_to_end.AddValue(strconv.Itoa(end_label))
	Add_tree(jmp_to_end)

	//temp, _ := strconv.Atoi(ifnode.GetValue())
	label_else:= &Label{label : else_label}
	Add_tree(label_else)

	// Keep iterating over all the elseif and basically just
	// call the if semantic check again

	a:=2
		
L1:
	for a=2 ; a < n.GetChildNos(); a++ {

		nt := n.GetChild(a)

		if nt.Ntype != compiler.Elseifn {
			break L1
		}

		if ret, err := Ast_ifstatement(nt, -1, end_label, depth+1); !ret {
			return false, err
		}
	}

	// Check to see if there is an else statement following this
	
	if a < n.GetChildNos() {

		et := n.GetChild(a)

		if ret, err := Ast_statementsequence(et.GetChild(0)); !ret {
			return false, err
		}
	}

	if depth <= 0 {
		label_ifend := &Label{label : end_label}
		Add_tree(label_ifend)
	}



	return true, "ok"	
}

// Function which checks if the given identifier is in the symb_table stack

func Ast_SymCheckIdent(ident string) (sym_ele compiler.Sym_table, ok bool) {

    for a:= len(st_stack)-1; a >= 0; a-- {
        st_ele := st_stack[a]

        sym_ele, ok = st_ele[ident]

        if ok  {
            return sym_ele, ok
        }

    }

    return compiler.Sym_table{}, false
}


// Function which returns the semantic type of the type given
func Ast_SymGetType(ident string) (s string, ok bool){

    // A very odd hack, because oberon allows internal DT to declared again

    switch ident {
        case "_#INTEGER" :
            return "_#INTEGER", true
        case "_#BOOLEAN" :
            return "_#BOOLEAN", true
        case "_#CHAR" :
            return "_#CHAR", true
        case "_ARRAY_#CHAR" :
            return "_ARRAY_#CHAR", true
        case "_ARRAY_#INTEGER" :
            return "_ARRAY_#INTEGER", true
        case "_ARRAY_#BOOLEAN" :
            return "_ARRAY_#BOOLEAN", true
    }

    ele, ok := Ast_SymCheckIdent(ident)

    // This call should ideally never fail
    if !ok {
        return "", ok
    }

    if ident == s {
        return s, true
    }

    //s, ok = Ast_SymGetType(s)

    if ele.GetNtype() != "_#INTEGER" && ele.GetNtype() != "_#BOOLEAN" && ele.GetNtype() != "_#CHAR" && 
                            ele.GetSize() >= 0 {
        s = "_ARRAY_"+ele.GetNtype()
    } else {
        s = ele.GetNtype()
    }

    return s, ok
}

func Ast_ArrayOf (atype string) (ntype string) {
    return strings.Split(atype, "_ARRAY")[1]
}

/* Function which returns the byte size */
func IR_GetByteSize (ntype string) (size int) {
	switch {
		case ntype == "_#INTEGER":
			return 4
		case ntype == "_#CHAR" :
			return 1
		case ntype == "_#BOOLEAN" :
			return 1
	}

	/* We should not have reached here if semantic analysis worked correctly. (It does reach here in some special cases and I rely on it :P ) */
	return 4
}

func GetLabel() (count int){
	temp := label_count

	label_count ++	

	return temp
}

/* This function accepts 2 irnodes, checks if one of them is CONST, if it is
 * then it just returns a 3 valued tuple in the following form. 
 * (true, const_irnode, the_other_irnode) */

func ContainsConst(irnode1 Node, irnode2 Node) (ret bool, cnode Node, onode Node) {
	if (irnode1 != nil && strings.Contains(irnode1.NodePrint(), "CONST")){
		cnode = irnode1
		onode = irnode2

		return true, cnode, onode
	} else if (irnode2 != nil && strings.Contains(irnode2.NodePrint(), "CONST")) {
		cnode = irnode2
		onode = irnode1

		return true, cnode, onode
	}

	/* If we are here it means that none of the branches had a CONST in it. */

	return false, nil, nil
}

func IsConst (str string) (ret bool){
	return strings.Contains(str, "CONST")
}

func IsAddrgp (str string) (ret bool){
	return strings.Contains(str, "ADDRGP")
}

func GetInt (irnode Node) (val int){
	temp, _ := strconv.Atoi(irnode.GetValue())

	return temp
}

/* When a node is passed to this function, after stripping of the 
 * indir, it checks if this is a array who's index is constant. 
 * If yes then it returns the following tuple -
 * (true, offset - the final thing, adrgp node - the register which contains 
 * the base address) */
 func OptimizeArray(irnode Node) (ret bool, offset int, addrgp Node) {
	/* Traverse the small IR tree to see if this is the type of tree 
	 * we can optimize for an array or not */

	offset = 0
	
	isadd := irnode
	if isadd.NodePrint() == "Addi" {
		left := isadd.GetChild(0)
		//right := isadd.GetChild(1)
	
		if left.NodePrint() == "ADD" {
			offset = offset + GetInt(isadd)

			ileft := left.GetChild(0)
			iright := left.GetChild(1)

			if (IsAddrgp(ileft.NodePrint()) && iright.NodePrint() == "Muli") {
				addrgp = ileft

				if (IsConst(iright.GetChild(0).NodePrint())) {
					offset = offset + (GetInt(iright.GetChild(0)) * GetInt(iright))

					/* This is the array pattern we were looking to optimize */
					return true, offset, addrgp
				} else {
					return false, -1, nil
				}
			} else {
				return false, -1, nil
			}
		} else {
			return false, -1, nil
		}

	} else {
		return false, -1, nil
	}
	return false, -1, nil
 }

 /* Function which traverses a boolean expression tree and optimizes it
  * to perform good shortcuiting. */
func OptimizeControl (irnode Node, endlabel int, invert bool) (ret bool) {

	op := irnode.NodePrint()

	ret = false

	if invert {
		op = InvertControl(op)
	}

	if op == "GTE" {
		ret = true
		cn := &Bge{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)
	} else if op == "LTE" {
		ret = true
		cn := &Ble{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)
	} else if op == "LT" {
		ret = true
		cn := &Blt{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)

	} else if op == "GT" {
		ret = true
		cn := &Bgt{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)

	} else if op == "EQ" {
		ret = true
		cn := &Beq{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)

	} else if (op == "NOT"){
		ret = true
		OptimizeControl(irnode.GetChild(0), endlabel, !invert)
	} else if op == "NE" {
		ret = true
		cn := &Bne{}
		cn.AddValue(strconv.Itoa(endlabel))
		cn.AddChild(irnode.GetChild(0))
		cn.AddChild(irnode.GetChild(1))

		Add_tree(cn)

	} else if op == "BAND" {
		ret = true
		temp := false
		newlabel := endlabel
		if invert {
			temp = true
			newlabel = GetLabel()
		}
		OptimizeControl(irnode.GetChild(0), newlabel, false)
		OptimizeControl(irnode.GetChild(1), endlabel, temp)

		if (invert) {
			l := &Label{}
			l.AddValue(strconv.Itoa(newlabel))
			Add_tree(l)
		}

	} else if op == "BOR" {
		ret = true
		newlabel := GetLabel()
		temp := false
		if invert {
			newlabel = endlabel
			temp = invert
		}
		OptimizeControl(irnode.GetChild(0), newlabel, true)
		OptimizeControl(irnode.GetChild(1), endlabel, temp)
		/* The new label to short-circuit to if the first expr evals true */
		if (!invert) {
			l := &Label{}
			l.AddValue(strconv.Itoa(newlabel))
			Add_tree(l)
		}

	/* } else if op == "BAndi" {
		OptimizeControl(irnode.GetChild(0), endlabel, false)
		c := &Const{}
		c.AddValue(irnode.GetValue())
		OptimizeControl(c, endlabel, false)
		*/
	} else if IsConst(op) {
		ret = true
		cn := &Jmp{}
		/* For AND */
		if irnode.GetValue() == "0" && !invert{
			cn.AddValue(strconv.Itoa(endlabel))
			Add_tree(cn)
		}
		/* For OR */
		if irnode.GetValue() == "1" && invert {
			cn.AddValue(strconv.Itoa(endlabel))
			Add_tree(cn)
		}
	} else if op == "INDIR" {
		ret = true
		/* If the invert flag is set then just apply the not operation on it */
		if (invert) {
			ccn := &IFn{}
			ccn.AddValue(strconv.Itoa(endlabel))
			ccn.AddChild(irnode)
			Add_tree(ccn)

		} else { 
			cn := &IF{}
			cn.AddValue(strconv.Itoa(endlabel))
			cn.AddChild(irnode)
			Add_tree(cn)
		}
	}

	/* Add for addrgp and constants (TRUE / FALSE) */

	return ret
}

/* Function which inverts the operations  for OR, 2 wrongs do make a right */
func InvertControl (op string) (newstr string) {
	if op == "GTE" {
		return "LT"
	} else if op == "LTE" {
		return "GT"
	} else if op == "LT" {
		return "GTE"
	} else if op == "GT" {
		return "LTE"
	} else if op == "EQ" {
		return "NE"
	} else if op == "NE" {
		return "EQ"
	}

	return  op

}

/* Function which returns true if the current statement can be 
 * optimized by applying the boolean branching rules. */
func CanBoolOptimize (op string) (ret bool) {
/* The optimization part generates more code (5 instructions
 * in case of the non optimized part (3 instructions). Hence this
 * has been removed from the optimization check. */
/*
	if op == "GTE" {
		return true
	} else if op == "LTE" {
		return true
	} else if op == "LT" {
		return true
	} else if op == "GT" {
		return true
	} else */
	if op == "EQ" {
		return true
	} else if op == "NE" {
		return true
	} else if op == "BAND" {
		return true
	} else if op == "BOR" {
		return true
	} else if op == "NOT" {
		return true
	}

	return  false
}

/* Function which when you pass an addrgp, returns the constant value inside
 * $fp+constant. This function must only be called on arrays, as it only
 * makes sense on arrays. Otherwise this function returns a -1.
 * Off course for any of this to make sense */
func ReturnConstant (addrgp Node) (constant int) {
	sym_ele, _ := Ast_SymCheckIdent(addrgp.GetValue())

	/* Only makes sense if this is a value parameter */
	if sym_ele.GetFormalVal() == 1 {
		return sym_ele.GetOffset();
	}
	
	return -1;
}
