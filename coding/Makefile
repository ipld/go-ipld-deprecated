
bin/multicodec:
	mkdir -p bin
	go get -d github.com/jbenet/go-multicodec/multicodec
	go build -o "$@" github.com/jbenet/go-multicodec/multicodec

bin/json2cbor:
	mkdir -p bin
	go get -d github.com/whyrusleeping/cbor/go/json2cbor
	go build -o "$@" github.com/whyrusleeping/cbor/go/json2cbor

bin/convert:
	mkdir -p bin
	cd bin; go build convert.go

bin/msgio:
	mkdir -p bin
	go get -d github.com/jbenet/go-msgio/msgio
	go build -o "$@" github.com/jbenet/go-msgio/msgio

json.testfile: bin/multicodec Makefile test1.json
	: >$@
	bin/multicodec header /mdagv1 >>$@
	bin/multicodec header /json >>$@
	cat test1.json >>$@

cbor.testfile: bin/multicodec bin/json2cbor Makefile test1.json
	./convert -i $< -o $@.tmp -c '/cbor'
	: >$@
	bin/multicodec header /mdagv1 >>$@
	bin/multicodec header /cbor >>$@
	bin/json2cbor -i test1.json -o - >>$@

protobuf.testfile: bin/multicodec bin/msgio
	bin/multicodec header /mdagv1 >$@
	bin/multicodec header /protobuf/msgio >>$@
	mkdir -p dir
	echo a >dir/a
	echo b >dir/b
	hash=`ipfs add -q -r dir | tail -n1` && \
		ipfs object get "$$hash" --enc=protobuf | bin/msgio wrap >>$@
