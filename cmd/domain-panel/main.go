package main

import (
	"errors"
	"log"
	"time"

	"git.cm/nb/domain-panel"
	"git.cm/nb/domain-panel/controller"
	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"
)

func init() {
	panel.DB.AutoMigrate(
		panel.User{},
		panel.Panel{},
		panel.Cat{},
		panel.Domain{},
		panel.Offer{},
	)
}

func main() {
	controller.Web()
	go updateWhois()
	select {}
}

func updateWhois() {
	var domains []panel.Domain
	for {
		panel.DB.Where("whois_update is NULL or DATEDIFF(now(),whois_update)>7", time.Now()).Find(&domains)
		for _, domain := range domains {
			result, err := whois.Whois(domain.Domain)
			if err == nil {
				var parsed parser.WhoisInfo
				parsed, err = parser.Parse(result)
				if err == nil {
					create, _ := parseTime(parsed.Registrar.CreatedDate)
					expire, _ := parseTime(parsed.Registrar.ExpirationDate)
					panel.DB.Model(&domain).UpdateColumns(panel.Domain{
						Registrar:   parsed.Registrar.RegistrarName,
						Create:      create,
						Expire:      expire,
						WhoisUpdate: time.Now(),
					})
					log.Println(domain)
				}
			}
			time.Sleep(time.Minute)
		}
		time.Sleep(time.Hour)
	}
}

var timeLayouts = []string{
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02.01.2006 15:04:05",
	time.RFC1123,     //= "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    //= "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RFC850,      //= "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC3339,     //= "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, //= "2006-01-02T15:04:05.999999999Z07:00"
}

func parseTime(t string) (tt time.Time, e error) {
	for _, layout := range timeLayouts {
		tt, e = time.Parse(layout, t)
		if e == nil {
			return
		}
	}
	e = errors.New("解析失败")
	return
}
