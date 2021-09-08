package server

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/examples/with-gin/config"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/models"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func Init() {
	config := config.GetConfig()

	router := gin.New()

	router.Use(gin.Recovery())

	// CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     append([]string{"content-type"}, supertokens.GetAllCORSHeaders()...),
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}))

	// Adding the SuperTokens middleware
	router.Use(func(c *gin.Context) {
		supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	})

	// Adding an API that requires session verification
	router.GET("/sessioninfo", verifySession(nil), sessioninfo)

	// starting the server
	err := router.Run(config.GetString("server.apiPort"))
	if err != nil {
		panic(err.Error())
	}
}

// This is a function that wraps the supertokens verification function
// to work the gin
func verifySession(options *models.VerifySessionOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		session.VerifySession(options, func(rw http.ResponseWriter, r *http.Request) {
			c.Request = c.Request.WithContext(r.Context())
			c.Next()
		})(c.Writer, c.Request)
	}
}

func sessioninfo(c *gin.Context) {
	session := session.GetSessionFromRequest(c.Request)
	if session == nil {
		c.JSON(500, "no session found")
		return
	}
	sessionData, err := session.GetSessionData()
	if err != nil {
		supertokens.ErrorHandler(err, c.Request, c.Writer)
		return
	}
	c.JSON(200, map[string]interface{}{
		"sessionHandle": session.GetHandle(),
		"userId":        session.GetUserID(),
		"jwtPayload":    session.GetJWTPayload(),
		"sessionData":   sessionData,
	})
}