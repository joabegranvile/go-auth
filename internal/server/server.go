package server

import (
	"encoding/json"
	"net/http"

	_ "github.com/joabegranvile/auth-go/docs"
	"github.com/joabegranvile/auth-go/internal/auth"
	"github.com/joabegranvile/auth-go/internal/rbac"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	auth *auth.Service
	rbac *rbac.Service
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(auth *auth.Service, rbac *rbac.Service) *Server {
	return &Server{auth: auth, rbac: rbac}
}

func (s *Server) Routes() {
	http.Handle(
		"/swagger/",
		httpSwagger.WrapHandler,
	)
	http.HandleFunc("/login", s.loginHandler)

	http.HandleFunc(
		"/admin",
		s.rbac.Middleware("admin", s.adminHandler),
	)

	http.HandleFunc("/", s.healthHandler)
}

// Admin godoc
// @Summary Painel admin
// @Description Endpoint restrito a admin
// @Tags admin
// @Success 200 {string} string
// @Failure 401 {string} string
// @Failure 403 {string} string
// @Router /admin [get]
func (s *Server) adminHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bem vindo ao Painel Secreto, admin"))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("auth-go Ok"))
}

// LoginHandler godoc
// @Summary Realizar Login
// @Description Autentica o usuário e retorna o token JWT
// @Tags auth
// @Accept json
// @Produce plain json
// @Param request body LoginRequest true "Credenciais de Login"
// @Success 200 {string} string "Token JWT"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.Username == "joao" && req.Password == "123" {
		token, _ := s.auth.Generate("joao", "admin")
		w.Write([]byte(token))
		return
	}
}
