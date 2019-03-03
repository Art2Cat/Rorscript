package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	re "regexp"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	start := time.Now()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var sm sync.Map

	res := washData("export.txt")
	clzs := toClasses(res)
	fmt.Println(len(clzs))
	for _, c := range clzs {
		if c != nil {
			fmt.Println(len(c.Methods))
		}
	}
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if info.Mode().IsRegular() && strings.Contains(info.Name(), ".java") {
			wg.Add(1)
			go find(path, sm, clzs)
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

	wg.Wait()

	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("total time: %v\n", duration.Seconds())
}

type Class struct {
	Name     string
	FullName string
	LibName  string
	Methods  []string
	usage    int
}

func newClass() *Class {
	return &Class{Methods: make([]string, 1)}
}

func washData(filename string) string {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	content, _ := ioutil.ReadAll(reader)
	var res string

	p := re.MustCompile(`(\([\d]+\susage\sfound\))`)
	res = p.ReplaceAllString(string(content), "")

	p = re.MustCompile(`(\([\d]+\susages\sfound\))`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(//\s?[\w.:?(\-)\s]+)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(@Deprecated)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(//\s?[\w.:?(\-)\s]+)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`Maven:\s(?P<lib>[\w.:_-]+)`)
	res = p.ReplaceAllString(res, "mvn: ${lib}\n")

	p = re.MustCompile(`(?P<java>[\w]+\.java)`)
	res = p.ReplaceAllString(res, "class: ${java}\n")

	p = re.MustCompile(`(?P<mtd>[\w]+\([\w<>?,\s.\[\]]*\))`)
	res = p.ReplaceAllString(res, "mtd: ${mtd}\n")

	res = addPKGTag(res)

	return res
}

func addPKGTag(content string) string {
	var builder strings.Builder
	pkgRE := re.MustCompile(`\s{12}(?P<pkg>[a-z0-9.:]*)`)
	for _, line := range strings.Split(content, "\n") {
		if !strings.Contains(line, "mvn:") &&
			!strings.Contains(line, "java:") &&
			!strings.Contains(line, "mtd:") {
			s := pkgRE.ReplaceAllString(line, "pkg: ${pkg}\n")
			builder.WriteString(s)
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", line))
		}
	}
	return builder.String()
}

func toClasses(content string) []*Class {
	var mvn string
	var pkg string
	var clz *Class
	mvnRe := re.MustCompile(`mvn:\s([\w.:_-]+)`)
	pkgRe := re.MustCompile(`pkg:\s([a-z0-9.]+)`)
	clzRe := re.MustCompile(`class:\s([\w]+)\.java`)
	mtdRe := re.MustCompile(`mtd:\s([a-zA-Z]+)\(`)
	classes := make([]*Class, 1)
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "mvn:") {
			if res := mvnRe.FindAllString(line, -1); res != nil {
				mvn = res[0]
				fmt.Println(mvn)
			}
		}
		if strings.Contains(line, "pkg:") {
			if res := pkgRe.FindAllString(line, -1); res != nil {
				pkg = res[0]
			}
		}
		if strings.Contains(line, "class:") {
			if res := clzRe.FindAllString(line, -1); res != nil {
				clz = newClass()
				clz.Name = res[0]
				clz.LibName = mvn
				clz.FullName = pkg + "." + clz.Name
				classes = append(classes, clz)
			}
		}
		if strings.Contains(line, "mtd:") {
			if res := mtdRe.FindAllString(line, -1); res != nil {
				if clz != nil {
					clz.Methods = append(clz.Methods, res[0])
				}
			}
		}
	}
	return classes
}

func find(filePath string, sm sync.Map, clzs []*Class) {
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	_, filename := filepath.Split(filePath)

	potential := make([]*Class, 1)
	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}
		for _, c := range clzs {
			if strings.Contains(string(line), c.FullName) {
				potential = append(potential, c)
			}
		}

		fmt.Printf("%s \n", line)
	}
	sm.Store(filename, potential)
}
