# CS208Project

run servers
```shell
docker-compose up --build
```

start requests
```shell
docker exec -it api-client curl api-client:8090/start
```

deploy docker services
```shell
sudo docker stack deploy --compose-file docker-stack.yml cs208
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