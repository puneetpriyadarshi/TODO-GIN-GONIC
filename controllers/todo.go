package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	orm "github.com/go-pg/pg/orm"
	guuid "github.com/google/uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	TITLE     string    `json:"title"`
	BODY      string    `json:"body"`
	COMPLETED string    `json:"completed"`
	CreatedAt time.Time `json : "updatedAt"`
}

func CreateTodoTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	createError := db.CreateTable(&Todo{}, opts)
	if createError != nil {
		log.Printf("error while creating a table, Reason : %v\n", createError)
		return createError
	}
	log.Printf("todo table created")
	return nil
}

var dbConnect *pg.DB

func InitiateDB(db *pg.DB) {
	dbConnect = db
}
func GetAllTodos(c *gin.Context) {
	var todos []Todo
	err := dbConnect.Model(&todos).Select()

	if err != nil {
		log.Printf("Error while getting all todos, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Todos",
		"data":    todos,
	})
	return
}

func CreateTodo(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	TITLE := todo.TITLE
	BODY := todo.BODY
	COMPLETED := todo.COMPLETED
	id := guuid.New().String()

	insertError := dbConnect.Insert(&Todo{
		ID:        id,
		TITLE:     TITLE,
		BODY:      BODY,
		COMPLETED: COMPLETED,
		CreatedAt: time.Now(),
	})
	if insertError != nil {
		log.Printf("Error while inserting new todo into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Todo created Successfully",
	})
	return
}
