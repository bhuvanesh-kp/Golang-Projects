package main

import (
	"fmt"
	"net/http"
	"strconv"
	"todo/handler"

	"github.com/gin-gonic/gin"
)

var (
	todoList *handler.Todos = handler.NewTodo()
)

func handleDelete(ctx *gin.Context) {
	postId := ctx.Param("id")

	id, err := strconv.Atoi(postId)
	if err != nil {
		ctx.JSON(http.StatusForbidden, map[string]any{
			"message": "Invalid postId",
		})
		return
	}

	response, err := todoList.HandleDelete(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]any{
			"message": err,
		})
	}

	ctx.JSON(http.StatusAccepted, response)
}

func main() {
	router := gin.Default()

	fmt.Println("Server starting at port 8000")

	// GET request to show all avalible todo lists
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, todoList.HandleGetAll())
	})

	// POST request to add a new todo list to the existing list
	router.POST("/add", func(ctx *gin.Context) {
		var req handler.Request
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": "Something wrong in posting a new request",
			})
			return
		}

		response := todoList.HandlePost(req)
		ctx.JSON(http.StatusAccepted, response)
	})

	// PUT request to update an exising Todo in the todo list
	router.PUT("/:id", func(ctx *gin.Context) {
		postId := ctx.Param("id")
		var req handler.Request

		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": "Something wrong in updating the post",
			})
			return
		}

		id, err := strconv.Atoi(postId)
		if err != nil {
			ctx.JSON(http.StatusForbidden, map[string]any{
				"message": "Invalid postId",
			})
			return
		}

		response, err := todoList.HandleUpdate(req, id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, map[string]any{
				"message": err,
			})
		}

		ctx.JSON(http.StatusAccepted, response)
	})

	//Delete an existing todo from the todos list
	router.DELETE("/:id", handleDelete)

	router.Run(":8000")
}
