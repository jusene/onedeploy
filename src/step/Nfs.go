package step

import (
	"deploy/utils"
	"github.com/spf13/viper"
)

func DeployNFS(config *viper.Viper) {
	host := config.Get("server.nfs").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := GetHostInfo(config, host)

	utils.SSHExec(host, info, "yum install -y nfs-utils && " +
		"mkdir -p /ddhome/arch/resource && " +
		"echo '/ddhome/arch/resource *(sync,rw,no_root_squash)' > /etc/exports && " +
		"systemctl enable nfs --now && " +
		"exportfs -rv")
}
