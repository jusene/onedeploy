package step

import (
	"github.com/spf13/viper"
	"deploy/temp"
	"deploy/utils"
)

func DeployGit(config *viper.Viper) {
	host := config.Get("server.git").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)

	utils.SSHExec(host, info, "yum install -y git httpd && " +
		"mkdir -p /ddhome/local/gitdata && " +
		"htpasswd -b -c -m 	/etc/httpd/conf/.httpd root DI_git12#$ && " +
		"git init --bare /ddhome/local/gitdata/configrepo.git && " +
		"chown -R apache.apache /ddhome/local/gitdata/ && " +
		"cd /ddhome/local/gitdata/configrepo.git;git config http.receivepack true")
	utils.SFTPut(host, info, temp.GitConf, "/etc/httpd/conf.d/git.conf")
	utils.SSHExec(host, info, "systemctl enable httpd --now")
}
