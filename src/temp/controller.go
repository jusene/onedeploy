package temp

var Controller = `
###
# The following values are used to configure the kubernetes controller-manager

# defaults from config and apiserver should be adequate

# Add your own!
KUBE_CONTROLLER_MANAGER_ARGS="--experimental-cluster-signing-duration=87600h0m0s --feature-gates=RotateKubeletServerCertificate=true --address=127.0.0.1  --service-cluster-ip-range=10.0.6.0/24  --cluster-name=kubernetes  --cluster-signing-cert-file=/etc/kubernetes/kubernetesTLS/ca.pem  --cluster-signing-key-file=/etc/kubernetes/kubernetesTLS/ca.key  --service-account-private-key-file=/etc/kubernetes/kubernetesTLS/ca.key  --root-ca-file=/etc/kubernetes/kubernetesTLS/ca.pem  --leader-elect=true  --cluster-cidr=10.1.0.0/16"
`
