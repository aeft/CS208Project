FROM nginx:latest

RUN rm /etc/nginx/conf.d/default.conf

# COPY load_balancer.conf /etc/nginx/conf.d/

RUN apt-get update && \
    apt install net-tools

# Install consul-template (adjust version as needed)
ENV CONSUL_TEMPLATE_VERSION=0.40.0
RUN apt-get update && \
    apt-get install -y wget unzip && \
    wget https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip && \
    unzip consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip && \
    mv consul-template /usr/local/bin/ && \
    rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip

# Copy template into container
COPY load_balancer.conf.ctmpl /etc/nginx/templates/load_balancer.conf.ctmpl

EXPOSE 80
