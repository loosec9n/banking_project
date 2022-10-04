package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves a new HTTP request
type Server struct {
	config     utils.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(config utils.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot make token %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//calling setup router
	server.setupRouter()

	return server, nil
}

// setupRouter setup the router for the project
func (server *Server) setupRouter() {

	router := gin.Default()

	//Users
	router.POST("/user", server.createUser)
	router.POST("/user/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))
	//Accounts
	authRoute.POST("/account", server.createAccount)
	authRoute.GET("/account/:id", server.getAccount)
	authRoute.GET("/accounts", server.listAccounts)

	//Money transfer
	authRoute.POST("/transfers", server.createTrasfer)

	server.router = router

}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse will be used by all api for reporting errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
