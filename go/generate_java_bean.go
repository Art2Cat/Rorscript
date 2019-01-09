package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type BeanXML struct {
	lock  *sync.Mutex
	beans []bean
}

type bean struct {
	BeanID   string     `xml:"id,attr"`
	RefClass string     `xml:"class,attr"`
	Parent   string     `xml:"parent,attr"`
	Fields   []property `xml:"property"`
}

type property struct {
	Name string `xml:"name,attr"`
	Ref  string `xml:"ref,attr"`
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

var beanXML = &BeanXML{
	lock:  &sync.Mutex{},
	beans: make([]bean, 5),
}

var wg sync.WaitGroup

func main() {
	fmt.Println("vim-go")
	start := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	findBeans(dir)
	wg.Wait()
	beans := beanXML.beans
	if len(beans) > 0 {
		beanToXML(filepath.Join(dir, "test.xml"), beans)
	} else {
		fmt.Println("empty beans")
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

func beanToXML(file string, beans []bean) {

	xmlWriter, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer xmlWriter.Close()

	for _, v := range beans {
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

func findBeans(dir string) {

	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if info.Mode().IsRegular() && strings.Contains(info.Name(), ".java") {
			wg.Add(1)
			go loadJava(path)
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

}

func loadJava(filePath string) {
	defer wg.Done()
	data, err := ioutil.ReadFile(filePath)
	check(err)
	b := new(bean)
	var fields []property
	pkgRex, err := regexp.Compile(`package (.+);`)
	checkPatternCompile(err)
	pkgResult := pkgRex.FindStringSubmatch(string(data))
	if len(pkgResult) > 0 {
		b.RefClass = pkgResult[1]
	}
	extRex, err := regexp.Compile(`(\w+)\sextends\s(\w+)`)
	checkPatternCompile(err)
	extResult := extRex.FindStringSubmatch(string(data))
	if len(extResult) > 0 {
		b.BeanID = firstLower(extResult[1])
		b.RefClass = b.RefClass + "." + extResult[1]
		if extResult[2] != "" {
			b.Parent = firstLower(extResult[2])
		}
	}
	fieldRex, err := regexp.Compile(`[\(\s]{1,3}(\w+)[\)\s]{1,4}BeanUtils[\.\s]+getBean[\(\s]+"(\w+)"[\s]*\)`)
	checkPatternCompile(err)
	fieldResult := fieldRex.FindAllStringSubmatch(string(data), -1)
	if len(fieldResult) > 0 {
		for _, v := range fieldResult {
			if len(v) > 0 {
				property := new(property)
				(*property).Name = firstLower(v[1])
				(*property).Ref = v[2]
				fields = append(fields, *property)

			}
		}
		b.Fields = fields
	}
	if len(b.BeanID) > 0 {
		beanXML.lock.Lock()
		beanXML.beans = append(beanXML.beans, *b)
		beanXML.lock.Unlock()
	}
}
