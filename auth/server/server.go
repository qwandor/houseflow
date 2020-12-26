package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gbaranski/houseflow/auth/types"
	"github.com/gbaranski/houseflow/auth/utils"

	"github.com/gbaranski/houseflow/auth/database"
)

// Server hold root server state
type Server struct {
	db     *database.Database
	router *gin.Engine
}

// NewServer creates server, it won't run till Server.Start
func NewServer(db *database.Database) *Server {
	return &Server{
		db:     db,
		router: gin.Default(),
	}
}

// Start starts server, this function is blocking
func (s *Server) Start() error {
	log.Println("Starting server at port 8080")
	s.router.POST("/login", s.onLogin)
	s.router.POST("/register", s.onRegister)
	s.router.POST("/someAction", s.onSomeAction)
	return s.router.Run(":8080")
}

func (s *Server) onLogin(c *gin.Context) {
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	dbUser, err := s.db.Mongo.GetUser(user.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, "Invalid email or password")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	//compare the user from the request, with the one we defined:
	if dbUser.Email != user.Email {
		c.JSON(http.StatusUnauthorized, "Invalid email or password")
		return
	}
	passmatch := utils.CheckPasswordHash(user.Password, dbUser.Password)
	if !passmatch {
		c.JSON(http.StatusUnauthorized, "Invalid email or password")
		return
	}

	tokens, err := utils.CreateTokens()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	err = s.db.Redis.CreateAuth(dbUser.ID, tokens)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	token := map[string]string{
		"access_token":  tokens.AccessToken.Token,
		"refresh_token": tokens.RefreshToken.Token,
	}
	c.JSON(http.StatusOK, token)
}

func (s *Server) onRegister(c *gin.Context) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id, err := s.db.Mongo.AddUser(user)
	if err != nil {
		if database.IsDuplicateError(err) {
			c.JSON(http.StatusConflict, "Account already exists")
			return
		}
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"id": id.String(),
	})
}

func (s *Server) onSomeAction(c *gin.Context) {
	strtoken := utils.ExtractToken(c.Request)
	fmt.Println("strtoken:", *strtoken)
	if strtoken == nil {
		c.JSON(http.StatusUnauthorized, "Authorization token not provided")
		return
	}
	_, claims, err := utils.VerifyToken(*strtoken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := s.db.Redis.FetchAuth(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"expires": claims.ExpiresAt,
		"ID":      claims.Id,
		"userID":  userID,
	})
}
