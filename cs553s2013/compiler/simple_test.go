package compiler

import (
	"testing"
)

func Test_err1(t *testing.T) {
	s :=
	`
	MODULE simple;

	VAR
		mychar : CHAR;

	PROCEDURE testproc (arr : ARRAY OF CHAR);
		BEGIN
			mychar := arr[3];
			
	END testproc;
	
		BEGIN
			mychar := "a";

		END simple.

`

	if ok, msg,_, _ := Checker(s); !ok {
		 t.Errorf(msg)
	}
}
