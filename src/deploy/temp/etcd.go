package temp

var ETCDTemp = `
#[member]
ETCD_NAME=etcd{{ .INDEX }}
ETCD_DATA_DIR="/ddhome/etcd/data"
ETCD_LISTEN_PEER_URLS="https://{{ .NODE }}:2380"
ETCD_LISTEN_CLIENT_URLS="https://{{ .NODE }}:2379"

#[cluster]
ETCD_INITIAL_ADVERTISE_PEER_URLS="https://{{ .NODE }}:2380"
ETCD_INITIAL_CLUSTER_TOKEN="etcd-cluster"
ETCD_ADVERTISE_CLIENT_URLS="https://{{ .NODE }}:2379"
ETCD_CLUSTER="{{ .ETCD }}"
`

var ETCDService = `
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target

[Service]
Type=notify
WorkingDirectory=/var/lib/etcd/
EnvironmentFile=-/etc/etcd/etcd.conf
User=etcd
# set GOMAXPROCS to number of processors
ExecStart=/usr/bin/etcd --name=${ETCD_NAME} --cert-file=/etc/etcd/etcdSSL/etcd.pem --key-file=/etc/etcd/etcdSSL/etcd-key.pem --peer-cert-file=/etc/etcd/etcdSSL/etcd.pem --peer-key-file=/etc/etcd/etcdSSL/etcd-key.pem --trusted-ca-file=/etc/etcd/etcdSSL/ca.pem --peer-trusted-ca-file=/etc/etcd/etcdSSL/ca.pem --initial-advertise-peer-urls=${ETCD_INITIAL_ADVERTISE_PEER_URLS} --listen-peer-urls=${ETCD_LISTEN_PEER_URLS} --listen-client-urls=${ETCD_LISTEN_CLIENT_URLS},http://127.0.0.1:2379 --advertise-client-urls=${ETCD_ADVERTISE_CLIENT_URLS} --initial-cluster-token=${ETCD_INITIAL_CLUSTER_TOKEN} --initial-cluster=${ETCD_CLUSTER} --initial-cluster-state=new --data-dir=${ETCD_DATA_DIR}
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`

type ETCDCONF struct {
	INDEX int
	NODE string
	ETCD string
}