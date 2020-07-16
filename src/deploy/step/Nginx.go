package step

import (
	"deploy/temp"
	"deploy/utils"
	"github.com/spf13/viper"
)

func DeployNginx(config *viper.Viper) {
	host := config.Get("server.nginx").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)
	utils.SSHExec(host, info, "yum install -y epel-release nginx && "+
		"mkdir -p /etc/nginx/vhosts && "+
		"mkdir -p /var/log/nginx && "+
		"chown -R nginx /var/log/nginx")

	utils.SFTPut(host, info, temp.NginxConf, "/etc/nginx/nginx.conf")
	utils.SFTPut(host, info, temp.ArchConf, "/etc/nginx/vhosts/arch.conf")

	utils.SSHExec(host, info, "mkdir -p /ddhome/project/bigdata/web/{user,admin} && "+
		"systemctl enable nginx --now")

	// 下载包
	userPkg := config.Get("package").(map[string]interface{})["userpkg"].(string)
	adminPkg := config.Get("package").(map[string]interface{})["adminpkg"].(string)
	userob := utils.NewDownload(config, userPkg)
	userob.Down()
	adminob := utils.NewDownload(config, adminPkg)
	adminob.Down()

	// 将前端包发送到服务器
	utils.SFTPutFile(host, info, "file/"+userPkg, "/ddhome/project/bigdata/"+userPkg)
	utils.SFTPutFile(host, info, "file"+adminPkg, "/ddhome/project/bigdata/"+adminPkg)

	utils.SSHExec(host, info, "yum install -y unzip && rm -rf /ddhome/project/bigdata/web/user/* &&"+
		" unzip /ddhome/project/bigdata/"+userPkg+" -d /ddhome/project/bigdata/web/user")
	utils.SSHExec(host, info, "rm -rf /ddhome/project/bigdata/web/admin/* &&"+
		" unzip /ddhome/project/bigdata/"+adminPkg+" -d /ddhome/project/bigdata/web/admin")
}
