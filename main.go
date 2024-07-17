package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mylukin/sensitive"
)

type Table struct {
	DictURL  string
	Filter   *sensitive.Filter
	Modified time.Time
}

var (
	tables       = make(map[string]*Table)
	mutex        = &sync.Mutex{}
	dataFilePath = "tables_data.json"
)

const globalTableName = "*"

func main() {
	loadTablesFromFile()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/add", addHandler)
	e.GET("/del", delHandler)
	e.GET("/filter", filterHandler)
	e.GET("/replace", replaceHandler)
	e.GET("/findin", findInHandler)
	e.GET("/findall", findAllHandler)
	e.GET("/validate", validateHandler)

	go monitorDictUpdates()
	go saveTablesToFilePeriodically()

	e.Logger.Fatal(e.Start(":8080"))
}

func addHandler(c echo.Context) error {
	table := c.QueryParam("table")
	dictURL := c.QueryParam("dict")

	if table == "" || dictURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table and dict parameters are required"})
	}

	mutex.Lock()
	defer mutex.Unlock()

	filter := sensitive.New()
	err := filter.LoadNetWordDict(dictURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to load dict"})
	}

	tables[table] = &Table{
		DictURL:  dictURL,
		Filter:   filter,
		Modified: time.Now(),
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "table added"})
}

func delHandler(c echo.Context) error {
	table := c.QueryParam("table")

	if table == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table parameter is required"})
	}

	mutex.Lock()
	defer mutex.Unlock()

	delete(tables, table)

	return c.JSON(http.StatusOK, map[string]string{"message": "table deleted"})
}

func filterHandler(c echo.Context) error {
	table := c.QueryParam("table")
	text := c.QueryParam("text")

	if table == "" || text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table and text parameters are required"})
	}

	mutex.Lock()
	tbl, exists := tables[table]
	globalTbl, globalExists := tables[globalTableName]
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	if globalExists {
		text = globalTbl.Filter.Filter(text)
	}
	filtered := tbl.Filter.Filter(text)
	return c.JSON(http.StatusOK, map[string]string{"filtered": filtered})
}

func replaceHandler(c echo.Context) error {
	table := c.QueryParam("table")
	text := c.QueryParam("text")
	to := c.QueryParam("to")

	if table == "" || text == "" || to == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table, text, and to parameters are required"})
	}

	mutex.Lock()
	tbl, exists := tables[table]
	globalTbl, globalExists := tables[globalTableName]
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	replacementRune := []rune(to)[0]
	if globalExists {
		text = globalTbl.Filter.Replace(text, replacementRune)
	}
	replaced := tbl.Filter.Replace(text, replacementRune)
	return c.JSON(http.StatusOK, map[string]string{"replaced": replaced})
}

func findInHandler(c echo.Context) error {
	table := c.QueryParam("table")
	text := c.QueryParam("text")

	if table == "" || text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table and text parameters are required"})
	}

	mutex.Lock()
	tbl, exists := tables[table]
	globalTbl, globalExists := tables[globalTableName]
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	if globalExists {
		found, word := globalTbl.Filter.FindIn(text)
		if found {
			return c.JSON(http.StatusOK, map[string]interface{}{"found": found, "word": word})
		}
	}
	found, word := tbl.Filter.FindIn(text)
	return c.JSON(http.StatusOK, map[string]interface{}{"found": found, "word": word})
}

func findAllHandler(c echo.Context) error {
	table := c.QueryParam("table")
	text := c.QueryParam("text")

	if table == "" || text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table and text parameters are required"})
	}

	mutex.Lock()
	tbl, exists := tables[table]
	globalTbl, globalExists := tables[globalTableName]
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	var words []string
	if globalExists {
		words = globalTbl.Filter.FindAll(text)
	}
	words = append(words, tbl.Filter.FindAll(text)...)
	return c.JSON(http.StatusOK, map[string]interface{}{"words": words})
}

func validateHandler(c echo.Context) error {
	table := c.QueryParam("table")
	text := c.QueryParam("text")

	if table == "" || text == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "table and text parameters are required"})
	}

	mutex.Lock()
	tbl, exists := tables[table]
	globalTbl, globalExists := tables[globalTableName]
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	if globalExists {
		valid, word := globalTbl.Filter.Validate(text)
		if !valid {
			return c.JSON(http.StatusOK, map[string]interface{}{"valid": valid, "word": word})
		}
	}
	valid, word := tbl.Filter.Validate(text)
	return c.JSON(http.StatusOK, map[string]interface{}{"valid": valid, "word": word})
}

func monitorDictUpdates() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mutex.Lock()
		for tableName, tbl := range tables {
			checkAndUpdateTableDict(tbl, tableName)
		}
		mutex.Unlock()
	}
}

func checkAndUpdateTableDict(tbl *Table, tableName string) {
	if tbl.DictURL == "" {
		return
	}

	resp, err := http.Head(tbl.DictURL)
	if err != nil {
		fmt.Printf("Failed to check dict for table %s: %v\n", tableName, err)
		return
	}

	lastModified := resp.Header.Get("Last-Modified")
	if lastModified == "" {
		return
	}

	modifiedTime, err := time.Parse(http.TimeFormat, lastModified)
	if err != nil {
		fmt.Printf("Failed to parse Last-Modified header for table %s: %v\n", tableName, err)
		return
	}

	if modifiedTime.After(tbl.Modified) {
		filter := sensitive.New()
		err := filter.LoadNetWordDict(tbl.DictURL)
		if err != nil {
			fmt.Printf("Failed to reload dict for table %s: %v\n", tableName, err)
			return
		}

		tbl.Filter = filter
		tbl.Modified = modifiedTime

		fmt.Printf("Dict for table %s updated\n", tableName)
	}
}

func saveTablesToFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	data := make(map[string]struct {
		DictURL  string
		Modified time.Time
	})

	for tableName, tbl := range tables {
		data[tableName] = struct {
			DictURL  string
			Modified time.Time
		}{
			DictURL:  tbl.DictURL,
			Modified: tbl.Modified,
		}
	}

	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataFilePath, fileData, 0644)
}

func loadTablesFromFile() error {
	mutex.Lock()
	defer mutex.Unlock()

	fileData, err := ioutil.ReadFile(dataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	data := make(map[string]struct {
		DictURL  string
		Modified time.Time
	})

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return err
	}

	for tableName, tblData := range data {
		filter := sensitive.New()
		err := filter.LoadNetWordDict(tblData.DictURL)
		if err != nil {
			fmt.Printf("Failed to load dict for table %s: %v\n", tableName, err)
			continue
		}

		tables[tableName] = &Table{
			DictURL:  tblData.DictURL,
			Filter:   filter,
			Modified: tblData.Modified,
		}
	}

	return nil
}

func saveTablesToFilePeriodically() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := saveTablesToFile()
		if err != nil {
			fmt.Printf("Failed to save tables to file: %v\n", err)
		}
	}
}
