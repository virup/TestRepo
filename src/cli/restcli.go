package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	pb "server/rpcdef"

	"github.com/fatih/color"
	"github.com/parnurzeal/gorequest"
	"github.com/urfave/cli"
)

var SubCommandPostSession = cli.Command{
	Name:        "postsession",
	Usage:       "post a session",
	Description: "",
	Action:      postSession,
	Flags:       []cli.Flag{},
}

var SubCommandGetSession = cli.Command{
	Name:        "getsession",
	Usage:       "Get session",
	Description: "",
	Action:      getGetSessionInfo,
	Flags:       []cli.Flag{SessionId},
}

var SubCommandList = cli.Command{
	Name:        "listsessions",
	Usage:       "List sessions",
	Description: "",
	Action:      listAllSessions,
	Flags:       []cli.Flag{},
}

func printHelp(c *cli.Context) {
	fmt.Printf("server host = %s\n", os.Getenv("SERVERIP"))
	//cli.HelpPrinter(c.App.Writer, common.AppHelpTemplate, c.App)
}

func getGetSessionInfo(c *cli.Context) {
	cmd := "getsession"
	sid := c.String(cmd)
	_, body, errs := gorequest.New().Get("http://" +
		os.Getenv("SERVERIP") +
		":8080/getsession/" + sid).End()
	if errs != nil {
		log.Fatal("getGetSession REST call error: ", errs)
	}
	fmt.Printf(" session response: %s\n", body)
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getNewSession() pb.SessionInfo {

	var si pb.SessionInfo
	t := time.Now()
	si.SessionTime = t.String()
	si.SessionType = "stype"
	si.InstructorID = randSeq(10)
	si.SessionDesc = "my session"
	return si
}

func getUrl() string {
	serv := os.Getenv("SERVERIP")
	if serv == "" {
		fmt.Printf("IP not set")
		os.Exit(-1)
	}
	return "http://" + serv + ":8080"
}
func postSession(c *cli.Context) {

	serverurl := getUrl()
	s := getNewSession()
	sJson, err := json.Marshal(s)
	if err != nil {
		fmt.Printf("Failed to marshal json %s", err)
	}
	fmt.Printf("Sending json : %s\n", sJson)
	/*
		request := gorequest.New()
		resp, body, errs := request.Post("http://"+os.Getenv("SERVERIP")+
			":8080/postsession").
			Set("Notes", "gorequst is coming!").
			SendString(sJson).
			End()
	*/
	contentReader := bytes.NewReader(sJson)
	req, err := http.NewRequest("POST", serverurl+"/postsession",
		contentReader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notes", "GoRequest is coming!")
	client := &http.Client{}
	resp, _ := client.Do(req)
	if err != nil {
		log.Fatal("postsession REST call error:%s ", err)
	}
	fmt.Printf(" postsession response: %v\n", resp)
}

// Rest call to get list of vCenter s.
func listAllSessions(c *cli.Context) {

	url := "http://" + os.Getenv("SERVERIP") + ":8080/getsessions"
	_, body, err := gorequest.New().Get(url).End()
	if err != nil {
		log.Fatal("listvms REST call error: ", err)
		return
	}

	databytes := []byte(body)
	fmt.Printf("reponse body = %v", body)

	var sessionList []pb.Session
	// Unmarshal string into structs.
	jerr := json.Unmarshal(databytes, &sessionList)
	if jerr != nil {
		log.Fatal("JSON unmarshal error: ", jerr)
		return
	}

	// Loop over structs and display them with some pretty-printing
	for _, s := range sessionList {
		color.Red("Info = %v", s.Info)
	}
	return
}

var Commands = []cli.Command{
	SubCommandList,
	SubCommandGetSession,
	SubCommandPostSession,
}

func main() {
	app := cli.NewApp()
	app.Name = "cli"
	//app.Usage = Usage
	app.Commands = Commands
	app.Flags = []cli.Flag{}

	app.Run(os.Args)
}
