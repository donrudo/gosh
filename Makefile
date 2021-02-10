default:
	./make.sh
	go build
	./gosh

plugin:
	./make.sh

install: 
	./make.sh
	go install

clean:
	rm ~/.gosh/plugins -rf
	rm ${GOPATH}/bin/gosh
	rm ./gosh