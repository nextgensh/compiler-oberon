A compiler for the oberon programming language.
Outputs Spim assembly code.

Oberon - https://en.wikipedia.org/wiki/Oberon_(programming_language)

Dependencies -

Please note that this compiler depends on the tokenizer and grammer which is in - 
github.com/proebsting/cs553s2013/lexer

Scince I do not control this repository, its contents might become unavaliable at any time.
In that case message me and I shall provide you with a copy of the grammer and tokernizer.

cs553s2013/mylexer -
	Lexical analyzer. The primary aim of this package was to call the tokenizer and 
	prepare for the syntax / semantic checking stage by feeding those stages with 
	tokens. Also checks for illegal characters.

cs553s2013/compiler -
	compiler.go - High level methods to start the syntax / semantic checker.
	Maintains the symbol table, defines the semantic meaning for the AST nodes, AST helper functions.
	
	parser.go - Create the AST from the output of the lexical analysis phase.

	semantic.go - Called by compiler.go and accepts the AST created by parser.
	Does complete semantic analysis on the AST.

cs553s2013/IR - Responsible for the generation of the Intermediate Representation using
	the AST. This also contains helper scripts for creating new nodes with the SPIM assembly code.
