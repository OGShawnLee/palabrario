Vagrant.configure("2") do |config|
  config.vm.box = "debian/bookworm64"
  config.vm.hostname = "palabrio-dev"

  # Desactivar la auto-actualización de Guest Additions (usamos rsync, no se necesitan)
  if Vagrant.has_plugin?("vagrant-vbguest")
    config.vbguest.auto_update = false
  end

  config.vm.synced_folder ".", "/home/vagrant/palabrio",
    type: "rsync",
    rsync__exclude: ["**/build/", "**/node_modules/", ".git/", "local-app/", "local-app-prototype/"]

  config.vm.network "forwarded_port", guest: 8000, host: 8000  # core-service
  config.vm.network "forwarded_port", guest: 8001, host: 8001  # realtime-service
  config.vm.network "forwarded_port", guest: 8002, host: 8002  # image-service
  config.vm.network "forwarded_port", guest: 5432, host: 5432  # postgres

  config.vm.provider "virtualbox" do |vb|
    vb.name = "palabrio-dev"
    vb.memory = "2048"
    vb.cpus = 2
  end

  config.vm.provision "shell", inline: <<-SHELL
    export DEBIAN_FRONTEND=noninteractive
    apt-get update
    apt-get install -y ca-certificates curl gnupg lsb-release

    # Añadir clave GPG oficial de Docker
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # Añadir repositorio de Docker
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

    apt-get update
    # Instalar Docker Engine y el plugin de Compose
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    systemctl enable --now docker
    usermod -aG docker vagrant
  SHELL
end
