FROM golang:latest

RUN apt update \
 && apt install -yq \
    curl \
 && curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

WORKDIR /project

COPY *.* /project/

RUN go get github.com/influxdata/influxdb-client-go/v2 \
 && go get github.com/BurntSushi/toml \
 && rm -rf /project

WORKDIR /project

# CMD go run main.go