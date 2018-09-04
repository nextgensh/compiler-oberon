package main

import (
	"github.com/proebsting/cs553s2013/lexer"
	"cs553s2013/mylexer"
	"fmt"
	"io/ioutil"
)

func main() {
	//make a new channel
	c := make(chan lexer.Token)

	filename := "testfile3.txt"

	b, err := ioutil.ReadFile(filename)

	if err != nil {
		// Don't panic
	}	

	    defer func(){
        if r := recover(); r != nil {
			fmt.Println("err");
        }   
    }() 
			

	go mylexer.Lexer(string(b), c)

	t := <- c

	for ; t.Enum() != lexer.EOF; t = <-c {
		fmt.Println(t.Enum(), t.Value(), t.Line(), t.Column(), t.Location())
	}
		fmt.Println(t.Enum(), t.Value(), t.Line(), t.Column(), t.Location())

}
