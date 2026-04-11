package handler

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Todos struct {
	TodoList map[int]Response
}

type Request struct{
	Title string `json:"title"`
	Description string `json:"description"`
}

type Response struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Time time.Time `json:"time"`
}

func NewTodo() *Todos {
	return &Todos{
		TodoList: make(map[int]Response),
	}
}

func (t *Todos) HandleGetAll() map[int]Response {
	return t.TodoList
}

func (t *Todos) HandlePost(req Request) Response{
	newId := rand.Intn(1000000)
	newTodo := Response{
		Id: newId,
		Title: req.Title,
		Description: req.Description,
		Time: time.Now(),
	}

	t.TodoList[newId] = newTodo
	return newTodo
}

func(t *Todos) HandleUpdate(req Request,id int) (Response, error){
	_, ok := t.TodoList[id]
	if !ok{
		fmt.Println("Id not found")
		return Response{}, errors.New("No Todo with ID found")
	} 

	updatedResponse := Response{
		Id: id,
		Title: req.Title,
		Description: req.Description,
		Time: time.Now(),
	}

	t.TodoList[id] = updatedResponse

	return updatedResponse, nil
}

func(t *Todos) HandleDelete(ctx context.Context) (Response, error){
	fmt.Print(ctx.Value("id"))
	return Response{}, nil
	id := 10
	res, ok := t.TodoList[id]
	if !ok{
		fmt.Println("Id not found")
		return Response{}, errors.New("No Todo with ID found")
	}

	delete(t.TodoList, id)
	return res, nil
}