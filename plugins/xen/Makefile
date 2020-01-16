all: ../xen.so

../%.so: %.go
	go build -buildmode=plugin -o $@ $<
