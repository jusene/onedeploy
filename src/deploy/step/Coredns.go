package step

import (
	"github.com/spf13/viper"
	"deploy/temp"
	"deploy/utils"
)

func DeployCoreDNS(config *viper.Viper) {
	host := config.Get("server.app.master").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := temp.GetHostInfo(config, host)

	corednsAtrr := new(temp.CoreDNSAttr)
	corednsAtrr.RESTRY = config.Get("registry.local").(map[string]interface{})["domain"].(string)

	corednsConf := utils.RendTemp(temp.CoreDNS, &corednsAtrr)
	utils.SFTPut(host, info, corednsConf, "/ddhome/k8s/coredns.yml")
	utils.SSHExec(host, info, "/usr/local/bin/kubectl apply -f /ddhome/k8s/coredns.yml")
}
