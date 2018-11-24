package main

import (
	"fmt"
	"os"

	"github.com/emmaly/anydesk"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// globals
	licenseID = kingpin.Flag("license", "AnyDesk License ID").Required().String()
	apiKey    = kingpin.Flag("apikey", "AnyDesk API Key").Required().String()

	// authtest
	cmdAuthTest = kingpin.Command("authtest", "Authentication Test")

	// sysinfo
	cmdSysInfo = kingpin.Command("sysinfo", "AnyDesk System Info")

	// client
	cmdClient = kingpin.Command("client", "Client Tools")

	// client - get
	cmdClientGet   = cmdClient.Command("get", "Get Single Client")
	argClientGetID = cmdClientGet.Arg("id", "Client ID").Required().Int()

	// client - list
	cmdClientList                = cmdClient.Command("list", "List Clients")
	flagClientListLimit          = cmdClientList.Flag("limit", "Max record count returned").Default("-1").Int()
	flagClientListOffset         = cmdClientList.Flag("offset", "Index of first item returned").Default("0").Int()
	flagClientListSort           = cmdClientList.Flag("sort", "Sort results by property name").Enum("id", "alias", "online")
	flagClientListOrder          = cmdClientList.Flag("order", "Sort direction").Default("desc").Enum("desc", "asc")
	flagClientListIncludeOffline = cmdClientList.Flag("includeoffline", "Include offline clients in results").Bool()

	// client - alias
	cmdClientAlias         = cmdClient.Command("alias", "Set Client Alias")
	argClientAliasClientID = cmdClientAlias.Arg("id", "Client ID").Required().Int()
	argClientAliasValue    = cmdClientAlias.Arg("alias", "New Alias").Required().String()

	// session
	cmdSession = kingpin.Command("session", "Session Tools")

	// session - get
	cmdSessionGet   = cmdSession.Command("get", "Get Single Session")
	argSessionGetID = cmdSessionGet.Arg("id", "Client ID").Required().Int()

	// session - list
	cmdSessionList           = cmdSession.Command("list", "List Sessions")
	argSessionListClientID   = cmdSessionList.Arg("cid", "Client ID").Int()
	flagSessionListLimit     = cmdSessionList.Flag("limit", "Max record count returned").Default("-1").Int()
	flagSessionListOffset    = cmdSessionList.Flag("offset", "Index of first item returned").Default("0").Int()
	flagSessionListSort      = cmdSessionList.Flag("sort", "Sort results by property name").Enum("from", "to", "start", "end", "duration")
	flagSessionListOrder     = cmdSessionList.Flag("order", "Sort direction").Default("desc").Enum("desc", "asc")
	flagSessionListDirection = cmdSessionList.Flag("direction", "Session connection direction").Default("inout").Enum("inout", "in", "out")

	// session - close
	cmdSessionClose          = cmdSession.Command("close", "Close an Open Session")
	argSessionCloseSessionID = cmdSessionClose.Arg("id", "Session ID").Required().Int()

	// session - comment
	cmdSessionComment          = cmdSession.Command("comment", "Set Comment on Session")
	argSessionCommentSessionID = cmdSessionComment.Arg("id", "Session ID").Required().Int()
	argSessionCommentText      = cmdSessionComment.Arg("text", "Comment Text").Required().String()
)

func main() {
	cmd := kingpin.Parse()

	ad, err := anydesk.New(*apiKey, *licenseID, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	switch cmd {
	case cmdAuthTest.FullCommand():
		authTest(ad)
	case cmdClientGet.FullCommand():
		clientGet(ad)
	case cmdClientList.FullCommand():
		clientList(ad)
	case cmdClientAlias.FullCommand():
		clientAlias(ad)
	case cmdSessionGet.FullCommand():
		sessionGet(ad)
	case cmdSessionList.FullCommand():
		sessionList(ad)
	case cmdSessionClose.FullCommand():
		sessionClose(ad)
	case cmdSessionComment.FullCommand():
		sessionComment(ad)
	}
}

func authTest(ad *anydesk.AnyDesk) {
	r, err := ad.AuthTest()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println(r.Result)
}

func clientGet(ad *anydesk.AnyDesk) {
	data, err := ad.Client(*argClientGetID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%+v\n", data)
}

func clientList(ad *anydesk.AnyDesk) {
	var sort string
	switch *flagClientListSort {
	case "id":
		sort = anydesk.SortClientID
	case "alias":
		sort = anydesk.SortAlias
	case "online":
		sort = anydesk.SortOnline
	}

	var order bool
	switch *flagClientListOrder {
	case "asc":
		order = true
	case "desc":
		order = false
	}

	opts := &anydesk.ClientsOptions{
		Limit:          *flagClientListLimit,
		Offset:         *flagClientListOffset,
		Sort:           sort,
		Order:          order,
		IncludeOffline: *flagClientListIncludeOffline,
	}

	data, err := ad.Clients(opts)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%+v\n", data)
}

func clientAlias(ad *anydesk.AnyDesk) {
	err := ad.ClientAlias(*argClientAliasClientID, *argClientAliasValue)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("OK!")
}

func sessionGet(ad *anydesk.AnyDesk) {
	data, err := ad.Session(*argSessionGetID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%+v\n", data)
}

func sessionList(ad *anydesk.AnyDesk) {
	var sort string
	switch *flagSessionListSort {
	case "alias":
		sort = anydesk.SortAlias
	case "online":
		sort = anydesk.SortOnline
	}

	var order bool
	switch *flagSessionListOrder {
	case "asc":
		order = true
	case "desc":
		order = false
	}

	opts := &anydesk.SessionsOptions{
		Limit:     *flagSessionListLimit,
		Offset:    *flagSessionListOffset,
		Sort:      sort,
		Order:     order,
		Direction: *flagSessionListDirection,
		// TimeAfter:  *flagSessionListTimeAfter,
		// TimeBefore: *flagSessionListTimeBefore,
	}

	if *argSessionListClientID > 0 {
		opts.ClientID = *argSessionListClientID
	}

	data, err := ad.Sessions(opts)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%+v\n", data)
}

func sessionClose(ad *anydesk.AnyDesk) {
	err := ad.SessionClose(*argSessionCloseSessionID)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("OK!")
}

func sessionComment(ad *anydesk.AnyDesk) {
	err := ad.SessionComment(*argSessionCommentSessionID, *argSessionCommentText)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("OK!")
}
