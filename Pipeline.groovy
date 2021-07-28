#!/usr/bin/env groovy
def nodes = [:]

nodesByLabel('master').each {
  nodes[it] = { ->
    node(it) {
      stage("preparation@${it}") {

        sh('whoami')
        sh('lsblk')

        dir('/home/centos/go/src/github.com/kong/kubernetes-ingress-controller') {
            checkout scm
        }
        
        sh('export GOPATH=/home/centos/go 
            && export GOBIN="/home/centos/go/bin" 
            && export GOROOT="/usr/local/go" 
            && export PATH=$PATH:$GOROOT/bin:$GOBIN:$GOPATH:/usr/local/bin/')

        sh('sudo chmod -R 777 /home/centos/go/src/github.com/kong/kubernetes-ingress-controller')
        sh(' creating test cluster ...')
        sh('cd /home/centos/go/src/github.com/kong/kubernetes-ingress-controller/railgun && make test.integration.cluster')

        sh('building docker iamge if not yet.')
        sh('docker build -t 477502 -f Dockerfile.Test .')

        sh('compose docker container.')
        sh('docker run --name execution -t -d -u 997:994 --volume-driver=nfs --network=host --privileged -v /home/centos/go:/home/centos/go -v /var/run/docker.sock:/var/run/docker.sock 477502:latest')

        sh('deploy controller into namesapce kong.')
        sh('''docker exec -i execution /bin/bash -c "cd /home/centos/go/src/github.com/kong/kubernetes-ingress-controller && kubectl apply -f deploy/single-v2/all-in-one-dbless.yaml''')
        
        sh('kick off test cases.')
        sh('''docker exec -i execution /bin/bash -c "cd /home/centos/go/src/github.com/kong/kubernetes-ingress-controller/railgun && 
        GO111MODULE=on TEST_DATABASE_MODE="off" GOFLAGS="-tags=performance_tests" go test -run "TestIngressPerformance" ./test/performance/ -v''')        

        }
    }
  }
}

parallel nodes
