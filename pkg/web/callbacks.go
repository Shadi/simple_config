package web

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)


type UpdatesNotifier interface {
  NotifyUpdate(namespace, key string);
  RegisterCallback(namespace, key, callback string);
}

type CallbacksHandler struct {
  data map[string]string
  updates chan string
  httpClient *http.Client
}

func GetCallbacksHandler(callbacksMapping map[string]string) UpdatesNotifier{
  h := &CallbacksHandler {
    updates: make(chan string, 100), //ToDo: make it configurable
    data: callbacksMapping,
    httpClient: &http.Client{
      Timeout: time.Second * 1,
    },
  }
  go h.watchUpdates()
  return h
}

func (c *CallbacksHandler) RegisterCallback(namespace, key, callback string) {
  c.data[namespace + key] = callback
}

func (c *CallbacksHandler) NotifyUpdate(namespace, key string) {
  url := c.data[namespace + key]
  if url == "" {
    return
  }
  c.updates <- url
  log.Debug().Msgf("notify %s for %s %s", c.data[namespace + key], namespace, key)
}

func (c *CallbacksHandler) watchUpdates() {
  for url := range c.updates {
		log.Debug().Msgf("Received update request: %s", url)
    c.callWebhook(url)
	}
}

func (c *CallbacksHandler) callWebhook(url string) {
  resp, err := c.httpClient.Get(url)
  if err != nil {
    log.Error().Err(err).Msgf("Error calling %s", url)
    return
  }
  defer resp.Body.Close()
  log.Debug().Msgf("%s resp %s", url, resp.Status)
}

