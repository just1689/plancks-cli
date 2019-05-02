package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/plancks-cloud/plancks-cli/model"
	"github.com/plancks-cloud/plancks-cli/pc"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/pretty"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

//Commands
var apply bool
var delete bool
var get bool
var install bool
var version bool
var project bool

//Flags
var filename string
var endpoint string
var object string

const versionID = "v1.3"

func main() {
	err := readFirst()
	if err != nil {
		logrus.Error(err)
		return
	}
	readFlags()

	if !apply && !delete && !get && !install && !version && !project {
		logrus.Error(errors.New("No command. Supported commands are apply, delete, get, install and version"))
		return
	}

	//TODO: load filename if "" from .plancks in ~/

	if (filename == "" && apply) || (filename == "" && delete) {
		logrus.Error(errors.New("No filename for command"))
		return
	}

	if (endpoint == "" && apply) || (endpoint == "" && delete) || (endpoint == "" && get) {
		endpoint = "http://127.0.0.1:6227"
		logrus.Println("Assuming endpoint http://127.0.0.1:6227")
	}

	/// End of checking

	if install {
		handleInstall()
		return
	}

	if apply {
		handleApply(&endpoint, &filename)
		return
	}

	if delete {
		handleDelete(&endpoint, &filename)
		return
	}

	if get {
		handleGet(&endpoint, &object)
		return
	}

	if version {
		handleVersion()
		return
	}

	if project {
		handleProject()
		return
	}

}

func handleProject() {
	//TODO: support using a filename provided
	b, err := ioutil.ReadFile("project.json")
	if err != nil {
		logrus.Error(err)
		return
	}

	//TODO: assume project v1 for now
	project := model.ProjectV1{}
	err = json.Unmarshal(b, &project)
	if err != nil {
		logrus.Error(err)
		return
	}

	c := exec.Command("git", "rev-parse", "--short", "HEAD")
	b, err = c.Output()
	if err != nil {
		logrus.Error(err)
		return
	}
	gitRevision := string(b)

	tag := fmt.Sprint(project.TeamName, "/", project.ProjectName, ":", gitRevision)

	//TODO: support supplied dockerfile name
	err = exec.Command("docker", "build", "-t", tag, ".").Run()
	if err != nil {
		logrus.Error(err)
		return
	}

	//////////

	b, err = ioutil.ReadFile(project.Service)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", project.Service, err)
	}

	s := model.Service{}
	s.Image = tag
	b, err = json.Marshal(s)
	if err != nil {
		log.Fatalf("Failed marshal serivce struct: %s'\n", err)
	}

	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(project.Endpoint, "/apply")
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}
	p := pretty.Pretty(b)
	p = pretty.Color(p, pretty.TerminalStyle)
	fmt.Println(string(p))

}

func handleVersion() {
	fmt.Println(versionID)
}

func readFirst() (err error) {
	if len(os.Args) < 2 {
		err = errors.New("Not enough arguments. Provide either apply or delete.")
		return
	}
	if os.Args[1] == "apply" || os.Args[1] == "a" {
		apply = true
	} else if os.Args[1] == "delete" || os.Args[1] == "d" {
		delete = true
	} else if os.Args[1] == "get" || os.Args[1] == "g" {
		get = true
	} else if os.Args[1] == "install" || os.Args[1] == "i" {
		install = true
	} else if os.Args[1] == "version" || os.Args[1] == "v" {
		version = true
	} else if os.Args[1] == "project" || os.Args[1] == "p" {
		project = true
	}
	return
}

func readFlags() {
	for i, s := range os.Args {
		if i < 2 {
			continue
		}
		f, v, e := split(s)
		if e != nil {
			continue
		}
		if f == "-f" || f == "-filename" {
			filename = v
			continue
		}
		if f == "-e" || f == "-endpoint" {
			endpoint = v
			continue
		}
		if f == "-o" || f == "-object" {
			object = v
			continue
		}
	}
}

func split(in string) (f, v string, err error) {
	if !strings.Contains(in, "=") {
		err = errors.New("No value")
		return
	}
	s := strings.Split(in, "=")
	f, v = s[0], s[1]
	return
}

func handleApply(endpoint, file *string) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", *file, err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/apply")
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}
	p := pretty.Pretty(b)
	p = pretty.Color(p, pretty.TerminalStyle)
	fmt.Println(string(p))
}

func handleDelete(endpoint, file *string) {
	b, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Failed to read file %s: '%s'\n", *file, err)
	}
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/delete")
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}
	p := pretty.Pretty(b)
	p = pretty.Color(p, pretty.TerminalStyle)
	fmt.Println(string(p))
}

func handleGet(endpoint, object *string) {
	//http://localhost:6227/route
	client := &http.Client{}
	client.Timeout = time.Second * 5

	uri := fmt.Sprint(*endpoint, "/", *object)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalf("http.NewRequest() failed with '%s'\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do() failed with '%s'\n", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll() failed with '%s'\n", err)
	}
	p := pretty.Pretty(b)
	p = pretty.Color(p, pretty.TerminalStyle)
	fmt.Println(string(p))

}

func handleInstall() {
	pc.Install()
}
