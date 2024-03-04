package frontend

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	ListenUri       string
	BackendEndpoint string
}

func StartFrontend(config Config) {
	r := gin.Default()
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	r.GET("/", getHandler)
	r.POST("/", postHandler)
	r.Run(config.ListenUri)
}

func loadTemplate() (*template.Template, error) {
	const indexFile = `<!doctype html>
<head>{{.Title}}</head>
<body>
	<form action="/" method="POST">
			<p>Check some colors</p>
			<label for="red">Red</label>
			<input type="checkbox" name="colors[]" value="red" id="red">
			<label for="green">Green</label>
			<input type="checkbox" name="colors[]" value="green" id="green">
			<label for="blue">Blue</label>
			<input type="checkbox" name="colors[]" value="blue" id="blue">
			<input type="submit">
	</form>
</body>`

	t, err := template.New("index").Parse(indexFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return t, nil
}

type myForm struct {
	Colors []string `form:"colors[]"`
}

func getHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"Title": "Your favourite programming language",
	})
}

func postHandler(c *gin.Context) {
	var fakeForm myForm
	c.ShouldBind(&fakeForm)
	c.JSON(200, gin.H{"color": fakeForm.Colors})
}
