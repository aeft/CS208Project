FROM openresty/openresty:alpine-fat

# 安装编译依赖
RUN apk add --no-cache wget unzip build-base

# 下载并编译安装 LuaRocks（版本可根据需要调整）
RUN wget https://luarocks.org/releases/luarocks-3.9.1.tar.gz && \
    tar zxpf luarocks-3.9.1.tar.gz && \
    cd luarocks-3.9.1 && \
    ./configure --with-lua=/usr/local/openresty/luajit --lua-suffix=jit && \
    make && make install && \
    cd .. && rm -rf luarocks-3.9.1*

# 使用编译好的 LuaRocks 安装 lua-resty-http 模块
RUN luarocks install lua-resty-http
RUN luarocks install lua-resty-core
RUN luarocks install lua-resty-balancer

# 确保日志目录存在，并移除默认配置
RUN mkdir -p /var/log/nginx && rm /etc/nginx/conf.d/default.conf

COPY load_balancer.conf /etc/nginx/conf.d/


EXPOSE 80
CMD ["/usr/local/openresty/bin/openresty", "-g", "daemon off;"]
