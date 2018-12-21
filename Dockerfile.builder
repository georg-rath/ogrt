FROM golang:1.11.1-stretch

RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential file

# reuse vendorize script to deploy dependencies for client
ADD client/vendorize /tmp/vendor/vendorize
RUN cd /tmp/vendor && ./vendorize 4 /usr/ && rm -rf /tmp/vendor

RUN go get -u github.com/golang/protobuf/protoc-gen-go && go get -u github.com/gogo/protobuf/protoc-gen-gogofaster

RUN chmod 770 /go && chmod 770 /go/bin

RUN mkdir /src
WORKDIR /src

