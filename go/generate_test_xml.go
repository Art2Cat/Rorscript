package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type step struct {
	File      string `xml:"file",attr`
	FileType  string `xml:"fileType",attr`
	Submitter string `xml:"submitterRole",attr`
	Submitted int    `xml:"submittedByUser",attr`
	Clear     bool   `xml:"clearData",attr`
}

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
	var steps []step
	walkPath(dir, steps)
	if len(steps) > 0 {
		toXML(filepath.Join(dir, "test1.xml"), steps)
	} else {
		fmt.Println("empty steps")
	}

	end := time.Now()
	duriation := end.Sub(start)
	fmt.Printf("total time: %v\n", duriation.Seconds())
}

func firstLower(str string) string {
	if str != "" {
		return strings.ToLower(string(str[0])) + str[1:]
	}
	return str
}

func toXML(file string, steps []step) {

	xmlWriter, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer xmlWriter.Close()

	for _, v := range steps {
		if data, err := xml.MarshalIndent(v, "", "\t"); err != nil {
			panic("xml.MarshalIndent FAILED: " + err.Error())
		} else {
			_, err := xmlWriter.Write(append(data, []byte("\n")...))
			if err != nil {
				panic(err)
			}
		}
	}
}

func walkPath(dir string, steps []step) {

	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if info.Mode().IsRegular() && strings.Contains(info.Name(), ".csv") {
			steps = append(steps, *analyzeFileName(info.Name()))
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

}

func analyzeFileName(fileName string) *step {
	lFileName := strings.ToLower(fileName)
	s := new(step)
	s.File = fileName
	s.Clear = false
	if strings.Contains(lFileName, "agent") {
		s.Submitter = "Agent"
		s.Submitted = 903
	} else if strings.Contains(lFileName, "lender") {
		s.Submitter = "Party"
		s.Submitted = 908
	} else if strings.Contains(lFileName, "vendor") {
		s.Submitter = "Vendor"
		s.Submitted = 907
	}

	if strings.Contains(lFileName, "opening") || strings.Contains(lFileName, "balance") {
		s.FileType = "OpeningBalance"
	} else if strings.Contains(lFileName, "trans") {
		s.FileType = "Transaction"
	} else if strings.Contains(lFileName, "master") {
		s.FileType = "Master"
		s.Clear = true
		if s.Submitted == 0 {
			s.Submitter = "Agent"
			s.Submitted = 903
		}
	} else if strings.Contains(lFileName, "fee") {
		s.FileType = "Fee"
	} else if strings.Contains(lFileName, "contract") {
		s.FileType = "Contract"
	}
	return s
}
