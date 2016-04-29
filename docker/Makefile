build:
	docker build -t factorio-server-manager .

run:
	docker run -d --name factorio-server-manager -v [yourpath]:/security -p80:80 -p443:443 -p34197:34197/udp -e "ADMIN_PASSWORD=jqkSnQS4rA" factorio-server-manager
