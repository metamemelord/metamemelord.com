package handlers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/metamemelord/portfolio-website/worker"
)

func refreshData(c *gin.Context) {
	log.Println("Force refresh called")
	worker.RefreshData()
	respond(c, http.StatusAccepted, nil, nil)
}

func getGithubReposHandler(c *gin.Context) {
	log.Println("[INFO] Github data accessed")
	c.JSON(http.StatusOK, worker.GetData().GithubData)
}

func getWordpressPostsHandler(c *gin.Context) {
	log.Println("[INFO] Wordpress data accessed")
	c.JSON(http.StatusOK, worker.GetData().WordpressData)
}

func getWordpressPostbyIDHandler(c *gin.Context) {
	id := c.Param("id")
	pid, err := strconv.Atoi(id)
	if err != nil || pid <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
	} else if pid > len(worker.GetData().WordpressData) {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		wpData := worker.GetData().WordpressData
		c.JSON(http.StatusOK, wpData[len(wpData)-pid])
	}
}

func getWordpressWebpage(c *gin.Context) {
	path := c.Param("path")
	log.Println(path)
	baseWordpressPath := "https://theanonymosopher.wordpress.com"

	finalPath := baseWordpressPath + path
	req, _ := http.NewRequest(http.MethodGet, finalPath, nil)
	req.Header = c.Request.Header.Clone()

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[Error]", err)
		c.Redirect(302, "/error")
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		if strings.ToLower(k) == "cache-control" {
			c.Header(k, "no-cache, must-revalidate, max-age=0")
			continue
		}
		c.Header(k, v[0])
	}
	io.Copy(c.Writer, resp.Body)
}
