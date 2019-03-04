package main

import (
	"bufio"
	"fmt"
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
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".java") {
			wg.Add(1)
			go find(path, &sm, res)
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

	wg.Wait()

	saveToCSV(sm, filepath.Join(dir, "result.csv"))

	end := time.Now()
	duration := end.Sub(start)
	fmt.Printf("total time: %v\n", duration.Seconds())
}

type class struct {
	name     string
	fullName string
	pkgName  string
	libName  string
	methods  map[string]int
	usage    int
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

	p := re.MustCompile(`(\([\d]+\sUsage\sfound\))`)
	res = p.ReplaceAllString(string(content), "")

	p = re.MustCompile(`(\([\d]+\susages\sfound\))`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(//\s?[\w.:?(\-)\s]+)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`([\d]*\s@Deprecated)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(//\s?[\w.:?(\-)\s]+)`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	res = p.ReplaceAllString(res, "")

	p = re.MustCompile(`Maven:\s(?P<lib>[\w.:_-]+)`)
	res = p.ReplaceAllString(res, "mvn: ${lib}")

	p = re.MustCompile(`(?P<java>[\w]+\.java)`)
	res = p.ReplaceAllString(res, "class: ${java}")

	p = re.MustCompile(`(?P<mtd>[\w]+\([\w<>?,\s.\[\]]*\))`)
	res = p.ReplaceAllString(res, "mtd: ${mtd}")

	res = addPKGTag(res)

	return res
}

func addPKGTag(content string) string {
	var builder strings.Builder
	pkgRE := re.MustCompile(`\s{12}(?P<pkg>[a-z0-9.:]*)`)
	for _, line := range strings.Split(content, "\n") {
		if !strings.Contains(line, "mvn:") &&
			!strings.Contains(line, "class:") &&
			!strings.Contains(line, "mtd:") {
			s := pkgRE.ReplaceAllString(line, "pkg: ${pkg}\n")
			builder.WriteString(s)
		} else {
			builder.WriteString(fmt.Sprintf("%s\n", line))
		}
	}
	return builder.String()
}

func toClasses(content string) []*class {
	var mvn string
	var pkg string
	var clz *class
	var classes []*class
	mvnRe := re.MustCompile(`mvn:\s([\w.:_-]+)`)
	pkgRe := re.MustCompile(`pkg:\s([a-z0-9.]+)`)
	clzRe := re.MustCompile(`class:\s([\w]+)\.java`)
	mtdRe := re.MustCompile(`mtd:\s([a-zA-Z]+)\(`)
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "mvn:") {
			if res := mvnRe.FindAllStringSubmatch(line, -1); res != nil {
				mvn = res[0][1]
			}
		}
		if strings.Contains(line, "pkg:") {
			if res := pkgRe.FindAllStringSubmatch(line, -1); res != nil {
				pkg = res[0][1]
			}
		}
		if strings.Contains(line, "class:") {
			if res := clzRe.FindAllStringSubmatch(line, -1); res != nil {
				clz = &class{methods: make(map[string]int)}
				clz.name = res[0][1]
				clz.libName = mvn
				clz.pkgName = pkg
				clz.fullName = clz.pkgName + "." + clz.name
				classes = append(classes, clz)
			}
		}
		if strings.Contains(line, "mtd:") {
			if res := mtdRe.FindAllStringSubmatch(line, -1); res != nil {
				if clz != nil {
					name := res[0][1]
					clz.methods[name] = 0
				}
			}
		}
	}
	return classes
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func find(filePath string, sm *sync.Map, res string) {

	clzs := toClasses(res)
	defer wg.Done()
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	_, filename := filepath.Split(filePath)

	var potential []*class

	cBytes, _ := ioutil.ReadAll(reader)
	content := string(cBytes)

	for _, c := range clzs {
		if strings.Contains(content, c.fullName) {
			potential = append(potential, c)
		}
	}

	for _, c := range potential {
		usage := 1
		for k, _ := range c.methods {
			f := string(k[0])
			count := 0
			if f == strings.ToUpper(f) {
				count = strings.Count(content, fmt.Sprintf("new %s", k))
			} else {
				count = strings.Count(content, k)
			}
			c.methods[k] += count
			usage += count
		}
		c.usage = usage
	}

	sm.Store(strings.TrimSuffix(filename, ".java"), potential)
}

func saveToCSV(sm sync.Map, output string) {

	var builder strings.Builder

	sm.Range(func(k, v interface{}) bool {
		vs := v.([]*class)
		for _, c := range vs {
			builder.WriteString(fmt.Sprintf("%s, %s, %s, %d\n", c.fullName, c.libName, k, c.usage))
			for n, count := range c.methods {
				var mtdName string
				if count == 0 {
					continue
					//mtdName = c.fullName
				} else {
					mtdName = c.fullName + "." + n
				}
				builder.WriteString(fmt.Sprintf("%s, %s, %s, %d\n", mtdName, c.libName, k, count))
			}
		}
		return true
	})
	err := ioutil.WriteFile(output, []byte(builder.String()), 0666)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(builder.String())
}
