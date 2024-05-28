package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"komikApi/controllers"
	"komikApi/initializers"
	"net/http"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.MigrateDatabase()
}

func main() {
	router := gin.Default()
	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "ngrok-skip-browser-warning", "Authorization"},
		AllowCredentials: true,
	}
	router.Use(cors.New(config))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})

	komikCast := router.Group("/komikCast")
	kiryuu := router.Group("/kiryuu")

	komikCast.GET("/daftar-komik", func(c *gin.Context) {
		order := c.Query("order")
		page := c.Query("page")

		response, err := controllers.GetDaftarKomik(order, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
	komikCast.GET("/fetch-data", controllers.GetDataHandler)
	komikCast.GET("/komik-info", controllers.GetKomikInfo)
	komikCast.GET("/search", controllers.SearchKomik)
	komikCast.GET("/genre", controllers.GetGenreInfo)
	komikCast.GET("/genre/komik", controllers.FetchComicsByGenre)

	kiryuu.GET("/daftar-komik", func(c *gin.Context) {
		page := c.Query("page")

		response, err := controllers.GetDaftarKomikKiryuu(page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
	kiryuu.GET("/search", controllers.SearchKomikKiryuu)
	kiryuu.GET("/komik-info", controllers.GetKomikInfoKiryuu)
	kiryuu.GET("/fetch-data", controllers.GetDataHandlerKiryuu)

	router.Run()
}
