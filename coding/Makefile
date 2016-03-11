
bin/multicodec:
	mkdir -p bin
	go get -d github.com/jbenet/go-multicodec/multicodec
	go build -o "$@" github.com/jbenet/go-multicodec/multicodec

bin/convert:
	cd bin; go build convert.go

bin/msgio:
	mkdir -p bin
	go get -d github.com/jbenet/go-msgio/msgio
	go build -o "$@" github.com/jbenet/go-msgio/msgio

json.testfile: bin/multicodec Makefile
	: >$@
	bin/multicodec header /multicodec >>$@
	bin/multicodec header /json >>$@
	echo '{"@codec":"/json","abc":{"mlink":"QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V"}}' >>$@

cbor.testfile: bin/multicodec json.testfile
	bin/multicodec header /multicodec >$@
	./convert -i $< -o $@.tmp -c '/cbor'
	cat $@.tmp >>$@
	rm -f $@.tmp

protobuf.testfile: bin/multicodec bin/msgio
	bin/multicodec header /mdagv1 >$@
	bin/multicodec header /protobuf/msgio >>$@
	mkdir -p dir
	echo a >dir/a
	echo b >dir/b
	hash=`ipfs add -q -r dir | tail -n1` && \
		ipfs object get "$$hash" --enc=protobuf | bin/msgio wrap >>$@
