
pb/bin/multicodec:
	$(MAKE) -C pb bin/multicodec

json.testfile: pb/bin/multicodec Makefile
	: >$@
	pb/bin/multicodec header /multicodec >>$@
	pb/bin/multicodec header /json >>$@
	echo '{"@codec":"/json","abc":{"mlink":"QmXg9Pp2ytZ14xgmQjYEiHjVjMFXzCVVEcRTWJBmLgR39V"}}' >>$@

