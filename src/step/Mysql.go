package step

import (
	"deploy/temp"
	"deploy/utils"
	"github.com/spf13/viper"
	"log"
	"os"
	"regexp"
	"strings"
)

func DeployMysql(config *viper.Viper) {
	host := config.Get("server.mysql").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := GetHostInfo(config, host)
	pkg := config.Get("package").(map[string]interface{})["mysql"].(string)
	utils.SFTPutFile(host, info, "bin/"+pkg, "/tmp/"+pkg)
	utils.SSHExec(host, info,"rm -rf /usr/local/mysql*")
	utils.SSHExec(host, info, "tar xf /tmp/"+pkg+" -C /usr/local/")
	utils.SSHExec(host, info, "ln -sf /usr/local/"+strings.Trim(pkg, ".tar.xz")+" /usr/local/mysql")
	utils.SSHExec(host, info, "yum install -y libaio-devel && " +
		"cp /usr/local/mysql/support-files/mysql.server /etc/init.d/mysqld && " +
		"mkdir -p /ddhome/local/mysql/data && " +
		"echo 'export PATH=$PATH:/usr/local/mysql/bin' >> /etc/profile && " +
		"if `id mysql &> /dev/null`;then echo '';else useradd mysql;fi")
	utils.SFTPut(host, info, temp.MySQLCNF, "/etc/my.cnf")
	ret := utils.SSHExec(host, info, "rm -rf /ddhome/local/mysql/data && /usr/local/mysql/bin/mysqld -I --basedir=/usr/local/mysql --datadir=/ddhome/local/mysql/data --user=mysql")
	log.Println(ret)
	utils.SSHExec(host, info,"service mysqld start")
	r := regexp.MustCompile(`root@localhost: (.*)`)
	v := r.FindAll([]byte(ret), 1)
	log.Println(string(v[0]))
	file, _ := os.OpenFile("tmp/mysql", os.O_CREATE | os.O_WRONLY, 0644)
	file.Write(v[0])
}
