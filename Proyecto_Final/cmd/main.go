package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// se van a usar las estructuras db.Create, db.Find, db.Delete, db.Save
// Estructura para estudiantes
type Estudiante struct {
	StudentID int    `json:"student_id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Group     string `json:"group"`
	Email     string `json:"email"`
}

// crear Tabla Materia
type Materia struct {
	Id_subject int    `json:"id_subject"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

// Estructura para calificaciones
type Calificacion struct {
	GradeID   int        `json:"grade_id" gorm:"primaryKey"`
	StudentID Estudiante `gorm:"foreignKey:Id_subject;references:Id_subject"` // Llave foránea que referencia a Student

	SubjectID Materia `gorm:"foreignKey:StudentID;references:StudentID"` // Llave foránea que referencia a Student

	Grade float64 `json:"grade"`
}

// Definiendo variable para la base de datos para no crear otro documento
var db *gorm.DB

func main() {

	//Creando conexion con tabla
	dsn := "root:test@tcp(127.0.0.1:3306)/proyecto_1?charset=utf8mb4&parseTime=True&loc=Local"
	var err error                                        //para ir guardando los errores
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //abriendo base de datos
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	//creando tabla si no existe usando la estructura de las 3 tablas
	db.AutoMigrate(&Estudiante{}, &Materia{}, &Calificacion{})

	fmt.Println("Conexión exitosa y tabla creada o actualizada.")

	router := gin.Default()

	//Tomar archivos de la carpeta template
	router.LoadHTMLGlob("templates/*")

	//revisar si está corriendo el servidor
	router.GET("/ping", func(c *gin.Context) { //se define una url
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//1. OBTENER MATERIAS
	//1.1MATERIAS
	router.GET("/", func(c *gin.Context) {

		materias := []Materia{}
		db.Find(&materias)

		c.HTML(200, "index.html", gin.H{
			"title":          "Main website",
			"total_materias": len(materias),
			"materias":       materias,
		})
	})

	//API URL obtener materias
	router.GET("/api/subjects", func(c *gin.Context) {
		materias := []Materia{}
		db.Find(&materias) // Obtener todos los usuarios
		c.JSON(200, materias)

	})

	//1.2 CALIFICACIONES

	router.GET("/api/grades/:grade_id/student/:student_id", func(c *gin.Context) {
		gradeID := c.Param("grade_id")
		studentID := c.Param("student_id")
		var calificacion Calificacion

		// Buscar la calificación específica por grade_id y student_id
		if err := db.Where("grade_id = ? AND student_id = ?", gradeID, studentID).First(&calificacion).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calificación no encontrada para el estudiante"})
			return
		}

		c.JSON(http.StatusOK, calificacion)
	})

	// Obtener todas las calificaciones de un estudiante
	router.GET("/api/grades/student/:student_id", func(c *gin.Context) {
		studentID := c.Param("student_id")
		var calificacion []Calificacion

		// Buscar todas las calificaciones de un estudiante
		if err := db.Where("student_id = ?", studentID).Find(&calificacion).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron calificaciones para este estudiante"})
			return
		}

		c.JSON(http.StatusOK, calificacion)
	})

	//2.CREAR
	//2.1 CREAR MATERIA
	router.POST("/api/subjects", func(c *gin.Context) {
		var materia Materia
		if c.BindJSON(&materia) == nil {

			// Crear un nuevo usuario en la base de datos
			result := db.Create(&materia)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear materia"})
				return
			}
			c.JSON(200, materia)

		} else {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
			return
		}

	})

	//2.2 CREAR CALIFICACIÓN
	router.POST("/api/grades", func(c *gin.Context) {
		var calificacion Calificacion
		if c.BindJSON(&calificacion) == nil {

			// Crear un nuevo usuario en la base de datos
			result := db.Create(&calificacion)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear nueva califacion"})
				return
			}
			c.JSON(200, calificacion)

		} else {
			c.JSON(400, gin.H{
				"error": "Invalid payload",
			})
			return
		}

	})

	//3. ELIMINAR
	//3.1 ELIMINAR MATERIA
	router.DELETE("/api/subjects:id", func(c *gin.Context) {
		id := c.Param("student_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		result := db.Delete(&Materia{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar Materia"})
			return
		}
		c.JSON(200, gin.H{
			"message": "Materia eliminada",
		})

		//c.JSON(201, gin.H{})

	})

	//3.2 ELIMINAR CALIFICACION
	router.DELETE("/api/grades/:grade_id", func(c *gin.Context) {
		id := c.Param("grade_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		result := db.Delete(&Materia{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar Materia"})
			return
		}
		c.JSON(200, gin.H{
			"message": "Materia eliminada",
		})

	})

	//3.ACTUALIZAR
	//3.1 ACTUALIZAR MATERIAS
	router.PUT("/api/subjects:id", func(c *gin.Context) {
		id := c.Param("student_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		var materia Materia
		if c.BindJSON(&materia) == nil {

			var materiaexistente Materia
			result := db.First(&materiaexistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{
					"error": "Materia no encontrada",
				})
				return

			}

			//Se actualiza materia
			materiaexistente.Name = materia.Name

			db.Save(&materiaexistente)

			c.JSON(200, gin.H{
				"message": "Materia actualizada",
			})

		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})

		}

	})

	//3.2 ACTUALIZAR CALIFICACIONES
	router.PUT("/api/grades/:grade_id", func(c *gin.Context) {
		id := c.Param("grade_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Invalid id",
			})
			return

		}

		var calificacion Calificacion
		if c.BindJSON(&calificacion) == nil {

			var calificaionexistente Calificacion
			result := db.First(&calificaionexistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{
					"error": "Calificaion no encontrada",
				})
				return

			}

			//Se actualiza materia
			calificaionexistente.Grade = calificacion.Grade

			db.Save(&calificaionexistente)

			c.JSON(200, gin.H{
				"message": "Calificación actualizada",
			})

		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})

		}

	})

	router.Run(":8001") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
