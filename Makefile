eveadm: *.go cmd/*.go
	go build github.com/itmo-eve/eveadm/

install: eveadm 
	go install github.com/itmo-eve/eveadm/

test:
	go test
