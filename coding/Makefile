
testfile: pb/testfile pb/bin/multicodec
	pb/bin/multicodec header /multicodec >testfile
	cat pb/testfile >>testfile

pb/testfile:
	$(MAKE) -C pb testfile

pb/bin/multicodec:
	$(MAKE) -C pb bin/multicodec

