VERSION=0.0.4

build:
	glu build linux,darwin

release: build
	glu release

install: build
	cp build/$(shell uname)/atlas /usr/local/bin/
	
.PHONY: build release install
