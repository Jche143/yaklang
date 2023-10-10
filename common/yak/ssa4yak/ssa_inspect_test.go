package ssa4yak

import (
	"fmt"
	"testing"
)

func TestSSA_SMOKING(t *testing.T) {
	prog := ParseSSA(`var a = 1;c= a + 1; a = c + 2;`)
	prog.Show()
	fmt.Println(prog.GetErrors().String())
}

func TestSSA_SMOKING_2(t *testing.T) {
	prog := ParseSSA(`var a = {}; a + 1; ; var a = 1;c= a + 1; a = c + 2;`)
	prog.Show()
	fmt.Println(prog.GetErrors().String())
}

func TestSSA_SMOKING_3(t *testing.T) {
	prog := ParseSSA(`var a = {"c": 1}; a + 1; ; var a = 1;c= a + 1; a = c + 2;`)
	prog.Show()
	fmt.Println(prog.GetErrors().String())
}

func TestSSA_SMOKING_4(t *testing.T) {
	prog := ParseSSA(`var a = []; a + 1; ; var a = 1;c= a + 1; a = c + 2;`)
	prog.Show()
	fmt.Println(prog.GetErrors().String())
}
