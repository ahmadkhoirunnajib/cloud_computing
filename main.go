package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"main/controllers"
	"main/models"
	"main/wa"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {

	// =====================================
	// Load ENV
	// =====================================
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// =====================================
	// Database
	// =====================================
	db := koneksi()

	db.AutoMigrate(
		&models.Suhu{},
		&models.User{},
		&models.Nilai{},
		&models.Dokumen{},
		&models.Pesan{},
	)

	// =====================================
	// Gin
	// =====================================
	r := gin.Default()

	r.Use(cors.Default())

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// =====================================
	// JWT
	// =====================================
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{

		Realm:       "cloud-computing",
		Key:         []byte("rahasia-jwt-cloud"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",

		Authenticator: controllers.UserLogin,

		PayloadFunc: func(data interface{}) jwtgo.MapClaims {

			if user, ok := data.(models.User); ok {

				return jwtgo.MapClaims{
					"id":       user.ID,
					"username": user.Username,
				}
			}

			return jwtgo.MapClaims{}
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal(err)
	}

	// =====================================
	// Public Route
	// =====================================

	r.POST("/login", authMiddleware.LoginHandler)
	r.POST("/register", controllers.UserTambah)
	r.POST("/logout", authMiddleware.LogoutHandler)

	// =====================================
	// Protected Route
	// =====================================

	auth := r.Group("/backend")
	auth.Use(authMiddleware.MiddlewareFunc())

	{
		auth.GET("/profile", controllers.Profile)

		// SUHU
		auth.GET("/suhu", controllers.Tampil)
		auth.POST("/suhu", controllers.Tambah)
		auth.PUT("/suhu", controllers.Ubah)
		auth.DELETE("/suhu", controllers.Hapus)

		// USER
		auth.GET("/user", controllers.UserTampil)
		auth.POST("/user", controllers.UserTambah)
		auth.PUT("/user", controllers.UserUbah)
		auth.DELETE("/user", controllers.UserHapus)

		// NILAI
		auth.GET("/nilai", controllers.NilaiTampil)
		auth.POST("/nilai", controllers.NilaiTambah)
		auth.PUT("/nilai", controllers.NilaiUbah)
		auth.DELETE("/nilai", controllers.NilaiHapus)

		// DRIVE
		auth.POST("/drive", controllers.DriveUpload)
		auth.GET("/drive", controllers.DriveTampil)
		auth.GET("/drive/:id", controllers.DriveUnduh)

		// SHEET
		auth.GET("/sheet", controllers.SheetTampil)
		auth.POST("/sheet", controllers.SheetTambah)
		auth.PUT("/sheet/update", controllers.SheetUpdate)

		//pesan
		auth.GET("/pesan", controllers.PesanTampil)
		auth.POST("/pesan", controllers.PesanTambah)
		auth.PUT("/pesan", controllers.PesanUbah)
		auth.DELETE("/pesan", controllers.PesanUbah)

	}

	// =====================================
	// Jalankan WhatsApp Bot
	// =====================================

	go wa.KonekWa()

	go func() {
		sig := make(chan os.Signal, 1)

		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

		<-sig

		log.Println("Menutup WhatsApp...")
		wa.StopWA()

		os.Exit(0)
	}()

	// =====================================
	// Jalankan Server
	// =====================================

	port := os.Getenv("PORT")

	if port == "" {
		port = "8111"
	}

	log.Println("===================================")
	log.Println("Server berjalan pada port :", port)
	log.Println("===================================")

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
