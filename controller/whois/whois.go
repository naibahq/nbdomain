package whois

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"
)

//Whois whois 查询
func Whois(c *gin.Context) {
	domain := c.Param("domain")
	if len(domain) < 4 {
		c.String(http.StatusForbidden, "域名格式不符合规范")
		return
	}
	result, err := whois.Whois(domain)
	if err == nil {
		var parsed parser.WhoisInfo
		parsed, err = parser.Parse(result)
		if err == nil {
			c.JSON(http.StatusOK, parsed)
			return
		}
	}
	log.Println("whois", err)
	c.Status(http.StatusNoContent)
}
