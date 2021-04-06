package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type todoModel struct {
	gorm.Model  // 包含四个字段 ID、CreatedAt、UpdatedAt、DeletedAt
	Title string `json:"title"`
	Completed int `json:"completed"`
}

// 处理返回的字段
type transformedTodo struct {
	ID uint `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

type SelectModel struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type SelectReturn struct {
	Options []SelectModel `json:"options"`
}

type TableModel struct {
	Id int `json:"id"`
	Browser string `json:"browser"`
	Engine string `json:"engine"`
	Grade string `json:"grade"`
	Platform string `json:"platform"`
	Version string `json:"version"`
	Idx string `json:"idx"`
	Entity string `json:"entity"`
}
type TableReturn struct {
	Count int `json:"count"`
	Rows []TableModel `json:"rows"`
}

type DetailModel struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Gender int `json:"gender"`
}

type CrudReturn struct {
	Count int `json:"count"`
	Items []TableModel `json:"items"`
}
var GradeArr = [5]string{"A", "B", "C", "D", "X"}


var db *gorm.DB

func init() {
	initDb()
}

func initDb()  {
	var err error
	db, err = gorm.Open("mysql", "root:13396095889@/gindemo?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logs.Error(err)
		panic("failed to connect database!")
	}
	// 迁移schema
	db.AutoMigrate(&todoModel{})
}

func main() {
	g := gin.Default()

	v1 := g.Group("/api/v1/todos")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", FetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	g.GET("/amis/mock/getOptions", GetSelectData)
	g.GET("/amis/mock/getTableData", GetTableData)
	g.GET("/amis/mock/getDetailData", GetDetailData)
	g.GET("/amis/mock/getCrudData", GetCrudData)
	g.GET("/amis/mock/getById", GetById)

	g.Run(":8099")

	//g.GET("/", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "index",
	//	})
	//})
	//g.GET("/test", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})
}

type Account struct {
	AccountName string `json:"AccountName"`
}
func GetById(c *gin.Context) {
	var data = Account{}
	data.AccountName = "xxxxxxxx"
	
	// 解决跨域问题
	c.Header("Access-Control-Allow-Origin", "http://localhost:3333")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg": "",
		"data": data,
	})
}
func GetSelectData(c *gin.Context) {
	var options []SelectModel
	var data = &SelectReturn{}

	options = append(options, SelectModel{
		Label: "A",
		Value: "a",
	})
	options = append(options, SelectModel{
		Label: "B",
		Value: "b",
	})
	options = append(options, SelectModel{
		Label: "C",
		Value: "c",
	})
	data.Options = options

	// 解决跨域问题
	c.Header("Access-Control-Allow-Origin", "http://localhost:3333")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg": "",
		"data": data,
	})
}
func GetTableData(c *gin.Context) {
	var data = &TableReturn{}
	var i int
	for i = 1; i <= 20; i++ {
		data.Rows = append(data.Rows, TableModel{
			Id:       i,
			Browser:  fmt.Sprintf("Browser-%d", i),
			Engine:   fmt.Sprintf("Engine-%d", i),
			Grade:    fmt.Sprintf("Grade-%d", i),
			Platform: fmt.Sprintf("Platform-%d", i),
			Version:  fmt.Sprintf("Version-%d", i),
		})
	}
	data.Count = i - 1

	// 解决跨域问题
	c.Header("Access-Control-Allow-Origin", "http://localhost:3333")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg": "",
		"data": data,
	})
}
func GetDetailData(c *gin.Context) {
	var data = DetailModel{
		Id:     3,
		Name:   "Stefan",
		Email:  "252556310@qq.com",
		Gender: 2,
	}

	// 解决跨域问题
	c.Header("Access-Control-Allow-Origin", "http://localhost:3333")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg": "",
		"data": data,
	})
}
func GetCrudData(c *gin.Context) {
	var data = CrudReturn{}
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	for i := (page - 1) * perPage + 1; i <= page * perPage; i++ {
		table := TableModel{
			Id:       i,
			Browser:  fmt.Sprintf("Browser-%d", i),
			Engine:   fmt.Sprintf("Engine-%d", i),
			Grade:    GradeArr[i % 5],
			Platform: fmt.Sprintf("Platform-%d", i),
			Version:  fmt.Sprintf("Version-%d", i),
		}

		if i % 2 == 0 {
			table.Idx = "A7LPTA0002XJ"
			table.Entity = "Account"
		} else {
			table.Idx = "A7LPTA014V6V"
			table.Entity = "Exhibition"
		}
		data.Items = append(data.Items, table)
	}
	data.Count = 100

	// 解决跨域问题
	c.Header("Access-Control-Allow-Origin", "http://localhost:3333")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg": "",
		"data": data,
	})
}

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "title is null",
		})
		return
	}
	todo := &todoModel{
		Title:     title,
		Completed: completed,
	}

	db.Save(todo)
	c.JSON(http.StatusCreated, gin.H{
		"status": http.StatusCreated,
		"message": "Todo item created successfully!",
		"resourceId": todo.ID,
	})
}
func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo
	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		}

		_todos = append(_todos, transformedTodo{
			ID:        item.ID,
			Title:     item.Title,
			Completed: completed,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _todos,
	})
}
func FetchSingleTodo(c *gin.Context) {
	var todo todoModel
	todoID, _ := strconv.Atoi(c.Param("id"))
	if todoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "id is wrong!",
		})
		return
	}

	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	}
	_todo := transformedTodo{
		ID:        todo.ID,
		Title:     todo.Title,
		Completed: completed,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data": _todo,
	})
}
func updateTodo(c *gin.Context) {
	var todo todoModel
	todoID, _ := strconv.Atoi(c.Param("id"))
	title := c.PostForm("title")
	if todoID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "id is wrong!",
		})
		return
	}
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"message": "title is null!",
		})
		return
	}

	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}

	db.Model(&todo).Update("title", title)
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Todo updated successfully!",
	})
}
func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"message": "No todo found!",
		})
		return
	}
	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"message": "Todo deleted successfully!",
	})
}