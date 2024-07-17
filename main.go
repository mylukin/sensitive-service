package main

import (
	"fmt"
	"net/http"
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

var tables = make(map[string]*Table)
var mutex = &sync.Mutex{}

func main() {
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
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
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
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	replacementRune := []rune(to)[0]
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
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
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
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
	}

	words := tbl.Filter.FindAll(text)
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
	mutex.Unlock()

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "table not found"})
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
			resp, err := http.Head(tbl.DictURL)
			if err != nil {
				fmt.Printf("Failed to check dict for table %s: %v\n", tableName, err)
				continue
			}

			lastModified := resp.Header.Get("Last-Modified")
			if lastModified == "" {
				continue
			}

			modifiedTime, err := time.Parse(http.TimeFormat, lastModified)
			if err != nil {
				fmt.Printf("Failed to parse Last-Modified header for table %s: %v\n", tableName, err)
				continue
			}

			if modifiedTime.After(tbl.Modified) {
				filter := sensitive.New()
				err := filter.LoadNetWordDict(tbl.DictURL)
				if err != nil {
					fmt.Printf("Failed to reload dict for table %s: %v\n", tableName, err)
					continue
				}

				tbl.Filter = filter
				tbl.Modified = modifiedTime

				fmt.Printf("Dict for table %s updated\n", tableName)
			}
		}
		mutex.Unlock()
	}
}
