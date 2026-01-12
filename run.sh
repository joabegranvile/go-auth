docker build --no-cache -t auth-go:v2 . && docker stack deploy -c docker-compose.yml auth-app
