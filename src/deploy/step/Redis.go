package step

import (
	"github.com/spf13/viper"
	"deploy/temp"
	"deploy/utils"
)

func DeployRedis(config *viper.Viper) {
	host := config.Get("server.redis").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)

	utils.SSHExec(host, info, "yum install -y redis")
	utils.SFTPut(host, info, temp.RedisConf, "/etc/redis.conf")
	utils.SSHExec(host, info, "systemctl enable redis --now")
}
