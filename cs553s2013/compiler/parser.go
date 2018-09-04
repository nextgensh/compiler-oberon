package compiler

import (
	"github.com/proebsting/cs553s2013/lexer"
	"strconv"
	"strings"
	//"fmt"
)

// Nothing to do for the AST
func TypeDeclaration(lex *lexType, topscope map[string]Sym_table) (ok bool, error string,
										te map[string]Sym_table) {	

	//var te map[string] int

	te = make(map[string]Sym_table)
	sz := 0
	var s string

	var msg string
	var type_name string
	var line int

	if lexer.IDENT == lex.Peek() {
		type_name = lex.Tok_value()
		type_name = "#" + type_name
		line = lex.Tok().Line()
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : pected ident"), make(map[string]Sym_table)
	}
	if lexer.EQUAL == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ="), make(map[string]Sym_table)
	}
	// This array type returns which type it is -
	// 
	if ok, msg, sz, s = ArrayType(lex); !ok {
		return ok, msg, make(map[string]Sym_table)
	}
	
	if ntype, found := topscope[s]; (found && ntype.ntype != "_#INTEGER" && 
													ntype.ntype != "_#BOOLEAN" &&
													ntype.ntype != "_#CHAR" &&
													ntype.size >= 0) {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Type nesting of `"+s+"` not allowed"),
												make(map[string]Sym_table)
	}

	if _, found := topscope[s]; !found {
		
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" :Undefined type `" +s+"`"), make(map[string]Sym_table)
	}

	te[type_name] = Sym_table{ntype : s, size : sz, Line : line}

	return true, "ok", te
}

func Parser_Type(lex *lexType) (ok bool, error string, s string) {
	switch lex.Peek() {
	case lexer.IDENT:
		s = lex.Tok_value()
		s = "#"+s
		lex.Advance() // IDENT
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"), ""
	}
	return true, "ok", s
}

func ArrayType(lex *lexType) (ok bool, error string, size int,
								s string) {
	var msg string

	if lexer.ARRAY == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ARRAY"), 0, ""
	}
	if lexer.INTEGER == lex.Peek() {
		     trim_str := strings.TrimRight(lex.Tok_value(), "H") 
            if trim_str == lex.Tok_value() {
				size, _ = strconv.Atoi(lex.Tok_value())
            } else {
                dec_val, _ := strconv.ParseInt(trim_str, 16, 64)
				size = int(dec_val)
            } 
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected integer"), 0, ""
	}
	if lexer.OF == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected OF"), 0, ""
	}
	if ok, msg, s = Parser_Type(lex); !ok {
		return ok, msg, 0, ""
	}
	return true, "ok", size, s
}

// Nothing to do for AST
func IdentList(lex *lexType) (ok bool, error string, te map[string]Sym_table) {

	te = make(map[string]Sym_table)

	if lexer.IDENT == lex.Peek() {
		te[lex.Tok_value()] = Sym_table{Line : lex.Tok().Line()} // By default this is blank string
		lex.Advance()
	} else {
		//return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident")
	}
L0:
	for {
		switch lex.Peek() {
		case lexer.COMMA:
			lex.Advance() // COMMA
			if lexer.IDENT == lex.Peek() {
				te[lex.Tok_value()] = Sym_table{Line : lex.Tok().Line()} // By default this is blank string
				lex.Advance()
			} else {
				//return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident")
			}
		default:
			break L0
		}
	}

	return true, "ok", te
}

//Nothing to do for AST
func VariableDeclaration(lex *lexType, topscope map[string]Sym_table) (ok bool, error string,
											te map[string]Sym_table) {

	var msg string

	if ok, msg, te = IdentList(lex); !ok {
		return ok, msg, te
	}
	if lexer.COLON == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected :"), make(map[string]Sym_table)
	}
	var s string
	if ok, msg, s = Parser_Type(lex); !ok {
			return ok, msg, make(map[string]Sym_table)
	} else {
		// Check if `s` is really a type
		if ntype, found := topscope[s]; !found || ntype.size < 0 {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" :Undefined type `" +s+"`"), make(map[string]Sym_table)
		}	
		for k, e := range te {
			te[k] = Sym_table{ntype: s, size : -1, Line : e.Line, offset : -1}
		}
	}
	
	return true, "ok", te
}

// AST done
func expression(lex *lexType) (ok bool, error string,
								n  Node) {

	var c1, c2, c3  Node
	var msg string

	if ok, msg, c1 = SimpleExpression(lex); !ok {
		return ok, msg,  Node{}
	}
	switch lex.Peek() {
	case lexer.HASH, lexer.GREATEREQUAL, lexer.EQUAL, lexer.GREATER, lexer.LESS, lexer.LESSEQUAL:
		if ok, msg, c2 = relation(lex); !ok {
			return ok, msg,  Node{}
		}
		if ok, msg, c3 = SimpleExpression(lex); !ok {
			return ok, msg,  Node{}
		}
	}

	if c2.Ntype == 0 {
		n = c1
	} else {
		n = c2
		n.AddChild(c1)
		n.AddChild(c3)
	}

	return true, "ok", n
}

// AST done
func relation(lex *lexType) (ok bool, error string,
								n  Node) {									
	switch lex.Peek() {
	case lexer.EQUAL:	
		n.Tok = lex.Tok()
		lex.Advance() // EQUAL
	case lexer.HASH:
		n.Tok = lex.Tok()
		lex.Advance() // HASH
	case lexer.LESS:
		n.Tok = lex.Tok()
		lex.Advance() // LESS
	case lexer.LESSEQUAL:
		n.Tok = lex.Tok()
		lex.Advance() // LESSEQUAL
	case lexer.GREATER:
		n.Tok = lex.Tok()
		lex.Advance() // GREATER
	case lexer.GREATEREQUAL:
		n.Tok = lex.Tok()
		lex.Advance() // GREATEREQUAL
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
	}

	n.Ntype =  Roperator

	return true, "ok", n
}

func SimpleExpression(lex *lexType) (ok bool, error string, 
									n  Node) {

	var c1,c2,c3,c4  Node
	var queue []  Node
	var msg string 

	switch lex.Peek() {
	case lexer.MINUS, lexer.PLUS:
		switch lex.Peek() {
		case lexer.PLUS:
			c1.Tok = lex.Tok()
			c1.Ntype =  Soperator 
			lex.Advance() // PLUS
		case lexer.MINUS:
			c1.Tok = lex.Tok()
			c1.Ntype =  Soperator 
			lex.Advance() // MINUS
		default:
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
		}
	}
	if ok, msg,c2 = term(lex); !ok {
		return ok, msg,  Node{}
	}
	
	if c1.Ntype != 0 {
		c1.AddChild(c2)
		queue = append(queue, c1)
	} else {
		queue = append(queue, c2)
	}

L1:
	for {
		switch lex.Peek() {
		case lexer.MINUS, lexer.OR, lexer.PLUS:
			if ok, msg, c3 = AddOperator(lex); !ok {
				return ok, msg,  Node{}
			}
			
			queue = append(queue, c3)	

			if ok, msg, c4 = term(lex); !ok {
				return ok, msg,  Node{}
			}

			queue = append(queue, c4)

		default:
			break L1
		}
	}

	for ; len(queue) > 1 ; {
		tl := queue[0]
		queue = queue[1:len(queue)]
		to := queue[0]
		queue = queue[1:len(queue)]
		tr := queue[0]
		to.AddChild(tl)
		to.AddChild(tr)
		queue[0] = to
	} 

	n = queue[0]

	return true, "ok", n
}

func AddOperator(lex *lexType) (ok bool, error string,
								n  Node) {
	switch lex.Peek() {
	case lexer.PLUS:
		n.Ntype =  Operator
		n.Tok = lex.Tok()
		lex.Advance() // PLUS
	case lexer.MINUS:
		n.Ntype =  Operator
		n.Tok = lex.Tok()
		lex.Advance() // MINUS
	case lexer.OR:
		n.Ntype =  Boperator
		n.Tok = lex.Tok()
		lex.Advance() // OR
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
	}

	return true, "ok", n
}

func term(lex *lexType) (ok bool, error string,
							n  Node) {
	
	var c1,c2,c3  Node
	var queue []  Node
	var msg string

	if ok, msg, c1 = factor(lex); !ok {
		return ok, msg,  Node{}
	}

	queue = append(queue, c1)

L2:
	for {
		switch lex.Peek() {
		case lexer.STAR, lexer.DIV, lexer.MOD , lexer.AMP:
			if ok, msg, c2 = MulOperator(lex); !ok {
				return ok, msg,  Node{}
			}

			queue = append(queue, c2)

			if ok, msg, c3 = factor(lex); !ok {
				return ok, msg,  Node{}
			}

			queue = append(queue, c3)

		default:
			break L2
		}
	}

	for ; len(queue) > 1 ; {
		tl := queue[0]
		queue = queue[1:len(queue)]
		to := queue[0]
		queue = queue[1:len(queue)]
		tr := queue[0]
		to.AddChild(tl)
		to.AddChild(tr)
		queue[0] = to
	} 

	n = queue[0]

	return true, "ok", n
}

func MulOperator(lex *lexType) (ok bool, error string,
									n  Node) {
	switch lex.Peek() {
	case lexer.STAR:
		n.Ntype =  Operator		
		n.Tok = lex.Tok()
		lex.Advance() // STAR
	case lexer.DIV:
		n.Ntype =  Operator		
		n.Tok = lex.Tok()
		lex.Advance() // DIV
	case lexer.MOD:
		n.Ntype =  Operator		
		n.Tok = lex.Tok()
		lex.Advance() // MOD
	case lexer.AMP:
		n.Ntype =  Boperator		
		n.Tok = lex.Tok()
		lex.Advance() // AMP
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
	}


	return true, "ok", n
}

func factor(lex *lexType) (ok bool, error string,
							n  Node) {

	var c2  Node
	var msg string

	switch lex.Peek() {
	case lexer.INTEGER:
		n.Ntype =  Integer
		n.Tok = lex.Tok()
		lex.Advance() // INTEGER
	case lexer.STRING:
		n.Ntype =  String
		n.Tok = lex.Tok()
		lex.Advance() // STRING
	case lexer.TRUE:
		n.Ntype =  Boolean
		n.Tok = lex.Tok()
		lex.Advance() // TRUE
	case lexer.FALSE:
		n.Ntype =  Boolean
		n.Tok = lex.Tok()
		lex.Advance() // FALSE
	case lexer.IDENT:
		if ok, msg, c2 = designator(lex); !ok {
			return ok, msg,  Node{}
		}
		n = c2
		switch lex.Peek() {
		case lexer.LPAREN:
			n.Ntype =  Procedure
			var c3  Node
			if ok, msg, c3 = ActualParameters(lex); !ok {
				return ok, msg,  Node{}
			}
			n.AddChild(c2)
			if c3.Ntype != 0 {
				n.AddChild(c3)
			}
		}
	case lexer.LPAREN:
		lex.Advance() // LPAREN
		if ok, msg, n = expression(lex); !ok {
			return ok, msg,  Node{}
		}
		if lexer.RPAREN == lex.Peek() {
			lex.Advance()
		} else {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected )"),  Node{}
		}
	case lexer.TILDE:
		n.Ntype =  Soperator
		n.Tok = lex.Tok()
		lex.Advance() // TILDE
		var c1   Node
		if ok, msg, c1 = factor(lex); !ok {
			return ok, msg,  Node{}
		}
		n.AddChild(c1)
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
	}
	return true, "ok", n
}

// AST Done
func designator(lex *lexType) (ok bool, error string,
								n  Node) {
	var temp  Node
	is_array := false
	var c1  Node
	var msg string

	if lexer.IDENT == lex.Peek() {
		temp.Tok = lex.Tok()		
		temp.Ntype =  Identifier

		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"),  Node{}
	}
L3:
	for {

		switch lex.Peek() {
		case lexer.LSQUARE:

			n.Ntype =  Indexed
			n.Tok = temp.Tok
			n.AddChild(temp)
			is_array = true

			if ok, msg, c1 = selector(lex); !ok {
				return ok, msg, c1
			}
			n.AddChild(c1)

		default:
			break L3
		}
	}
	
	if !is_array  {
		n = temp
	}

	return true, "ok", n
}

func selector(lex *lexType) (ok bool, error string,
								n  Node) {
	var msg string

	if lexer.LSQUARE == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ["),  Node{}
	}
	if ok, msg, n = expression(lex); !ok {
		return ok, msg,  Node{}
	}
	if lexer.RSQUARE == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ]"),  Node{}
	}
	return true, "ok", n
}

func ExpList(lex *lexType) (ok bool, error string,
								n  Node) {
	n.Ntype =  Explist

	var c1  Node
	var msg string

	if ok, msg, c1 = expression(lex); !ok {
		return ok, msg,  Node{}
	}

	n.Tok = c1.Tok

	n.AddChild(c1)
	var c2  Node
L4:
	for {
		switch lex.Peek() {
		case lexer.COMMA:
			lex.Advance() // COMMA
			if ok, msg, c2 = expression(lex); !ok {
				return ok, msg,  Node{}
			}
			n.AddChild(c2)
		default:
			break L4
		}
	}
	return true, "ok", n
}

func ActualParameters(lex *lexType) (ok bool, error string, 
										n  Node) {

	var msg string

	if lexer.LPAREN == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ("),  Node{}
	}
	switch lex.Peek() {
	case lexer.MINUS, lexer.FALSE, lexer.INTEGER, lexer.NIL, lexer.TILDE, lexer.PLUS, lexer.TRUE, lexer.LPAREN, lexer.IDENT, lexer.STRING:
		if ok, msg, n = ExpList(lex); !ok {
			return ok, msg,  Node{}
		}
	}
	if lexer.RPAREN == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected )"),  Node{}
	}
	return true, "ok", n
}

func statement(lex *lexType) (ok bool, error string,
								n  Node) {

	var msg string

	switch lex.Peek() {
	case lexer.IF, lexer.FOR, lexer.IDENT, lexer.WHILE:
		switch lex.Peek() {
		case lexer.IDENT:
			if ok, msg, n = AssignOrCall(lex); !ok {
				return ok, msg,  Node{}
			}
		case lexer.IF:
			if ok, msg, n = IfStatement(lex); !ok {
				return ok, msg,  Node{}
			}
		case lexer.WHILE:
			if ok, msg, n = WhileStatement(lex); !ok {
				return ok, msg,  Node{}
			}
		case lexer.FOR:
			if ok, msg, n = ForStatement(lex); !ok {
				return ok, msg,  Node{}
			}
		default:
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
		}
	}
	return true, "ok", n
}

// AST Done
func AssignOrCall(lex *lexType) (ok bool, error string, 
									n  Node) {

	var c1, c2  Node
	var msg string

	if ok, msg, c1  = designator(lex); !ok {
		return ok, msg,  Node{}
	}
	switch lex.Peek() {
	case lexer.COLONEQUAL:
		
		n.Ntype =  Assoperator
		n.Tok = lex.Tok()

		lex.Advance() // COLONEQUAL
		if ok, msg, c2 = expression(lex); !ok {
			return ok, msg, c2
		}
	case lexer.LPAREN:

		n.Ntype =  Procedure
		n.Tok = c1.Tok
	
		if ok, msg, c2 = ActualParameters(lex); !ok {
			return ok, msg, c2
		}
	default:
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : unexpected token"),  Node{}
	}

	n.AddChild(c1)
	n.AddChild(c2)

	return true, "ok", n
}

func StatementSequence(lex *lexType) (ok bool, error string,
										n  Node) {

	n.Ntype =  Statseq

	var c1,c2  Node
	var msg string

	if ok, msg, c1 = statement(lex); !ok {
		return ok, msg,  Node{}
	}
	n.Tok = c1.Tok
	n.AddChild(c1)
L5:
	for {
		switch lex.Peek() {
		case lexer.SEMICOLON:
			lex.Advance() // SEMICOLON
			if ok, msg, c2 = statement(lex); !ok {
				return ok, msg,  Node{}
			}
			n.AddChild(c2)
		default:
			break L5
		}
	}
	return true, "ok", n
}

func IfStatement(lex *lexType) (ok bool, error string,
									n  Node) {
	n.Ntype =  Ifn
	n.Tok = lex.Tok()	

	var c1, c2  Node
	var msg string

	if lexer.IF == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected IF"),  Node{}
	}
	if ok, msg, c1 = expression(lex); !ok {
		return ok, msg,  Node{}
	}
	if lexer.THEN == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected THEN"),  Node{}
	}
	if ok, msg, c2 = StatementSequence(lex); !ok {
		return ok, msg,  Node{}
	}

	n.AddChild(c1)
	n.AddChild(c2)
L6:
	for {

		var ei1, ei2, ei3  Node

		ei1.Ntype =  Elseifn
		ei1.Tok = lex.Tok()

		switch lex.Peek() {
		case lexer.ELSIF:
			lex.Advance() // ELSIF
			if ok, msg, ei2 = expression(lex); !ok {
				return ok, msg,  Node{}
			}
			if lexer.THEN == lex.Peek() {
				lex.Advance()
			} else {
				return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected THEN"),  Node{}
			}
			if ok, msg, ei3 = StatementSequence(lex); !ok {
				return ok, msg,  Node{}
			}
			ei1.AddChild(ei2)
			ei1.AddChild(ei3)
			n.AddChild(ei1)
		default:
			break L6
		}
	}
	switch lex.Peek() {
	case lexer.ELSE:
		var e1, e2  Node
		e1.Ntype =  Elsen
		lex.Advance() // ELSE
		if ok, msg, e2 = StatementSequence(lex); !ok {
			return ok, msg,  Node{}
		}
		e1.AddChild(e2)
		n.AddChild(e1)
	}
	if lexer.END == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected END"),  Node{}
	}
	return true, "ok", n
}

func WhileStatement(lex *lexType) (ok bool, error string,
									n  Node) {
	n.Ntype =  Whilen	

	var c1, c2  Node
	var msg string

	if lexer.WHILE == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected WHILE"),  Node{}
	}
	if ok, msg, c1 = expression(lex); !ok {
		return ok, msg,  Node{}
	}
	if lexer.DO == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected DO"),  Node{}
	}
	if ok, msg, c2 = StatementSequence(lex); !ok {
		return ok, msg,  Node{}
	}

	n.AddChild(c1)
	n.AddChild(c2)

//L7:
//	for {
//		switch lex.Peek() {
//		case lexer.ELSIF:
//			lex.Advance() // ELSIF
//			if ok, msg, t1 := expression(lex); !ok {
//				return ok, msg,  Node{}
//			}
//			if lexer.DO == lex.Peek() {
//				lex.Advance()
//			} else {
//				return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected DO"),  Node{}
//			}
//			if ok, msg, t1 := StatementSequence(lex); !ok {
//				return ok, msg,  Node{}
//			}
//		default:
//			break L7
//		}
//	}
	if lexer.END == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected END"),  Node{}
	}
	return true, "ok", n
}

func ForStatement(lex *lexType) (ok bool, error string,
									n  Node) {

	n.Ntype =  Forn

	var c1, c2, c3, c4, c5, c6  Node
	var msg string

	if lexer.FOR == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected FOR"),  Node{}
	}
	if lexer.IDENT == lex.Peek() {
		c1.Ntype =  Identifier
		c1.Tok = lex.Tok()
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"),  Node{}
	}
	if lexer.COLONEQUAL == lex.Peek() {
		c2.Ntype =  Assoperator
		c2.Tok = lex.Tok()
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected :="),  Node{}
	}
	if ok, msg, c3 = expression(lex); !ok {
		return ok, msg,  Node{}
	}
	c2.AddChild(c1)
	c2.AddChild(c3)
	n.AddChild(c2)
	if lexer.TO == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected TO"),  Node{}
	}
	if ok, msg, c4 = expression(lex); !ok {
		return ok, msg,  Node{}
	}
	n.AddChild(c4)
	switch lex.Peek() {
	case lexer.BY:
		lex.Advance() // BY
		if lexer.INTEGER == lex.Peek() {
			c5.Ntype =  Forstep
			c5.Tok = lex.Tok()
			n.AddChild(c5)
			lex.Advance()
		} else {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected integer"),  Node{}
		}
	}
	if lexer.DO == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected DO"),  Node{}
	}
	if ok, msg, c6 = StatementSequence(lex); !ok {
		return ok, msg,  Node{}
	}
	n.AddChild(c6)
	if lexer.END == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected END"),  Node{}
	}
	return true, "ok", n
}

func ProcedureDeclaration(lex *lexType, topscope map[string]Sym_table) (ok bool, error string,
											n  Node,
											t map[string]map[string]Sym_table) {

	var msg string
	var pname string	
	var te, tee map[string]Sym_table
	var ret_type string
	var type_arr []string
	var type_val []int
	var names_var []string

	t = make(map[string]map[string]Sym_table)
	if ok, msg, pname, type_arr, type_val, ret_type, tee, names_var = ProcedureHeading(lex, topscope); !ok {
		return ok, msg,  Node{}, make(map[string]map[string]Sym_table)
	}
	// Quickly check if the return type of this procedure is pre-defined
	// type in the scope above (module)

	if e, found := topscope[ret_type]; len(ret_type)>0 && (!found || e.size < 0){
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" :Type `"+ret_type+"` is undefined\n"), Node{}, make(map[string]map[string]Sym_table)
	}

	if lexer.SEMICOLON == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ;"),  Node{}, 
												make(map[string]map[string]Sym_table)
	}
	if ok, msg, n, te = ProcedureBody(lex, topscope); !ok {
		return ok, msg,  Node{}, make(map[string]map[string]Sym_table)
	}

	n.void = pname
	
	// Check for duplicates
	for k, _ := range tee {
		if e, found := te[k]; found {
			return false, lex.ErrorMessage(strconv.Itoa(e.Line)+" : Variable "+ k  +" redeclared\n"),  Node{}, 
																	make(map[string]map[string]Sym_table)
		}
		te[k] = tee[k]
	}
	
	te["#"+pname] = Sym_table {proc_ret : ret_type, arg_type : type_arr, arg_val : type_val, argument_names : names_var}

	t[pname] = te

	if lexer.IDENT == lex.Peek() {
		if pname != "_"+lex.Tok().Value() {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Procedure end ident does not match "),  Node{},
											make(map[string]map[string]Sym_table)
			
		}			
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"),  Node{},
										make(map[string]map[string]Sym_table)
	}
	return true, "ok", n, t
}

func ProcedureHeading(lex *lexType, topscope map[string]Sym_table) (ok bool, error string, s string, type_arr []string, val_arr []int, ret_type string, tee map[string]Sym_table, var_names []string) {

	tee = make(map[string]Sym_table)

	var te map[string]Sym_table
	var msg string

	var var_names_temp []string

	if lexer.PROCEDURE == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected PROCEDURE"), "", nil, nil, "", make(map[string]Sym_table), var_names
	}
	if lexer.IDENT == lex.Peek() {
		s = "_"+lex.Tok_value()
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"), "", nil, nil, "", make(map[string]Sym_table), var_names
	}
	if lexer.LPAREN == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ("), "", nil, nil, "", make(map[string]Sym_table), var_names
	}
	switch lex.Peek() {
	case lexer.VAR, lexer.IDENT:
		if ok, msg, te, var_names = FPSection(lex, topscope); !ok {
			return ok, msg, "", nil, nil, "", make(map[string]Sym_table), var_names
			}
		for k, e := range te {

			if e.size >= 0 {
				type_arr = append(type_arr, "_ARRAY_"+e.ntype)
			} else {
				type_arr = append(type_arr, e.ntype)
			}

			val_arr = append(val_arr, e.formal_val)

			if _, found := tee[k]; found {
				return false, lex.ErrorMessage(strconv.Itoa(e.Line)+" : Variable "+ k  +" redeclared\n"), "", nil, nil, "",
																			make(map[string]Sym_table), var_names
			}

			tee[k] = e
		}
	L8:
		for {
			switch lex.Peek() {
			case lexer.SEMICOLON:
				lex.Advance() // SEMICOLON
				if ok, msg, te, var_names_temp = FPSection(lex, topscope); !ok {
					return ok, msg, "", nil, nil, "", make(map[string]Sym_table), var_names
				}
				for j:=0; j < len(var_names_temp); j++ {
					var_names = append(var_names, var_names_temp[j])
				}
				for k, e := range te {

					if e.size >= 0 {
						type_arr = append(type_arr, "_ARRAY_"+e.ntype)
					} else {
						type_arr = append(type_arr, e.ntype)
					}

			val_arr = append(val_arr, e.formal_val)

			if _, found := tee[k]; found {
				return false, lex.ErrorMessage(strconv.Itoa(e.Line)+" : Variable "+ k  +" redeclared\n"), "", nil, nil, "",
																			make(map[string]Sym_table), var_names
			}

					tee[k] = e
				}
			default:
				break L8
			}
		}
	}
	if lexer.RPAREN == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected )"), "", nil, nil, "", make(map[string]Sym_table), var_names
	}
	switch lex.Peek() {
	case lexer.COLON:
		lex.Advance() // COLON
		if lexer.IDENT == lex.Peek() {
			ret_type = lex.Tok().Value()
			lex.Advance()
		} else {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"), "", nil, nil, "", make(map[string]Sym_table), var_names
		}
	}
	return true, "ok", s, type_arr, val_arr, ret_type, tee, var_names
}

func ProcedureBody(lex *lexType, topscope map[string]Sym_table) (ok bool, error string,
									n  Node, te map[string]Sym_table) {

	var c2, r2  Node
	var msg string

	if ok, msg, _, te, _ = DeclarationSequence(lex, topscope); !ok {
		return ok, msg,  Node{}, make(map[string]Sym_table)
	}
	// This child should usually be  Node{}, if it is not then it is an error
	switch lex.Peek() {
	case lexer.BEGIN:
		lex.Advance() // BEGIN
		if ok, msg, c2 = StatementSequence(lex); !ok {
			return ok, msg,  Node{}, make(map[string]Sym_table)
		}
		n.AddChild(c2)
	}
	switch lex.Peek() {
	case lexer.RETURN:
		lex.Advance() // RETURN
		if ok, msg, r2 = expression(lex); !ok {
			return ok, msg,  Node{}, make(map[string]Sym_table)
		}
		n.AddChild(r2)
	}
	if lexer.END == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected END"),  Node{}, 
														make(map[string]Sym_table)
	}
	return true, "ok", n, te
}

func DeclarationSequence(lex *lexType, topscope map[string]Sym_table) (ok bool, error string,
											n  Node, 
											te map[string]Sym_table,
											t map[string]map[string]Sym_table) {

	var c1  Node
	var tte map[string]Sym_table
	var tt map[string]map[string]Sym_table

	n.Ntype =  Declseq
	var msg string

	te = make(map[string]Sym_table)
	t = make(map[string]map[string]Sym_table)

	topmap := make(map[string]Sym_table)
		for k, e := range topscope {
			topmap[k] = e	
		}

	// No need to include this in the AST
	switch lex.Peek() {
	case lexer.TYPE:
		lex.Advance() // TYPE
			
	L9:
		for {
			switch lex.Peek() {
			case lexer.IDENT:
				if ok, msg,  tte= TypeDeclaration(lex, topmap); !ok {
					return ok, msg,  Node{}, make(map[string]Sym_table),
													make(map[string]map[string]Sym_table)
				}
				for k,e := range tte {
					if _, v := te[k]; v {
						return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Type "+ k  +" redeclared\n"),  Node{}, 
																make(map[string]Sym_table),
																make(map[string]map[string]Sym_table)
					
					}
				te[k] = e
				topmap[k] = e
				}
				if lexer.SEMICOLON == lex.Peek() {
					lex.Advance()
				} else {
					return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ;"),  Node{}, 
														make(map[string]Sym_table),
														make(map[string]map[string]Sym_table)
				}
			default:
				break L9
			}
		}
	}
	
	// No need to include this in the AST
	switch lex.Peek() {
	case lexer.VAR:
		lex.Advance() // VAR
		
	L10:
		for {
			switch lex.Peek() {
			case lexer.IDENT:
				if ok, msg, tte = VariableDeclaration(lex, topmap); !ok {
					return ok, msg,  Node{}, make(map[string]Sym_table),
													make(map[string]map[string]Sym_table)
				}
				for k,e := range tte {
					if _, v := te[k]; v {
						return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Variable "+ k  +" redeclared\n"),  Node{}, 
																make(map[string]Sym_table),
																make(map[string]map[string]Sym_table)
					
					}
					te[k] = e
					topmap[k] = e
				}
				if lexer.SEMICOLON == lex.Peek() {
					lex.Advance()
				} else {
					return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ;"),  Node{}, 
															make(map[string]Sym_table),
															make(map[string]map[string]Sym_table)
				}
			default:
				break L10
			}
		}
	}

L11:
	for {
		switch lex.Peek() {
		case lexer.PROCEDURE:
			if ok, msg, c1, tt = ProcedureDeclaration(lex, topmap); !ok {
				return ok, msg,  Node{}, make(map[string]Sym_table),
											make(map[string]map[string]Sym_table)
			}
			for k, e := range tt {
				if _, found := t[k]; found && k!="_CHR" && k!="_WRITE" && k!="_COPY"&& k!="_ORD"{
					return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : `"+k+"` redeclared"), Node{}, make(map[string]Sym_table),
																	make(map[string]map[string]Sym_table)
				}
				t[k] = e
			}
			n.AddChild(c1)
			if lexer.SEMICOLON == lex.Peek() {
				lex.Advance()
			} else {
				return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ;"),  Node{}, 
															make(map[string]Sym_table),
														make(map[string]map[string]Sym_table)
			}
		default:
			break L11
		}
	}
	return true, "ok", n, te,t 
}

func FPSection(lex *lexType, topscope map[string]Sym_table) (ok bool, error string, te map[string]Sym_table, var_names[]string) {
	var msg string

	te = make(map[string]Sym_table)

	vtype := 1
	switch lex.Peek() {
	case lexer.VAR:
		vtype = 0
		lex.Advance() // VAR
	}
	if lexer.IDENT == lex.Peek() {
		te[lex.Tok().Value()] = Sym_table{Line : lex.Tok().Line(),
										formal_val : vtype}
		/* Also add to the variable name to an ordered slice */
		var_names = append(var_names, lex.Tok().Value())
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"), make(map[string]Sym_table), var_names
	}
L12:
	for {
		switch lex.Peek() {
		case lexer.COMMA:
			lex.Advance() // COMMA
			if lexer.IDENT == lex.Peek() {
				if _, found := te[lex.Tok().Value()]; found {
					return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Variable "+ lex.Tok().Value()  +" redeclared\n"),  
																	make(map[string]Sym_table), var_names
				}
				te[lex.Tok().Value()] = Sym_table{Line : lex.Tok().Line(),
										formal_val : vtype}
				var_names = append(var_names, lex.Tok().Value())
				lex.Advance()
			} else {
				return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"), make(map[string]Sym_table), var_names
			}
		default:
			break L12
		}
	}
	if lexer.COLON == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected :"), make(map[string]Sym_table), var_names
	}
	sz := -1
	nt := ""
	if ok, msg, sz, nt = FormalType(lex); !ok {
		return ok, msg, make(map[string]Sym_table), var_names
	}

	if e, found := topscope[nt]; len(nt)>0 && (!found || e.size < 0){
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" :Type `"+nt+"` is undefined\n"),make(map[string]Sym_table), var_names
	}

	for k, e := range te {
		if sz >= 0 {
			nt = "_ARRAY_"+nt
			sz = -1
		}
		te[k] = Sym_table{ntype : nt, size: sz, Line : e.Line, formal_val : e.formal_val}
	}

	return true, "ok", te, var_names
}

func FormalType(lex *lexType) (ok bool, error string, size int, ntype string) {
	size = -1

	switch lex.Peek() {
	case lexer.ARRAY:
		lex.Advance() // ARRAY
		if lexer.OF == lex.Peek() {
			size = 0
			lex.Advance()
		} else {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected OF"), 0 , ""
		}
	}
	if lexer.IDENT == lex.Peek() {
		ntype = lex.Tok().Value()
		ntype = "#" + ntype
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"), 0, ""
	}
	return true, "ok", size, ntype
}

func module(lex *lexType) (ok bool, error string,
							n  Node, t map[string]map[string]Sym_table) {

	var c1, c2, c3  Node
	var te map[string]Sym_table
	var te1 map[string]Sym_table

	te = make(map[string]Sym_table)
	mname := ""

	// Merge the global and module table into one because 
	// we have only 1 module to worry about

	// note to self = There is a reason why this is 2, do not try to change it
	te["#INTEGER"] = Sym_table {ntype: "_#INTEGER", size : 0}
	te["#BOOLEAN"] = Sym_table {ntype: "_#BOOLEAN", size : 0}
	te["#CHAR"] = Sym_table {ntype: "_#CHAR", size : 1}

	canary := "LIFEHACKERFFERS"

	/* Predefined procedures - They all have 2 underscores in front of them rather than just one */
	te["__LEN"] = Sym_table {proc_ret : "_#INTEGER" , arg_type : []string{"_ARRAY_#CHAR|_ARRAY_#INTEGER|_ARRAY_#BOOLEAN"}, arg_val : []int{1}, visible : 1, argument_names : []string{canary+"a"}, offset : 16} 
	te[canary+"a"] = Sym_table {ntype : "_ARRAY_#CHAR", offset : 12}
	te["__CHR"] = Sym_table {proc_ret : "_#CHAR" , arg_type : []string{"_#INTEGER"}, arg_val : []int{1}, visible : 1, argument_names : []string{canary+"b"}, offset : 16} 
	te[canary+"b"] = Sym_table {ntype : "_#INTEGER", offset : 12}
	te["__COPY"] = Sym_table {arg_type : []string{"_ARRAY_#CHAR", "_ARRAY_#CHAR"}, arg_val : []int{1,0}, visible : 1, argument_names : []string{canary+"c",canary+"d"}, offset : 32} 
	te[canary+"c"] = Sym_table {ntype : "_ARRAY_#CHAR", offset : 12}
	te[canary+"d"] = Sym_table {ntype : "_ARRAY_#CHAR", offset : 16}
	te["__WRITE"] = Sym_table {arg_type : []string{"_ARRAY_#CHAR"}, arg_val : []int{1}, visible : 1, argument_names : []string{canary+"e"}, offset : 28} 
	te[canary+"e"] = Sym_table {ntype : "_ARRAY_#CHAR", offset : 12}
	te["__ORD"] = Sym_table {proc_ret : "_#INTEGER", arg_type : []string{"_#CHAR|_#BOOLEAN"}, arg_val : []int{1}, visible : 1, argument_names : []string{canary+"f"}, offset : 16} 
	te[canary+"f"] = Sym_table {ntype : "_#CHAR", offset : 12}

	t = make(map[string]map[string]Sym_table)

	n.Ntype =  Module
	var msg string

	if lexer.MODULE == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected MODULE"),  Node{},
										make(map[string]map[string]Sym_table)
	}
	if lexer.IDENT == lex.Peek() {
		mname = lex.Tok().Value()
		n.Tok = lex.Tok()
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"),  Node{},
										make(map[string]map[string]Sym_table)
	}
	if lexer.SEMICOLON == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ;"),  Node{},
										make(map[string]map[string]Sym_table)
	}
	if ok, msg, c1, te1, t = DeclarationSequence(lex, te); !ok {
		return ok, msg,  Node{}, make(map[string]map[string]Sym_table)
	}
	
	for k, v := range te1 {
		te[k] = v 
	}

	t["module"] = te

	// Bring the procedure in the module's scope.
	for k, _ := range t {
		if k != "module" {
			if _, found := t["module"][k]; found && k!="_CHR" && k!="_WRITE" && k!="_COPY" && k!="_ORD" {
				return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" :redeclared "+ k), Node{},
											make(map[string]map[string]Sym_table)
			}
			t["module"][k] = t[k]["#"+k]
		}
	}

	n.AddChild(c1)
	switch lex.Peek() {
	case lexer.BEGIN:
		lex.Advance() // BEGIN
		if ok, msg, c2 = StatementSequence(lex); !ok {
			return ok, msg,  Node{}, make(map[string]map[string]Sym_table)
		}
		n.AddChild(c2)
	}
	if lexer.END == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected END"),  Node{},
											make(map[string]map[string]Sym_table)
	}
	if lexer.IDENT == lex.Peek() {
		c3.Ntype =  Moduleename
		c3.Tok = lex.Tok()
		if mname != lex.Tok().Value() {
			return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : Module end ident does not match "),  Node{},
											make(map[string]map[string]Sym_table)
		}
		n.AddChild(c3)
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ident"),  Node{},
											make(map[string]map[string]Sym_table)
	}
	if lexer.DOT == lex.Peek() {
		lex.Advance()
	} else {
		return false, lex.ErrorMessage(strconv.Itoa(lex.Tok().Line())+" : expected ."),  Node{},
											make(map[string]map[string]Sym_table)
	}
	return true, "ok", n, t
}
