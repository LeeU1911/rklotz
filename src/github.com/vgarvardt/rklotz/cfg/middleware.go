package cfg

import (
	"net/http"

	"gopkg.in/flosch/pongo2.v3"
	. "github.com/gin-gonic/gin"
)

func Pongo2() HandlerFunc {
	return func(c *Context) {
		c.Next()

		templateName, templateNameExists := c.Get("template")
		templateNameValue, isString := templateName.(string)

		if templateNameExists && isString {
			var template = pongo2.Must(pongo2.FromFile(templateNameValue))
			err := template.ExecuteWriter(getContext(c.Get("data")), c.Writer)
			if err != nil {
				http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func getContext(templateData interface{}, templateDataExists bool) pongo2.Context {
	if templateData == nil || !templateDataExists {
		return nil
	}
	contextData, isMap := templateData.(map[string]interface{})
	if isMap {
		return contextData
	}
	return nil
}
