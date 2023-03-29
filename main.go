package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"log"
)

type UrlInfo struct {
	ID int64
	ShortUrl string
	MainUrl string
}

func getUrls(c *gin.Context) {
	db := initializeDB()

	rows, err := db.Query("Select * from urls")
	if (err != nil) {
		c.String(http.StatusInternalServerError, "Something went wrong")
		return
	}
	defer rows.Close()

	var urlInfos []UrlInfo

	for rows.Next() {
		var urlInfo UrlInfo
		if err := rows.Scan(&urlInfo.ID, &urlInfo.ShortUrl, &urlInfo.MainUrl); err != nil {
			fmt.Println("error")
			c.String(http.StatusInternalServerError, "Something went wrong")
			return
		}

		urlInfos = append(urlInfos, urlInfo)

		fmt.Print("%v", urlInfos)

	}

	if err := rows.Err(); err != nil {
        c.String(http.StatusInternalServerError, "last error")
		return
    }

    c.JSON(http.StatusOK, urlInfos)
	
}

func createUrl(c *gin.Context) {
	var urlInfo UrlInfo

	if err := c.BindJSON(&urlInfo); err != nil {
		fmt.Printf("%v", err)
		return
	}

	db := initializeDB()

	_, err := db.Exec("Insert into urls (id, short_url, main_url) values(?, ?, ?)", urlInfo.ID, urlInfo.ShortUrl, urlInfo.MainUrl)

	if (err != nil) {
		fmt.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when inserting in database.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "url created successfully.",
	})
}

func updateUrl(c *gin.Context) {
	var urlInfo UrlInfo

	if err := c.BindJSON(&urlInfo); err != nil {
		fmt.Printf("%v", err)
		return
	}

	db := initializeDB()

	_, err := db.Exec("update urls set short_url = ?, main_url = ? where id = ?", urlInfo.ShortUrl, urlInfo.MainUrl, c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when updating the database.",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}


func deleteUrl(c *gin.Context) {
	db := initializeDB()

	_, err := db.Exec("delete from urls where id = ?", c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred when deleting from database.",
		})

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func showUrl(c *gin.Context) {
	db := initializeDB()

	var urlInfo UrlInfo

	row := db.QueryRow("Select * from urls where id = ?", c.Param("id"))
	
	if err := row.Scan(&urlInfo.ID, &urlInfo.ShortUrl, &urlInfo.MainUrl); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "url not found",
			})
			return
		}
		fmt.Printf("Error occurred %v", err)
		return
	}

	c.IndentedJSON(http.StatusOK, urlInfo)
}

func initializeDB() (db *sql.DB) {
	cfg := mysql.Config{
		User: "root",
		Net: "tcp",
		Addr: "127.0.0.1:3306",
		DBName: "golang",
	}

	var err error

	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
        log.Fatal(err)
    }

	pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")

	return
}

func main() {
	router := gin.Default()

	router.GET("/urls", getUrls)
	router.POST("/urls", createUrl)
	router.PUT("/urls/:id", updateUrl)
	router.DELETE("/urls/:id", deleteUrl)
	router.GET("/urls/:id", showUrl)

	router.Run("localhost:8080")
}

