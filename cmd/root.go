package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

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

type Column struct {
	col1 string
	col2 string
	col3 string
}

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
		reader1, err := Open(config.File1)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", config.File1, err)
		}
		reader2, err := Open(config.File2)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", config.File2, err)
		}

		columns, err := CompareLines(reader1, reader2, config)
		if err != nil {
			return fmt.Errorf("comparison error: %w", err)
		}

		PrintColumns(columns, config)

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
func CompareLines(reader1, reader2 *bufio.Reader, config Config) ([]Column, error) {
	var columns []Column

	for {
		line1, err1 := readLine(reader1)
		line2, err2 := readLine(reader2)

		if err1 != nil && err1 != io.EOF {
			return nil, err1
		}
		if err2 != nil && err2 != io.EOF {
			return nil, err2
		}

		if line1 == "" && line2 == "" {
			break
		}

		if config.Insensitive {
			line1 = strings.ToLower(line1)
			line2 = strings.ToLower(line2)
		}

		col := Column{}
		switch {
		case line1 == line2:
			col.col3 = line1
		case line2 == "":
			col.col1 = line1
		case line1 == "":
			col.col2 = line2
		default:
			col.col1 = line1
			col.col2 = line2
		}

		columns = append(columns, col)
	}

	return columns, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	return strings.TrimSpace(line), err
}

func PrintColumns(columns []Column, config Config) {
    for _, col := range columns {
        var output strings.Builder
        
        if !config.HideCol1 && col.col1 != "" {
            output.WriteString(col.col1)
        }
        output.WriteString(config.Delimiter)
        
        if !config.HideCol2 && col.col2 != "" {
            output.WriteString(col.col2)
        }
        output.WriteString(config.Delimiter)
        
        if !config.HideCol3 && col.col3 != "" {
            output.WriteString(col.col3)
        }
        
        if output.Len() > len(config.Delimiter)*2 {
            fmt.Println(strings.TrimSuffix(output.String(), config.Delimiter))
        }
    }
}