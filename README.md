# global-ipam

```
kubectl apply -f deploy/
```

## on linux system

```
mkdir -p /opt/cni/bin/ && cd  /opt/cni/bin/

## install plugin
git clone https://github.com/containernetworking/plugins.git
cd plugins
## install plugin

```

## install cnitool

```
git clone https://github.com/containernetworking/cni.git
cd cni/cnitool/
go build -o /usr/local/bin/cnitool cnitool.go

cnitool

# output
cnitool: Add, check, or remove network interfaces from a network namespace
cnitool add   <net> <netns>
cnitool check <net> <netns>
cnitool del   <net> <netns>
```

git clone https://github.com/yametech/global-ipam.git
cd global-ipam 
rm -rf /usr/local/bin/global-ipam && go build -o /usr/local/bin/global-ipam cmd/cni/main.go

```
cat >/etc/cni/net.d/10-macvlan-global-ipam.conf  << "EOF"
{
    "name": "macvlan-global-ipam",
    "type": "macvlan",
    "cniVersion": "0.4.0",
    "master": "eth0",
    "ipam": {
        "name": "global-ipam",
        "type": "global-ipam",
        "subnet": "10.211.55.0/24",
        "rangeStart": "10.211.55.30",
        "rangeEnd": "10.211.55.50",
        "routes": [{ "dst": "0.0.0.0/0" }],
        "gateway": "10.211.55.1"
    }
}
EOF


export CNI_PATH=/opt/cni/bin/
# delete ns a
ip netns delete a

# if not exists create
```

ip netns add a
cnitool add macvlan-global-ipam /var/run/netns/a
cnitool del macvlan-global-ipam /var/run/netns/a

```

# check ns ip addr
```

ip netns exec a ip addr

```

# etcd
export ETCDCTL_API=3
etcdctl --endpoints=10.200.100.200:42379 get /global-ipam-etcd-cni --prefix

```
