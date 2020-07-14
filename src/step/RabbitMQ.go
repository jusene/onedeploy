package step

import (
	"deploy/utils"
	"github.com/spf13/viper"
)

func DeployRabbitMQ(config *viper.Viper) {
	host := config.Get("server.rabbitmq").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := GetHostInfo(config, host)

	utils.SSHExec(host, info, "yum install -y epel-release")
	utils.SSHExec(host, info, "yum install -y rabbitmq-server && " +
		"systemctl enable rabbitmq-server --now")
	utils.SSHExec(host, info, "rabbitmqctl add_user rabbitadmin rabbitadmin && " +
		"rabbitmqctl add_vhost /arch/prod && " +
		"rabbitmqctl set_permissions -p /arch/prod rabbitadmin '.*' '.*' '.*' && " +
		"rabbitmqctl set_user_tags rabbitadmin administrator && " +
		"rabbitmq-plugins enable rabbitmq_management" )
}
