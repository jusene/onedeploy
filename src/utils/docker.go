package utils

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
	"strings"
)

func DockerLogin(username, passwd, registry string) {
	cmd := exec.Command("docker", "login", registry, "-u", username, "-p", passwd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal("docker login ", registry, err)
	}
	log.Println(out.String())
}

func DockerPull(image, tag string) {
	target := []string{image, tag}
	cmd := exec.Command("docker", "pull", strings.Join(target, ":"))
	// StdoutPipe Start后命令与标准输出相关的管道，wait方法获取命令结束后会关闭这个管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	// 创建一个流来读取管道内的内容，一行一行读
	reader := bufio.NewReader(stdout)

	for {
		// 以换行符作为一行结尾
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		log.Print(line)
	}
	cmd.Wait()
}

func DockerTagAndPush(src, target, tag string) {
	s := []string{src, tag}
	t := []string{target, tag}
	cmd := exec.Command("docker", "tag", strings.Join(s, ":"), strings.Join(t, ":"))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	log.Println(out.String())

	cmd2 := exec.Command("docker", "push", strings.Join(t, ":"))
	stdout, err := cmd2.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd2.Start()

	// 创建一个流来读取管道内的内容，一行一行读
	reader := bufio.NewReader(stdout)
	for {
		// 以换行符作为一行结尾
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		log.Print(line)
	}
	cmd2.Wait()
}
