VERSION=0.0.2

build:
	glu build linux,darwin

release:
	glu release

install: build
	cp build/$(shell uname)/atlas /usr/local/bin/
	
.PHONY: build release install
