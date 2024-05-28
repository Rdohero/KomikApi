package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kiryuu2 "komikApi/controllers/kiryuu"
	komikCast2 "komikApi/controllers/komikCast"
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

	router.Static("/images", "images/")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})

	komikCast := router.Group("/komikCast")
	kiryuu := router.Group("/kiryuu")

	komikCast.GET("/daftar-komik", func(c *gin.Context) {
		order := c.Query("order")
		page := c.Query("page")

		response, err := komikCast2.GetDaftarKomik(order, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
	komikCast.GET("/fetch-data", komikCast2.GetDataHandler)
	komikCast.GET("/komik-info", komikCast2.GetKomikInfo)
	komikCast.GET("/search", komikCast2.SearchKomik)
	komikCast.GET("/genre", komikCast2.GetGenreInfo)
	komikCast.GET("/genre/komik", komikCast2.FetchComicsByGenre)

	kiryuu.GET("/daftar-komik", func(c *gin.Context) {
		page := c.Query("page")

		response, err := kiryuu2.GetDaftarKomikKiryuu(page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	})
	kiryuu.GET("/search", kiryuu2.SearchKomikKiryuu)
	kiryuu.GET("/komik-info", kiryuu2.GetKomikInfoKiryuu)
	kiryuu.GET("/fetch-data", kiryuu2.GetDataHandlerKiryuu)

	router.Run()
}
