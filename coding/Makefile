
pb/bin/multicodec:
	$(MAKE) -C pb bin/multicodec

bin/convert:
	cd bin; go build convert.go

json.testfile: pb/bin/multicodec Makefile
	: >$@
	pb/bin/multicodec header /multicodec >>$@
	pb/bin/multicodec header /json >>$@
	echo '{"@codec":"/json","abc":{"mlink":"QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V"}}' >>$@

cbor.testfile: json.testfile
	pb/bin/multicodec header /multicodec >$@
	./convert -i $< -o $@.tmp -c '/cbor'
	cat $@.tmp >>$@
	rm -f $@.tmp

