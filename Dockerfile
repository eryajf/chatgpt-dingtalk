FROM golang:1.17.10 AS builder

# ENV GOPROXY      https://goproxy.io

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o chatgpt-dingtalk .

FROM centos:centos7
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/ .
RUN chmod +x chatgpt-dingtalk && cp config.dev.json config.json && yum -y install vim net-tools telnet wget curl && yum clean all

CMD ./chatgpt-dingtalk