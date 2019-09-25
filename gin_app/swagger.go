package gin_app

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/johnnyeven/libtools/env"
)

func getSwaggerJSON() []byte {
	file, err := os.Open("./swagger.json")
	if err != nil {
		return []byte{}
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}
	}

	return data
}

func WithSwagger() gin.HandlerFunc {
	swaggerJSON := []byte{}

	if !env.IsOnline() {
		swaggerJSON = getSwaggerJSON()
	}

	return func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json", swaggerJSON)
	}
}
