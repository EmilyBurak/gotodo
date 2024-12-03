package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to the task list",
	Long:  `Add a task to the task list`,
	Run: func(cmd *cobra.Command, args []string) {

		pomNeeded := cmd.Flag("number").Value.String()

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
			err = writer.Write([]string{"ID", "Task", "Status", "Deleted", "Pomodoros Completed", "Pomodoros Needed"})
			if err != nil {
				panic(err)
			}
		}

		rows := len(data)

		if rows == 0 {
			rows = 1
		}

		// Append to front of the slice the new task ID
		args = append(append([]string{fmt.Sprintf("%d", rows)}, args[0]), "Pending", "false", "0", pomNeeded)
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
	addCmd.Flags().StringP("number", "n", "", "pomodoros needed to complete the task")
}
