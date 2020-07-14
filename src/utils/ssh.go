package utils

import (
	"bytes"
	"deploy/step"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"strings"
	"time"
)

func sshConnect(host string, hostinfo step.HostInfo)  *ssh.Client {
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
	//defer sshClient.Close()
	return sshClient
}


func SSHExec(host string, hostinfo step.HostInfo, command string) string {
	sshClient := sshConnect(host, hostinfo)
	defer sshClient.Close()
	// 创建ssh-session
	session, err1 := sshClient.NewSession()
	if err1 != nil {
		log.Fatal("创建ssh session失败", err1)
	}
	defer session.Close()

	var stdOut, stdErr bytes.Buffer
	session.Stdout = &stdOut
	session.Stderr = &stdErr

	session.Run(command)
	if stdErr.String() != "" && strings.Contains(stdErr.String(), "Created symlink") {
		log.Print(stdErr.String())
	} else if stdErr.String() != "" && strings.Contains(stdErr.String(), "Signature") {
		log.Print(stdErr.String())
	} else if stdErr.String() != "" && strings.Contains(stdErr.String(), "Generating") {
		log.Print(stdErr.String())
	} else if stdErr.String() != "" && strings.Contains(stdErr.String(), "mysqld") {
		return stdErr.String()
	} else if stdErr.String() != "" && strings.Contains(stdErr.String(), "Adding password") {
		log.Print(stdErr.String())
	} else if stdErr.String() != "" {
		log.Fatal("err: ", stdErr.String())
	}
	log.Println(stdOut.String())
	return stdOut.String()
}
