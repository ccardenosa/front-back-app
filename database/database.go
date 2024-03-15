package database

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	ListenUri string
}

type languageType struct {
	Name string `json:"name"`
}

var languages = make(map[string]int)

type developerType struct {
	Name              string `json:"name"`
	FavouriteLanguage string `json:"favourite_lang"`
}

var developers = make(map[string]string)

func StartDatabase(config Config) {

	r := gin.Default()
	r.GET("/languages", listLangHandler)
	r.POST("/language", postLangHandler)
	r.GET("/developers", listDevelHandler)
	r.GET("/developers/:name", getDevelHandler)
	r.POST("/developer", postDevelHandler)
	r.Run(config.ListenUri)
}

func listLangHandler(c *gin.Context) {
	c.JSON(200, gin.H{"Languages": languages})
}

func postLangHandler(c *gin.Context) {
	var h languageType
	c.ShouldBindJSON(&h)
	_, exists := languages[h.Name]
	if !exists {
		languages[h.Name] = 0
		c.String(200, "Success")
	} else {
		c.String(401, "Not Found")
	}
}

func listDevelHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"Developers": developers,
	})
}

func getDevelHandler(c *gin.Context) {
	type Developer struct {
		Name string `uri:"name"`
	}

	var devel Developer
	if err := c.ShouldBindUri(&devel); err != nil {
		c.JSON(400, gin.H{"msg": err})
	} else {
		_, exists := developers[devel.Name]
		if exists {
			c.JSON(200, gin.H{"name": devel.Name, "favourite_lang": developers[devel.Name]})
		} else {
			c.String(401, "Not Found")
		}
	}
}

func postDevelHandler(c *gin.Context) {
	var h developerType
	c.ShouldBindJSON(&h)
	_, exists := developers[h.Name]
	if !exists {
		_, exists := languages[h.FavouriteLanguage]
		if exists {
			developers[h.Name] = h.FavouriteLanguage
			languages[h.FavouriteLanguage]++
			c.String(200, "Success")
			return
		}
	}
	c.String(401, "Not Found")
}
