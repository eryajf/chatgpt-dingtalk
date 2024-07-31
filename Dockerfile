FROM golang:1.22.5-alpine3.20 AS builder

# ENV GOPROXY      https://goproxy.io

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o chatgpt-dingtalk .

FROM alpine:3.20

ARG TZ="Asia/Shanghai"

ENV TZ ${TZ}

RUN mkdir /app && apk upgrade \
    && apk add bash tzdata \
    && ln -sf /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone

WORKDIR /app
COPY --from=builder /app/ .
RUN chmod +x chatgpt-dingtalk && cp config.example.yml config.yml

CMD ./chatgpt-dingtalk