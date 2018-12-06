package main

import (
	"bytes"
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

/* TCP network type */
const netTCP = "tcp"

var (
	wg        sync.WaitGroup
	upgradedb bool
)

func init() {
	flag.BoolVar(&upgradedb, "upgradedb", false, "Upgrade DataBase.")
	flag.BoolVar(&upgradedb, "u", false, "Upgrade DataBase.")

}

func main() {
	flag.Parse()

	if len(os.Args) < 1 {
		log.Fatalf(`Use %s in CMD or PowerShell,
		1. normal: %s code_dir
		2. upgrade database: %s -u code_dir db_schema
		`, os.Args[0])
		os.Exit(-1)
	}

	codeDir := os.Args[1]
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("directory don't exists")
		os.Exit(-1)
	}
	fmt.Println("start archive files...")
	sourceCodePath := filepath.Join(dir, codeDir)
	files := listCompressFiles(sourceCodePath)
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
		if len(os.Args) < 2 {
			log.Fatalf(`Use %s in CMD or PowerShell,
			1. normal: %s code_dir
			2. upgrade database: %s -u code_dir db_schema
			`, os.Args[0])
			os.Exit(-1)
		}
		fmt.Println("start upgrade database...")
		dbupgradePath := filepath.Join(sourceCodePath, "bin", "deploy")
		fmt.Printf("cd: %s\n", dbupgradePath)
		os.Chdir(dbupgradePath)

		cmd := exec.Command("dbupgrade.bat", "192.168.1.44", "markit", "markit@123", os.Args[2])
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
		Password: "abc@123",
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
	client, err = ssh.Dial(netTCP, info.getAddress(), config)
	if err != nil {
		log.Fatalf("create ssh client failed: %v ", err)
		os.Exit(-1)
	}
	defer client.Close()

	// copy newly code to the server
	fmt.Println("start copy files to server...")
	fmt.Printf("cd: %s\n", dir)
	os.Chdir(dir)
	err = copyFileToRemote(client, zipFilePath, zipFileName)
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

	unzipCode := fmt.Sprintf("unzip %s -d %s", zipFileName, targetPath)
	cmds = append(cmds, unzipCode)

	rmZipFile := "rm -f " + zipFileName
	cmds = append(cmds, rmZipFile)

	deploy := fmt.Sprintf("cd /home/lp_app/lp_live && ./deploy.sh %s", codeDir)
	cmds = append(cmds, deploy)

	for _, cmd := range cmds {
		wg.Add(1)
		fmt.Println("execute -----> " + cmd)
		err = executeRemoteCmd(client, cmd)
		if err != nil {
			log.Fatalf("execute command %s failed: %v", cmd, err)
			os.Exit(-1)
		}
	}
	wg.Wait()
	fmt.Println("Done")
}

type sshInfo struct {
	Hostname string // hostname
	Port     int64  // port
	Username string // username
	Password string // password
}

func (sshInfo *sshInfo) getAddress() string {
	return fmt.Sprintf("%s:%d", sshInfo.Hostname, sshInfo.Port)
}

func listCompressFiles(dir string) []string {

	files := make([]string, 0)
	walkErr := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Prevent panic by handling failure accessing a path%q: %v\n", dir, err)
			return err
		}
		if !strings.Contains(path, "target") && !strings.Contains(path, ".svn") {
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

func executeRemoteCmd(client *ssh.Client, cmd string) error {
	defer wg.Done()
	outputChannel := make(chan error)

	go func() {
		session, err := client.NewSession()
		if err != nil {
			log.Fatalf("create session failed: %v", err)
			outputChannel <- err
		}
		defer session.Close()
		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf

		err = session.Run(cmd)
		if err != nil {
			log.Fatalf("run cmd %s failed: %v", cmd, err)
			outputChannel <- err
		} else {
			fmt.Println(fmt.Sprintf("%s :\n %s", cmd, stdoutBuf.String()))
			outputChannel <- nil
		}
	}()
	return <-outputChannel
}
