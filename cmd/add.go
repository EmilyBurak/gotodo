package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the task list",
	Long:  `Add a task to the task list`,
	Run: func(cmd *cobra.Command, args []string) {

		var (
			name        string
			isPomNeeded bool
			confirm     bool
		)

		pomNeeded := ""

		file, err := os.Open("tasks.csv")
		newFile := false
		if err != nil {
			// If the file does not exist, create it
			file, err = os.Create("tasks.csv")
			if err != nil {
				panic(err)
			}
		}

		defer file.Close()

		data, err := csv.NewReader(file).ReadAll()
		if err != nil {
			panic(err)
		}

		// If the file is empty, create a new task list
		if len(data) == 0 {
			newFile = true
			fmt.Println("Creating new task list")
		}

		file.Close() // Close the file after reading

		// Reopen the file in append mode
		file, err = os.OpenFile("tasks.csv", os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)

		if newFile {
			// Write the header if the file is empty
			err = writer.Write([]string{"ID", "Task", "Status", "Deleted", "Working Sessions Completed", "Working Sessions Needed"})
			if err != nil {
				panic(err)
			}
		}

		rows := len(data)

		if rows == 0 {
			rows = 1
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Input task name").
					Prompt("Task name: ").
					Value(&name),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Do you want to set a number of Pomodoro working sessions needed for this task?").
					Affirmative("Yes").
					Negative("No").
					Value(&isPomNeeded),
			),
		)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if isPomNeeded {
			form = huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Input number of Pomodoro working sessions needed").
						Value(&pomNeeded),
				),
			)

			err = form.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		form = huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					TitleFunc(func() string {
						return "Do you want to add this task?" + "\n" + "Task name: " + name
					}, &name).
					Affirmative("Yes").
					Negative("No").
					Value(&confirm),
			),
		)

		err = form.Run()
		if err != nil {
			log.Fatal(err)
		}

		if !confirm {
			fmt.Println("Task not added")
			return
		}

		args = append(append([]string{fmt.Sprintf("%d", rows)}, name), "Pending", "false", "0", pomNeeded)

		err = writer.Write(args)
		if err != nil {
			panic(err)
		}

		writer.Flush()

		if err := writer.Error(); err != nil {
			panic(err)
		}
		fmt.Println("Task ID " + strings.Join(args[:len(args)-4], " ") + " added")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("sessions", "s", "", "Pomodoro working sessions needed to complete the task")
}
