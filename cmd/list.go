package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks that are not marked as done.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Need to get the value of the all flag!
		All, _ := cmd.Flags().GetBool("all")

		file, err := os.Open("tasks.csv")
		if err != nil {
			fmt.Println("No tasks to list")
			return
		}
		defer file.Close()

		csvReader := csv.NewReader(file)

		headers, err := csvReader.Read()
		if err != nil {
			panic(err)
		}

		// Print the headers
		fmt.Printf("%-3s %-30s %-10s %-3s\n", headers[0], headers[1], headers[2], headers[3])

		// Loop through the records and print the tasks that are not done or deleted, unless the all flag is set
		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			if !All && (record[2] == "done" || record[3] == "true") {
				continue
			}
			fmt.Printf("%-3s %-30s %-10s %-3s\n", record[0], record[1], record[2], record[3])
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("all", "a", false, "List all tasks")
}
