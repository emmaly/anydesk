package main

import (
	"fmt"
	"os"

	"github.com/emmaly/anydesk"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s <apiKey> <licenseID>\n", os.Args[0])
		os.Exit(1)
	}

	a, _ := anydesk.New(os.Args[1], os.Args[2], nil)

	err := a.AuthTest()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Seems OK!")
}
