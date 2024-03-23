package gotodo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []item

func (t *Todos) Add(task string) {
	todo := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now().UTC(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

func (t *Todos) Complete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index")
	}

	ls[index-1].CompletedAt = time.Now().UTC()
	ls[index-1].Done = true

	return nil
}

func (t *Todos) Delete(index int) error {
	ls := *t
	if index <= 0 || index > len(ls) {
		return errors.New("Invalid index")
	}

	*t = append(ls[:index-1], ls[index:]...)

	return nil
}

func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (t *Todos) Print() {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done"},
			{Align: simpletable.AlignCenter, Text: "CreatedAt"},
			{Align: simpletable.AlignCenter, Text: "CompletedAt"},
		},
	}

	var cells [][]*simpletable.Cell
	for index, itm := range *t {
		index++

		task := blue(itm.Task)
		done := blue("No")
		if itm.Done {
			task = green(fmt.Sprintf("\u2705 %s", itm.Task))
			done = green("Yes")
		}

		cells = append(cells, *&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", index)},
			{Text: task},
			{Text: done},
			{Text: itm.CreatedAt.Format(time.RFC822)},
			{Text: itm.CompletedAt.Format(time.RFC822)},
		})
	}

	footerText := red(fmt.Sprintf("You have %d pending TODO's.", t.PendingCount()))
	if t.PendingCount() == 0 {
		footerText = green("You have completed all the TODO's.")
	}
	table.Body = &simpletable.Body{Cells: cells}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: footerText},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

func (t *Todos) PendingCount() int {
	pending := 0
	for _, itm := range *t {
		if !itm.Done {
			pending++
		}
	}

	return pending
}
