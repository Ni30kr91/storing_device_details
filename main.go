package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DEVICES struct {
	ID                int       `json:"id"`
	TYPE              string    `json:"type"`
	BROWSER           string    `json:"browser"`
	BROWSER_VERSION   string    `json:"browser_version"`
	CREATED_AT        time.Time `json:"created_at"`
	SCREEN_RESOLUTION string    `json:"screen_resolution"`
}

var (
	DB *sql.DB
)

func main() {
	createDBConnection()
	defer DB.Close()
	r := gin.Default()
	r.Use(CORSMiddleware())
	setupRouters(r)
	r.Run()
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func setupRouters(r *gin.Engine) {

	r.POST("/device", PostDevices)
	r.GET("/get/device", GetDevice)

}
func PostDevices(c *gin.Context) {
	reqBody := DEVICES{}
	err := c.Bind(&reqBody)
	if err != nil {
		res := gin.H{
			"error": "Invalid Request Body",
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, res)
		return

	}
	reqBody.CREATED_AT = time.Now()

	fmt.Println(reqBody)
	res, err := DB.Exec(`INSERT INTO "device_details" ("type", "browser", "browser_version", "created_at", "screen_resolution")
	VALUES ($1, $2, $3, $4, $5)`, reqBody.TYPE, reqBody.BROWSER, reqBody.BROWSER_VERSION, reqBody.CREATED_AT, reqBody.SCREEN_RESOLUTION)
	if err != nil {
		fmt.Println("err inserting data: ", err)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return

	}
	lastInsID, err := res.LastInsertId()
	fmt.Println("errr: ", err)
	reqBody.ID = int(lastInsID)
	fmt.Println("res: ", lastInsID)
	c.JSON(http.StatusOK, reqBody)
	c.Writer.Header().Set("Content-Type", "application/jason")
}

func GetDevice(c *gin.Context) {

	rows, err := DB.Query("SELECT id, type, browser, browser_version, created_at, screen_resolution FROM device_details order by id desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	devices := []DEVICES{}
	for rows.Next() {
		device := DEVICES{}

		err := rows.Scan(&device.ID, &device.TYPE, &device.BROWSER, &device.BROWSER_VERSION, &device.CREATED_AT, &device.SCREEN_RESOLUTION)

		if err != nil {
			fmt.Println(err)
		}
		devices = append(devices, device)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	res := gin.H{
		"data": devices,
	}

	c.JSON(http.StatusOK, res)
}
