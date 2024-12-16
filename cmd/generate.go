package cmd

import (
	"fmt"
	"os"

	"github.com/metal-stack/sonic-configdb-utils/values"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config_db.json",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		if input == "" {
			fmt.Println("missing input values; please provide an input file via --input flag")
			os.Exit(1)
		}

		bytes, err := os.ReadFile(input)
		if err != nil {
			fmt.Printf("failed to read input file, %v\n", err)
			os.Exit(1)
		}

		values, err := values.UnmarshalValues(bytes)
		if err != nil {
			fmt.Printf("failed to parse input file, %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%#v\n", values)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input", "i", "", "input file to generate the config_db.json from")
}
