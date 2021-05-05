package main

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type CsvData struct {
	ID   int `gorm:"primary_key"`
	Num  int
	Data string
}

func csvInit(io string) {
	db, err := gorm.Open("sqlite3", "/tmp/csvgorm_"+io+".db")
	if err != nil {
		panic("failed connect databases\n")
	}
	defer db.Close()
	db.LogMode(true)
	//db.AutoMigrate(&CsvData{})
}

func dbinput(num int, data string, io string) {
	db, err := gorm.Open("sqlite3", "/tmp/csvgorm_"+io+".db")
	if err != nil {
		panic("failed connect databases\n")
	}
	defer db.Close()
	db.Create(&CsvData{Num: num, Data: data})
}

func dbinput_struct(csvData string) {
	//reader := csv.NewReader(strings.NewReader(csvData))
	reader := csv.NewReader(strings.NewReader(csvData))
	for i := 1; ; i++ {

		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		var v string
		for _, v = range line {
			dbinput(i, v, "in")
		}
	}
}

func csv_all(io string) []CsvData {
	db, err := gorm.Open("sqlite3", "/tmp/csvgorm_"+io+".db")
	if err != nil {
		panic("failed connect databases\n")
	}
	defer db.Close()
	var csvData []CsvData
	db.Find(&csvData)
	return csvData
}

func csvDelete(io string) {
	db, err := gorm.Open("sqlite3", "/tmp/csvgorm_"+io+".db")
	if err != nil {
		panic("failed connect databases(csvDelete)\n")
	}
	db.Debug().Delete(CsvData{ID: 0})
	//defer db.Close()
}

func main() {
	route := gin.Default()
	route.LoadHTMLGlob("templates/*")
	csvInit("in")

	route.GET("/", func(c *gin.Context) {
		csvData := csv_all("in")

		c.HTML(http.StatusOK, "index.html", gin.H{
			"csvData": csvData,
		})
	})

	route.POST("/w", func(c *gin.Context) {
		var csvData string
		csvDelete("in")
		csvData = c.PostForm("csv")
		dbinput_struct(csvData)

		c.Redirect(302, "/")
	})

	route.POST("/big", func(c *gin.Context) {
		csvData := csv_all("in")
		csvData_c := csv_all("in")
		for num, datas := range csvData_c {
			data := datas.Data
			csvData_c[num].Data = strings.ToUpper(data)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"csvData":   csvData,
			"csvData_c": csvData_c,
		})
	})

	route.POST("/small", func(c *gin.Context) {
		csvData := csv_all("in")
		csvData_c := csv_all("in")
		for num, datas := range csvData_c {
			data := datas.Data
			csvData_c[num].Data = strings.ToLower(data)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"csvData":   csvData,
			"csvData_c": csvData_c,
		})
	})

	route.POST("/grep", func(c *gin.Context) {
		csvData := csv_all("in")
		csvData_c := csv_all("in")
		word := c.PostForm("grep")
		for num, datas := range csvData_c {
			data := datas.Data
			if strings.Contains(data, word) {
			} else {
				csvData_c[num].Data = ""
			}
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"csvData":   csvData,
			"csvData_c": csvData_c,
		})
	})

	route.POST("/sed", func(c *gin.Context) {
		csvData := csv_all("in")
		csvData_c := csv_all("in")
		before := c.PostForm("before")
		after := c.PostForm("after")
		for num, datas := range csvData_c {
			data := datas.Data
			csvData_c[num].Data = strings.Replace(data, before, after, -1)
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"csvData":   csvData,
			"csvData_c": csvData_c,
		})
	})

	route.POST("/delete", func(c *gin.Context) {
		csvDelete("in")
		c.Redirect(302, "/")
	})

	route.Run()
}
