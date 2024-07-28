package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oficialrivas/sgi/config"
	_ "github.com/oficialrivas/sgi/docs" // Importa tu documentación de Swagger generada aquí
	"github.com/oficialrivas/sgi/routers"
	swaggerFiles "github.com/swaggo/files"    // Archivos estáticos para Swagger
	ginSwagger "github.com/swaggo/gin-swagger" // Gin-swagger para la documentación de la API
)

// @title API 
// @description API SISTEMA DE EXPLORACION POPULAR.
// @version 1.0
// @host 10.51.16.147:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configs.ConnectToDB() // Establece la conexión a la base de datos

	r := gin.Default()

	// Configura CORS para permitir todas las solicitudes
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Configura el endpoint para Swagger UI utilizando gin-swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Servir archivos estáticos desde la carpeta 'static'
	r.Static("/static", "./static")

	// Configura tus rutas aquí
	routes.SetupRouter(r) // Esta función ahora configura las rutas directamente

	// Inicia el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}
	r.Run(":" + port)
}
