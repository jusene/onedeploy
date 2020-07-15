package temp

var FlannelConf = `
# Flanneld configuration options

# etcd url location.  Point this to the server where etcd runs
FLANNEL_ETCD_ENDPOINTS="{{ .ETCD }}"

# etcd config key.  This is the configuration key that flannel queries
# For address range assignment
FLANNEL_ETCD_PREFIX="/k8s/network"

# Any additional options that you want to pass
FLANNEL_OPTIONS="--etcd-cafile=/etc/etcd/etcdSSL/ca.pem  --etcd-certfile=/etc/etcd/etcdSSL/etcd.pem  --etcd-keyfile=/etc/etcd/etcdSSL/etcd-key.pem"
{{ end }}
`
