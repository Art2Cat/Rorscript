package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func checkPatternCompile(e error) {
	if e != nil {
		log.Fatal("An error occurred on compile rex pattern")
		panic(e)
	}
}

func main() {
	fmt.Println("vim-go")
	start := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if info.Mode().IsRegular() && strings.Contains(info.Name(), ".xml") && !strings.Contains(info.Name(), "output") {

			replace(path)
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

	end := time.Now()
	duriation := end.Sub(start)
	fmt.Printf("total time: %v\n", duriation.Seconds())
}

func replaceAll(data string, ptn string, val string) string {
	fromRe, err := regexp.Compile(ptn)
	checkPatternCompile(err)
	result := fromRe.ReplaceAllString(data, val)
	return result
}

func replace(filePath string) {
	output := strings.Replace(filePath, "xml", "output", -1)
	if _, err := os.Stat(output); os.IsNotExist(err) {
		data, err := ioutil.ReadFile(filePath)
		check(err)
		fromRe, err := regexp.Compile(`((from|FROM)\s([a-zA-z_]+))`)
		checkPatternCompile(err)
		result := string(data)
		result = fromRe.ReplaceAllString(result, "FROM dbo.$3")
		insertPtn := `(([insertINSERT]+)\s(into|INTO)\s([a-zA-Z_]+))`
		result = replaceAll(result, insertPtn, "INSERT INTO dbo.$4")
		joinPtn := `((join|JOIN)\s([a-zA-Z_]+))`
		result = replaceAll(result, joinPtn, "JOIN dbo.$3")
		updatePtn := `((update|UPDATE)\s([a-zA-Z_]+))`
		result = replaceAll(result, updatePtn, "UPDATE dbo.$3")

		err = ioutil.WriteFile(output, []byte(result), 0664)
		check(err)
	}
}
