PREFIX=/usr/local
DESTDIR=/usr/local
BINDIR=${PREFIX}/bin

BUILDDIR=build

APPS = image_tool
all: $(APPS)


$(BUILDDIR)/image_tool: $(wildcard apps/image_tool/*.go viewer/*.go)

$(BUILDDIR)/%:
	@mkdir -p $(dir $@)
	go build ${GOFLAGS} -o $@ ./apps/$*

$(APPS): %: $(BUILDDIR)/%

$(APPS): %: $(BUILDDIR)/%

clean:
	rm -fr $(BUILDDIR)

.PHONY: install clean all
.PHONY: $(APPS)

install: $(APPS)
	install -m 755 -d ${DESTDIR}${BINDIR}
	for APP in $^ ; do install -m 755 ${BUILDDIR}/$$APP ${DESTDIR}${BINDIR}/$$APP${EXT} ; done
