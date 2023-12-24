package marketplace

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (controller *MarketplaceController) GetRepoStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Make a GET request to the GitHub API's repository endpoint
		resp, err := http.Get("https://api.github.com/repos/{owner}/{repo}")
		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "an internal server error occurred",
			})
			return
		}
		defer resp.Body.Close()

		// Parse the response into a map
		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "an internal server error occurred",
			})
			return
		}

		// Extract the statistics from the map
		files := data["size"].(float64)
		lines := data["lines"].(float64)
		commits := data["commits"].(float64)
		contributors := data["contributors"].(float64)

		// Create a new map with the statistics
		stats := map[string]float64{
			"files":       files,
			"lines":       lines,
			"commits":     commits,
			"contributors": contributors,
		}

		// Return the statistics in the HTTP response
		c.JSON(http.StatusOK, stats)
	}
}
