package IR

import (
	"testing"
)

func Test_err1(t *testing.T) {
	s :=

`
MODULE simple

END simple;
`

	if ok, msg := IRGen(s); !ok {
		 t.Errorf(msg)
	}
}
