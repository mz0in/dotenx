package marketplace

import (
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (controller *MarketplaceController) GetFunction() gin.HandlerFunc {
	return func(c *gin.Context) {

		var accountId string
		if a, ok := c.Get("accountId"); ok {
			accountId = a.(string)
		}
		functionName := c.Param("function_name")
		function, err := controller.Service.GetFunction(functionName)
		if err != nil && err.Error() == "function not found" {
			c.JSON(http.StatusOK, gin.H{
				"access": true,
				"exist":  false,
			})
			return
		} else if err != nil {
			logrus.Error(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "an internal server error occurred",
			})
			return
		}

		if function.AccountId != accountId {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "you haven't access to get information about this function",
			})
			return
		}

		repoStatsHandler := controller.GetRepoStats()
		repoStatsHandler(c)

		var stats map[string]float64
		if err := json.NewDecoder(c.Request.Body).Decode(&stats); err != nil {
			logrus.Error("Failed to decode repository stats: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to decode repository stats"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access":   true,
			"exist":    true,
			"function": function,
			"stats":    stats,
		})
	}

}
