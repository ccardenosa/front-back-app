package frontend

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
)

type Config struct {
	ListenUri       string
	BackendEndpoint string
}

var cnf Config

func getTemplatesPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Unable to get the current filename")
	}
	dirname := filepath.Join(filepath.Dir(filename), "html")
	log.Println("Templates are stored at: ", dirname)
	return dirname
}

func generateTemplatesNamePathMap() map[string]string {

	basenames := make(map[string]string)

	dir := getTemplatesPath()

	// Walk the directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if it's a file
		if !info.IsDir() {
			// Get the basename without extension
			basename := filepath.Base(path)
			basenameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))

			// Add to the map
			basenames[basenameWithoutExt] = basename
		}

		return nil
	})

	if err != nil {
		log.Printf("Error walking the directory: %v\n", err)
		panic(err)
	}

	return basenames
}

var tmplMap map[string]string = generateTemplatesNamePathMap()

func StartFrontend(config Config) {

	cnf = config

	r := gin.Default()
	r.LoadHTMLGlob(getTemplatesPath() + "/*.tmpl")
	r.GET("/", getHandler)
	r.POST("/", postHandler)
	r.Run(config.ListenUri)
}

func getLanguages() []string {

	var languages []string

	resp, err := http.Get("http://" + cnf.BackendEndpoint + "/languages")
	if err != nil {
		log.Println(err)
		return languages
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Not OK Response status:", resp.Status)
	} else {
		log.Println("Response status:", resp.Status)
		var langs interface{}
		err = json.NewDecoder(resp.Body).Decode(&langs)
		if err != nil {
			panic(err)
		}
		log.Printf("Backend response is: %s", langs)
		for _, v := range langs.(map[string]any) {
			for k, _ := range v.(map[string]any) {
				log.Printf("Found '%v' of type '%T'", k, k)
				languages = append(languages, k)
			}
		}
	}

	return languages
}

func getHandler(c *gin.Context) {

	availableLanguages := getLanguages()

	c.HTML(http.StatusOK, tmplMap["index"], gin.H{
		"Title": "Your favourite programming language",
		"Lang":  availableLanguages,
	})
}

type newDevelForm struct {
	Name    string `form:"name"`
	Lang    string `form:"fav_language"`
	NewLang string `form:"new_language"`
}

type developerType struct {
	Name              string `json:"name"`
	FavouriteLanguage string `json:"favourite_lang"`
}

type langCount struct {
	Name  string
	Count float64
}

func postHandler(c *gin.Context) {
	var f newDevelForm
	var dev developerType
	c.ShouldBind(&f)
	dev.Name = f.Name
	if f.Lang == "Other" {
		dev.FavouriteLanguage = f.NewLang
	} else {
		dev.FavouriteLanguage = f.Lang
	}
	log.Println("Form Fields", dev)
	data, _ := json.Marshal(dev)
	ret := postNewDeveloper(data)
	c.HTML(http.StatusOK, tmplMap["languages-ranking"], gin.H{
		"Title":         "Favourites programming languages",
		"LanguagesList": ret,
	})
}

func postNewDeveloper(newDeveloper []byte) (languages []langCount) {

	ct := "application/json"
	resp, err := http.Post("http://"+cnf.BackendEndpoint+"/developer", ct, bytes.NewBuffer(newDeveloper))
	if err != nil {
		log.Println(err)
		return languages
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Not OK Response status:", resp.Status)
	} else {
		log.Println("Response status:", resp.Status)
		var langs map[string]any
		err = json.NewDecoder(resp.Body).Decode(&langs)
		if err != nil {
			panic(err)
		}
		log.Printf("Backend response: %s", langs)
		for _, v := range langs {
			for k, vv := range v.(map[string]any) {
				log.Printf("Found '%v' of type '%T'", k, k)
				l := langCount{
					Name:  k,
					Count: vv.(float64),
				}
				languages = append(languages, l)
			}
		}
	}

	return languages
}
