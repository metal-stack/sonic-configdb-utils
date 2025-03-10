package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/metal-stack/sonic-configdb-utils/configdb"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a config_db.json",
	Run: func(cmd *cobra.Command, args []string) {
		configDir, _ := cmd.Flags().GetString("sonic-config-dir")

		sonicPlatform, _ := cmd.Flags().GetString("platform")
		if sonicPlatform == "" {
			var err error
			sonicPlatform, err = retrieveSonicPlatform(configDir)
			if err != nil {
				fmt.Printf("failed to retrieve sonic platform: %v\n", err)
				os.Exit(1)
			}
		}

		deviceDir, _ := cmd.Flags().GetString("device-dir")
		platformFile := fmt.Sprintf("%s/%s/platform.json", deviceDir, sonicPlatform)

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

		inputFile, _ := cmd.Flags().GetString("input-file")
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

		configDB, err := configdb.GenerateConfigDB(values, platform)
		if err != nil {
			fmt.Printf("failed to generate config, %v\n", err)
			os.Exit(1)
		}

		configBytes, err := json.MarshalIndent(configDB, "", "  ")
		if err != nil {
			fmt.Printf("failed to serialize json, %v\n", err)
			os.Exit(1)
		}

		outputFileName, _ := cmd.Flags().GetString("output-file-name")
		outputFile := fmt.Sprintf("%s/%s", configDir, outputFileName)

		err = os.WriteFile(outputFile, configBytes, 0644)
		if err != nil {
			fmt.Printf("failed to write file, %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().String("device-dir", "/usr/share/sonic/device", "directory that holds all vendor specific files")
	generateCmd.Flags().StringP("input-file", "i", "sonic-config.yaml", "path to input file to generate the config_db.json from")
	generateCmd.Flags().StringP("output-file-name", "o", "config_db.json", "output file name")
	generateCmd.Flags().StringP("platform", "p", "", "sonic platform")
	generateCmd.Flags().String("sonic-config-dir", "/etc/sonic", "where to store the generated config_db.json")
}

func retrieveSonicPlatform(configDir string) (string, error) {
	sonicEnvFile := fmt.Sprintf("%s/sonic-environment", configDir)
	f, err := os.Open(sonicEnvFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); strings.Contains(line, "PLATFORM=") {
			sonicPlatform, _ := strings.CutPrefix(line, "PLATFORM=")
			return sonicPlatform, nil
		}
	}

	return "", fmt.Errorf("no platform information found")
}
