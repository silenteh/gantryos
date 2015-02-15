NETWORKING

This modules contains:
- implementation of SDN for supporting VLAN and GRE between containers on different hosts


- How to install openvswitch on ubuntu
apt-get install build-essential libtool

wget http://openvswitch.org/releases/openvswitch-2.3.1.tar.gz
tar zxvf openvswitch-2.3.1.tar.gz && cd openvswitch-2.3.1
./boot.sh
./configure --with-linux=/lib/modules/`uname -r`/build
make -j && sudo make install
sudo make modules_install
sudo modprobe gre
sudo modprobe openvswitch
sudo modprobe libcrc32c


ovsdb-tool create /usr/local/etc/openvswitch/conf.db /usr/local/share/openvswitch/vswitch.ovsschema

ovsdb-server --remote=punix:/usr/local/var/run/openvswitch/db.sock \
--remote=db:Open_vSwitch,Open_vSwitch,manager_options \
--pidfile --detach --log-file

example for remote port
ovsdb-server --remote=ptcp:6633 --remote=db:Open_vSwitch,Open_vSwitch,manager_options --pidfile --detach --log-file

echo "openvswitch " >> /etc/modules
echo "gre" >> /etc/modules
echo "libcrc32c" >> /etc/modules