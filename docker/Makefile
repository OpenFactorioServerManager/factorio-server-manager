build:
	docker build -t factorio-server-manager .

run:
	docker run -d -p 8080:8080 -p 34197:34197/udp -v /home/user/.factorio:/factorio factorio-server-manager
