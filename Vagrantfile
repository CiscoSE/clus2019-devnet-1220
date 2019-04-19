Vagrant.configure(2) do |config|
  # collector (Ubuntu 16.04)
  config.vm.define "collector" do |config|
    config.vm.box = "ubuntu/xenial64"
    config.vm.provision "shell", path: "collector.sh"
    config.vm.provider "virtualbox" do |vb|
      vb.name = "collector"
    config.vm.hostname = "collector"
    end
    config.vm.network :private_network, virtualbox__intnet: "link1", ip: "192.0.2.1"
  end
  # XR (6.4.2)
  config.vm.define "xr" do |config|
    config.vm.box =  "iosxrv-fullk9-x64.snapshot.6.4.2"
    config.vm.provision "file", source: "router1.xr", destination: "/home/vagrant/config"
    config.vm.provision "shell" do |s|
        s.path =  "apply_config.sh"
        s.args = ["/home/vagrant/config"]
    end
    config.vm.provider "virtualbox" do |vb|
      vb.name = "xr"
    end
    config.vm.network :private_network, virtualbox__intnet: "link1", auto_config: false
    config.vm.network :private_network, virtualbox__intnet: "link2", auto_config: false
  end
 end

