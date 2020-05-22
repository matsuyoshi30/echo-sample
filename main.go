package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
)

func setupDB() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", "/tmp/echo-sample-db.bin")
	if err != nil {
		return nil, err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Todo{}, "todos").SetKeys(true, "id")
	if err = dbmap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}
	return dbmap, nil
}

type Controller struct {
	dbmap *gorp.DbMap
}

// Todo is a struct to hold unit of request and response.
type Todo struct {
	Id        int64     `json:"id" db:"id,primarykey,autoincrement"`
	Taskname  string    `json:"taskname" form:"taskname" db:"taskname,notnull,size:200"`
	Completed bool      `json:"completed" form:"completed" db:"completed"`
	Created   time.Time `json:"created" db:"created,notnull"`
}

// InsertTodo is GET handler to return record.
func (controller *Controller) GetTodo(c echo.Context) error {
	var todo Todo
	// fetch record specified by parameter id
	err := controller.dbmap.SelectOne(&todo, "SELECT * FROM todos WHERE id = $1", c.Param("id"))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Logger().Error("SelectOne: ", err)
			return c.String(http.StatusBadRequest, "SelectOne: "+err.Error())
		}
		return c.String(http.StatusNotFound, "Not Found")
	}
	return c.JSON(http.StatusOK, todo)
}

// ListTodos is GET handler to return records.
func (controller *Controller) ListTodos(c echo.Context) error {
	var todos []Todo
	// fetch last 10 records
	_, err := controller.dbmap.Select(&todos, "SELECT * FROM todos ORDER BY created desc LIMIT 10")
	if err != nil {
		c.Logger().Error("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

// InsertTodo is POST handler to insert record.
func (controller *Controller) InsertTodo(c echo.Context) error {
	var todo Todo
	// bind request to Todo struct
	if err := c.Bind(&todo); err != nil {
		c.Logger().Error("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	// insert record
	if err := controller.dbmap.Insert(&todo); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("inserted todo: %v", todo.Id)
	return c.NoContent(http.StatusCreated)
}

// UpdateTodo is PUT handler to update record
func (controller *Controller) UpdateTodo(c echo.Context) error {
	var todo Todo
	// fetch record specified by parameter id
	err := controller.dbmap.SelectOne(&todo, "SELECT * FROM todos WHERE id = $1", c.Param("id"))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Logger().Error("SelectOne: ", err)
			return c.String(http.StatusBadRequest, "SelectOne: "+err.Error())
		}
		return c.String(http.StatusNotFound, "Not Found")
	}
	// update record
	todo.Completed = !todo.Completed
	if _, err := controller.dbmap.Update(&todo); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("inserted todo: %v", todo.Id)
	return c.NoContent(http.StatusCreated)
}

// RemoveTodo is DELETE handler to delete record
func (controller *Controller) RemoveTodo(c echo.Context) error {
	var todo Todo
	// fetch record specified by parameter id
	err := controller.dbmap.SelectOne(&todo, "SELECT * FROM todos WHERE id = $1", c.Param("id"))
	if err != nil {
		if err != sql.ErrNoRows {
			c.Logger().Error("SelectOne: ", err)
			return c.String(http.StatusBadRequest, "SelectOne: "+err.Error())
		}
		return c.String(http.StatusNotFound, "Not Found")
	}
	// delete record
	if _, err := controller.dbmap.Delete(&todo); err != nil {
		c.Logger().Error("Delete: ", err)
		return c.String(http.StatusBadRequest, "Delete: "+err.Error())
	}
	c.Logger().Infof("deleted todo: %v", todo.Id)
	return c.NoContent(http.StatusCreated)
}

func main() {
	dbmap, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	defer dbmap.Db.Close()
	controller := &Controller{dbmap: dbmap}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/todos/:id", controller.GetTodo)
	e.GET("/api/todos", controller.ListTodos)
	e.POST("/api/todos", controller.InsertTodo)
	e.PUT("/api/todos/:id", controller.UpdateTodo)
	e.DELETE("/api/todos/:id", controller.RemoveTodo)
	e.Static("/", "static/")

	e.Logger.Fatal(e.Start(":1323"))
}
