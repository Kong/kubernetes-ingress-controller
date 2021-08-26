# sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
# sudo chmod +x /usr/local/bin/docker-compose
# sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose

# docker-compose up -d
sudo add-apt-repository ppa:deadsnakes/ppa
sudo apt-get update
sudo apt-get install python3.7
sudo apt-get update
sudo apt-get install python3-setuptools
sudo easy_install3 pip
sudo pip3 install awscli==1.18.115
sudo pip3 install botocore==1.17.63


cd dolphin-kic
echo "install enterprise test tool."
sudo make install
echo "install target cluster"
dolphin cluster create using flavor kind cluster_name --verbose
echo "configure test environment"
dolphin cluster envsetup cluster-name --external_cluster=kind --verbose --kubeconfig=/hom/ubuntu/.kube/config
echo "run basic integration tests cases"
dolphin load --ingress_type=all --ingress_number=1 --export_type=csv --kubeconfig=/hom/ubuntu/.kube/config