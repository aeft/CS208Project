FROM openresty/openresty:alpine-fat

RUN mkdir -p /var/log/nginx && rm /etc/nginx/conf.d/default.conf

# COPY load_balancer.conf /etc/nginx/conf.d/

RUN apk add --no-cache net-tools wget unzip luarocks build-base git openssl-dev curl

RUN ln -s /usr/local/openresty/luajit/bin/luarocks /usr/local/bin/luarocks

RUN luarocks install lua-resty-http
RUN luarocks install lua-resty-core
RUN luarocks install lua-resty-balancer
RUN luarocks install lua-resty-openssl

RUN luarocks install lua-resty-prometheus

# Install consul-template (adjust version as needed)
ENV CONSUL_TEMPLATE_VERSION=0.40.0
RUN apk add --no-cache wget unzip && \
    wget https://releases.hashicorp.com/consul-template/${CONSUL_TEMPLATE_VERSION}/consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip && \
    unzip consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip && \
    mv consul-template /usr/local/bin/ && \
    rm consul-template_${CONSUL_TEMPLATE_VERSION}_linux_amd64.zip

# Copy template into container
COPY load_balancer.conf.ctmpl /etc/nginx/templates/load_balancer.conf.ctmpl

EXPOSE 80
