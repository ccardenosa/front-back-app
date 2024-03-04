package backend

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	ListenUri        string
	DababaseEndpoint string
}

var languages = []string{
	"golang",
	"python",
	"C",
	"C++",
	"Rust",
}

type developer struct {
	Name          string `form:"name"`
	FavouriteLang string `form:"favourite_language"`
}

func StartBackend(config Config) {

	r := gin.Default()
	r.GET("/languages", listLangHandler)
	r.POST("/language", postLangHandler)
	r.GET("/developers", listDevelHandler)
	r.GET("/developers/:developer", getDevelHandler)
	r.POST("/developers/:developer", postDevelHandler)
	r.Run(config.ListenUri)
}

func listLangHandler(c *gin.Context) {
	c.JSON(200, gin.H{"languages": languages})
}

func postLangHandler(c *gin.Context) {
}

func listDevelHandler(c *gin.Context) {
	dv := developer{
		Name:          "ecascaz",
		FavouriteLang: "C++",
	}
	c.JSON(200, gin.H{
		"Developers": []developer{dv},
	})
	// 	"Developer": [
	// 		developer{
	// 			"Name": "ecascaz",
	// 			"Favourite Lang": "C++",
	// 		}
	// 	],
	//  )
}

func getDevelHandler(c *gin.Context) {
	c.JSON(200, gin.H{"Developer": developer{"ecascaz", "C++"}})
}

func postDevelHandler(c *gin.Context) {
	var devel developer
	c.ShouldBind(&devel)
	c.JSON(200, gin.H{})
}
