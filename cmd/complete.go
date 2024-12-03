package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete",
	Short: "Complete a task",
	Long:  `Complete a task by providing the task ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("tasks.csv")
		if err != nil {
			fmt.Println("No tasks to complete")
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

		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		taskFound := false
		// Loop through the records to find the task ID and mark it as done
		for i := 0; ; i++ {
			record, err := csvReader.Read()
			if err == io.EOF { // End of file
				if taskFound {
					break
				}
				// Task ID not found in records, so return
				log.Println("Task ID not found")
				csvWriter.Flush()
				os.Remove("temp.csv")
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			// Convert the task ID to string to compare with the record
			if record[0] == strconv.Itoa(taskID) {
				if record[3] == "true" {
					// Task already deleted=> Do not write the record
					fmt.Println("Task already deleted")
					csvWriter.Flush()
					os.Remove("temp.csv")
					return
				}
				if record[2] == "done" {
					// Task already completed=> Do not write the record
					fmt.Println("Task already completed")
					csvWriter.Flush()
					os.Remove("temp.csv")
					return
				}
				// Mark the task as done and break to write the record
				record[2] = "done"
				taskFound = true
			}
			err = csvWriter.Write(record)
			if err != nil {
				log.Fatal(err)
			}
			csvWriter.Flush()
			if csvWriter.Error() != nil {
				log.Fatal(csvWriter.Error())
			}
		}

		// Rename the temp file to tasks.csv to update the tasks
		err = os.Rename("temp.csv", "tasks.csv")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Task of task ID " + args[0] + " completed")
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
