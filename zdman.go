package main

import (
	"fmt"
	"gozdman/data"
	"gozdman/forms"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/assets", "./vol/assets")
	router.LoadHTMLGlob("vol/templates/*.html")

	router.GET("/", func(c *gin.Context) {
		ip := c.Request.RemoteAddr // attempt detecting localhost or distant server:
		islocal := strings.HasPrefix(ip, "127") || strings.HasPrefix(ip, "[::1]")
		var content = ""
		if islocal {
			content = "Select a menu item hereabove to get started "
		} else {
			content = "Caution: all modifications to the database will be discarded at the end of the fly.io session"
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": content,
		})
	})
	router.GET("/size", dicsize)
	router.GET("/getziform", getziform)
	router.GET("/getpyform", getpyform)
	router.GET("/addziform", addziform)
	router.GET("/selupdate", selupdate)
	router.GET("/updatezi", updatezi)
	router.GET("/getdelete", getdelete)
	router.GET("/showlast", showlast)
	router.GET("/listdic", listdic)
	router.GET("/remove", remove)
	router.POST("/listzi", listzi)
	router.POST("/listpy", listpy)
	router.POST("/addzi", addzi)
	router.POST("/seldelete", seldelete)
	router.PUT("/doupdate/:id", doupdate)
	router.DELETE("delete/:id", delete)

	router.Run(":8080")
}

func dicsize(c *gin.Context) {
	len, time := data.Dicsize()
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("The dictionary presently contains "+len+" entries ; last updated on "+time))
}

func getziform(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Ziform()))
}

func getpyform(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Pyform()))
}

func addziform(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Addziform()))
}

func selupdate(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Selupdate()))
}

func getdelete(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Seldelete()))
}

func showlast(c *gin.Context) {
	c.Data(http.StatusOK, "text/html/json; charset=utf-8", []byte(data.Printlast()))
}

func listdic(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data.Printdiclist()))
}

func remove(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Form canceled."))
}

func listzi(c *gin.Context) {
	var newZi data.Hanzi
	if err := c.Bind(&newZi); err != nil {
		log.Println("Erreur listzi:" + err.Error())
		return
	}
	carac := newZi.Carac
	log.Println("Reçu :" + carac)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data.Listforzi(carac)))
}

func listpy(c *gin.Context) {
	var newPy data.PinYin
	if err := c.Bind(&newPy); err != nil {
		log.Println("Erreur listpy:" + err.Error())
		return
	}
	pinyin := newPy.Pinyin
	log.Println("Reçu :" + pinyin)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data.Listforpy(pinyin)))
}

func updatezi(c *gin.Context) {
	var newId data.Ident
	if err := c.Bind(&newId); err != nil {
		log.Println("Error binding in updatezi:" + err.Error())
		return
	}
	someZi, err := data.Getforid(newId.Id)
	if err != nil {
		log.Println("No record for ident " + strconv.Itoa(newId.Id))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("Error, %v", err)))
		return
	}
	log.Println("bien reçu " + someZi.Pinyin_ton)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Updateziform(someZi)))
}

func addzi(c *gin.Context) {
	var newZi data.DBzi
	if err := c.Bind(&newZi); err != nil {
		log.Println("Erreur addzi:" + err.Error())
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(fmt.Sprintf("Error, %v", err)))
		return
	}
	log.Println("Got " + newZi.Pinyin_ton)
	id, err := data.DBaddzi(newZi)
	if err != nil || id == 0 {
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(fmt.Sprintf("Error, %v", err)))
		log.Println(err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte("Created record n°"+strconv.FormatInt(id, 10)))
}

func doupdate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var zi data.DBidzi
	if err := c.Bind(&zi); err != nil {
		log.Println("Erreur doupdate:" + err.Error())
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data.DBupdatezi(id, zi)))
}

func seldelete(c *gin.Context) {
	var newId data.Ident
	if err := c.Bind(&newId); err != nil {
		log.Println("Erreur updatezi:" + err.Error())
		return
	}
	// check that record wih id newId exists
	zi, err := data.Getforid(newId.Id)
	if err != nil {
		log.Println("No record for id " + strconv.Itoa(newId.Id))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf("Error, %v", err)))
		return
	}
	// display delete confirmation dialog
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(forms.Confdelete(zi)))
}

func delete(c *gin.Context) {
	aid := c.Param("id")
	id, _ := strconv.Atoi(aid)
	// check that record wih id newId exists
	_, err := data.Getforid(id)
	if err != nil {
		log.Println("No record for id " + strconv.Itoa(id))
		c.Data(http.StatusBadRequest, "text/html; charset=utf-8", []byte(fmt.Sprintf("Error, %v", err)))
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(data.DBdelete(aid)))
}
