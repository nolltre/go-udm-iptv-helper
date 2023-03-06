# go-udm-iptv-helper

This is a small web server whose only purpose is to restart the udm-iptv service [udm-iptv]

It can serve a web site by HTTP or HTTPS, depending on the options given.

## Building
Prerequisites: `go`

Cross-compile on your PC by doing:
```bash
$ GOOS=linux GOARCH=arm64 go build .
```

## Usage
```bash
$ ./go-udm-iptv-helper
Usage of ./go-udm-iptv-helper:
  -cert string
    	path to certificate
  -key string
    	path to certificate key
  -port int
    	port to listen to (default 12345)
```

[udm-iptv]: https://github.com/fabianishere/udm-iptv
