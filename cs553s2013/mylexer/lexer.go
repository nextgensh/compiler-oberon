// Author : Shravan Aras < shravanaras@cs.arizona.edu>
// Date : 15/1/13

package mylexer

import (
	"github.com/proebsting/cs553s2013/lexer"
)

// This structure does the book keeping were the our lexical 
// analyser is currently looking at.
type Lexim struct {
	value string	// The token string value
	location int	// File offset were the token was found
	line int		// Line number were this token can be found
	column int 		// The column number in the line
	enum int		// The integer value for this token
}

// Possible states the input can be in
const (
	start = 0
	ident = 1 
	literal = 2
	literal_real = 3
	literal_hex = 4
	literal_ed = 5
	literal_pm = 6
	comment = 7
	str = 8
)

func Lexer(input string, output chan <- lexer.Token){

	char_a := byte('a')
	char_z := byte('z')
	char_A := byte('A')
	char_Z := byte('Z')
	char_F := byte('F')
	char_0 := byte('0')
	char_9 := byte('9')

	// takes care of matching string "
	str_q := 0
	// takes care of matching (*
	cmt_q := 0

	filelen := len(input)

	line, column := 1, 1

	state := start

	// The current token being parsed
	buf := ""

	a := 0

	var eof_l Lexim

	prev_state := start // Only used in the case of comments
	start_line := 1 // Only used in case of comment errs

	enu := lexer.INTEGER

	for a=0; a < filelen; a ++ {

		switch state {
			
		case start:				
			buf = ""
			if (byte(input[a]) >= char_a && byte(input[a]) <= char_z) ||
					(byte(input[a]) >= char_A && byte(input[a]) <= char_Z) {
				
				buf = buf + string(input[a])
				state = ident
			} else if input[a] == '"'{

				state = str 	// Put this in the string state
				str_q = 1
				buf = buf + string(input[a])
			
		    } else if ((a + 1) < filelen ) && input[a] == '(' &&
											input[a+1] == '*' {
				// We are inside a comment
				cmt_q = 1
				prev_state = start
				state = comment
				start_line = line
				a = a + 1
				column = column + 1
			} else if (input[a] >= char_0 && input[a] <= char_9) {
				buf = buf + string(input[a])
				state = literal
				
			} else { // This might be an operator
				if (a + 1) < filelen {
					buf = string(input[a]) + string(input[a+1])
				} else {
					buf = string(input[a])
				}

				if e, l, s := lexer.Operator(buf); s == true {

					var t Lexim

					switch {
					case l == 1:
						t = Lexim{string(input[a]), a, line, column, e}
					case l == 2:
						t = Lexim{buf, a, line, column, e}
					default:
						// Not a good place to be in
					}
				
					output <- &t
					a = a + (l - 1) // Some adjustment if it is 2 char operator
					column = column + (l - 1)
				} else {
					// This is illegal character 
					if input[a] != ' ' && input[a] != '\t' &&
											input[a] != '\n' &&
											input[a] != '\r' {
						eof_l = Lexim{string(input[a]), a, line, 
														column, 
														lexer.ERROR}
						goto end
					}
				}
			}

		case ident:
			if (byte(input[a]) >= char_a && byte(input[a]) <= char_z) ||
				(byte(input[a]) >= char_A && byte(input[a]) <= char_Z) || 
				(byte(input[a]) >= char_0 && byte(input[a]) <= char_9) {
				
				buf = buf + string(input[a])
			} else {
				// Was this a keyword ?
				if e, s := lexer.Keyword(buf); s == true {
					t := Lexim{buf, a - len(buf), line, column - len(buf), e}
					output <- &t
				} else { // This is just a normal ident
					t := Lexim{buf, a - len(buf), line, column - len(buf), lexer.IDENT}
					output <- &t
				}

				// Adjust some variables for the next iteration
				a = a - 1
				column = column - 1
				state = start
			}

		case literal, literal_real, literal_hex, literal_ed, literal_pm:


			if (input[a] >= char_0 && input[a] <= char_9){
				buf = buf + string(input[a])
			} else if input[a] == '.' && state == literal {
				// a real number, maybe
				buf = buf + string(input[a])
				state = literal_real
				enu = lexer.REAL
			} else if (state == literal || state == literal_hex) && 
													input[a] >= char_A && 
													input[a] <= char_F {
				// A hex number, maybe
				buf = buf + string(input[a])
				state = literal_hex
				enu = lexer.INTEGER
			} else if state == literal_real && (input[a] == 'E' || input[a] == 'D'){
				buf = buf + string(input[a])
				state = literal_ed
				enu = lexer.REAL
			} else if state == literal_ed && (input[a] == '+' || input[a] == '-')  {
				buf = buf + string(input[a])
				state = literal_pm
				enu = lexer.REAL
			} else {

				if (state == literal_hex || state == literal) && input[a] == 'H' {
					buf = buf + string(input[a])
					column = column + 1
					a = a + 1
				} else if (state == literal || state == literal_hex) && input[a] == 'X' {
					buf = buf + string(input[a])
					column = column + 1
					a = a + 1
					enu = lexer.STRING
					
				} else if state == literal_hex {
					eof_l = Lexim{buf, a - len(buf), line, column - len(buf), 
																lexer.ERROR}
					goto end
				}

				// Are we in a pm or ed state, then we must end with an interger.
				if len(buf) > 0 && (state == literal_ed || state == literal_pm) {
					if buf[len(buf)-1] >= char_0 && 
										buf[len(buf)-1] <= char_9 {
						//accepted, do nothing. 
					}else {
						// error, throw out an error statement.
						eof_l = Lexim{buf, a - len(buf), line, column - len(buf), 
															lexer.ERROR}
						goto end
					}
				}

				t := Lexim{buf, a - len(buf), line, column - len(buf), enu}
				output <- &t
	
				a = a - 1
				column = column - 1
				state = start
				enu = lexer.INTEGER
			}

		case comment:

			if (a + 1) < filelen && input[a] == '(' && input[a + 1] == '*' {
				cmt_q = cmt_q + 1
				a = a + 1
				column = column + 1
			} else if (a + 1) < filelen && input[a] == '*' &&
														input[a +1] == ')' {
				cmt_q = cmt_q - 1
				a = a + 1
				column = column + 1
			}
			
			if cmt_q == 0 {
				state = prev_state
			}

		case str:
			if input[a] == '"' {
				str_q = str_q - 1	
				buf = buf + string(input[a])
				t := Lexim{buf, a - len(buf), line, column - len(buf) + 1, lexer.STRING}
				output <- &t
				state = start
			} else {
				buf = buf + string(input[a])
			}
		}

		column = column + 1

		if input[a] == '\n' {
	
			if str_q > 0 {
				// A new-line character in string is not allowed
				goto strerr			
			}

			line = line + 1
			column = 1
			buf = ""	
		}

	}

	// In case some panic is raised
	defer func(){
		if r := recover(); r != nil {
			eof_l = Lexim{"error", a - len(buf), start_line, 
								column - len(buf), lexer.ERROR}
			output <- &eof_l
		}
	}()

	// If there is an unmatched string, throw an error
strerr:
	if str_q > 0 {
		eof_l = Lexim{"string error", a - len(buf), line, 
							column - len(buf), lexer.ERROR}
		goto end
	} else if cmt_q > 0 {
		eof_l = Lexim{"bad comment", a - len(buf), start_line, 
							column - len(buf), lexer.ERROR}
		goto end
	}

	eof_l = Lexim{"<EOF>", a , line, column, lexer.EOF}

end:
	output <- &eof_l
}

// Implement all the methods described inside the interface

func (l *Lexim) Value() string {
	return l.value
}

func (l *Lexim) Location() int {
	return l.location
}

func (l *Lexim) Line() int {
	return l.line
}

func (l *Lexim) Column() int {
	return l.column
}

func (l *Lexim) Enum() int {
	return l.enum
}
