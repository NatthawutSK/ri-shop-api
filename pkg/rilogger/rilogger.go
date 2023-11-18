package rilogger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/NatthawutSK/ri-shop/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type IRiLogger interface {
	Print() IRiLogger
	Save()
	setQuery(c *fiber.Ctx)
	setBody(c *fiber.Ctx)
	setResponse(res any)
}

type RiLogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitRiLogger(c *fiber.Ctx, res any) IRiLogger {
	log := &RiLogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		Path:       c.Path(),
		StatusCode: c.Response().StatusCode(),
	}
	log.setQuery(c)
	log.setBody(c)
	log.setResponse(res)
	return log
}

// Print implements IRiLogger.
func (l *RiLogger) Print() IRiLogger {
	utils.Debug(l)
	return l

}

// Save implements IRiLogger.
func (l *RiLogger) Save() {
	data := utils.Output(l)

	fileName := fmt.Sprintf("./assets/logs/rilogger_%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	file.WriteString(string(data) + "\n")
}

// setBody implements IRiLogger.
func (l *RiLogger) setBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(&body); err != nil {
		log.Printf("error parsing body: %v", err)
	}

	switch l.Path{
		case "v1/users/signup":
			l.Body = "HAHA XD"
		default:
			l.Body = body
	}
}

// setQuery implements IRiLogger.
func (l *RiLogger) setQuery(c *fiber.Ctx) {
	var query any
	if err := c.BodyParser(&query); err != nil {
		log.Printf("error parsing query: %v", err)
	}
	l.Query = query
	
}

// setResponse implements IRiLogger.
func (l *RiLogger) setResponse(res any) {
	l.Response = res
}