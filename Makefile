PLUGINS=test rkt xen

.PHONY: plugins $(PLUGINS)

install: eveadm plugins
	go install github.com/itmo-eve/eveadm/

eveadm: *.go cmd/*.go
	go build github.com/itmo-eve/eveadm/

plugins: $(PLUGINS)

$(PLUGINS):
	make -C plugins/$@
