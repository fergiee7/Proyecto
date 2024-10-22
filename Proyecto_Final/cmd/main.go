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
	//Email      string `json:"email"`
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
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "Main.html", gin.H{
			"title": "Menú de Acciones",
		})
	})

	// 1. OBTENER
	//1.1OBTENER MATERIAS
	/*router.GET("/", func(c *gin.Context) {
		materias := []Materia{}
		db.Find(&materias)
		c.HTML(200, "index.html", gin.H{
			"title":          "Main website",
			"total_materias": len(materias),
			"materias":       materias,
		})
	})*/
	//Obtener todos los estudiantes
	router.GET("/api/students/:student_id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "XEstudiantes.html", gin.H{})
		studentID := c.Param("student_id")
		var estudiante []Estudiante

		//buscar todos los estudiantes
		if err := db.Where("student_id = ?", studentID).Find(&estudiante).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron datos para este estudiante"})
			return
		}
		c.JSON(http.StatusOK, estudiante)
	})
	// Obtener materias
	router.GET("/api/subjects/:id_subject", func(c *gin.Context) {
		c.HTML(http.StatusOK, "XMaterias.html", gin.H{})
		//materias := []Materia{}
		subjectID := c.Param("id_subject")
		var materias []Materia

		//buscar todos los estudiantes
		if err := db.Where("id_subject = ?", subjectID).Find(&materias).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron materias"})
			return
		}
		c.JSON(http.StatusOK, materias)
		/*db.Find(&materias)
		c.JSON(200, materias)*/
	})

	// 1.2 OBTENER CALIFICACIONES
	//Calificaion especifica por grade_id y strudent_id
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

	//obtener calificaciones  de un estudiante
	router.GET("/api/grades/student/:student_id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "XCalificaciones.html", gin.H{})
		studentID := c.Param("student_id")
		var calificacion []Calificacion

		//buscar todas las calificaciones
		if err := db.Where("student_id = ?", studentID).Find(&calificacion).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron calificaciones para este estudiante"})
			return
		}
		c.JSON(http.StatusOK, calificacion)
	})

	// 2. CREAR
	//Crear estudiante.
	router.POST("/api/students", func(c *gin.Context) {
		c.HTML(http.StatusOK, "LEstudiantes.html", gin.H{})
		var estudiante Estudiante
		if err := c.BindJSON(&estudiante); err == nil {
			result := db.Create(&estudiante)
			if result.Error != nil {
				c.JSON(500, gin.H{"error": "Error al crear Estudiante"})
				return
			}
			c.JSON(200, estudiante)
		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})
	// CREAR MATERIA
	router.POST("/api/subjects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "LMaterias.html", gin.H{})
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
		c.HTML(http.StatusOK, "LCalificaciones.html", gin.H{})
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
	//Actualizar estudiante
	router.PUT("/api/students/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "MEstudiantes.html", gin.H{})
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}

		var estudiante Estudiante
		if err := c.BindJSON(&estudiante); err == nil {
			var estudianteExistente Estudiante
			result := db.First(&estudianteExistente, idParsed)
			if result.Error != nil {
				c.JSON(404, gin.H{"error": "Estudiante no encontrado"})
				return
			}

			estudianteExistente.Name = estudiante.Name
			db.Save(&estudianteExistente)
			c.JSON(200, gin.H{"message": "Alumno actualizado"})
		} else {
			c.JSON(400, gin.H{"error": "Invalid payload"})
		}
	})
	//3.1ACTUALIZAR MATERIA
	router.PUT("/api/subjects/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "MMaterias.html", gin.H{})
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
		c.HTML(http.StatusOK, "MCalificaciones.html", gin.H{})
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
		c.HTML(http.StatusOK, "KMaterias.html", gin.H{})
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
		c.HTML(http.StatusOK, "KCalificaciones.html", gin.H{})
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
	//Eliminar estudiante.
	router.DELETE("/api/students/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "KEstudiantes.html", gin.H{})
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid id"})
			return
		}
		result := db.Delete(&Calificacion{}, idParsed)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error al eliminar estudiante"})
			return
		}
		c.JSON(200, gin.H{"message": "Estudiante eliminado"})
	})
	router.Run(":8001")
}
