package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Start, please type your name")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Println("Hello,", text)
	fmt.Println("type your age")
	text, _ = reader.ReadString('\n')
	fmt.Println("Your age is", text)
}
