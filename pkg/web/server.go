package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shadi/simple_config/pkg/data"
)

type StoreBackedServer struct {
  db data.Storage
}

type HttpServer interface {
  ServeRequests()
}

func GetHttpServer(db data.Storage) HttpServer {
  return &StoreBackedServer{db}
}

func (h *StoreBackedServer) ServeRequests() {
  r := gin.Default()

  
  r.GET("/ping", func(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  r.POST("/put", h.putProperty)
  r.GET("/get", h.getProperty)
  
  r.Run()
}

func (h *StoreBackedServer)putProperty(c *gin.Context) {
  p := &data.Property{}
  
  if err := c.BindJSON(p); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "Problem with request body, needs namespace, key, value",
    })
    log.Debug().Err(err).Msgf("Property: %v",p)
    return 
  }

  if p.Key == "" || p.Value == "" || p.Namespace == "" {
    c.JSON(http.StatusBadRequest, gin.H{
    "message": "request should have properties namespace, key, value",
    })
    return
  }

  if err := h.db.SetProperty(p.Namespace, p.Key, p.Value, p.Callback); err != nil {
    log.Err(err).Msg("Error saving Property")
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"Message": "property updated"})
}

func (h *StoreBackedServer) getProperty(c *gin.Context) {
  n := c.DefaultQuery("namespace", "")
  k := c.DefaultQuery("key", "")
  if n == "" || k == "" {
    c.JSON(http.StatusBadRequest, gin.H{
    "message": "request should have properties namespace, key",
    })
    return
  }
  v, err := h.db.GetProperty(n, k);
  if err != nil {
    log.Err(err).Msg("Error saving Property")
    c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal error"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"Value": v})
}
