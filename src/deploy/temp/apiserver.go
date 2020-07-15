package temp

var APIServerTem = `
###
## kubernetes system config
##
## The following values are used to configure the kube-apiserver
##
#
## The address on the local server to listen to.
KUBE_API_ADDRESS="--advertise-address={{ .MASTER }} --bind-address={{ .MASTER }} --insecure-bind-address={{ .MASTER }}"
#
## The port on the local server to listen on.
KUBE_API_PORT="--port=8080"
#
## Port minions listen on
KUBELET_PORT="--kubelet-port=10250"
#
## Comma separated list of nodes in the etcd cluster
KUBE_ETCD_SERVERS="--etcd-servers={{ .ETCD  }}"
#
## Address range to use for services
KUBE_SERVICE_ADDRESSES="--service-cluster-ip-range=10.0.6.0/24"
#
## default admission control policies
KUBE_ADMISSION_CONTROL="--admission-control=ServiceAccount,NamespaceLifecycle,NamespaceExists,LimitRanger,ResourceQuota,NodeRestriction"
#
## Add your own!
KUBE_API_ARGS="--authorization-mode=Node,RBAC  --runtime-config=rbac.authorization.k8s.io/v1beta1  --kubelet-https=true  --token-auth-file=/etc/kubernetes/token.csv  --service-node-port-range=30000-32767  --tls-cert-file=/etc/kubernetes/kubernetesTLS/apiserver.pem  --tls-private-key-file=/etc/kubernetes/kubernetesTLS/apiserver.key  --client-ca-file=/etc/kubernetes/kubernetesTLS/ca.pem  --service-account-key-file=/etc/kubernetes/kubernetesTLS/ca.key  --storage-backend=etcd3  --etcd-cafile=/etc/etcd/etcdSSL/ca.pem  --etcd-certfile=/etc/etcd/etcdSSL/etcd.pem  --etcd-keyfile=/etc/etcd/etcdSSL/etcd-key.pem  --enable-swagger-ui=true  --apiserver-count=3  --audit-log-maxage=30  --audit-log-maxbackup=3  --audit-log-maxsize=100  --audit-log-path=/var/lib/audit.log  --event-ttl=1h"
`

type APIServerAttr struct {
	ETCD string
	MASTER string
}