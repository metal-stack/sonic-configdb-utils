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
		inputFile, _ := cmd.Flags().GetString("input-file")
		outputFile, _ := cmd.Flags().GetString("output-file")
		deviceDir, _ := cmd.Flags().GetString("device-dir")

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

		platformIdentifier := values.DeviceMetadata.Platform
		platformFile := fmt.Sprintf("%s/%s/platform.json", deviceDir, platformIdentifier)

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

		configDB, err := configdb.GenerateConfigDB(values, platform, values.DeviceMetadata)
		if err != nil {
			fmt.Printf("failed to generate config, %v\n", err)
			os.Exit(1)
		}

		configDBBytes, err := json.MarshalIndent(configDB, "", "  ")
		if err != nil {
			fmt.Printf("failed to serialize json, %v\n", err)
			os.Exit(1)
		}

		err = os.WriteFile(outputFile, configDBBytes, 0644) //nolint:gosec
		if err != nil {
			fmt.Printf("failed to write file, %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input-file", "i", "sonic-config.yaml", "path to input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output-file", "o", "config_db.json", "path to output file")
	generateCmd.Flags().String("device-dir", "/usr/share/sonic/device", "directory which holds all device-specific files")
}
