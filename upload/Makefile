# paste-upload
# See LICENSE for copyright and license details.
.POSIX:

PREFIX ?= /usr
GO ?= go
GOFLAGS ?= -buildvcs=false
RM ?= rm -f

all: paste-upload

paste-upload:
	$(GO) build $(GOFLAGS) .

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f upload $(DESTDIR)$(PREFIX)/bin/paste-upload
	chmod 755 $(DESTDIR)$(PREFIX)/bin/paste-upload

uninstall:
	$(RM) $(DESTDIR)$(PREFIX)/bin/paste-upload

clean:
	$(RM) upload

run:
	go run -race .

watch:
	fd -e go -e tmpl | entr -rcs "go run -race ."

lint:
	fd -e go | entr -c golangci-lint run

.PHONY: all paste-upload install uninstall clean run watch
