default: test

TOPTARGETS := test tidy

SUBDIRS := sweet sweet/factories/testcontainers/redis sweet/factories/databases/redisfactory

$(TOPTARGETS): $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)


.PHONY: $(TOPTARGETS) $(SUBDIRS)
