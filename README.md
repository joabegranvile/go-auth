# auth-go

API de autenticação simples em Go com demonstração de RBAC (Role-Based Access Control).

## Descrição

Projeto de exemplo que implementa autenticação via JWT e um middleware RBAC simples. O objetivo é fornecer uma base mínima para ter um RBAC funcional e extensível.

## Objetivo

Ter um servidor HTTP que:
- Autentique usuários e emita tokens JWT.
- Proteja rotas com middleware de autenticação.
- Aplique verificação de papel/role via middleware RBAC.

## Estrutura relevante do projeto

- `main.go` - inicialização, montagem dos serviços e leitura de segredos/variáveis.
- `internal/auth/jwt.go` - geração e parsing de JWT.
- `internal/auth/middleware.go` - middleware que popula o contexto com as claims do usuário.
- `internal/rbac/middleware.go` - middleware que valida se o usuário tem o `role` requerido.
- `internal/db/postgres.go` - conexão e migração básica do Postgres.
- `docker-compose.yml` - compose para rodar `app` + `db` (usa secret externo `pg_password`).

## Requisitos

- Go (conforme `go.mod`): go 1.25.5
- Docker & Docker Compose (se optar por rodar com containers)

## Quickstart — Execução local (modo dev)

1) Prepare o segredo do Postgres (o binário lê `/run/secrets/pg_password`):

```bash
sudo mkdir -p /run/secrets
echo -n "sua_senha_postgres" | sudo tee /run/secrets/pg_password
sudo chmod 644 /run/secrets/pg_password
```

2) Configure variáveis de ambiente necessárias:

```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_NAME=postgres
# main.go lê a senha via /run/secrets/pg_password (veja passo 1)
```

3) Instale dependências e rode:

```bash
go mod download
go run main.go
```

O servidor escutará em `:8080`.

## Rodando com Docker Compose

O `docker-compose.yml` incluído pressupõe que exista um secret Docker chamado `pg_password` e uma rede `traefik_public` externa. Para testes locais sem Traefik, você pode adaptar o compose:

- Criar secret no Docker Swarm (exemplo):

```bash
echo -n "sua_senha_postgres" | docker secret create pg_password -
```

- Em seguida:

```bash
docker compose up --build
```

Se preferir não usar secrets, edite `docker-compose.yml` para usar `POSTGRES_PASSWORD` diretamente (apenas para desenvolvimento).

## Endpoints principais

- `POST /login` — aceita JSON `{ "username": "joao", "password": "123" }`. Retorna token JWT (no código atual o usuário `joao` com senha `123` gera um token com role `admin`).
- `GET /admin` — endpoint protegido por RBAC, exemplo de uso: `rbac.Middleware("admin", handler)`.
- `GET /` — health check.
- `GET /swagger/` — UI Swagger gerada (documentação).

Exemplo para obter token e acessar rota protegida:

```bash
# Obter token
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"joao","password":"123"}'

# Supondo que a resposta seja o token em plain text
TOKEN=<token_recebido>

# Acessar rota admin
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/admin
```

## Como funciona o RBAC aqui

1. `internal/auth/jwt.go` gera tokens com `username` e `role` nas claims.
2. `internal/auth/middleware.go` valida o header `Authorization: Bearer <token>`, faz `Parse` do JWT e coloca as `claims` no contexto da request sob a chave `user`.
3. `internal/rbac/middleware.go` recupera `claims` do contexto e compara `claims.Role` com o `requiredRole` passado ao middleware; se diferente, retorna `403 Forbidden`.

Fluxo mínimo: login -> receber JWT (role=admin) -> chamar rota protegida com header Authorization -> RBAC valida o papel.

## Pontos a melhorar / próximos passos

- Tornar o `JWT secret` configurável (atualmente hardcoded como `super-secret` em `main.go`).
- Substituir armazenamento de usuários hardcoded por uma tabela `users` (o `internal/db/postgres.go` já cria uma tabela simples).
- Implementar hashing de senhas (bcrypt) e endpoints de CRUD de usuários/roles.
- Implementar políticas mais flexíveis (roles hierárquicos, permissões por recurso).
- Adicionar testes automatizados e pipeline CI.

## Links úteis no projeto

- [main.go](main.go) — bootstrap e leitura de segredos/variáveis.
- [internal/auth/jwt.go](internal/auth/jwt.go) — geração e parsing JWT.
- [internal/auth/middleware.go](internal/auth/middleware.go) — middleware de autenticação.
- [internal/rbac/middleware.go](internal/rbac/middleware.go) — middleware RBAC.
- [internal/db/postgres.go](internal/db/postgres.go) — conexão e migração Postgres.
- [docker-compose.yml](docker-compose.yml) — exemplifica como orquestrar `app` + `db`.

---

Se quiser, eu posso:

- ajustar `main.go` para ler `JWT_SECRET` de environment/secret;
- adicionar endpoints de gerenciamento de usuários (registro/login com DB);
- criar um `docker-compose.override.yml` para facilitar execução local sem secrets externos.

Arquivo criado para referência e documentação do projeto.
# go-auth
