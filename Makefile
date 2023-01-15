default: test

test: sweet
	$(MAKE) -C $<

.PHONY: test
