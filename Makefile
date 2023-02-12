default: test

TOPTARGETS := test tidy

SUBDIRS := sweet sweet/factories/tc

$(TOPTARGETS): $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)


.PHONY: $(TOPTARGETS) $(SUBDIRS)
