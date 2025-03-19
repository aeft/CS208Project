(cd ./api/client \
&& docker build -t xor01/cs208-api-client:latest . \
&& docker push xor01/cs208-api-client:latest)

(cd ./api/server \
&& docker build -t xor01/cs208-api-server:latest . \
&& docker push xor01/cs208-api-server:latest)

(cd ./nginx \
&& docker build -t xor01/cs208-nginx:latest . \
&& docker push xor01/cs208-nginx:latest)

echo "Built and Pushed all images to Docker Hub"
