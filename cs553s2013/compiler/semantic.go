// Note0 -> Remember to code "re-declared error"
// Note1 -> All types are defined as _type
//			Ident != types as idents can never have _ in them

package compiler

import (
	"fmt"
	"strings"
	"strconv"
)

var st_stack []map[string]Sym_table

// n -> AST tree node
// st -> symbol table pointer

func Ast_module (n *Node, st map[string]map[string]Sym_table) (ok bool, err string){
	// Load the module and predefined sym table onto the stack

	st_stack = st_stack[0:0]
	st_ele := st["module"]
	st_stack = append(st_stack, st_ele)

	pdecseq := n.getChild(0)

	ok = true
	
	if pdecseq.Ntype == Declseq {
		for a:=0; a < len(pdecseq.children); a++ {
			decseq := pdecseq.getChild(a)
			ok, err := Ast_Procedure(decseq, st)
			if !ok {
				return false, err
			}
		}
	}

	if len(n.children) >=2 {
		decseq := n.getChild(1)
		if decseq.Ntype == Statseq {
			ok, err = Ast_statementsequence(decseq)
		}
	}

	// Assert - If we are here there then we only have the module name to match

	//At the end of the day check if the module name matches the one
	//used to close the module

	return ok, err
}

func Ast_Procedure (n *Node, st map[string]map[string]Sym_table) (ok bool, err string) {
	// Load the scope into the current scope stack
	name := n.void

	st_ele := st[name]
	st_stack = append(st_stack, st_ele)

	proc_header := st["module"][name]

	proc_header.visible = 1
	st["module"][name] = proc_header

	fall_through := false
	fall_next := 1

	for k, e := range st[name]{
		e.offset = -1
		st[name][k] = e
	}


	if len(n.children) >= 1 {
		if n.getChild(0).Ntype != Statseq {
			fall_through = true
			fall_next = 0
		} else {
			ok, err := Ast_statementsequence(n.getChild(0))
			if !ok {
				return false, err
			}
		}
	}

	if len(n.children) >= 2 || fall_through {
		ntype, ok, err := Ast_expression(n.getChild(fall_next), false)
		if !ok {
			return false, err
		}
		ret_type := proc_header.proc_ret;

		st_stack = st_stack[:len(st_stack)-1]

		if ret_type == "INTEGER" || ret_type == "BOOLEAN" || ret_type == "CHAR" {
			// Assert - No need to check if the ret_type is an actual type. 
			// It has already been tested
			ret_type, _ = Ast_SymGetType(ret_type)
		}

		if ntype == "_ARRAY_#CHAR" || ntype == "_#CHAR" {
			temp_type, _ := Ast_SymGetType(ret_type)
			if temp_type == "_ARRAY_#CHAR" {
				ret_type = ntype
			}
		}

		if ret_type != ntype {
			return false, fmt.Sprintf("Return type, found %s expected %s\n",
															ntype,
															ret_type)
		}

	// Cruel hack but in the interest of less time this will have to do
		if _, found := st[name][ret_type]; found {
			return false, fmt.Sprintf("Return type `%s` does not match the outer scope `%s`\n",
														ret_type,
														ret_type)
		}

	} else {
		// Make sure there was no return statement
		if len(proc_header.proc_ret) > 1  {
			return false, fmt.Sprintf("Expected return statement for %s", name)
		}
	}

	/* Insert all the offsets into the symbol table. */
	offset := 0

	/* Adjust for the return address also. */
	offset = offset + 4
	/* Adjust for the $fp also. */
	offset = offset + 4
	/* Adjust for the $sp also. */
	offset = offset + 4

	/* Offset calculation for arguments passed. */
	for i:=0; i < len(proc_header.argument_names); i++ {
		arg_sym_table, _ := Ast_SymCheckIdent(proc_header.argument_names[i])
		arg_sym_table.argument_param = true	/* This is a argument parameter */ 
		/* type */
		arg_type_temp, _ := Ast_SymGetType(proc_header.argument_names[i])
		/* The literal type */
		arg_type, _ := Ast_SymGetType(arg_type_temp)

		if status := Ast_IsArray(arg_type); status {
			/* Make the offset adjustment for arrays */
			/* Well arrays are always passed by refernce, hence this is just
			 * going to be the size of the array. */
			arg_sym_table.offset = offset
			offset = offset + 4
			/* This is because all arrays are passed by reference and we the value
			 * stored inside the stack to be the address not the stack cell address */
			arg_sym_table.formal_val = 0
		} else {
			/* Well this is just a normal variable */
			if(arg_sym_table.formal_val == 1) {
				/* Then this is a value parameter */
				arg_sym_table.offset = offset
				offset = offset + Ast_GetByteSize(arg_type)
				for ; (offset % 4) != 0; {
					offset = offset + 1
				}
				/* Be nice and leave the stack aligned at word boundary */
			} else {
				/* Its passed by reference so only the address will be here */
				arg_sym_table.offset = offset
				offset = offset + 4				
			}
		}

		/* The offset stored in the procedure header entry indicates
		 * the total size of the activation frame */ 			 
		 proc_header.offset = offset
		 st["module"][name] = proc_header

		/* Reflect the change in the actual symbol table */
		st[name][proc_header.argument_names[i]] = arg_sym_table
	}
	

	/* Offset calculation for local variables */	
	for k,e := range st[name] {

		if e.offset >= 0 || e.size >=0 {
			continue
		}

		arg_sym_table := e
		arg_sym_table.argument_param = false
		/* type */
		arg_type_temp, _ := Ast_SymGetType(k)
		/* The literal type */
		arg_type, _ := Ast_SymGetType(arg_type_temp)

		if status := Ast_IsArray(arg_type); status {
			/* Make the offset adjustment for arrays */
			/* We need to get the stack address so must treat these as val parameters */
			arg_sym_table.formal_val = 1 
			offset = offset + 4	 /* This is just a safety adjustment to leave some buffer space */
			array_of := Ast_ArrayOf(arg_type)
			/* Sym table entry for the actual type */
			type_sym_table, _ := Ast_SymCheckIdent(arg_type_temp)
			offset = offset + (type_sym_table.size * Ast_GetByteSize(array_of))
			/* Get this offset to the next closest multiple of 4 (word boundary) */
			for ; (offset % 4 != 0); {
				offset = offset + 1
			}
			offset = offset + 4				/* Make space for the array length */
			arg_sym_table.offset = offset	/* Arrays are stored in reverse on the stack */
			offset = offset + 4
			/* We leave the stack space aligned on word boundary */
		} else {
			/* Well this is just a normal variable */
			arg_sym_table.offset = offset
			arg_sym_table.formal_val = 1
			offset = offset + Ast_GetByteSize(arg_type)
			for ; (offset % 4 != 0); {
				offset = offset + 1
			}
			/* Be nice and leave the stack aligned at word boundary */
		}

		/* The offset stored in the procedure header entry indicates
		 * the total size of the activation frame */ 			 
		 proc_header.offset = offset
		 st["module"][name] = proc_header

		/* Reflect the change in the actual symbol table */
		st[name][k] = arg_sym_table
	}

	// Assert - If we are in a procedure then the stack size can never be more than 2
	if len(st_stack) >= 2 {
		st_stack = st_stack[:len(st_stack)-1]
	}

	return true, "ok"
}

func Ast_statementsequence (n *Node) (ok bool, err string){
	// We iterate over all the children 
	
	for a:=0; a < n.getChildNos(); a++ {
		c := n.getChild(a)
		ntype := c.Ntype

		switch ntype {
			case Assoperator :
				ok, err = Ast_assignorcall(c)
			case Ifn :
				ok, err = Ast_ifstatement(c)
			case Whilen :
				ok, err = Ast_whilestatement(c)
			case Forn :
				ok, err = Ast_forstatement(c)
			case Procedure :
				_, ok, err = Ast_expression(c, false)
				// Even if this returns true a proper
				// function can never be called with its return value ignored.
				if (!ok) {
					return false, err 
				}
		}

		if !ok {
			return ok, err
		}
	}

	return true, "ok"
}

func Ast_forstatement (n *Node) (v bool, err string) {
	
	cexpr := n.getChild(0)

	cont := cexpr.getChild(0)
	
	ctype, ok, err := Ast_expression(cont, false)
	if !ok {
		return false, err
	}

	ctype, _ = Ast_SymGetType(ctype)

	if ctype != "_#INTEGER" {
		return false, fmt.Sprintf("%d : Control variable must be an integer\n", cont.Tok.Line())
	}

	// Assert - if the control variable is integer then the assign
	// statement will force an int or give an error.

	ok, err = Ast_assignorcall(cexpr)
	if !ok {
		return false, err
	}

	toexpr := n.getChild(1)

	// The to expr must also eval to an integer
	ttype, ok, err := Ast_expression(toexpr, false)
	if !ok {
		return false, err
	}

	// Assert - This expr will eval to scalar integer so no need to 
	// translate

	ttype,_ = Ast_SymGetType(ttype)

	if ttype != "_#INTEGER" {
		return false, fmt.Sprintf("%d : In TO expression expecting INTEGER found %s\n",
										toexpr.Tok.Line(),
										ttype)
	}

	// Is the next node a state seq or a BY node

	statorby := n.getChild(2)
	
	if statorby.Ntype == Forstep {
		// We don't really need to do anything here, 
		// since the parser gaurentees that this is an integer literal
		statorby = n.getChild(3)
	} 

	ok, err = Ast_statementsequence(statorby)
	if !ok {
		return false, err	
	}

	return true, "ok"
}

// While AST is a sub set of IF, so more than happy to kick the can down
func Ast_whilestatement (n *Node) (v bool, err string) {
	return Ast_ifstatement(n)
}

func Ast_ifstatement(n *Node) (v bool, err string){
	
	expc := n.getChild(0)
	statc := n.getChild(1)

	exp_type, rete, err := Ast_expression(expc, false)

	if !rete {
		return false, err
	}

	exp_type, _ = Ast_SymGetType(exp_type)

	exp_type = Ast_Literal2Normal(exp_type)

	if exp_type != "_#BOOLEAN" {
		return false, fmt.Sprintf("%d : Expecting BOOLEAN expression found %s\n",
																expc.Tok.Line(),
																exp_type)
	}

	rets, err := Ast_statementsequence(statc)

	if !rets {
		return false, err
	}

	// Keep iterating over all the elseif and basically just
	// call the if semantic check again

	a:=2
		
L1:
	for a=2 ; a < n.getChildNos(); a++ {
		nt := n.getChild(a)

		if nt.Ntype != Elseifn {
			break L1
		}

		if ret, err := Ast_ifstatement(nt); !ret {
			return false, err
		}
	}

	// Check to see if there is an else statement following this
	
	if a < n.getChildNos() {
		et := n.getChild(a)

		if ret, err := Ast_statementsequence(et.getChild(0)); !ret {
			return false, err
		}
	}
	
	return true, "ok"	
}

func Ast_assignorcall(n *Node) (v bool, err string){

	if n.Ntype == Assoperator {

		lc := n.getChild(0)
		lr := n.getChild(1)

		//lt := lc.Tok
		
		//lc_sym, ret := Ast_SymCheckIdent(lt.Value())
		lctype, ret, err := Ast_expression(lc, true)
		backup_lctype := lctype

		if !ret {
			return false, err 
		}
	
		if lc.Tok != nil {
			sym_ele, _ := Ast_SymCheckIdent(lc.Tok.Value())
			/* We must check if this is an array */
			sym_ele_temp, _ := Ast_SymGetType(sym_ele.ntype)
			if sym_ele.formal_val == 1 && Ast_IsArray(sym_ele_temp){
				return false, fmt.Sprintf("%d : Cannot assign to formal variable %s\n",
															lc.Tok.Line(), lc.Tok.Value())
			}
		}

		ntype, ret, err := Ast_expression(lr, false)

		if !ret {
			return false, err	
		}

		ntype, _ = Ast_SymGetType(ntype)
		backup_ntype := ntype

		ntype = Ast_Literal2Normal(ntype)

		// Need to handle things a little differently if we are dealing with
		// character arrays

		lntype := ""


		if ntype == "_ARRAY_#CHAR" || ntype == "_#CHAR" {
			lntype, _ = Ast_SymGetType(lctype)
			if lntype == "_#CHAR" {
				goto L1
			}else if lntype != "_ARRAY_#CHAR" {
				return false, fmt.Sprintf("%d : Cannot assign string literal to %s\n",
													n.Tok.Line(),
													lntype)
			}
			
			if ntype == "_#CHAR" {
				ntype = "_ARRAY_#CHAR"
			}
			
			// Assert - If there is a variable entry in sym_table then there will
			// be a type entry.
			sym_ele, _ := Ast_SymCheckIdent(lctype)
			
			  new_string, stat := strconv.Unquote(lr.Tok.Value())

            if stat != nil {
                /* Now it is 0AX */
                trim_str := strings.TrimRight(n.Tok.Value(), "X")
                dec_str, _ := strconv.ParseInt(trim_str, 16, 64) 
                new_string = fmt.Sprintf("%c", dec_str)
            } 

			if (len(new_string) > sym_ele.size){
				return false, fmt.Sprintf("%d : String literal exceeds size of ARRAY %s\n",
																n.Tok.Line(),
																lc.Tok.Value())
			}
			
		} else {
			//lctype, _ = Ast_SymGetType(lctype)
			//lntype = lctype
			ntype = backup_ntype
			lntype = backup_lctype
			if ntype == "_#INTEGER" || ntype == "_#BOOLEAN" || ntype == "_#CHAR" {
				lntype, _ = Ast_SymGetType(lntype)
			}
		}	

		L1:
		if ntype != lntype {
			// Types of the left and right hand side do not match
			
			nt := n.Tok
			return false, fmt.Sprintf("%d : Type mismatch. Cannot assign %s to %s\n", nt.Line(),
													ntype, lctype)
		}
		
	}

	return true, "ok"
}

// This functions applies semantic rules on a expression and returns its
// final type

func Ast_expression (n *Node, left bool) (ntype string, ok bool, err string){
	
	var sym_ele Sym_table

	switch n.Ntype {
		case Operator : 
			ntype, ok, err = Ast_BinaryOp(n)
			if !ok {
				return "", false, err
			}
		case Integer :
			ntype = "_#INTEGER"
			ok = true
		case String :
			new_string, stat := strconv.Unquote(n.Tok.Value())

            if stat != nil {
                /* Now it is 0AX */
                trim_str := strings.TrimRight(n.Tok.Value(), "X")
                dec_str, _:= strconv.ParseInt(trim_str, 16, 64) 
                new_string = fmt.Sprintf("%c", dec_str)
            } 
			if len(new_string) > 1 && stat == nil{
				ntype = "_ARRAY_#CHAR"
			} else {
				ntype = "_#CHAR"
			}
			ok = true
		case Boolean :
			ntype = "_#BOOLEAN"
			ok = true
		case Procedure :
			
			if n.getChildNos() < 2 {
				return "", false, fmt.Sprintf("%d : %s is not a procedure or invalid number of arguments\n",
													n.getChild(0).Tok.Line(),
													n.getChild(0).Tok.Value())
			}

			proc_name := n.getChild(0)
			sys_module := st_stack[0]
			proc_header, found := sys_module["_"+proc_name.Tok.Value()]

			found2 := found;

			/* If this is not found lets try again after applying the prebuilt conversion */
			if !found || proc_header.visible != 1 {
				proc_header, found2 = sys_module[ConvertPrebuild("_"+proc_name.Tok.Value())]
			}

			if !found2 || proc_header.visible != 1 {
				return "", false, fmt.Sprintf("%d : Proc `%s` undeclared\n", 
															proc_name.Tok.Line(),
															proc_name.Tok.Value())
			}

			expcount := 0

			if n.getChildNos() <= 1 {
				expcount = 0
			 } else {
				expcount = n.getChild(1).getChildNos()
			 }

			if len(proc_header.arg_type) != expcount {
				return "", false, fmt.Sprintf("%d : Not enough argument, expected %d found %d\n",
															proc_name.Tok.Line(),
															len(proc_header.arg_type),
															expcount)
			}
			//Assert - We have the correct number of arguments
			var explist *Node
			if expcount > 0 {
				explist = n.getChild(1)
			}
			for i:=0; i < len(proc_header.arg_type); i++ {		
				wantedtype_str := proc_header.arg_type[i]
				wantedtype_arr := strings.Split(wantedtype_str, "|")

				exptype, ok, err := Ast_expression(explist.getChild(i), false)
				exptype, _ = Ast_SymGetType(exptype)
				isvaltype := proc_header.arg_val[i]
			
				if !ok {
					return "", false, err
				}
				
				temp_exptype, _ := Ast_SymGetType(exptype)
	
				nodetype := explist.getChild(i).Ntype 	

				if nodetype != Identifier && nodetype != Indexed  {
					if isvaltype == 0{
						return "", false, fmt.Sprintf("%d : Cannot pass absolute param %s to variable parameter\n",
																explist.getChild(i).Tok.Line(),
																explist.getChild(i).Tok.Value())
					}
				}
	
				for in:=0; in < len(wantedtype_arr); in++ { 
					wantedtype :=  wantedtype_arr[in]
	
					wantedtype_temp, _  := Ast_SymGetType(wantedtype)
					if Ast_IsArray(wantedtype_temp) {
						wantedtype = wantedtype_temp
					}	

					if wantedtype == "_ARRAY_#INTEGER" && temp_exptype == "_ARRAY_#INTEGER" {
						exptype = "_ARRAY_#INTEGER" 
					}

					if wantedtype == "_ARRAY_#BOOLEAN" && temp_exptype == "_ARRAY_#BOOLEAN" {
						exptype = "_ARRAY_#BOOLEAN" 
					}

					if wantedtype == "_ARRAY_#CHAR" && temp_exptype == "_ARRAY_#CHAR" && exptype != "_ARRAY_#CHAR"{
						exptype = "_ARRAY_#CHAR" 
						goto L2
					}
			
					switch exptype {
						case "_#INTEGER", "_#BOOLEAN", "_#CHAR" :
							wantedtype, _ = Ast_SymGetType(wantedtype) 
					}

					if (exptype == "_ARRAY_#CHAR" || exptype == "_#CHAR") {
						/* SH.hack
						if isvaltype == 0 {
							return "", false, fmt.Sprintf("%d : Cannot pass STRING literal %s to non-actual parameter\n",
																	explist.getChild(i).Tok.Line(),
																	explist.getChild(i).Tok.Value())
						}
						*/
						// Assert - This type must be present, else it would have been caugh 
						// before only
						tempwantedtype, _ := Ast_SymGetType(wantedtype)
						if tempwantedtype != "_ARRAY_#CHAR" {
						} else {
							if wantedtype != "_ARRAY_#CHAR" {
								symele, _ := Ast_SymCheckIdent(wantedtype)
								size := symele.size
								if (len(explist.getChild(i).Tok.Value())-2) > size {
									return "", false, fmt.Sprintf("%d : String literal exceeds size of ARRAY \n",
																		explist.getChild(i).Tok.Line())
								}
							}  
							wantedtype = tempwantedtype
						}
						if wantedtype == "_ARRAY_#CHAR" && exptype == "_#CHAR" {
							exptype = "_ARRAY_#CHAR"
						}
					}
					if exptype == wantedtype {
						goto L2
					}
					
				}		
					return "", false, fmt.Sprintf("%d : Argument type mismatch, expected %s found %s\n",
															proc_name.Tok.Line(),
															wantedtype_str,
															exptype)
			}

			L2:
			
			ntype = proc_header.proc_ret
			ok = true

		case Identifier :
			sym_ele, ok = Ast_SymCheckIdent(n.Tok.Value())

			if !ok {
				// Undefined identifier
				return "", false, fmt.Sprintf("%d : %s is not declared\n", n.Tok.Line(), n.Tok.Value())
			}

			ntype = sym_ele.ntype

		case Soperator :
			ntype, ok, err = Ast_expression(n.getChild(0), false)

			whattype := "_#INTEGER"

			if n.Tok.Value() == "~" {
				whattype = "_#BOOLEAN"
			}

			if !ok {				
				return "", false, err			
			}

			ntype, _ = Ast_SymGetType(ntype)

			if ntype !=  whattype{
				return "", false, fmt.Sprintf("%d : Cannot use %s with type %s\n", n.Tok.Line(),
										n.Tok.Value(), ntype)
			}
			ntype = whattype
		case Roperator :
			lntype, ret, err := Ast_expression(n.getChild(0), false)
			if !ret {
				return "", false, err
			}
			rntype, ret, err:= Ast_expression(n.getChild(1), false)
			if !ret {
				return "", false, err
			}

			lntype, _ = Ast_SymGetType(lntype)
			rntype, _ = Ast_SymGetType(rntype)

			if (n.Tok.Value() == "=" || n.Tok.Value() == "#") && 
				(lntype == "_#INTEGER" || lntype == "_#CHAR" || lntype == "_#BOOLEAN" || lntype == "_ARRAY_#CHAR") &&  
				(rntype == "_#INTEGER" || rntype == "_#CHAR" || rntype == "_#BOOLEAN" || rntype == "_ARRAY_#CHAR") {
				ntype = "_#BOOLEAN"
				ok = true
			} else if (n.Tok.Value() != "=" || n.Tok.Value() != "#") && 
				(lntype == "_#INTEGER" || lntype == "_#CHAR" || lntype == "_ARRAY_#CHAR") &&  
				(rntype == "_#INTEGER" || rntype == "_#CHAR" || rntype == "_ARRAY_#CHAR") {
				ntype = "_#BOOLEAN"
				ok = true
			} else {
				// This is an error.
				return "", false, fmt.Sprintf("%d : %s cannot be used with %s and %s\n", 
										n.Tok.Line(), n.Tok.Value(),
										lntype, rntype)
			}

			// Assert -> if we are here it means the acceptable lntype and rntype
			// have been seen. But are they the same ?
			
			if n.getChild(0).Ntype == String && n.getChild(1).Ntype == String {
				lntype = "_ARRAY_#CHAR"
				rntype = "_ARRAY_#CHAR"
			}

			if lntype != rntype {
				return "", false, fmt.Sprintf("%d : %s cannot be used with %s and %s\n", 
										n.Tok.Line(), n.Tok.Value(),
										lntype, rntype)
			}

		case Boperator :
			lntype, ret, err := Ast_expression(n.getChild(0), false)
			if !ret {
				return "", false, err
			}
			rntype, ret, err := Ast_expression(n.getChild(1), false)

			if !ret {
				return "", false, err
			}

			lntype, _ = Ast_SymGetType(lntype)
			rntype, _ = Ast_SymGetType(rntype)

			if lntype == "_#BOOLEAN" && rntype == "_#BOOLEAN" {
				ntype = "_#BOOLEAN"
				ok = true
			} else {
				return "", false, fmt.Sprintf("%d : %s cannot be used with %s and %s\n", 
										n.Tok.Line(), n.Tok.Value(),
										lntype, rntype)
			}

		case Indexed :
			lc := n.getChild(0)
			lr := n.getChild(1)

	if lc.Tok != nil && left {
				sym_ele, _ := Ast_SymCheckIdent(lc.Tok.Value())
				if sym_ele.formal_val == 1 && Ast_IsArray(sym_ele.ntype){
					return "", false, fmt.Sprintf("%d : Cannot assign to formal variable %s\n",
																lc.Tok.Line(), lc.Tok.Value())
				}
			}
			
			sym_ele, ret := Ast_SymCheckIdent(lc.Tok.Value()) 

			if !ret {
				return "", false, fmt.Sprintf("%d : %s is not declared\n", lc.Tok.Line(), lc.Tok.Value())
				
			}

			arr_type, _ := Ast_SymGetType(sym_ele.ntype)
			
			if ok := Ast_IsArray(arr_type); !ok {
				return "", false, fmt.Sprintf("%d : Invalid indexing operation, %s is not an array\n",
																	lc.Tok.Line(),
																	lc.Tok.Value())
			}

			exptype, ret, err := Ast_expression(lr, false)

			if !ret {
				return "", false, err
			}

			exptype, _ = Ast_SymGetType(exptype)

			if exptype != "_#INTEGER" {
				return "", false, fmt.Sprintf("%d : Array index must be INTEGER found %s\n", 
														lr.Tok.Line(),
														exptype)
			}

			
			// Assert - the ident is an array and the index is a valid integer expression
			ntype = Ast_ArrayOf(arr_type)	
			ok = true
	}

	return ntype, ok, "ok"
}

// Binary operators
// This function should have been merged with Ast_expression,
// it has been kept here only for legacy reasons.
func Ast_BinaryOp (n *Node) (ntype string, ok bool, err string) {

	lc := n.getChild(0)
	lr := n.getChild(1)

	lntype, ret, err := Ast_expression(lc, false)
	if !ret {
		return "", false, err
	}
	rntype, ret, err := Ast_expression(lr, false)
	if !ret {
		return "", false, err
	}

	lntype, _ = Ast_SymGetType(lntype)
	rntype, _ = Ast_SymGetType(rntype)

	if lntype != "_#INTEGER" || rntype != "_#INTEGER" {
		lt := n.Tok
		return "", false, fmt.Sprintf("%d : Cannot apply %s to %s and  %s\n",
										lt.Line(), lt.Value(),
										lntype,
										rntype)
	}

	return "_#INTEGER", true, "ok"
}

// Function which checks if the given identifier is in the symb_table stack

func Ast_SymCheckIdent(ident string) (sym_ele Sym_table, ok bool) {

	for a:= len(st_stack)-1; a >= 0; a-- {
		st_ele := st_stack[a]

		sym_ele, ok = st_ele[ident]

		if ok  {
			return sym_ele, ok
		}
		
	}

	return Sym_table{}, false
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

	if ele.ntype != "_#INTEGER" && ele.ntype != "_#BOOLEAN" && ele.ntype != "_#CHAR" &&
							ele.size >= 0 {
		s = "_ARRAY_"+ele.ntype
	} else {
		s = ele.ntype
	}

	return s, ok
}

// function to convert from literal type string to normal type string
func Ast_Literal2Normal (s string) (rs string) {
	switch s {
		case "_#INTEGER" :
			return "_#INTEGER"
		case "_#BOOLEAN" :
			return "_#BOOLEAN"
		case "_#CHAR" :
			return "_#CHAR"
		case "_ARRAY_#CHAR" :
			return "_ARRAY_#CHAR"
	}

	return s
}

func Ast_IsArray (ident string) (ok bool) {
	return strings.Contains(ident, "_ARRAY")	
}

func Ast_ArrayOf (atype string) (ntype string) {
	return strings.Split(atype, "_ARRAY")[1]
}

/* Function which returns the byte size */
func Ast_GetByteSize (ntype string) (size int) {
    switch {
        case ntype == "_#INTEGER":
            return 4
        case ntype == "_#CHAR" :
            return 1
        case ntype == "_#BOOLEAN" :
            return 1
    }   

    /* We should not have reached here if semantic analysis worked correctly. */ 
    return 4
}

/* Function which does mapping */
func ConvertPrebuild (name string) (convert string) {
	if name == "_WRITE" {
		return "__WRITE";
	} else if name == "_ORD" {
		return "__ORD";
	} else if name == "_CHR" {
		return "__CHR";
	} else if name == "_COPY" {
		return "__COPY";
	} else {
		return "__LEN"
	}		

	/* Otherwise we just return what we got */
	return name;
}
