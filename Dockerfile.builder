FROM golang:1.11.1-stretch

RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential file

RUN wget https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protobuf-all-3.6.1.tar.gz && \
  tar xf protobuf-all-3.6.1.tar.gz && \
  cd protobuf-3.6.1 && \
  ./configure --prefix=/usr && \
  make install -j4 && \
  rm -rf protobuf-3.6.1

RUN wget https://github.com/protobuf-c/protobuf-c/releases/download/v1.3.1/protobuf-c-1.3.1.tar.gz && \
  tar xf protobuf-c-1.3.1.tar.gz && \
  cd protobuf-c-1.3.1 && \
  ./configure --prefix=/usr && \
  make install -j4 && \
  rm -rf protobuf-c-1.3.1

RUN go get -u github.com/golang/protobuf/protoc-gen-go && go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

RUN chmod 770 /go && chmod 770 /go/bin

RUN mkdir /src
WORKDIR /src

