# CS208Project

## run local
```shell
docker-compose up --build
```

## run in CloudLab

Step 1: docker swarm init
```shell
# install docker in all nodes
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# in node0
sudo docker swarm init --advertise-addr $(hostname -I | awk '{print $2}')

# execute the output command from node0 on other nodes, it is like:
sudo docker swarm join --token ...

# output the join command again in node0 (run in node0)
sudo docker swarm join-token worker
```

Step 2: add shortname (it is convenient to refer later)
```shell
for ID in $(sudo docker node ls --format '{{.ID}}'); do
  HOST=$(sudo docker node inspect $ID --format '{{.Description.Hostname}}')
  SHORT=$(echo $HOST | cut -d. -f1)
  sudo docker node update --label-add shortname=$SHORT $ID
done
```

Step 3: deploy docker services
```shell
# before executing, copy docker-stack.yml and prometheus/prometheus.yml to node0

sudo docker stack deploy --compose-file docker-stack.yml cs208
```

## Useful Commands

list info in all nodes
```shell
for node in $(sudo docker node ls --format '{{.Hostname}}'); do   echo "=== Node: $node ===";   sudo docker node ps "$node"; done
```

delete all docker services/container/network
```shell
sudo docker service rm $(sudo docker service ls -q)
sudo docker container prune -f
sudo docker network prune -f
```

execute a command inside a container
```shell
sudo docker exec <ContainerID> curl api-client:8090/start
```

```shell
sudo docker ps
sudo docker service ls
```

## Test different load balance algorithms (configs)
1. Modify/Replace load_balancer.conf.ctmpl / trigger consul reload
```shell
# modify local ctmpl

sudo docker cp ./nginx/load_balancer.conf.ctmpl $(sudo docker ps | grep nginx | awk '{print $1}'):/etc/nginx/templates/load_balancer.conf.ctmpl

sudo docker exec -it $(sudo docker ps | grep nginx | awk '{print $1}') consul-template -consul-addr=consul:8500 -template='/etc/nginx/templates/load_balancer.conf.ctmpl:/etc/nginx/conf.d/load_balancer.conf:nginx -s reload || true' -once
```

2. (optional) examine the file in the container
```shell
sudo docker exec -it $(sudo docker ps | grep nginx | awk '{print $1}') cat /etc/nginx/templates/load_balancer.conf.ctmpl

sudo docker exec -it $(sudo docker ps | grep nginx | awk '{print $1}') cat /etc/nginx/conf.d/load_balancer.conf
```

3. start requsts
```shell
sudo docker exec -it $(sudo docker ps | grep nginx | awk '{print $1}') curl api-client:8090/start
```



## Scale service
```shell
sudo apt install python3-pip
pip3 install docker

sudo python scale_service.py cs208_api-server-normal 1 --delay 30
```

## Stress service 
```shell
sudo docker exec -it <ContainerID> stress-ng --cpu 1 --cpu-load 100 --timeout 300
```
