package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task from the task list",
	Long:  `Delete a task from the task list by providing the task ID or name.`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("task").Value.String()
		id := cmd.Flag("ID").Value.String()

		if name == "" && id == "0" {
			fmt.Println("Please provide the task ID or name to delete the task")
			return
		}

		file, err := os.Open("tasks.csv")
		if err != nil {
			fmt.Println("No tasks to delete")
			return
		}
		tempFile, err := os.Create("temp.csv")
		if err != nil {
			fmt.Println("No tasks to list")
			return
		}

		// Close the files after the function completes
		defer file.Close()
		defer tempFile.Close()

		csvReader := csv.NewReader(file)
		csvWriter := csv.NewWriter(tempFile)

		taskID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		rows, err := csvReader.ReadAll()
		if err != nil {
			fmt.Println("Error reading task csv file")
			return
		}

		recordTaskNames := make([]string, 0)

		for _, row := range rows {
			if len(row) != 0 {
				if id, err := strconv.Atoi(row[0]); err == nil && id == taskID {
					recordTaskNames = append(recordTaskNames, row[1])
					if row[3] == "true" {
						// Task already deleted=> Do not write the record
						fmt.Println("Task already deleted")
						csvWriter.Flush()
						os.Remove("temp.csv")
						return
					}
					if row[0] == strconv.Itoa(taskID) {
						// Mark as deleted
						row[3] = "true"
					}
				}
				if row[0] != strconv.Itoa(taskID) {
					recordTaskNames = append(recordTaskNames, row[1])
					searchResult := fuzzy.Find(name, recordTaskNames)
					if len(searchResult) != 0 {
						log.Println("Search result: ", searchResult)
						row[3] = "true"
						// clear recordTaskNames to avoid deleting the same task multiple times
						// maybe due to fuzzy search, multiple tasks are found, this will delete all of them, TODO: fix
						recordTaskNames = []string{}
					}
				}
			}
			err = csvWriter.Write(row)
			if err != nil {
				log.Fatal(err)
			}
		}

		csvWriter.Flush()
		if csvWriter.Error() != nil {
			log.Fatal(csvWriter.Error())
		}

		// Rename the temp file to tasks.csv to update the tasks
		err = os.Rename("temp.csv", "tasks.csv")
		if err != nil {
			log.Fatal(err)
		}

		if id != "0" {
			fmt.Println("Task of task ID " + id + " deleted")
		} else {
			fmt.Println("Task " + name + " deleted")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().IntP("ID", "i", 0, "ID of the task")
	deleteCmd.Flags().StringP("task", "t", "", "Name of task to delete")
}
