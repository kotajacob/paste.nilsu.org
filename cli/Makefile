# pcf
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr/local
GO ?= go
GOFLAGS ?= -o pcf
RM ?= rm -f

all: pcf

pcf:
	$(GO) build $(GOFLAGS)

clean:
	$(RM) pcf

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f pcf $(DESTDIR)$(PREFIX)/bin/pcf
	chmod 755 $(DESTDIR)$(PREFIX)/bin/pcf

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/pcf

.DEFAULT_GOAL := all

.PHONY: all pcf clean install uninstall
