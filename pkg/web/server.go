package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shadi/simple_config/pkg/data"
)

type HttpServer interface {
  ServeRequests()
}

type StoreBackedServer struct {
  db data.Storage
  notifier UpdatesNotifier
}

type PutProperty struct {
  Key string `json:"key"`
  Value string `json:"value"`
  Namespace string `json:"namespace"`
}
type WebhookProperty struct {
  Key string `json:"key"`
  Namespace string `json:"namespace"`
  Callback string `json:"callback"`
}


func GetHttpServer(db data.Storage, notifier UpdatesNotifier) HttpServer {
  return &StoreBackedServer{db, notifier}
}

func (h *StoreBackedServer) ServeRequests() {
  r := gin.New()
  r.Use(gin.Recovery())
  
  r.GET("/ping", func(c *gin.Context){
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })

  r.POST("/v1/put", h.putProperty)
  r.GET("/v1/get", h.getProperty)
  r.POST("/v1/webhook", h.registerWebhook)
  
  r.Run()
}

func (h *StoreBackedServer) putProperty(c *gin.Context) {
  p := &PutProperty{}
  
  if err := c.BindJSON(p); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "problem with request body, needs namespace, key, value",
    })
    return 
  }

  if p.Key == "" || p.Value == "" || p.Namespace == "" {
    c.JSON(http.StatusBadRequest, gin.H{
    "message": "request should have properties namespace, key, value",
    })
    return
  }

  if err := h.db.SetProperty(p.Namespace, p.Key, p.Value); err != nil {
    log.Err(err).Msg("Error saving Property")
    c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
    return
  }

  go h.notifier.NotifyUpdate(p.Namespace, p.Key)

  c.JSON(http.StatusOK, gin.H{"message": "property updated"})
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
    log.Err(err).Msgf("Error Getting Property %s %s", n, k)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
    return
  }

  c.JSON(http.StatusOK, gin.H{"value": v})
}

func (h *StoreBackedServer) registerWebhook(c *gin.Context){
  p := &WebhookProperty{}
  if err := c.BindJSON(p); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "Problem with request body, needs namespace, key, callback",
    })
    return 
  }
  if p.Namespace == "" || p.Key == "" || p.Callback == "" {
    c.JSON(http.StatusBadRequest, gin.H{
    "error": "request should have properties namespace, key, callback",
    })
    return
  }
  err := h.db.RegisterCallback(p.Namespace, p.Key, p.Callback)
  if err != nil {
    log.Err(err).Msg("failed to register callback")
    c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
    return
  }

  h.notifier.RegisterCallback(p.Namespace, p.Key, p.Callback)

  c.JSON(http.StatusOK, gin.H{"message": "callback registered"})
}
