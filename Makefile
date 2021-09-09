build:
	docker build -t demo-service .

run:
	docker run -p 8080:8080 -it -v /home/username/demo-service:/opt/demo-service  demo-service

clean:
	docker rmi -f demo-service
