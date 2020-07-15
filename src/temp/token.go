package temp

var Token = `
c6c26805e0569291a57f9c99d2551a3a,kubelet-bootstrap,10001,"system:kubelet-bootstrap"
`

var Instruct = `
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: system:certificates.k8s.io:certificatesigningrequests:selfnodeserver
rules:
  - apiGroups: ["certificates.k8s.io"]
    resources: ["certificatesigningrequests/selfnodeserver"]
    verbs: ["create"]
`
