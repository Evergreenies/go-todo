package main

import (
	"flag"
	"fmt"
	"os"

	gotodo "github.com/evergreenies/go-todo"
)

const (
	todoFile = ".todo.json"
)

func main() {
	add := flag.Bool("add", false, "Add a new todo.")
	message := flag.String("message", "", "Todo title")
	complete := flag.Int("complete", 0, "Mark a todo as completed.")
	del := flag.Int("delete", 0, "Delete a todo.")
	listTodos := flag.Bool("list", false, "Listing all todos.")
	flag.Parse()

	todos := &gotodo.Todos{}
	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		if *message == "" {
			fmt.Fprintln(os.Stdout, "You must provide some message to add as TODO title.")
			os.Exit(1)
		}
		todos.Add(*message)
		err := todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	case *listTodos:
		todos.Print()
	default:
		fmt.Fprintln(os.Stdout, "Invalid command")
		os.Exit(0)
	}
}
