package main

import (
	"compress/flate"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"time"

	"github.com/mholt/archiver"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	wg        sync.WaitGroup
	upgradedb bool
)

func init() {
	flag.BoolVar(&upgradedb, "upgradedb", false, "Upgrade DataBase.")
	flag.BoolVar(&upgradedb, "u", false, "Upgrade DataBase.")

}

func checkArgs(args []string, size int) bool {
	return len(args) < size
}

func main() {
	start := time.Now()
	flag.Parse()
	argSize := 1

	if upgradedb == true {
		argSize = 2
	}

	if checkArgs(flag.Args(), argSize) {
		log.Fatalf(`Use %[1]s in CMD or PowerShell,
		1. normal: %[1]s code_dir
		2. upgrade database: %[1]s -u code_dir db_schema
		3. needs help: %[1]s -help 
		`, os.Args[0])
		os.Exit(-1)
	}

	codeDir := flag.Arg(0)
	fmt.Println(codeDir)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("directory don't exists")
		os.Exit(-1)
	}
	os.Chdir(codeDir)
	cmd := exec.Command("mvn", "clean")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
	fmt.Printf("%s\n", stdoutStderr)

	os.Chdir(dir)

	fmt.Println("start archive files...")
	sourceCodePath := filepath.Join(dir, codeDir)

	files := make([]string, 0)
	files = append(files, sourceCodePath)
	fmt.Printf("total: %d files.\n", len(files))
	zipFileName := fmt.Sprintf("%s.zip", codeDir)
	zipFilePath := filepath.Join(dir, zipFileName)

	if _, err = os.Stat(zipFilePath); !os.IsNotExist(err) {
		fmt.Printf("%s exists.\n", zipFilePath)
		os.Remove(zipFilePath)
		fmt.Println("removed " + zipFilePath)
	}

	err = archiveFiles(files, zipFilePath)
	if err != nil {
		log.Fatalf("compress file failed: %v", err)
		os.Exit(-1)
	}
	fmt.Println("archive files done.")

	if upgradedb == true {
		// update database
		fmt.Println("start upgrade database...")
		dbupgradePath := filepath.Join(sourceCodePath, "bin", "deploy")
		fmt.Printf("cd: %s\n", dbupgradePath)
		os.Chdir(dbupgradePath)

		cmd := exec.Command("dbupgrade.bat", "192.168.1.44", "root", "password", flag.Arg(1))
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
			os.Exit(-1)
		}
		fmt.Printf("%s\n", stdoutStderr)
		fmt.Println("upgrade database done.")
	}

	info := &sshInfo{
		Hostname: "192.168.0.59",
		Port:     22,
		Username: "root",
		Password: "password",
	}

	config := &ssh.ClientConfig{
		User: info.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(info.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// optional tcp connect timeout
		Timeout: 15 * time.Second,
	}

	// connect
	var client *ssh.Client
	client, err = ssh.Dial("tcp", info.getAddress(), config)
	if err != nil {
		log.Fatalf("create ssh client failed: %v ", err)
		os.Exit(-1)
	}
	defer client.Close()

	// copy newly code to the server
	fmt.Println("start copy files to server...")
	fmt.Printf("cd: %s\n", dir)
	os.Chdir(dir)
	fmt.Println(zipFilePath)
	dst := fmt.Sprintf("/home/lp_app/lp_live/%s", zipFileName)
	err = copyFileToRemote(client, zipFilePath, dst)
	if err != nil {
		log.Fatalf("copy file to remote failed: %v", err)
		os.Exit(-1)
	}
	fmt.Println("copy files to server done.")

	// generate cmds
	var cmds = make([]string, 0)

	targetPath := fmt.Sprintf("/home/lp_app/lp_live/%s", codeDir)
	rmOldCode := "rm -rf " + targetPath
	cmds = append(cmds, rmOldCode)

	unzipCode := fmt.Sprintf("unzip /home/lp_app/lp_live/%s -d /home/lp_app/lp_live/", zipFileName)
	cmds = append(cmds, unzipCode)

	rmZipFile := "rm -f /home/lp_app/lp_live/" + zipFileName
	cmds = append(cmds, rmZipFile)

	build := fmt.Sprintf("cd /home/lp_app/lp_live/%s && mvn clean install -Dmaven.test.skip=true", codeDir)
	cmds = append(cmds, build)

	deploy := fmt.Sprintf("cd /home/lp_app/lp_live/ && ./deploy.sh %s", codeDir)
	cmds = append(cmds, deploy)
	cmds = append(cmds, "exit")

	executeRemoteCmd(client, cmds)
	fmt.Println("Done")
	end := time.Now()
	duriation := end.Sub(start)
	fmt.Printf("total time: %v\n", duriation.Seconds())
}

type sshInfo struct {
	Hostname string // hostname
	Port     int64  // port
	Username string // username
	Password string // password
}

func (s *sshInfo) getAddress() string {
	return fmt.Sprintf("%s:%d", s.Hostname, s.Port)
}

func listCompressFiles(dir string) []string {

	files := make([]string, 0)
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if !strings.Contains(path, "target") && !strings.Contains(path, ".svn") && !strings.Contains(path, "archive") && !strings.Contains(path, "db-scripts") {
			files = append(files, path)
		}
		return nil
	})
	if walkErr != nil {
		log.Fatalf("Error walking the path %q: %v\n", dir, walkErr)
	}

	return files
}

func archiveFiles(files []string, dest string) error {
	z := archiver.Zip{
		CompressionLevel:       flate.DefaultCompression,
		MkdirAll:               true,
		SelectiveCompression:   true,
		ContinueOnError:        false,
		OverwriteExisting:      false,
		ImplicitTopLevelFolder: false,
	}

	return z.Archive(files, dest)
}

func copyFileToRemote(client *ssh.Client, sourcePath, destPath string) error {

	// create new SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Fatalf("create sftp client failed: %v ", err)
		return err
	}
	defer sftpClient.Close()

	// create destination file
	dstFile, err := sftpClient.Create(destPath)
	if err != nil {
		log.Fatalf("create remote file failed: %v ", err)
		return err
	}
	defer dstFile.Close()

	// create source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		log.Fatalf("open %s failed: %v", sourcePath, err)
		return err
	}

	// copy source file to destination file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatalf("copy file failed: %v", err)
		return err
	}
	fmt.Printf("%d bytes copied\n", bytes)
	return nil
}

func executeRemoteCmd(client *ssh.Client, cmds []string) {
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment to store output in variable
	// var b bytes.Buffer
	// sess.Stdout = &b
	//sess.Stderr = &amp;b

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}

	// send the commands
	for _, cmd := range cmds {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
