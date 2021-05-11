package handler

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-starter-project/pkg/config"
	"go-starter-project/pkg/derror"
)

func HandlerReturnError(c *gin.Context, err error) {
	var derr derror.Derror

	switch err := err.(type) {
	case derror.Derror:
		derr = err
	case *derror.Derror:
		derr = *err
	default:
		derr = derror.E(err).Log()
	}

	var reply interface{}
	if config.Conf.Debug {
		reply = derr
	} else {
		reply = derr.Strip()
	}

	if derr.HTTPCode != 0 {
		c.JSON(derr.HTTPCode, reply)
	} else if derr.ErrCode == derror.ErrCodeUnauthorized {
		c.JSON(http.StatusUnauthorized, reply)
	} else if derr.ErrCode == derror.ErrCodeNotFound {
		c.JSON(http.StatusNotFound, reply)
	} else if derr.ErrCode == derror.ErrCodeForbidden {
		c.JSON(http.StatusForbidden, reply)
	} else if derr.ErrCode == derror.ErrCodeInputError {
		c.JSON(http.StatusBadRequest, reply)
	} else {
		c.JSON(http.StatusInternalServerError, reply)
	}
}

func HandlerReturnData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func HandlerRegisterSwagger(r *gin.Engine) {
	if !config.Conf.Debug {
		return
	}

	basicAuth := gin.BasicAuth(map[string]string{
		config.Conf.Docs.Username: config.Conf.Docs.Password,
	})

	url := ginSwagger.URL("swagger.yml")
	r.GET("/docs", func(c *gin.Context) {
		redirect := `<meta http-equiv="refresh" content="0; url=docs/index.html">`
		c.Data(http.StatusOK, "text/html", []byte(redirect))
	})
	r.GET("/docs/*filepath", basicAuth, func(c *gin.Context) {
		filepath := c.Param("filepath")
		switch filepath {
		case "/swagger.yml":
			c.File("./swagger.yml")
		case "/":
			redirect := `<meta http-equiv="refresh" content="0; url=index.html">`
			c.Data(http.StatusOK, "text/html", []byte(redirect))
		default:
			ginSwagger.WrapHandler(swaggerFiles.Handler, url)(c)
		}
	})
}

func HandlerRegisterVersion(r *gin.Engine) {
	r.StaticFile("/version", "./version")
}

func HandlerRegisterCORS(r *gin.Engine) {
	corsConf := cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowAllOrigins = config.Conf.CORS.AllowAll
	corsConf.AllowWebSockets = true
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	if !config.Conf.CORS.AllowAll {
		corsConf.AllowOrigins = config.Conf.CORS.AllowedDomains
	}
	r.Use(cors.New(corsConf))
}

func HandlerRegisterHealthCheck(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func HandlerRegisterCORSRG(r *gin.RouterGroup) {
	corsConf := cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowAllOrigins = config.Conf.CORS.AllowAll
	corsConf.AllowWebSockets = true
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	if !config.Conf.CORS.AllowAll {
		corsConf.AllowOrigins = config.Conf.CORS.AllowedDomains
	}
	r.Use(cors.New(corsConf))
}

func HandlerRegister404(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			derr := derror.E(derror.ErrEndpointNotFound).SetCode(derror.ErrCodeNotFound).
				SetDebug("endpoint not found").
				SetExtraInfo("path", c.Request.URL.Path).
				Log()
			HandlerReturnError(c, derr)
		} else {
			c.Data(http.StatusNotFound, "text/plain", []byte(http.StatusText(http.StatusNotFound)))
		}
	})
}
