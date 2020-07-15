package temp

var ConfigTmp = `
###
# kubernetes system config
#
# The following values are used to configure various aspects of all
# kubernetes services, including
#
#   kube-apiserver.service
#   kube-controller-manager.service
#   kube-scheduler.service
#   kubelet.service
#   kube-proxy.service
# logging to stderr means we get it in the systemd journal
# 表示错误日志记录到文件还是输出到stderr。
KUBE_LOGTOSTDERR="--logtostderr=true"

# journal message level, 0 is debug
# 日志等级。设置0则是debug等级
KUBE_LOG_LEVEL="--v=0"

# Should this cluster be allowed to run privileged docker containers
# 允许运行特权容器。
KUBE_ALLOW_PRIV="--allow-privileged=true"

# How the controller-manager, scheduler, and proxy find the apiserver
# 设置master服务器的访问
KUBE_MASTER="--master=http://{{- .MASTER -}}:8080"
`
