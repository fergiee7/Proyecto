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

	// ** Rutas para el menú principal **
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Menú de Acciones",
		})
	})

	// 1. OBTENER
	// 1.1 OBTENER MATERIAS
	router.GET("/api/subjects", func(c *gin.Context) {
		materias := []Materia{}
		db.Find(&materias)
		c.JSON(http.StatusOK, materias)
	})

	// 1.2 OBTENER CALIFICACIONES
	// Calificación específica por grade_id y student_id
	router.GET("/api/grades/:grade_id/student/:student_id", func(c *gin.Context) {
		gradeID := c.Param("grade_id")
		studentID := c.Param("student_id")
		var calificacion Calificacion

		// Convertir los parámetros a enteros
		gradeIDInt, err := strconv.Atoi(gradeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid grade ID"})
			return
		}
		studentIDInt, err := strconv.Atoi(studentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}

		// Buscar la calificación
		if err := db.Where("grade_id = ? AND student_id = ?", gradeIDInt, studentIDInt).First(&calificacion).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Calificación no encontrada"})
			return
		}
		c.JSON(http.StatusOK, calificacion)
	})

	// Obtener calificaciones de un estudiante
	router.GET("/api/grades/student/:student_id", func(c *gin.Context) {
		studentID := c.Param("student_id")
		var calificaciones []Calificacion

		// Buscar todas las calificaciones
		if err := db.Where("student_id = ?", studentID).Find(&calificaciones).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No se encontraron calificaciones para este estudiante"})
			return
		}
		c.JSON(http.StatusOK, calificaciones)
	})

	// 1.3 OBTENER ESTUDIANTES
	router.GET("/api/students", func(c *gin.Context) {
		estudiantes := []Estudiante{}
		db.Find(&estudiantes)
		c.JSON(http.StatusOK, estudiantes)
	})

	// Obtener Estudiante específico
	router.GET("/api/students/:student_id", func(c *gin.Context) {
		studentID := c.Param("student_id")
		var estudiante Estudiante
		if err := db.First(&estudiante, studentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
			return
		}
		c.JSON(http.StatusOK, estudiante)
	})

	// 2. CREAR
	// 2.1 CREAR MATERIA
	router.POST("/api/subjects", func(c *gin.Context) {
		var materia Materia
		if err := c.BindJSON(&materia); err == nil {
			result := db.Create(&materia)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear materia"})
				return
			}
			c.JSON(http.StatusCreated, materia)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 2.2 CREAR CALIFICACION
	router.POST("/api/grades", func(c *gin.Context) {
		var calificacion Calificacion
		if err := c.BindJSON(&calificacion); err == nil {
			// Verificar si el estudiante existe
			var estudiante Estudiante
			if err := db.First(&estudiante, calificacion.StudentID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Estudiante no encontrado"})
				return
			}

			// Verificar si la materia existe
			var materia Materia
			if err := db.First(&materia, calificacion.SubjectID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Materia no encontrada"})
				return
			}

			result := db.Create(&calificacion)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear calificación"})
				return
			}

			c.JSON(http.StatusCreated, calificacion)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 2.3 CREAR ESTUDIANTE
	router.POST("/api/students", func(c *gin.Context) {
		var estudiante Estudiante
		if err := c.BindJSON(&estudiante); err == nil {
			result := db.Create(&estudiante)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear estudiante"})
				return
			}
			c.JSON(http.StatusCreated, estudiante)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 3. ACTUALIZAR
	// 3.1 ACTUALIZAR MATERIA
	router.PUT("/api/subjects/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		var materia Materia
		if err := c.BindJSON(&materia); err == nil {
			var materiaExistente Materia
			result := db.First(&materiaExistente, idParsed)
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Materia no encontrada"})
				return
			}

			materiaExistente.Name = materia.Name
			db.Save(&materiaExistente)
			c.JSON(http.StatusOK, gin.H{"message": "Materia actualizada"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 3.2 ACTUALIZAR CALIFICACION
	router.PUT("/api/grades/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		var calificacion Calificacion
		if err := c.BindJSON(&calificacion); err == nil {
			var calificacionExistente Calificacion
			result := db.First(&calificacionExistente, idParsed)
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Calificación no encontrada"})
				return
			}
			calificacionExistente.Grade = calificacion.Grade
			db.Save(&calificacionExistente)
			c.JSON(http.StatusOK, gin.H{"message": "Calificación actualizada"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 3.3 ACTUALIZAR ESTUDIANTE
	router.PUT("/api/students/:student_id", func(c *gin.Context) {
		id := c.Param("student_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		var estudiante Estudiante
		if err := c.BindJSON(&estudiante); err == nil {
			var estudianteExistente Estudiante
			result := db.First(&estudianteExistente, idParsed)
			if result.Error != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Estudiante no encontrado"})
				return
			}

			estudianteExistente.Name = estudiante.Name
			estudianteExistente.Group = estudiante.Group
			estudianteExistente.Email = estudiante.Email
			db.Save(&estudianteExistente)
			c.JSON(http.StatusOK, gin.H{"message": "Estudiante actualizado"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		}
	})

	// 4. ELIMINAR
	// 4.1 ELIMINAR MATERIA
	router.DELETE("/api/subjects/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		result := db.Delete(&Materia{}, idParsed)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar materia"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Materia eliminada"})
	})

	// ELIMINAR CALIFICACION
	router.DELETE("/api/grades/:id", func(c *gin.Context) {
		id := c.Param("id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		result := db.Delete(&Calificacion{}, idParsed)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar calificación"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Calificación eliminada"})
	})

	// 4.3 ELIMINAR ESTUDIANTE
	router.DELETE("/api/students/:student_id", func(c *gin.Context) {
		id := c.Param("student_id")
		idParsed, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}

		result := db.Delete(&Estudiante{}, idParsed)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar estudiante"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Estudiante eliminado"})
	})

	router.Run(":8001")
}
