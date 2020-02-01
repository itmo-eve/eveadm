eveadm: *.go cmd/*.go
	go build github.com/itmo-eve/eveadm/

install: eveadm 
	go install github.com/itmo-eve/eveadm/

test: test_help test_func test-rkt test_xen

test_help:
	go test test_utils.go help_test.go

test_func:
	LANG=C go test test_utils.go func_test.go

test_rkt: eveadm
	go test test_utils.go rkt_test.go

test_xen: eveadm
	go test test_utils.go xen_test.go
