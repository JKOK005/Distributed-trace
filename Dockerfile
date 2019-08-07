FROM golang:latest

RUN apt-get update && apt-get -y install git unzip build-essential autoconf libtool

RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

RUN git clone https://github.com/edenhill/librdkafka.git && \
    cd librdkafka && \
    ./configure --prefix /usr && \
    make && \
    make install

RUN go get google.golang.org/grpc
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/samuel/go-zookeeper/zk
RUN go get -u gopkg.in/confluentinc/confluent-kafka-go.v1/kafka

ENV MAIN_APP /go/src
WORKDIR ${MAIN_APP}
COPY . Distributed-trace
RUN mkdir -p Distributed-trace/api/proto/v1/
RUN protoc -I Distributed-trace/api/proto/v1/ --go_out=plugins=grpc:Distributed-trace/pkg/api/proto/ Distributed-trace/api/proto/v1/messages.proto
RUN GOBIN=/${MAIN_APP}/Distributed-trace/bin go install Distributed-trace/main.go

EXPOSE 4000
CMD ["Distributed-trace/bin/main"]