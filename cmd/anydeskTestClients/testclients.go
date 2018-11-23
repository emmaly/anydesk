package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/emmaly/anydesk"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("%s <apiKey> <licenseID>\n", os.Args[0])
		os.Exit(1)
	}

	a, err := anydesk.New(os.Args[1], os.Args[2], nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if len(os.Args) >= 4 {
		i, err := strconv.ParseInt(os.Args[3], 10, 32)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		client, err := a.Client(int(i))
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("%+v\n", client)
	} else {
		data, err := a.Clients(true, 0, 0, "cid", false)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("%+v\n", data)
	}
}
