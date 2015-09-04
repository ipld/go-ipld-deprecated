PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO) test.pb

%.pb.go: %.proto
	protoc --gogo_out=. --proto_path=../../../../../../:/usr/local/opt/protobuf/include:. $<

clean:
	rm -f *.pb.go
	rm -f *.go

testfile: bin/multicodec bin/msgio
	bin/multicodec header /mdagv1 >testfile
	bin/multicodec header /protobuf/msgio >>testfile
	hash=`ipfs add -q -r . | tail -n1` && \
		ipfs object get "$$hash" --enc=protobuf | bin/msgio wrap >>testfile

bin/multicodec:
	mkdir -p bin
	go get -d github.com/jbenet/go-multicodec/multicodec
	go build -o "$@" github.com/jbenet/go-multicodec/multicodec

bin/msgio:
	mkdir -p bin
	go get -d github.com/jbenet/go-msgio/msgio
	go build -o "$@" github.com/jbenet/go-msgio/msgio
