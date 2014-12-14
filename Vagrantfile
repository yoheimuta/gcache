$script_n1 = <<SCRIPT

echo "Setting hostname..."
sudo hostname 172.20.20.10
echo "Done setting hostname"

echo "Installing latest redis-server..."
sudo apt-get update
sudo apt-get install -y python-software-properties
sudo add-apt-repository -y ppa:chris-lea/redis-server
sudo apt-get update
sudo apt-cache showpkg redis-server
sudo apt-get install -y redis-server
sudo redis-server -v
sudo sed -i -e"s/bind 127.0.0.1/bind 172.20.20.10/" /etc/redis/redis.conf
sudo redis-cli shutdown
sudo redis-server /etc/redis/redis.conf
echo "Done installing latest redis-server"

echo "Installing go..."
wget https://raw.github.com/kwmt/goinstall/master/goinstall.sh
wget https://raw.github.com/kwmt/goinstall/master/gosetting.sh
sudo chmod +x goinstall.sh gosetting.sh
sudo ./goinstall.sh vagrant
/usr/local/go/bin/go version
source ~/.bashrc
echo "Done Installing go"

echo "Setting up project..."
cd /vagrant
export GOPATH="/vagrant/.godeps:/vagrant/"
unset GOBIN
/usr/local/go/bin/go install gcache
echo "Done Setting up project"

SCRIPT

$script_n2 = <<SCRIPT

echo "Setting hostname..."
sudo hostname 172.20.20.11
echo "Done setting hostname"

SCRIPT

$script_n3 = <<SCRIPT

echo "Setting hostname..."
sudo hostname 172.20.20.12
echo "Done setting hostname"

SCRIPT

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
config.vm.box = "precise64"
config.vm.box_url = "http://files.vagrantup.com/precise64.box"

config.vm.define "n1" do |n1|
n1.vm.provision "shell", inline: $script_n1
n1.vm.network "private_network", ip: "172.20.20.10"
end

config.vm.define "n2" do |n2|
n2.vm.provision "shell", inline: $script_n2
n2.vm.network "private_network", ip: "172.20.20.11"
end

config.vm.define "n3" do |n3|
n3.vm.provision "shell", inline: $script_n3
n3.vm.network "private_network", ip: "172.20.20.12"
end
end
