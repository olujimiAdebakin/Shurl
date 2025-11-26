package main

import (
	"fmt"
	"github.com/olujimiAdebakin/Shurl/initializers"
	
)


func init() {
	initializers.LoadEnvVariables();
	initializers.ConnectToDB();
}

func main() {
	fmt.Println("Welcome to Shurl API!!!!")
}