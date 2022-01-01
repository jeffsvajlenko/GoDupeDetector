/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"GoDupeDetector/internal/detection"
	"GoDupeDetector/internal/parsing"
	"GoDupeDetector/internal/printer"
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// detectCmd represents the detect command
var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Runs clone detection.",
	Long:  `Runs clone detection.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")

		finfo, err := os.Stat(input)
		if finfo == nil {
			return errors.New("Input directory does not exist or is invalid.")
		}
		if !finfo.IsDir() {
			return errors.New("Input must be a directory.")
		}

		finfo, err = os.Stat(output)
		if finfo != nil {
			return errors.New("Output file already exists or is a directory.")
		}
		outf, err := os.Create(output)
		if err != nil {
			return fmt.Errorf("Failed to create output file with error: %s", err.Error())
		}
		defer outf.Close()

		fmt.Print("Collecting go files...")
		files, err := parsing.FileList(input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(" %d files collected.\n", len(files))

		fmt.Print("Parsing...")
		pset, _ := parsing.Parse(files)
		fmt.Println(" Done!")

		fmt.Print("Detecting...")
		cset, _ := detection.DetectClones(pset, 0.70)
		fmt.Printf(" %d clones detected.\n", len(cset.Clones))

		fmt.Println("Outputting...")
		bout := bufio.NewWriter(outf)
		err = printer.PrintCloneReport(pset, cset, bout)

		if err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// detectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// detectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	detectCmd.Flags().String("input", "", "Input source folder.")
	detectCmd.MarkFlagRequired("input")
	detectCmd.Flags().String("output", "", "Output report file.")
	detectCmd.MarkFlagRequired("output")
}
