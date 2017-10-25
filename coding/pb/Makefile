PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO) ipld.pb.go

%.pb.go: %.proto
	go get github.com/gogo/protobuf/protoc-gen-gogo
	go get github.com/gogo/protobuf/protoc-gen-gofast
	protoc --gogo_out=. --proto_path=../../../../../:/usr/include:/usr/local/opt/protobuf/include:. $<
	patch $@ <$@.patch

clean:
	rm -f *.pb.go
	rm -f *.go

testfile: ../bin/multicodec ../bin/msgio
	../bin/multicodec header /mdagv1 >testfile
	../bin/multicodec header /protobuf/msgio >>testfile
	hash=`ipfs add -q -r . | tail -n1` && \
		ipfs object get "$$hash" --enc=protobuf | ../bin/msgio wrap >>testfile

../bin/msgio:
	$(MAKE) -C .. bin/msgio

../bin/multicodec:
	$(MAKE) -C .. bin/multicodec
