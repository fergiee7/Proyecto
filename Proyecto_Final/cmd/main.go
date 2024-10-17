package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Estudiante struct {
	IDStudent int    `json:"student_id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Group     string `json:"group"`
	Email     string `json:"email"`
}

type Materia struct {
	Id_subject int    `json:"id_subject" gorm:"primaryKey"`
	Name       string `json:"name"`
	Email      string `json:"email"`
}

type Calificacion struct {
	GradeID   int     `json:"grade_id" gorm:"primaryKey"`
	StudentID int     `json:"student_id"`
	SubjectID int     `json:"subject_id"`
	Grade     float64 `json:"grade"`

	Estudiante Estudiante `gorm:"foreignKey:StudentID;references:IDStudent"`
	Materia    Materia    `gorm:"foreignKey:SubjectID;references:Id_subject"`
}

var db *gorm.DB

func main() {
	dsn := "root:test@tcp(127.0.0.1:3306)/proyecto_1?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error al conectar a la base de datos:", err)
		return
	}

	db.AutoMigrate(&Estudiante{}, &Materia{}, &Calificacion{})
	fmt.Println("Conexión exitosa y tabla creada o actualizada.")

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// 1. OBTENER
	//1.1OBTENER MATERIAS
	router.GET("/", func(c *gin.Context) {
		materias := []Materia{}
		db.Find(&materias)
		c.HTML(200, "index.html", gin.H{
			"title":          "Main website",
			"total_materias": len(materias),
			"materias":       materias,
		})
	})

	// Obtener materias
	router.GET("/api/subjects", func(c *gin.Context) {
		materias := []Materia{}
		db.Find(&materias)
		c.JSON(200, materias)
	})

	// 1.2 OBTENER CALIFICACIONES
	router.GET("/api/grades/:grade_id/student/:student_id", func(c *gin.Context) {
		gradeID := c.Param("grade_id")
		studentID := c.Param("student_id")
		var calificacion Calificacion
		if err := db.Where("grade_id = ? AND student_id = ?", gradeID, studentID).First(&calificacion).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calificación no encontrada"})
			return
		}
		c.JSON(http.StatusOK, calificacion)
	})

	// 2. CREAR
	// CREAR MATERIA
	router.POST("/api/subjects", func(c *gin.Context) {
		var materia Materia
		if err := c.BindJSON(&materia); err == nil {
			result := db.Create(&materia)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear materia"})
				return
			}
			c.JSON(200, materia)
		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})
	//2.2 CREAR  CALIFICACION
	router.POST("/api/grades", func(c *gin.Context) {
		var calificacion Calificacion
		if err := c.BindJSON(&calificacion); err == nil {
			result := db.Create(&calificacion)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear calificación"})
				return
			}
			c.JSON(200, calificacion)

		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})

	// 3. ACTUALIZAR
	//3.1ACTUALIZAR MATERIA
	router.PUT("/api/subjects/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}

		var materia Materia
		if err := c.BindJSON(&materia); err == nil {
			var materiaExistente Materia
			result := db.First(&materiaExistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{"error": "Materia no encontrada"})
				return
			}

			materiaExistente.Name = materia.Name
			db.Save(&materiaExistente)
			c.JSON(200, gin.H{"message": "Materia actualizada"})
		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})
	//3.2 ACTUALIZAR CALIFICACION
	router.PUT("/api/grades/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}

		var calificacion Calificacion
		if err := c.BindJSON(&calificacion); err == nil {
			var calificacionExistente Calificacion
			result := db.First(&calificacionExistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{"error": "Calificacion no encontrada"})
				return
			}
			calificacionExistente.Grade = calificacion.Grade
			db.Save(&calificacionExistente)
			c.JSON(200, gin.H{"message": "Calificacion actualizada"})
		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})

	// 4. ELIMINAR
	//4.1 ELIMINAR MATERIA
	router.DELETE("/api/subjects/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}

		result := db.Delete(&Materia{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar materia"})
			return
		}
		c.JSON(200, gin.H{"message": "Materia eliminada"})
	})

	//ELIMINAR CALIFICACION
	router.DELETE("/api/grades/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
		result := db.Delete(&Calificacion{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar calificacion"})
			return
		}
		c.JSON(200, gin.H{"message": "Calificacion eliminada"})
	})

	router.Run(":8001")
}
