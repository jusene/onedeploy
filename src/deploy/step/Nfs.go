package step

import (
	"github.com/spf13/viper"
	"deploy/utils"
	"deploy/temp"
)

func DeployNFS(config *viper.Viper) {
	host := config.Get("server.nfs").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)

	utils.SSHExec(host, info, "yum install -y nfs-utils && " +
		"mkdir -p /ddhome/arch/resource && " +
		"echo '/ddhome/arch/resource *(sync,rw,no_root_squash)' > /etc/exports && " +
		"systemctl enable nfs --now && " +
		"exportfs -rv")
}
