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
	kubePkg := config.Get(fmt.Sprintf("package")).(map[string]interface{})["master"].([]interface{})
	nodePkg := config.Get(fmt.Sprintf("package")).(map[string]interface{})["node"].([]interface{})
	dockerPkg := config.Get(fmt.Sprintf("package")).(map[string]interface{})["docker"].([]interface{})
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
	apiserverAttr := new(temp.APIServerAttr)
	apiserverAttr.ETCD = etcdCluster
	apiserverAttr.MASTER = kubeMaster[0].(string)

	apiserverConf := utils.RendTemp(temp.APIServerTem, &apiserverAttr)
	configConf := utils.RendTemp(temp.ConfigTmp, &apiserverAttr)
	masterinfo := GetHostInfo(config, kubeMaster[0].(string))
	utils.SSHExec(kubeMaster[0].(string), masterinfo, "mkdir -p /etc/kubernetes/kubernetesTLS && mkdir -p /ddhome/k8s")
	utils.SFTPut(kubeMaster[0].(string), masterinfo, configConf, "/etc/kubernetes/config")
	utils.SFTPut(kubeMaster[0].(string),  masterinfo, apiserverConf, "/etc/kubernetes/apiserver")
	utils.SFTPut(kubeMaster[0].(string),  masterinfo, temp.Controller, "/etc/kubernetes/controller-manager")
	utils.SFTPut(kubeMaster[0].(string),  masterinfo, temp.Scheduler, "/etc/kubernetes/scheduler")

	// 证书分发
	for _, f := range filepem {
		srcFile := fmt.Sprintf("tmp/kubernetes/%s", f)
		destFile := fmt.Sprintf("/etc/kubernetes/kubernetesTLS/%s", f)
		utils.SFTPutFile(kubeMaster[0].(string), masterinfo, srcFile, destFile)
	}

	// 二进制上传
	for _, f := range kubePkg {
		srcFile := fmt.Sprintf("bin/%s", f.(string))
		destFile := fmt.Sprintf("/usr/local/bin/%s", f.(string))
		utils.SFTPutFile(kubeMaster[0].(string), masterinfo, srcFile, destFile)
	}
	utils.SSHExec(kubeMaster[0].(string), masterinfo, "chmod +x /usr/local/bin/kube*")

	// service 文件上传
	utils.SFTPut(kubeMaster[0].(string), masterinfo, temp.Apiservice, "/usr/lib/systemd/system/kube-apiserver.service")
	utils.SFTPut(kubeMaster[0].(string), masterinfo, temp.Schedulerservice, "/usr/lib/systemd/system/kube-scheduler.service")
	utils.SFTPut(kubeMaster[0].(string), masterinfo, temp.Controllerservice, "/usr/lib/systemd/system/kube-controller-manager.service")


	utils.SFTPut(kubeMaster[0].(string), masterinfo, temp.Token, "/etc/kubernetes/token.csv")
	utils.SFTPut(kubeMaster[0].(string), masterinfo, temp.Instruct, "/tmp/tls-instruct.yml")

	// 设置admin用户集群参数
	cmd := fmt.Sprintf("/usr/local/bin/kubectl config set-cluster kubernetes --certificate-authority=/etc/kubernetes/kubernetesTLS/ca.pem --embed-certs=true --server=https://%s:6443", host[0])
	utils.SSHExec(kubeMaster[0].(string), masterinfo, cmd)
	utils.SSHExec(kubeMaster[0].(string), masterinfo, "/usr/local/bin/kubectl config set-credentials admin --client-certificate=/etc/kubernetes/kubernetesTLS/admin.pem --client-key=/etc/kubernetes/kubernetesTLS/admin.key --embed-certs=true")
	utils.SSHExec(kubeMaster[0].(string), masterinfo, "/usr/local/bin/kubectl config set-context kubernetes --cluster=kubernetes --user=admin")
	utils.SSHExec(kubeMaster[0].(string), masterinfo, "/usr/local/bin/kubectl config use-context kubernetes")

	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"systemctl daemon-reload && "+
			"systemctl enable kube-apiserver --now && "+
			"systemctl enable kube-controller-manager --now && "+
			"systemctl enable kube-scheduler --now")

	// 生成证书循环注册及信任node
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl apply -f /tmp/tls-instruct.yml && "+
			"/usr/local/bin/kubectl create clusterrolebinding node-client-auto-approve-csr --clusterrole=system:certificates.k8s.io:certificatesigningrequests:nodeclient --user=kubelet-bootstrap && "+
			"/usr/local/bin/kubectl create clusterrolebinding node-client-auto-renew-crt --clusterrole=system:certificates.k8s.io:certificatesigningrequests:selfnodeclient --group=system:nodes && "+
			"/usr/local/bin/kubectl create clusterrolebinding node-server-auto-renew-crt --clusterrole=system:certificates.k8s.io:certificatesigningrequests:selfnodeserver --group=system:nodes")

	// 生成kube-proxy.kubeconfig
	cmd1 := fmt.Sprintf("/usr/local/bin/kubectl config set-cluster kubernetes --certificate-authority=/etc/kubernetes/kubernetesTLS/ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig", host[0])
	utils.SSHExec(kubeMaster[0].(string), masterinfo, cmd1)
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config set-credentials kube-proxy --client-certificate=/etc/kubernetes/kubernetesTLS/proxy.pem --client-key=/etc/kubernetes/kubernetesTLS/proxy.key --embed-certs=true --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig")
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config set-context default --cluster=kubernetes --user=kube-proxy --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig")
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config use-context default --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig")

	// 生成kubelet.kubeconifg
	cmd2 := fmt.Sprintf("/usr/local/bin/kubectl config set-cluster kubernetes --certificate-authority=/etc/kubernetes/kubernetesTLS/ca.pem --embed-certs=true --server=https://%s:6443 --kubeconfig=/etc/kubernetes/bootstrap.kubeconfig", host[0])
	utils.SSHExec(kubeMaster[0].(string), masterinfo, cmd2)
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config set-credentials kubelet-bootstrap --token=c6c26805e0569291a57f9c99d2551a3a --kubeconfig=/etc/kubernetes/bootstrap.kubeconfig")
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config set-context default --cluster=kubernetes --user=kubelet-bootstrap --kubeconfig=/etc/kubernetes/bootstrap.kubeconfig")
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl config use-context default --kubeconfig=/etc/kubernetes/bootstrap.kubeconfig")
	utils.SSHExec(kubeMaster[0].(string), masterinfo,
		"/usr/local/bin/kubectl create --insecure-skip-tls-verify clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap")

	// 收集kubeconfig
	utils.SFTPFetchFile(kubeMaster[0].(string), masterinfo,
		"/etc/kubernetes/bootstrap.kubeconfig", "tmp/kubernetes/bootstrap.kubeconfig")
	utils.SFTPFetchFile(kubeMaster[0].(string), masterinfo,
		"/etc/kubernetes/kube-proxy.kubeconfig", "tmp/kubernetes/kube-proxy.kubeconfig")

	// 部署node节点
	var nodeWg sync.WaitGroup
	for _, node := range kubeNode {
		nodeWg.Add(1)
		go deployNode(config, node.(string), nodePkg, dockerPkg, configConf, &nodeWg)
	}
	nodeWg.Done()

	// 部署flannel网络
	var flannelWg sync.WaitGroup
	for _, node := range kubeAll {
		flannelWg.Add(1)
		go deployFlannel(config, node.(string), etcdCluster, &flannelWg)
	}
	flannelWg.Done()
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

var filepem = []string{
	"ca.pem",
	"ca.key",
	"apiserver.key",
	"apiserver.pem",
	"admin.key",
	"admin.pem",
	"proxy.key",
	"proxy.pem",
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

func deployNode(config *viper.Viper, node string, nodePkg, dockerPkg []interface{}, configConf string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("**** %s 节点", node)
	// 上传node二进制
	info := GetHostInfo(config, node)
	for _, pkg := range nodePkg {
		utils.SFTPut(node, info, "bin/"+pkg.(string), "/usr/local/bin/"+pkg.(string))
	}

	// 上传docker安装包
	var dockerSlice []string
	for _, docker := range dockerPkg {
		dockerSlice = append(dockerSlice, docker.(string))
		utils.SFTPut(node, info, "bin/"+docker.(string), "/tmp/"+docker.(string))
	}

	// 安装docker
	utils.SSHExec(node, info, "chmod +x /usr/local/bin/kube* && "+
		"mkdir -p /etc/kubernetes/kubernetesTLS && "+
		"cd /tmp/; yum localinstall -y "+strings.Join(dockerSlice, " "))

	utils.SFTPutFile(node, info, "tmp/kubernetes/bootstrap.kubeconfig", "/etc/kubernetes/bootstrap.kubeconfig")
	utils.SFTPutFile(node, info, "tmp/kubernetes/kube-proxy.kubeconfig", "/etc/kubernetes/kube-proxy.kubeconfig")

	// 分发证书
	for _, f := range filepem {
		srcFile := fmt.Sprintf("tmp/kubernetes/%s", f)
		destFile := fmt.Sprintf("/etc/kubernetes/kubernetesTLS/%s", f)
		utils.SFTPutFile(node, info, srcFile, destFile)
	}

	// 生成配置文件
	nodeAttr := new(temp.KubeletAttr)
	nodeAttr.NODE = node
	nodeAttr.RESTRY = config.Get("registry.local").(map[string]interface{})["domain"].(string)

	kubeletconf := utils.RendTemp(temp.KubeletTmp, nodeAttr)
	utils.SFTPut(node, info, kubeletconf, "/etc/kubernetes/kubelet")
	utils.SFTPut(node, info, configConf, "/etc/kubernetes/config")
	utils.SFTPut(node, info, temp.Proxy, "/etc/kubernetes/proxy")
	utils.SFTPut(node, info, temp.Kubeletservice, "/usr/lib/systemd/system/kubelet.service")
	utils.SFTPut(node, info, temp.Proxyservice, "/usr/lib/systemd/system/kube-proxy.service")

	utils.SSHExec(node, info,
		"swapoff -a && " +
			"yum install -y nfs-utils && "+
			"systemctl enable docker --now && "+
			"systemctl enable kubelet --now && "+
			"systemctl enable kube-proxy --now")

	cert := fmt.Sprintf("/etc/docker/certs.d/%s", viper.Get("registry.local").(map[string]interface{})["domain"].(string))
	utils.SSHExec(node, info, "mkdir -p "+cert)
	utils.SFTPutFile(node, info, "tmp/harbor/ca.crt", cert+"/ca.crt")
}

func deployFlannel(config *viper.Viper, node string, etcdCluster string, wg *sync.WaitGroup) {
	defer wg.Done()
	flannelAttr := new(temp.FlannelAttr)
	flannelAttr.ETCD = etcdCluster

	flannelConf := utils.RendTemp(temp.FlannelConf, &flannelAttr)
	info := GetHostInfo(config, node)
	utils.SSHExec(node, info, "yum install -y flannel && mkdir -p /ddhome/local/docker/data")
	utils.SFTPut(node, info, temp.FlannelJSON, "/root/flannel-config.json")
	utils.SFTPut(node, info, flannelConf, "/etc/sysconfig/flanneld")
	utils.SFTPut(node, info, temp.Dockerservice, "/usr/lib/systemd/system/docker.service")
	utils.SSHExec(node, info, "etcdctl --ca-file=/etc/etcd/etcdSSL/ca.pem  --cert-file=/etc/etcd/etcdSSL/etcd.pem   --key-file=/etc/etcd/etcdSSL/etcd-key.pem set /k8s/network/config < /root/flannel-config.json && "+
		"systemctl enable flanneld --now && "+
		"iptables -P FORWARD ACCEPT")

	// 准备etcd证书
	utils.SSHExec(node, info, "mkdir -p /etc/etcd/etcdSSL")
	utils.SFTPutFile(node, info, "tmp/etcd/etcd-key.pem", "/etc/etcd/etcdSSL/etcd-key.pem")
	utils.SFTPut(node, info, "tmp/etcd/ca.pem", "/etc/etcd/etcdSSL/ca.pem")
	utils.SFTPut(node, info, "tmp/etcd/etcd.pem", "/etc/etcd/etcdSSL/etcd.pem")
	utils.SSHExec(node, info, "if `id etcd &> /dev/null`;then echo '';else useradd etcd;fi && chown -R etcd /etc/etcd/etcdSSL")

	// 启动flannel
	utils.SSHExec(node, info, "systemctl enable flanneld --now && "+
		"systemctl daemon-reload && "+
		"systemctl restart docker && "+
		"iptables -P FORWARD ACCEPT")
}
