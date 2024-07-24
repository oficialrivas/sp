package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/oficialrivas/sgi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	// Cargar el archivo .env
	err := godotenv.Load() // Carga las variables de entorno desde el archivo .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ConnectToDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("Database connection successfully established")
	DB = db

	
	// Migrar las tablas
	err = db.AutoMigrate(
		&models.User{}, &models.IIO{}, &models.Persona{}, &models.Vehiculo{},
		&models.Empresa{}, &models.Direccion{}, &models.Pasaporte{}, &models.Visa{},
		&models.Documento{}, &models.Caso{}, &models.Correo{}, &models.Redes{}, &models.TemporaryAccess{}, &models.Nacionalidad{}, &models.Mensaje{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Crear índices de texto completo después de la migración
	createFullTextIndexes(db)
}



func createFullTextIndexes(db *gorm.DB) {
	// Verificar si la tabla persona existe
	if !db.Migrator().HasTable(&models.Persona{}) {
		log.Fatalf("Table 'persona' does not exist")
		return
	}

	// Crear índice de texto completo para la tabla persona
	err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_persona_fulltext
		ON persona
		USING gin (to_tsvector('spanish', nombre || ' ' || apellido || ' ' || cedula || ' ' || correo))
	`).Error
	if err != nil {
		log.Fatalf("Failed to create fulltext index on persona: %v", err)
	}

	log.Println("Fulltext index on persona created successfully")
}

