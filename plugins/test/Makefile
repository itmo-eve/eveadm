all: ../test.so

../%.so: %.go
	go build -buildmode=plugin -o $@ $<
