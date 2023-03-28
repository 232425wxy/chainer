package cenc

import (
	"fmt"
	"testing"
)

func TestNormalReset(t *testing.T) {
	fmt.Printf("\x1b[%dmred message\x1b[0m\nnon-read message\n", 31)
}

func TestNormalAndBold(t *testing.T) {
	// cyanNormalColor := ColorCyan.Normal()
	// fmt.Print(cyanNormalColor)
	// fmt.Println("normal cyan color test:1")
	// fmt.Println("normal cyan color test:2")

	// magentaBoldColor := ColorMagenta.Bold()
	// fmt.Print(magentaBoldColor)
	// fmt.Println("bold magenta xcolor test:1")
	// fmt.Println("bold magenta xcolor test:2")
	cyanNormalColor := ColorCyan.Normal()
	fmt.Print(cyanNormalColor)
	fmt.Println("normal cyan color test:1")
    fmt.Print(ResetColor())
	fmt.Println("normal cyan color test:2")

	magentaBoldColor := ColorMagenta.Bold()
	fmt.Print(magentaBoldColor)
	fmt.Println("bold magenta xcolor test:1")
    fmt.Print(ResetColor())
	fmt.Println("bold magenta xcolor test:2")
}