package temp

var Proxy = `
###
# kubernetes proxy config

# defaults from config and proxy should be adequate

# Add your own!
KUBE_PROXY_ARGS="--kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig --cluster-cidr=172.16.0.0/16"
`
