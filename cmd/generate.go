package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/metal-stack/sonic-configdb-utils/configdb"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config_db.json",
	Run: func(cmd *cobra.Command, args []string) {
		platformFile, _ := cmd.Flags().GetString("platform")
		platformBytes, err := os.ReadFile(platformFile)
		if err != nil {
			fmt.Printf("failed to read platform.json file: %v\n", err)
			os.Exit(1)
		}

		platform, err := p.UnmarshalPlatformJSON(platformBytes)
		if err != nil {
			fmt.Printf("failed to parse platform.json: %v\n", err)
			os.Exit(1)
		}

		inputFile, _ := cmd.Flags().GetString("input")
		inputBytes, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("failed to read input file, %v\n", err)
			os.Exit(1)
		}

		values, err := values.UnmarshalValues(inputBytes)
		if err != nil {
			fmt.Printf("failed to parse input file, %v\n", err)
			os.Exit(1)
		}

		configDBFile, _ := cmd.Flags().GetString("output")
		configDBBytes, err := os.ReadFile(configDBFile)
		if err != nil {
			fmt.Printf("failed to read current config file, %v\n", err)
			os.Exit(1)
		}

		currentConfig, err := configdb.UnmarshalConfigDB(configDBBytes)
		if err != nil {
			fmt.Printf("failed to parse current config file, %v\n", err)
			os.Exit(1)
		}

		configDB, err := configdb.GenerateConfigDB(values, platform, currentConfig.DeviceMetadata)
		if err != nil {
			fmt.Printf("failed to generate config, %v\n", err)
			os.Exit(1)
		}

		configBytes, err := json.MarshalIndent(configDB, "", "  ")
		if err != nil {
			fmt.Printf("failed to serialize json, %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(configDBFile, configBytes, 0644) //nolint:gosec
		if err != nil {
			fmt.Printf("failed to write file, %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input", "i", "sonic-config.yaml", "path to input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output", "o", "/etc/sonic/config_db.json", "path to output file")
	generateCmd.Flags().StringP("platform", "p", "/usr/share/sonic/platform/platform.json", "path to platform.json")
}
