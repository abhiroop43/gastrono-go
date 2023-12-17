package main

import (
	_ "gastrono-go/docs"
	"gastrono-go/middleware"
	"gastrono-go/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

//	@title			Gastrono Go Backend API
//	@version		1.0
//	@description	This is the Backend API for Gastrono Go.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Abhiroop Santra
//	@contact.url	https://www.abhiroopsantra.dev/
//	@contact.email	abhiroop.santra@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// log to custom file

	LogFile := "./tmp/gastronogo-" + time.Now().Format(time.DateOnly) + ".log"
	// open log file
	logFile, err := os.OpenFile(LogFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	// Set log out put and enjoy :)
	log.SetOutput(logFile)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.ForwardedByClientIP = true
	if os.Getenv("APP_ENV") == "development" {
		// gin set release mode
		gin.SetMode(gin.DebugMode)
		err = router.SetTrustedProxies([]string{"127.0.0.1"})
	} else {
		gin.SetMode(gin.ReleaseMode)
		err = router.SetTrustedProxies([]string{"0.0.0.0"})
	}

	routes.UserRoutes(router)
	router.Use(middleware.AuthenticationMiddleware())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.TableRoutes(router)
	routes.OrderRoutes(router)
	routes.OrderItemRoutes(router)
	routes.InvoiceRoutes(router)

	log.Println("[INFO] Starting Gastrono-Go Backend API on port " + port)
	err = router.Run(":" + port)
	if err != nil {
		return
	}
}
