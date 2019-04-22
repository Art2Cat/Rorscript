package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"regexp"

	"github.com/tkanos/gonfig"
)

// Config json file
type Config struct {
	/* Url */
	Url string
}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// load BT trackers list url
	config := Config{}
	err = gonfig.GetConf(filepath.Join(dir, "config.json"), &config)
	if err != nil {
		log.Panic(err)
	}

	// get BT Trackers List
	lists := getBTTrackersList(config.Url)

	// join list to string
	btTrackers := strings.Join(lists, ",")

	confPath := filepath.Join(dir, "aria2.conf")
	aria2Config := loadAria2Config(confPath)

	if !strings.Contains(aria2Config, btTrackers) {
		data := replace(aria2Config, btTrackers)
		saveAria2Config(confPath, data)
	}
	exe := filepath.Join(dir, "aria2c.exe")
	// if aria2.exe not exist then panic
	if _, er := os.Stat(exe); os.IsNotExist(er) {
		log.Panic(err)
	}

	run := exec.Command("aria2c.exe", "--conf-path=aria2.conf")
	out, err := run.Output()
	if err != nil {
		panic(err)
	}
	log.Println(out)
	log.Panicln("start success!")
}

func loadAria2Config(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	return string(data)
}

func saveAria2Config(path string, data string) {

	err := ioutil.WriteFile(path, []byte(data), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func replace(data, new string) string {

	p := regexp.MustCompile(`(bt-tracker=.*)`)
	res := p.ReplaceAllString(data, "bt-tracker="+new)
	fmt.Println(res)
	return res
}

func getBTTrackersList(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	trackers := make([]string, 0)
	for _, v := range strings.Split(string(body), "\n") {
		if v != "" {
			trackers = append(trackers, v)
		}
	}
	// log.Println(string(body))
	return trackers
}
