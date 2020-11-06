package firstservice_test

import (
	"fmt"
	. "go-to-do/firstservice"
	"testing"
)


func TestFunc_Marshall(t *testing.T) {
	var tasks Tasks

	tasks.Tasks = append(tasks.Tasks, *NewTask("task_1", "22"))
	tasks.Tasks = append(tasks.Tasks, *NewTask("task_2", "22"))

	b, err := tasks.Marshall()
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	fmt.Println(string(b))
}
