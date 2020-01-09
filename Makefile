default:
	./make.sh
	go install

test:
	bash make.sh
	go build ./
	./gosh