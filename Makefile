default: test

TOPTARGETS := test tidy

SUBDIRS := sweet sweet/factories/testcontainers/redis

$(TOPTARGETS): $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)


.PHONY: $(TOPTARGETS) $(SUBDIRS)
