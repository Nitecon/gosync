Vagrant.configure("2") do |config|
    config.vm.box = "centos"
    config.vm.box_url = "https://github.com/2creatives/vagrant-centos/releases/download/v6.5.3/centos65-x86_64-20140116.box"
    config.vm.synced_folder "./", "/data/go/src/github.com/Nitecon/gosync", create:true

    config.vm.define "node1" do |node1|
        node1.vm.network :private_network, ip: "10.0.1.100"
        node1.vm.hostname = "node1.example.com"
        node1.vm.provision :shell, :path => ".vagrant_setup/node_setup.sh"
        node1.vm.box = "centos"
    end

    config.vm.define "node2" do |node2|
        node2.vm.network :private_network, ip: "10.0.1.101"
        node2.vm.provision :shell, :path => ".vagrant_setup/node_setup.sh"
        node2.vm.hostname = "node2.example.com"
        node2.vm.box = "centos"
    end

    config.vm.define "db" do |db|
        db.vm.network :private_network, ip: "10.0.1.105"
        db.vm.hostname = "db.example.com"
        db.vm.provision :shell, :path => ".vagrant_setup/db_setup.sh"
        db.vm.box = "centos"
    end
end
