package main 

import ("os"
		"fmt"
		"io/ioutil"
		"cs553s2013/IR"
		)

func main(){
	bytes, _ := ioutil.ReadAll(os.Stdin)
	_, code:= IR.IRGen(string(bytes))
	fmt.Println(code)
}
