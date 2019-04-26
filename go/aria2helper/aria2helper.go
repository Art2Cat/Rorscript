package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"regexp"
)

// Config json file
type Config struct {
	TrackerURL string `json:"trackerUrl"`
	ConfigURL  string `json:"confUrl"`
}

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// load BT trackers list url
	config := Config{}
	filepath.Join(dir, "config.json")
	data, err := ioutil.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Panic(err)
	}

	// get BT Trackers List
	lists := getBTTrackersList(config.TrackerURL)

	// join list to string
	btTrackers := strings.Join(lists, ",")

	confPath := filepath.Join(dir, "aria2.conf")
	// if aria2.conf does not exist, download it from config repository
	if _, er := os.Stat(confPath); os.IsNotExist(er) {
		downloadFile(confPath, config.ConfigURL)
	}
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

	// if aria2.session and aria2.log do not exists, create one
	sessionFile := filepath.Join(dir, "aria2.session")
	if _, er := os.Stat(sessionFile); os.IsNotExist(er) {
		createFile(sessionFile)
	}

	logFile := filepath.Join(dir, "aria2.log")
	if _, er := os.Stat(logFile); os.IsNotExist(er) {
		createFile(logFile)
	}

	// start aria2c.exe
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

func createFile(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Panic(err)
	}
	f.Close()
}

func replace(data, new string) string {
	p := regexp.MustCompile(`(bt-tracker=.*)`)
	res := p.ReplaceAllString(data, "bt-tracker="+new)
	return res
}

func getBTTrackersList(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

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

func downloadFile(filepath string, url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Panic(err)
	}
}
