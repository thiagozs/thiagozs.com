build:
	docker build -t thiagozs/thiagozs.com .

run:
	docker run --name thiagozscom -d -p 8080:80 thiagozs/thiagozs.com:latest

clean:
	docker stop thiagozscom
	docker rm thiagozscom