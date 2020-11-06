package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metamemelord/portfolio-website/model"
	"github.com/metamemelord/portfolio-website/pkg/core"
	"github.com/metamemelord/portfolio-website/pkg/worker"
)

func addProxy(c *gin.Context) {
	proxyItem := model.ProxyItem{}
	if err := c.BindJSON(&proxyItem); err != nil {
		log.Println(err)
		respond(c, http.StatusBadRequest, nil, err)
	} else {
		_, err := worker.AddProxyItem(c.Request.Context(), &proxyItem)
		if err != nil {
			respond(c, http.StatusServiceUnavailable, nil, err)
		} else {
			respond(c, http.StatusCreated, nil, nil)
		}
	}
}

func resolveProxy(c *gin.Context) {
	routingKey := c.Param(core.ROUTING_KEY)
	pathToForward := c.Param(core.PATH_LABEL)
	if target, code, err := worker.ResolveProxyItem(routingKey, pathToForward, c.Request.URL.RawQuery); err == nil {
		c.Redirect(code, target)
		c.Abort()
	} else {
		log.Println("Error while redirecting", err)
	}
}

func deleteProxy(c *gin.Context) {
	routingKey := c.Param(core.ROUTING_KEY)
	if err := worker.DeleteProxyItem(routingKey); err == nil {
		respond(c, http.StatusOK, nil, nil)
	} else {
		respond(c, http.StatusServiceUnavailable, nil, err)
	}
}
