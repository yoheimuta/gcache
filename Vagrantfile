# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
config.vm.box = "precise64"
config.vm.box_url = "http://files.vagrantup.com/precise64.box"

config.vm.define "n1" do |n1|
n1.vm.network "private_network", ip: "172.20.20.10"
end

config.vm.define "n2" do |n2|
n2.vm.network "private_network", ip: "172.20.20.11"
end

config.vm.define "n3" do |n3|
n3.vm.network "private_network", ip: "172.20.20.12"
end
end
