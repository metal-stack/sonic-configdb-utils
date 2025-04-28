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
	RunE: func(cmd *cobra.Command, args []string) error {
		sonicEnvFile, _ := cmd.Flags().GetString("env-file")
		env, err := p.GetEnvironment(sonicEnvFile)
		if err != nil {
			return fmt.Errorf("failed to get environment information:%w", err)
		}

		platformIdentifier := env.Platform
		deviceDir, _ := cmd.Flags().GetString("device-dir")
		platformFile := fmt.Sprintf("%s/%s/platform.json", deviceDir, platformIdentifier)

		platformBytes, err := os.ReadFile(platformFile)
		if err != nil {
			return fmt.Errorf("failed to read platform.json file:%w", err)
		}

		platform, err := p.UnmarshalPlatformJSON(platformBytes)
		if err != nil {
			return fmt.Errorf("failed to parse platform.json:%w", err)
		}

		inputFile, _ := cmd.Flags().GetString("input")
		inputBytes, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read input file:%w", err)
		}

		values, err := values.UnmarshalValues(inputBytes)
		if err != nil {
			return fmt.Errorf("failed to parse input file:%w", err)
		}

		configDB, err := configdb.GenerateConfigDB(values, platform, env)
		if err != nil {
			return fmt.Errorf("failed to generate config:%w", err)
		}

		configDBBytes, err := json.MarshalIndent(configDB, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to serialize json:%w", err)
		}

		output, _ := cmd.Flags().GetString("output")
		err = os.WriteFile(output, configDBBytes, 0644) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to write file:%w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input", "i", "sonic-config.yaml", "path to input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output", "o", "/etc/sonic/config_db.json", "path to output file")
	generateCmd.Flags().String("device-dir", "/usr/share/sonic/device", "directory which holds all device-specific files")
	generateCmd.Flags().StringP("env-file", "e", "/etc/sonic/sonic-environment", "sonic-environment file holding platform information")
}
