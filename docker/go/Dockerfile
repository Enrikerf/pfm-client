FROM debian:buster AS base

RUN apt-get update && apt-get install -y curl

FROM base AS dev

ARG GO_VERSION
WORKDIR /tmp

RUN apt-get update && apt-get install -y build-essential vim fish git sudo \
    && groupadd -g 1000 docker-user \
    && useradd -d /home/docker-user -s /bin/bash -u 1000 -g 1000 docker-user \
    && usermod -aG sudo docker-user && echo "docker-user:1234" | sudo chpasswd \
    && mkdir /home/docker-user \
    && chown -R docker-user:docker-user /home/docker-user

USER docker-user
WORKDIR /home/docker-user

RUN URL=https://storage.googleapis.com/golang/go1.${GO_VERSION}.linux-amd64.tar.gz && \
    curl ${URL} -o go.tar.gz &&  \
    ls -lrta && tar -zxvf go.tar.gz && \
    rm -rf go.tar.gz

ENV GOPATH /home/docker-user/go
ENV PATH $PATH:/home/docker-user/go/bin
ENV CGO_ENABLED 0

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN ls -s /home/docker-user/go/bin/dlv /usr/local/bin

WORKDIR /app

FROM base AS compile

ARG GO_VERSION
WORKDIR /usr/local
RUN URL=https://storage.googleapis.com/golang/go1.${GO_VERSION}.linux-amd64.tar.gz \
        && curl ${URL} -o go.tar.gz  \
        && tar -zxf go.tar.gz  \
        && rm -rf go.tar.gz

ENV GOPATH /usr/local/go
ENV PATH $PATH:/usr/local/go/bin
# If you enable this, then gcc is needed to debug your app
ENV CGO_ENABLED 0

WORKDIR /app
COPY . .
RUN go build -ldflags "-s -w" -o final.sh

FROM debian:buster AS prod

WORKDIR /
COPY --from=compile /app/.env  /
COPY --from=compile /app/final.sh /
CMD ["./final.sh"]