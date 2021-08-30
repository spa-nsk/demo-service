build:
	GOOS=linux GOARCH=amd64 go build 
	go build  .
	docker build -t demo-service .

run:
	docker run -p 8080:50051 demo-service

clean:
	rm ./demo-srvice
	docker rmi -f demo-service
