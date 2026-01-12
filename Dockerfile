
# Etapa 1 : Builder(Compilar o codigo)
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o meu-app .

# ETAPA 2: Runner (Imagem final levissima)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/meu-app .

# Permissao para executar
RUN chmod +x ./meu-app
CMD ["./meu-app"]

