# global-ipam

## on linux system

mkdir -p /opt/cni/bin/ && cd  /opt/cni/bin/

## install plugin

git clone https://github.com/containernetworking/plugins.git
cd plugins
安装。。。。


## install cnitool

git clone https://github.com/containernetworking/cni.git
cd cnitool/
go install

cnitool

# output
cnitool: Add, check, or remove network interfaces from a network namespace
cnitool add   <net> <netns>
cnitool check <net> <netns>
cnitool del   <net> <netns>




git clone https://github.com/yametech/global-ipam.git


cat >/etc/cni/net.d/10-macvlan-global-ipam.conf  << "EOF"
{
    "name": "macvlan-global-ipam",
    "type": "macvlan",
    "master": "ens192",
    "ipam": {
        "name": "global-ipam",
        "type": "global-ipam",
        "etcdConfig": {
        "etcdURL": "http://10.200.100.200:42379"
        },
        "subnet": "10.22.0.0/16",
        "rangeStart": "10.22.0.2",
        "rangeEnd": "10.22.0.254",
        "routes": [{ "dst": "0.0.0.0/0" }]
    }
}
EOF


export CNI_PATH=/opt/cni/bin/

ip netns delete a && ip netns add a && cnitool add macvlan-global-ipam /var/run/netns/a
