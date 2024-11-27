/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
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

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a task from the task list",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		taskID, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}

		readComplete := false

		// Loop through the records to find the task ID and mark it as done
		for i := 0; ; i++ {
			record, err := csvReader.Read()
			if err == io.EOF {
				if readComplete {
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
			if record[0] == strconv.Itoa(taskID) {
				if record[3] == "true" {
					// Task already deleted=> Do not write the record
					fmt.Println("Task already deleted")
					csvWriter.Flush()
					os.Remove("temp.csv")
					return
				}
				if record[0] == strconv.Itoa(taskID) {

					fmt.Println("Deleting record...")
					record[3] = "true"
					readComplete = true
				}
			}
			// Mark the task as done and break to write the record
			// csvWriter.Flush()
			// 	os.Remove("temp.csv")
			err = csvWriter.Write(record)
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

		fmt.Println("Task of task ID " + args[0] + " deleted")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
