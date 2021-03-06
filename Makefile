SHELL=/bin/bash -o pipefail

.bin/gui:
	go build -o .bin/gui ./gui

.bin/obs-auto-livestream:
	go build -o .bin/obs-auto-livestream ./

.PHONY: build
build: clean .bin/gui .bin/obs-auto-livestream

.PHONY: run
run: build
	.bin/obs-auto-livestream

.PHONY: clean
clean:
	rm -rf .bin
