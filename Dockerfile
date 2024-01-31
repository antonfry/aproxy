FROM debian:bullseye-slim
WORKDIR /
COPY cmd/aproxy/aproxy-linux /aproxy
COPY conf/aproxy.yml /aproxy.yml
EXPOSE 8080
CMD ["/aproxy", "-conf", "/aproxy.yml"]