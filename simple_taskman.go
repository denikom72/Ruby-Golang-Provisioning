package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Task represents a task with a title, description, and completion status.
type Task struct {
	Title       string
	Description string
	Completed   bool
}

// TaskList manages a list of tasks.
type TaskList struct {
	Tasks []*Task
}

// NewTask creates a new task.
func NewTask(title, description string) *Task {
	return &Task{
		Title:       title,
		Description: description,
	}
}

// NewTaskList creates a new task list.
func NewTaskList() *TaskList {
	return &TaskList{
		Tasks: []*Task{},
	}
}

// AddTask adds a task to the list.
func (tl *TaskList) AddTask(title, description string) {
	task := NewTask(title, description)
	tl.Tasks = append(tl.Tasks, task)
}

// ListTasks lists all tasks.
func (tl *TaskList) ListTasks() {
	for i, task := range tl.Tasks {
		status := " "
		if task.Completed {
			status = "X"
		}
		fmt.Printf("%d. [%s] %s - %s\n", i+1, status, task.Title, task.Description)
	}
}

// MarkTaskCompleted marks a task as completed.
func (tl *TaskList) MarkTaskCompleted(index int) {
	if index >= 0 && index < len(tl.Tasks) {
		tl.Tasks[index].Completed = true
	}
}

// EditTaskDescription edits the description of a task.
func (tl *TaskList) EditTaskDescription(index int, description string) {
	if index >= 0 && index < len(tl.Tasks) {
		tl.Tasks[index].Description = description
	}
}

// DeleteTask deletes a task.
func (tl *TaskList) DeleteTask(index int) {
	if index >= 0 && index < len(tl.Tasks) {
		tl.Tasks = append(tl.Tasks[:index], tl.Tasks[index+1:]...)
	}
}

// SaveTasksToFile saves tasks to a text file.
func (tl *TaskList) SaveTasksToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, task := range tl.Tasks {
		completed := "false"
		if task.Completed {
			completed = "true"
		}
		_, err := fmt.Fprintf(writer, "%s,%s,%s\n", task.Title, task.Description, completed)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadTasksFromFile loads tasks from a text file.
func (tl *TaskList) LoadTasksFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		title := parts[0]
		description := parts[1]
		completed := parts[2] == "true"
		task := NewTask(title, description)
		task.Completed = completed
		tl.Tasks = append(tl.Tasks, task)
	}

	return nil
}

func main() {
	taskList := NewTaskList()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Task Manager Menu:")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Mark Task as Completed")
		fmt.Println("4. Edit Task Description")
		fmt.Println("5. Delete Task")
		fmt.Println("6. Save Tasks to File")
		fmt.Println("7. Load Tasks from File")
		fmt.Println("8. Exit")
		fmt.Print("Choose an option: ")

		var choice int
		if _, err := fmt.Scan(&choice); err != nil {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Enter task title: ")
			scanner.Scan()
			title := scanner.Text()
			fmt.Print("Enter task description: ")
			scanner.Scan()
			description := scanner.Text()
			taskList.AddTask(title, description)
		case 2:
			taskList.ListTasks()
		case 3:
			fmt.Print("Enter the task number to mark as completed: ")
			var index int
			if _, err := fmt.Scan(&index); err == nil {
				taskList.MarkTaskCompleted(index - 1)
			}
		case 4:
			fmt.Print("Enter the task number to edit description: ")
			var index int
			if _, err := fmt.Scan(&index); err == nil {
				fmt.Print("Enter new description: ")
				scanner.Scan()
				description := scanner.Text()
				taskList.EditTaskDescription(index - 1, description)
			}
		case 5:
			fmt.Print("Enter the task number to delete: ")
			var index int
			if _, err := fmt.Scan(&index); err == nil {
				taskList.DeleteTask(index - 1)
			}
		case 6:
			fmt.Print("Enter the filename to save tasks to: ")
			scanner.Scan()
			filename := scanner.Text()
			if err := taskList.SaveTasksToFile(filename); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 7:
			fmt.Print("Enter the filename to load tasks from: ")
			scanner.Scan()
			filename := scanner.Text()
			if err := taskList.LoadTasksFromFile(filename); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case 8:
			os.Exit(0)
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}
