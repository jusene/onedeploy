package step

import (
	"deploy/temp"
	"deploy/utils"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
)

func DeployK8S(config *viper.Viper, tag string) {
	var kubeAll []interface{}
	etcd := config.Get(fmt.Sprintf("server.%s.etcd", tag)).(map[string]interface{})["ip"].([]interface{})
	kubeMaster := config.Get(fmt.Sprintf("server.%s.master", tag)).(map[string]interface{})["ip"].([]interface{})
	kubeNode := config.Get(fmt.Sprintf("server.%s.node", tag)).(map[string]interface{})["ip"].([]interface{})
	kubeAll = append(kubeAll, kubeMaster...)
	kubeAll = append(kubeAll, kubeNode...)

	// etcd
	genEtcd(config, etcd)
    genKubernetes(config, kubeMaster)
	var etcdWg sync.WaitGroup
	for index, host := range etcd {
		etcdWg.Add(1)
		go deployEtcd(config, index, etcd, host.(string), &etcdWg)
	}
	etcdWg.Wait()

	// master
	var etcdSlice []string
	for _, host := range etcd {
		etcdUrl := fmt.Sprintf("https://%s:2379", host)
		etcdSlice = append(etcdSlice, etcdUrl)
	}
	etcdCluster := strings.Join(etcdSlice, ",")

}

func genEtcd(config *viper.Viper, hosts []interface{}) {
	host := config.Get("server.app.master").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := GetHostInfo(config, host)

	log.Println("**** 生成etcd证书")
	etcdAttr := new(temp.EtcdAttr)
	etcdAttr.ETCDLIST = hosts

	etcdSSL := utils.RendTemp(temp.ETCDSSL, etcdAttr)
	utils.SSHExec(host, info, "mkdir -p /tmp/etcd")
	utils.SFTPut(host, info, etcdSSL, "/tmp/etcd/openssl.cnf")
	createCA := fmt.Sprintf("cd /tmp/etcd;openssl genrsa -out ca.key 2048;openssl req -x509 -new -nodes -key ca.key -days 36500 -out ca.pem -subj \"/C=CN/ST=ZJ/L=HZ/O=ZJHW/OU=ARCH/CN=etcd-ca\"")
	createEtcd := fmt.Sprintf("cd /tmp/etcd;openssl genrsa -out etcd-key.pem;openssl req -new -key etcd-key.pem -subj \"/CN=etcd-client/O=system:masters\" -out etcd.csr;openssl x509 -in etcd.csr -req -CA ca.pem -CAkey ca.key -CAcreateserial -extensions v3_req_etcd -extfile openssl.cnf -out etcd.pem -days 36500")
	utils.SSHExec(host, info, createCA + " && " + createEtcd)

	// 收集证书
	utils.SFTPFetchFile(host, info, "/tmp/etcd/ca.pem", "tmp/etcd/ca.pem")
	utils.SFTPFetchFile(host, info, "/tmp/etcd/etcd-key.pem", "tmp/etcd/etcd-key.pem")
	utils.SFTPFetchFile(host, info, "/tmp/etcd/etcd.pem", "tmp/etcd/etcd.pem")
}

func genKubernetes(config *viper.Viper, master []interface{}) {
	host := config.Get("server.app.master").(map[string]interface{})["ip"].([]interface{})[0].(string)
	info := GetHostInfo(config, host)

	log.Print("**** 生成kubernetes证书")
	utils.SSHExec(host, info, "mkdir -p /tmp/kubernetes && " +
		"cd /tmp/kubernetes;openssl genrsa -out ca.key 2048 && openssl req -x509 -new -nodes -key ca.key -days 36500 -out ca.pem -subj \"/CN=kubernetes/O=k8s\"")
	sslAttr := new(temp.KubeS)
	sslAttr.NODE = master[0].(string)

	kubeSSL := utils.RendTemp(temp.KubeSSL, sslAttr)
	utils.SFTPut(host, info, kubeSSL, "/tmp/kubernetes/openssl.cnf")
	utils.SSHExec(host, info, "cd /tmp/kubernetes;openssl genrsa -out apiserver.key 2048;openssl req -new -key apiserver.key -out apiserver.csr -subj \"/CN=kubernetes/O=k8s\" -config openssl.cnf;openssl x509 -req -in apiserver.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out apiserver.pem -days 36500 -extensions v3_req -extfile openssl.cnf")
	utils.SSHExec(host, info, "cd /tmp/kubernetes;openssl genrsa -out admin.key 2048;openssl req -new -key admin.key -out admin.csr -subj \"/CN=admin/O=system:masters/OU=System\";openssl x509 -req -in admin.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out admin.pem -days 36500")
	utils.SSHExec(host, info, "cd /tmp/kubernetes;openssl genrsa -out proxy.key 2048;openssl req -new -key proxy.key -out proxy.csr -subj \"/CN=system:kube-proxy\";openssl x509 -req -in proxy.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out proxy.pem -days 36500")

	filepem := []string{
		"ca.pem",
		"ca.key",
		"apiserver.key",
		"apiserver.pem",
		"admin.key",
		"admin.pem",
		"proxy.key",
		"proxy.pem",
	}

	for _, f := range filepem {
		srcFile := fmt.Sprintf("/tmp/kubernetes/%s", f)
		destFile := fmt.Sprintf("tmp/kubernetes/%s", f)
		utils.SFTPFetchFile(host, info, srcFile, destFile)
	}
}

func deployEtcd(config *viper.Viper, idx int, hosts []interface{}, host string, wg *sync.WaitGroup) {
	defer wg.Done()
	info := GetHostInfo(config, host)
	utils.SSHExec(host, info, "yum install -y etcd && " +
		"mkdir -p /etc/etcd/etcdSSL && "+
		"mkdir -p /ddhome/etcd/data && "+
		"chown -R etcd /ddhome/etcd/data")

	var etcdList []string
	for idx, host := range hosts {
		etcdUrl := fmt.Sprintf("etcd%d=https://%s:2038", idx, host.(string))
		etcdList = append(etcdList, etcdUrl)
	}

	etcdCluster := strings.Join(etcdList, ",")
	etcdAttr := new(temp.ETCDCONF)
	etcdAttr.NODE = host
	etcdAttr.INDEX = idx
	etcdAttr.ETCD = etcdCluster

	etcdConf := utils.RendTemp(temp.ETCDTemp, etcdAttr)
	utils.SFTPut(host, info, etcdConf, "/etc/etcd/etcd.conf")
	utils.SFTPut(host, info, temp.ETCDService, "/usr/lib/systemd/system/etcd.service")
	utils.SFTPutFile(host, info, "tmp/etcd/etcd-key.pem", "/etc/etcd/etcdSSL/etcd-key.pem")
	utils.SFTPutFile(host, info, "tmp/etcd/ca.pem", "/etc/etcd/etcdSSL/ca.pem")
	utils.SFTPutFile(host, info, "tmp/etcd/etcd.pem", "/etc/etcd/etcdSSL/etcd.pem")
	utils.SSHExec(host, info, "chown -R etcd /etc/etcd/etcdSSL && " +
		"systemctl daemon-reload && " +
		"systemctl enable etcd --now")
}