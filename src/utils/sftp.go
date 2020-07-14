package utils

import (
	"deploy/step"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func sftpConnect(host string, hostinfo step.HostInfo) (*sftp.Client, error) {
	conf := &ssh.ClientConfig{
		User:            hostinfo.User,
		Auth:            []ssh.AuthMethod{ssh.Password(hostinfo.Pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, hostinfo.Port)
	sshClient, err := ssh.Dial("tcp", addr, conf)
	if err != nil {
		log.Fatal("创建ssh client失败", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatal("创建sftp client失败", err)
	}
	return sftpClient, err
}

func SFTPutFile(host string, hostinfo step.HostInfo, src, dest string) {
	sftpClient, err := sftpConnect(host, hostinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// src 文件
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal("打开源文件错误:", err)
	}
	defer srcFile.Close()

	// dest 文件
	destFile, err := sftpClient.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()


	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		log.Fatal(err)
	}
	destFile.Write(ff)
	log.Print("拷贝文件", src, "完成")
}

func SFTPut(host string, hostinfo step.HostInfo, src, dest string) {
	sftpClient, err := sftpConnect(host, hostinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// dest 文件
	destFile, err := sftpClient.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	destFile.Write([]byte(src))
}

func SFTPFetchFile(host string, hostinfo step.HostInfo, src, dest string) {
	sftpClient, err := sftpConnect(host, hostinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	// src 文件
	srcFile, err := sftpClient.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// dest 文件
	destFile, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	defer destFile.Close()

	if _, err := srcFile.WriteTo(destFile); err != nil {
		log.Fatal(err)
	}
	log.Print("拉取文件", src, "完成")
}
