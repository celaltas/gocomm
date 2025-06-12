package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Config holds all command arguments and options
type Config struct {
	File1       string
	File2       string
	HideCol1    bool
	HideCol2    bool
	HideCol3    bool
	Insensitive bool
	Delimiter   string
}

var (
	hideCol1    bool
	hideCol2    bool
	hideCol3    bool
	insensitive bool
	delimiter   string
)

var rootCmd = &cobra.Command{
	Use:   "gocomm [OPTION]... FILE1 FILE2",
	Short: "Compare sorted files FILE1 and FILE2 line by line.",
	Long: `Compare sorted files FILE1 and FILE2 line by line.

With no options, produce three-column output:
  column 1: lines only in FILE1
  column 2: lines only in FILE2
  column 3: lines in both FILE1 and FILE2

Options:
  -1                      suppress column 1 (lines unique to FILE1)
  -2                      suppress column 2 (lines unique to FILE2)
  -3                      suppress column 3 (lines that appear in both files)
  -i, --insensitive       case-insensitive comparison
  -d, --output-delimiter  separate columns with STR

Examples:
  gocomm file1.txt file2.txt
  gocomm -12 file1.txt file2.txt
  gocomm -i -d "|" -12 file1.txt file2.txt`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file1, file2 := args[0], args[1]
		if file1 == "-" && file2 == "-" {
			return fmt.Errorf("Both input files cannot be STDIN (\"-\")")
		}
		config := Config{
			File1:       file1,
			File2:       file2,
			HideCol1:    hideCol1,
			HideCol2:    hideCol2,
			HideCol3:    hideCol3,
			Insensitive: insensitive,
			Delimiter:   delimiter,
		}
		_, err := Open(config.File1)
		if err != nil {
			return err
		}
		_, err = Open(config.File2)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	rootCmd.SilenceUsage = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&hideCol1, "suppress1", "1", false, "suppress column 1 (lines unique to FILE1)")
	rootCmd.Flags().BoolVarP(&hideCol2, "suppress2", "2", false, "suppress column 2 (lines unique to FILE2)")
	rootCmd.Flags().BoolVarP(&hideCol3, "suppress3", "3", false, "suppress column 3 (lines that appear in both files)")
	rootCmd.Flags().BoolVarP(&insensitive, "insensitive", "i", false, "case-insensitive comparison")
	rootCmd.Flags().StringVarP(&delimiter, "output-delimiter", "d", "\t", "separate columns with STR")
}

func Open(fileName string) (*bufio.Reader, error) {
	if fileName == "-" {
		return bufio.NewReader(os.Stdin), nil
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	return reader, nil
}
