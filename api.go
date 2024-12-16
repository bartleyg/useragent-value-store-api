package useragentvaluestore

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const VERSION = "/v1"

type Api struct {
	mu      sync.RWMutex
	kvstore map[string]string
}

func NewApi() *gin.Engine {
	api := Api{
		kvstore: make(map[string]string),
	}
	return api.setupRoutes()
}

func (api *Api) setupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET(VERSION+"/value", api.Get)
	r.POST(VERSION+"/value", api.Upsert)
	r.DELETE(VERSION+"/value", api.Delete)

	r.GET("/", api.HomePage)

	return r
}

// Get
func (api *Api) Get(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")

	api.mu.RLock()
	defer api.mu.RUnlock()

	// validate userAgent exists in store
	value, found := api.kvstore[userAgent]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("userAgent '%v' has nothing stored", userAgent)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userAgent": userAgent,
		"value":     value,
	})
}

// Upsert (Insert or Update)
func (api *Api) Upsert(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")

	api.mu.Lock()
	defer api.mu.Unlock()

	// read raw data as value from request body
	bytesValue, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value := string(bytesValue)
	api.kvstore[userAgent] = value

	c.JSON(http.StatusOK, gin.H{
		"userAgent": userAgent,
		"value":     value,
	})
}

// Delete
func (api *Api) Delete(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")

	api.mu.Lock()
	defer api.mu.Unlock()

	// validate userAgent exists in store
	_, found := api.kvstore[userAgent]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("userAgent '%v' has nothing stored", userAgent)})
		return
	}

	delete(api.kvstore, userAgent)

	c.Status(http.StatusOK)
}

// HomePage shows the current user-agent and how many user-agents have stored an item.
func (api *Api) HomePage(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")

	api.mu.RLock()
	defer api.mu.RUnlock()

	numUserAgents := len(api.kvstore)

	c.JSON(http.StatusOK, gin.H{
		"userAgent":      userAgent,
		"numUserAgents:": numUserAgents,
	})
}
