# CS208Project

run local
```shell
docker-compose up --build
```

run in CloudLab

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

list info in all nodes
```shell
for node in $(sudo docker node ls --format '{{.Hostname}}'); do   echo "=== Node: $node ===";   sudo docker node ps "$node"; done
```

delete all docker services
```shell
sudo docker service rm $(sudo docker service ls -q)
```

execute a command inside a container
```shell
sudo docker exec <ContainerID> curl api-client:8090/start
```

```shell
sudo docker ps
sudo docker service ls
```