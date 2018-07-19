build:
	go build -ldflags="-s" -o kubeusr github.com/etiennecoutaud/kubeusr

install:
	go install

clean:
	rm -f kubeusr

.PHONY: build install
