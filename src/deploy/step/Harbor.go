package step

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"deploy/temp"
	"deploy/utils"
)

func DeployHarbor(config *viper.Viper) {
	host := config.Get("server.harbor").(map[string]interface{})["ip"].(string)
	info := temp.GetHostInfo(config, host)
	pkg := config.Get("package").(map[string]interface{})["harbor"].(string)
	dockerPkg := config.Get("package").(map[string]interface{})["docker"].(map[string]interface{})

	// 注册证书
	genHarbor(config)
	utils.SSHExec(host, info, "yum install -y docker-compose && mkdir -p /ddhome/ && rm -f /tmp/"+pkg)
	// 上传docker安装包
	var dockerSlice []string
	for _, docker := range dockerPkg {
		utils.SFTPut(host, info, "bin/"+docker.(string), "/tmp/"+docker.(string))
		dockerSlice = append(dockerSlice, docker.(string))
	}

	utils.SFTPutFile(host, info, "bin/"+pkg, "/tmp/"+pkg)
	utils.SSHExec(host, info, "tar xf /tmp/"+pkg+" -C /ddhome")
	utils.SSHExec(host, info,"cd /tmp/; yum localinstall -y "+strings.Join(dockerSlice, " "))
	utils.SFTPutFile(host, info, "tmp/harbor/harbor.crt", "/ddhome/harbor/harbor.crt")
	utils.SFTPutFile(host, info, "tmp/harbor/harbor.key", "/ddhome/harbor/harbor.key")

	// 生成配置文件
	harborAttr := new(temp.HarborAttr)
	harborAttr.DOMAIN = config.Get("registry.local").(map[string]interface{})["domain"].(string)
	harborCfg := utils.RendTemp(temp.HarborTmpl, harborAttr)
	utils.SFTPut(host, info, harborCfg, "/ddhome/harbor/harbor.cfg")
	utils.SFTPut(host, info, temp.DockerCom, "/ddhome/harbor/docker-compose.yml")
	utils.SSHExec(host, info, "systemctl enable docker --now && " +
		"/ddhome/harbor/prepare && " +
		"/ddhome/harbor/install.sh")
}


func genHarbor(config *viper.Viper) {
	host := config.Get("server.app.master").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)
	harborDomain := config.Get("registry.local").(map[string]interface{})["domain"].(string)
	log.Print("**** 创建Harbor证书")
	createCA := fmt.Sprintf("cd /tmp/harbor;openssl genrsa -out ca.key 2048;openssl req -x509 -new -nodes -key ca.key -days 36500 -out ca.pem -subj \"/C=CN/ST=ZJ/L=HZ/O=ZJHW/OU=ARCH/CN=ca.%s\"",
		strings.Join([]string{strings.Split(harborDomain, ".")[0], strings.Split(harborDomain, ".")[1]}, "."))
	createHarbor := fmt.Sprintf("cd /tmp/harbor;openssl genrsa -out harbor.key 2048;openssl req -new -key harbor.key -out harbor.csr -subj \"/C=CN/ST=ZJ/L=HZ/O=ZJHW/OU=ARCH/CN=%s\";openssl x509 -req -in harbor.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out harbor.crt -days 36500",
		harborDomain)
	utils.SSHExec(host, info, "mkdir -p /tmp/harbor && "+
		createCA + " && " + createHarbor)

	// 收集证书
	utils.SFTPFetchFile(host, info, "/tmp/harbor/ca.pem", "tmp/harbor/ca.crt")
	utils.SFTPFetchFile(host, info, "/tmp/harbor/harbor.key", "tmp/harbor/harbor.key")
	utils.SFTPFetchFile(host, info, "/tmp/harbor/harbor.crt", "tmp/harbor/harbor.crt")

}