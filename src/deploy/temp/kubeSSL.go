package temp

var KubeSSL = `
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = kubernetes
DNS.2 = kubernetes.default
DNS.3 = kubernetes.default.svc
DNS.4 = kubernetes.default.svc.cluster
DNS.5 = kubernetes.default.svc.cluster.local
DNS.6 = k8s_master
IP.1 = 10.0.6.1              # ClusterServiceIP 地址
IP.2 = {{ .NODE }}           # master IP地址
IP.3 = 10.0.6.200            # kubernetes DNS IP地址
`

type KubeS struct {
	NODE string
}