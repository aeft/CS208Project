FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY load_balancer.conf /etc/nginx/conf.d/

EXPOSE 80
