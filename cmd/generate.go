package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/metal-stack/sonic-configdb-utils/configdb"
	"github.com/metal-stack/sonic-configdb-utils/values"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config_db.json",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile, _ := cmd.Flags().GetString("input")
		if inputFile == "" {
			fmt.Println("missing input values; please provide an input file via --input flag")
			os.Exit(1)
		}

		bytes, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("failed to read input file, %v\n", err)
			os.Exit(1)
		}

		values, err := values.UnmarshalValues(bytes)
		if err != nil {
			fmt.Printf("failed to parse input file, %v\n", err)
			os.Exit(1)
		}

		config := configdb.GenerateConfigDB(values)

		bytes, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Printf("failed to serialize json, %v\n", err)
			os.Exit(1)
		}

		fileInfo, err := os.Lstat(inputFile)
		if err != nil {
			fmt.Printf("failed to retrieve file info for %s, %v\n", inputFile, err)
			os.Exit(1)
		}

		outputFile, _ := cmd.Flags().GetString("output")
		if outputFile == "" {
			outputFile = filenameWithoutExtension(inputFile) + ".json"
		}

		err = os.WriteFile(outputFile, bytes, fileInfo.Mode())
		if err != nil {
			fmt.Printf("failed to write file, %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input", "i", "", "input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output", "o", "", "output file")
}

func filenameWithoutExtension(name string) string {
	return strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
}
