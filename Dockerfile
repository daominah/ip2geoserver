FROM daominah/go114

COPY go.mod /go.mod
COPY go.sum /go.sum
RUN cd / && go mod download

ENV APP_DIR=/go/src/app
WORKDIR ${APP_DIR}
COPY . ${APP_DIR}
RUN cd ${APP_DIR}/cmd/echo_ip_httpsvr && go build

CMD ["bash", "-c", "${APP_DIR}/cmd/echo_ip_httpsvr/echo_ip_httpsvr"]
