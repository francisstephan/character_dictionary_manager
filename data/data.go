package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// _ "github.com/mattn/go-sqlite3"

type Hanzi struct {
	Carac string `form:"carac" binding:"required,len=1"`
}
type PinYin struct {
	Pinyin string `form:"pinyin"`
}
type Ident struct {
	Id int `form:"id"`
}

type DBzi struct {
	Pinyin_ton string `form:"pinyin_ton" binding:"required"`
	Unicode    string `form:"unicode" binding:"required,len=4"`
	Sens       string `form:"sens" binding:"required"`
}

type DBidzi struct {
	Id         int    `form:"id"`
	Pinyin_ton string `form:"pinyin_ton" binding:"required"`
	Unicode    string `form:"unicode" binding:"required,len=4"`
	Sens       string `form:"sens" binding:"required"`
}

type Zi struct {
	Id         int
	Pinyin_ton string
	Unicode    string
	Hanzi      string
	Sens       string
}

type Dico []Zi

func litdic(where string) (dic Dico) { // read dictionary with WHERE clause
	dic = make(Dico, 0, 10) // initialize dic with size 0, capacity 10
	db, err := sql.Open("sqlite3", "vol/zidian.db")
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	defer db.Close()

	query := "SELECT id, pinyin_ton, unicode, sens FROM pyhz"
	if where != "" {
		query += " WHERE " + where
	}
	// log.Println("requête = ", query)

	rows, err := db.Query(query)
	if err != nil {
		log.Println("erreur=", err.Error())
		return nil
	}
	var tempZi Zi
	for rows.Next() {
		rows.Scan(&tempZi.Id, &tempZi.Pinyin_ton, &tempZi.Unicode, &tempZi.Sens)
		intrune, _ := strconv.ParseInt(tempZi.Unicode, 16, 64)
		tempZi.Hanzi = string(rune(intrune))
		dic = append(dic, tempZi)
	}
	return dic
}

func Printzi(zi Zi) string {
	return fmt.Sprintf(
		"<tr><td>%5d</td><td>%7s</td><td>%4s</td><td>%s</td><td>%s</td></tr>",
		zi.Id, zi.Pinyin_ton, zi.Unicode, zi.Hanzi, zi.Sens)
}

func printdic(sousdic Dico, py string) string { // prepare browser display of partial dictionary
	if len(sousdic) == 0 {
		return "No result for " + py + " query"
	}
	retour := "<p>Result for query '" + py + "' :</p>" +
		"<table><tr><td>  Id   </td><td> Pinyin </td><td> Unicode </td><td> Character </td><td> Translation</td></tr>"
	for _, zi := range sousdic {
		retour += Printzi(zi)
	}
	return retour + "</table>"
}

func Printlast() string { // print last entry in dictionary
	dic := litdic("")
	zi := dic[len(dic)-1]
	return "Last entry: <table>" + Printzi(zi) + "</table>"
}

func Printdiclist() string {
	dic := litdic("")
	return printdic(dic, "whole dictionary")
}

func Dicsize() (string, string) { // get dictionary size = number of entries
	dic := litdic("")
	if dic != nil {
		file, _ := os.Stat("vol/zidian.db")
		time := file.ModTime()
		return strconv.Itoa(len(dic)), time.String()[0:10]
	}
	return "zero", "err"
}

func Listforzi(zi string) string {
	r := []rune(zi)
	where := "unicode='" + fmt.Sprintf("%X", r[0]) + "' ORDER BY pinyin_ton "
	dic := litdic(where)
	return (printdic(dic, zi))
}

func Listforpy(py string) string {
	last := py[len(py)-1]
	numeric := "01234"
	var where string
	if !strings.Contains(numeric, string(last)) {
		where = "pinyin_ton='" + py + "0' OR pinyin_ton='" + py + "1' OR pinyin_ton='" + py + "2' OR pinyin_ton='" + py +
			"3' OR pinyin_ton='" + py + "4' ORDER BY pinyin_ton , unicode"
	} else {
		where = "pinyin_ton='" + py + "' ORDER BY unicode"
	}
	dic := litdic(where)
	return (printdic(dic, py))
}

func Getforid(id int) (Zi, error) {
	where := "id='" + strconv.Itoa(id) + "'"
	tempZi := new(Zi)
	dic := litdic(where)
	if len(dic) == 0 {
		return *tempZi, fmt.Errorf("Record " + strconv.Itoa(id) + " does not exist")
	}
	return dic[0], nil
}

func DBaddzi(newZi DBzi) (int64, error) {
	db, err := sql.Open("sqlite3", "vol/zidian.db")
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	defer db.Close()
	var i int32
	var p, u, s string
	query := "SELECT * FROM pyhz WHERE pinyin_ton='" + newZi.Pinyin_ton + "' AND unicode='" + newZi.Unicode + "';"
	errDeja := db.QueryRow(query).Scan(&i, &p, &u, &s)
	if errDeja != nil && errDeja != sql.ErrNoRows {
		return 0, errDeja
	}
	if errDeja == nil {
		return 0, fmt.Errorf("Record " + newZi.Pinyin_ton + " , " + newZi.Unicode + " already exists")
	}

	result, err := db.Exec("INSERT INTO pyhz (pinyin_ton, unicode, sens) VALUES (?, ?, ?)", newZi.Pinyin_ton, newZi.Unicode, newZi.Sens)

	if err != nil {
		log.Println("erreur=", err.Error())
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddZ i: %v", err)
	}
	// Return the new album's ID.
	return id, nil
}

func DBupdatezi(id int, zi DBidzi) string {
	log.Printf("Update entry n° %d", id)
	db, err := sql.Open("sqlite3", "vol/zidian.db")
	var mess string
	if err != nil {
		mess = err.Error()
		log.Println(err.Error())
		return mess
	}
	defer db.Close()
	_, err2 := db.Exec("UPDATE pyhz SET pinyin_ton = ?, unicode = ?, sens = ? WHERE id = ?",
		zi.Pinyin_ton, zi.Unicode, zi.Sens, id)
	if err2 != nil {
		mess = err2.Error()
		log.Println(err.Error())
		return mess
	}
	return "Update success for entry id " + strconv.Itoa(id)
}

func DBdelete(aid string) string {
	id, _ := strconv.ParseInt(aid, 10, 32)
	db, err := sql.Open("sqlite3", "vol/zidian.db")
	var mess string
	if err != nil {
		mess = err.Error()
		log.Println(err.Error())
		return mess
	}
	defer db.Close()
	_, err2 := db.Exec("DELETE FROM pyhz WHERE id= ? ", id)
	if err2 != nil {
		mess = err2.Error()
		log.Println(err.Error())
		return mess
	}
	return "delete success"
}

func init() { // initialize the data package with a test on database availability
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(version)
}
