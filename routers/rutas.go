package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/oficialrivas/sgi/controllers"
	"github.com/oficialrivas/sgi/middleware"
)

func SetupRouter(r *gin.Engine) {
	// Rutas públicas (sin protección JWT)
	r.POST("/signup", controllers.CreateUser)
	r.POST("/login", controllers.Login)
	r.POST("/webhook", controllers.WebhookHandler)
	r.POST("/gestion", controllers.GetRecordsByAreaAndPeriod)
	r.POST("/gestion/por-area", controllers.GetRecordsCountByAreaAndPeriod)
	r.POST("/websms", controllers.WebsmsHandler)  
	r.POST("/gestion/area-modalidad", controllers.GetRecordsCountByAreaAndModalidad)
	r.POST("/gestion/user", controllers.GetRecordsByUserAndPeriod)
    r.POST("/gestion/user-area-modalidad", controllers.GetRecordsCountByUserAndModalidad)
	r.POST("/generate-token", controllers.GenerateToken) 
			
	// Endpoints protegidos con JWT
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired()) // Middleware de autenticación JWT
	{
		// CRUD para User
		protected.GET("/users/:id", middleware.RoleRequired("admin", "superuser"), controllers.GetUser)
		protected.GET("/users", middleware.RoleRequired("admin"), controllers.GetUsers)
		protected.PUT("/users/:id", middleware.RoleRequired("admin"), controllers.UpdateUser)
		protected.DELETE("/users/:id", middleware.RoleRequired("admin"), controllers.DeleteUser)
		protected.GET("/users/:id/otp-setup", middleware.RoleRequired("admin"), controllers.SetupOTP)
		protected.PUT("/users/:id/password", middleware.RoleRequired("admin"), controllers.UpdatePassword)
		protected.GET("/users/:id/messages", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.GetMensajesByUserID)
		

		// CRUD para Caso
		protected.POST("/casos", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateCaso)
		protected.GET("/casos/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetCasoByID)
		protected.PUT("/casos/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateCaso)
		protected.DELETE("/casos/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteCaso)
		protected.PUT("/casos/valorar/:id", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.ValorarCaso)

		// CRUD para Documento
		protected.POST("/documentos", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateDocumento)
		protected.GET("/documentos/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetDocumentoByID)
		protected.PUT("/documentos/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateDocumento)
		protected.DELETE("/documentos/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteDocumento)

		// CRUD para Pasaporte
		protected.POST("/pasaportes", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreatePasaporte)
		protected.GET("/pasaportes/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPasaporteByID)
		protected.PUT("/pasaportes/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdatePasaporte)
		protected.DELETE("/pasaportes/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeletePasaporte)

		// CRUD para Persona
		protected.POST("/personas", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreatePersona)
		protected.GET("/personas/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersona)
		protected.PUT("/personas/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdatePersona)
		protected.DELETE("/personas/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeletePersona)
		protected.GET("/personas/cedula/:cedula", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonaByCedula)
		protected.GET("/personas/pasaporte/:pasaporte", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonaByPasaporte)
		protected.GET("/personas/nombre/:nombre", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonaByNombre)
		protected.GET("/personas/nacionalidad/:nacionalidad", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonasByNacionalidad)
		protected.GET("/personas", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonas)
		protected.GET("/personas/cedulas", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetPersonasByCedula)
		protected.GET("/personas/search", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.SearchPersonas)

		// CRUD para Vehiculo
		protected.POST("/vehiculos", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateVehiculo)
		protected.GET("/vehiculos/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetVehiculoByID)
		protected.GET("/vehiculos/search", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetVehiculoByMatricula)
		protected.PUT("/vehiculos/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateVehiculo)
		protected.DELETE("/vehiculos/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteVehiculo)

		// Endpoints para Empresa
		protected.POST("/empresas", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateEmpresa)
		protected.GET("/empresas/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetEmpresaByID)
		protected.GET("/empresas/search", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetEmpresaByRIF)
		protected.PUT("/empresas/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateEmpresa)
		protected.DELETE("/empresas/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteEmpresa)

		// Endpoints para Dirección
		protected.POST("/direcciones", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateDireccion)
		protected.GET("/direcciones/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetDireccionByID)
		protected.PUT("/direcciones/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateDireccion)
		protected.DELETE("/direcciones/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteDireccion)

		// CRUD para Visa
		protected.POST("/visas", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateVisa)
		protected.GET("/visas/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetVisaByID)
		protected.PUT("/visas/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateVisa)
		protected.DELETE("/visas/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteVisa)

		// CRUD para IIO
		protected.POST("/iios", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateIIO)
		protected.GET("/iios/:id", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.GetIIO)
		protected.PUT("/iios/:id", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.UpdateIIO)
		protected.DELETE("/iios/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteIIO)
		protected.GET("/iios", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetIIOs)
		protected.GET("/iios/filter", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.GetIIOs)

		// Configuración de acceso temporal
		protected.POST("/configuracion/acceso-temporal", middleware.RoleRequired("admin"), controllers.GrantTemporaryAccess)
		protected.POST("/configuracion/area", middleware.RoleRequired("admin"), controllers.AddArea)
		protected.PUT("/configuracion/area", middleware.RoleRequired("admin"), controllers.UpdateArea)
		protected.DELETE("/configuracion/area", middleware.RoleRequired("admin"), controllers.RemoveArea)

		// CRUD para Correo
		protected.POST("/correos", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateCorreo)
		protected.GET("/correos/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetCorreoByID)
		protected.PUT("/correos/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateCorreo)
		protected.DELETE("/correos/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteCorreo)

		// CRUD para Mensajes
		protected.GET("/mensajes/:id", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.GetMensaje)
		protected.PUT("/mensajes/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.UpdateMensaje)
		protected.DELETE("/mensajes/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteMensaje)
		protected.GET("/mensajes", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.GetMensajes)
		protected.GET("/mensajes/filter", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.FilterMensajes)
		protected.PUT("/mensajes/:id/procesado", middleware.RoleRequired("admin", "superuser", "analyst"), controllers.UpdateMensajeStatus)
		protected.POST("/mensajes", middleware.RoleRequired("admin", "superuser", "analyst", "user"), controllers.CreateMensaje)
		protected.POST("/create-and-send-mensaje", middleware.RoleRequired("admin", "superuser", "analyst", "user"), controllers.CreateAndSendMensaje)
		protected.POST("/send-mensaje-to-user/:user_id", middleware.RoleRequired("admin", "superuser", "analyst", "user"), controllers.SendMensajeToUser)

		// CRUD para Redes
		protected.POST("/redes", middleware.RoleRequired("admin", "superuser", "user"), controllers.CreateRedes)
		protected.GET("/redes/:id", middleware.RoleRequired("admin", "superuser", "analyst"), middleware.AreaCheck(), controllers.GetRedesByID)
		protected.PUT("/redes/:id", middleware.RoleRequired("admin", "superuser"), middleware.AreaCheck(), controllers.UpdateRedes)
		protected.DELETE("/redes/:id", middleware.RoleRequired("admin"), middleware.AreaCheck(), controllers.DeleteRedes)
	}
}
