package temp

var KubeletTmp = `
# kubelet (minion) config
#
## The address for the info server to serve on (set to 0.0.0.0 or "" for all interfaces)
KUBELET_ADDRESS="--address=0.0.0.0"
#
## The port for the info server to serve on
KUBELET_PORT="--port=10250"
#
## You may leave this blank to use the actual hostname
KUBELET_HOSTNAME="--hostname-override={{ .NODE }}"
#
## location of the api-server
KUBELET_CONFIG="--kubeconfig=/etc/kubernetes/kubelet.kubeconfig"
#
## pod infrastructure container
KUBELET_POD_INFRA_CONTAINER="--pod-infra-container-image={{ .RESTRY }}/library/pause-amd64:3.1"
#
## Add your own!
KUBELET_ARGS="--cluster-dns=10.0.6.200  --serialize-image-pulls=false  --bootstrap-kubeconfig=/etc/kubernetes/bootstrap.kubeconfig  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig  --cert-dir=/etc/kubernetes/kubernetesTLS  --cluster-domain=cluster.local.  --hairpin-mode=promiscuous-bridge"
`

type KubeletAttr struct {
	NODE string
	RESTRY string
}
