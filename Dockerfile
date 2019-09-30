#源镜像
FROM golang:latest
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
WORKDIR $GOPATH/src/serve
COPY . $GOPATH/src/serve
RUN go build .
EXPOSE 9002
ENTRYPOINT  ["./serve"]