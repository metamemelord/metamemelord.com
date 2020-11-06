package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/metamemelord/portfolio-website/model"
	"github.com/metamemelord/portfolio-website/pkg/core"
	"github.com/metamemelord/portfolio-website/pkg/worker"
)

func addRedirection(c *gin.Context) {
	redirectionItem := model.RedirectionItem{}
	if err := c.BindJSON(&redirectionItem); err != nil {
		log.Println(err)
		respond(c, http.StatusBadRequest, nil, err)
	} else {
		_, err := worker.AddRedirectionItem(c.Request.Context(), &redirectionItem)
		if err != nil {
			respond(c, http.StatusServiceUnavailable, nil, err)
		} else {
			respond(c, http.StatusCreated, nil, nil)
		}
	}
}

func resolveRedirection(c *gin.Context) {
	routingKey := c.Param(core.ROUTING_KEY)
	pathToForward := c.Param(core.PATH_LABEL)
	if target, code, err := worker.ResolveRedirectionItem(routingKey, pathToForward, c.Request.URL.RawQuery); err == nil {
		c.Redirect(code, target)
		c.Abort()
	} else {
		log.Println("Error while redirecting", err)
	}
}

func deleteRedirection(c *gin.Context) {
	routingKey := c.Param(core.ROUTING_KEY)
	if err := worker.DeleteRedirectionItem(routingKey); err == nil {
		respond(c, http.StatusOK, nil, nil)
	} else {
		respond(c, http.StatusServiceUnavailable, nil, err)
	}
}
