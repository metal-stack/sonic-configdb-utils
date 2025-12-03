package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/metal-stack/sonic-configdb-utils/configdb"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
	v "github.com/metal-stack/sonic-configdb-utils/version"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config_db.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, _ := cmd.Flags().GetString("input-file")
		outputFile, _ := cmd.Flags().GetString("output-file")
		deviceDir, _ := cmd.Flags().GetString("device-dir")
		sonicEnvFile, _ := cmd.Flags().GetString("env-file")
		sonicVersionFile, _ := cmd.Flags().GetString("version-file")

		env, err := p.GetEnvironment(sonicEnvFile)
		if err != nil {
			return fmt.Errorf("failed to get environment information:%w", err)
		}

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

		platformFile := fmt.Sprintf("%s/%s/platform.json", deviceDir, env.Platform)
		platformBytes, err := os.ReadFile(platformFile)
		if err != nil {
			return fmt.Errorf("failed to read platform.json file:%w", err)
		}

		platform, err := p.UnmarshalPlatformJSON(platformBytes)
		if err != nil {
			return fmt.Errorf("failed to parse platform.json:%w", err)
		}

		versionBytes, err := os.ReadFile(sonicVersionFile)
		if err != nil {
			return fmt.Errorf("failed to read version file:%w", err)
		}

		version, err := v.GetVersion(versionBytes)
		if err != nil {
			return fmt.Errorf("failed to parse version file:%w", err)
		}

		configDB, err := configdb.GenerateConfigDB(values, platform, env, version)
		if err != nil {
			return fmt.Errorf("failed to generate config:%w", err)
		}

		configDBBytes, err := json.MarshalIndent(configDB, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to serialize json:%w", err)
		}

		err = os.WriteFile(outputFile, configDBBytes, 0644) //nolint:gosec
		if err != nil {
			return fmt.Errorf("failed to write file:%w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringP("input-file", "i", "sonic-config.yaml", "path to input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output-file", "o", "config_db.json", "path to output file")
	generateCmd.Flags().String("device-dir", "/usr/share/sonic/device", "directory which holds all device-specific files")
	generateCmd.Flags().StringP("env-file", "e", "/etc/sonic/sonic-environment", "sonic-environment file holding platform information")
	generateCmd.Flags().StringP("version-file", "v", "/etc/sonic/sonic_version.yml", "sonic version file")
}
