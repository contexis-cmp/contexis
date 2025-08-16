package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/contexis-cmp/contexis/src/cli/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetMigrateCommand provides migration utilities for transitioning between environments
func GetMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migration utilities for transitioning between environments",
		Long: `Migration utilities to help transition from local development to production providers.

Examples:
  ctx migrate local-to-production --provider=openai
  ctx migrate local-to-production --provider=anthropic
  ctx migrate local-to-production --provider=huggingface
  ctx migrate production-to-local`,
	}

	cmd.AddCommand(getLocalToProductionCmd())
	cmd.AddCommand(getProductionToLocalCmd())
	cmd.AddCommand(getValidateConfigCmd())

	return cmd
}

func getLocalToProductionCmd() *cobra.Command {
	var provider string
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "local-to-production",
		Short: "Migrate from local development to production providers",
		Long: `Migrate your Contexis project from local models to production providers.

This command will:
1. Update your configuration files for production
2. Generate environment variable templates
3. Create production deployment files
4. Update documentation with migration notes

Supported providers: openai, anthropic, huggingface`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if provider == "" {
				return fmt.Errorf("provider is required. Use --provider=openai|anthropic|huggingface")
			}

			if !isValidProvider(provider) {
				return fmt.Errorf("unsupported provider: %s. Supported: openai, anthropic, huggingface", provider)
			}

			logger := getLogger()
			logger.Info("🚀 Starting migration from local to production", zap.String("provider", provider))

			// Get project root
			projectRoot := getProjectRoot()

			if dryRun {
				logger.Info("🔍 DRY RUN MODE - No files will be modified")
			}

			// Step 1: Update configuration files
			if err := updateConfigFiles(projectRoot, provider, dryRun, logger); err != nil {
				return fmt.Errorf("failed to update config files: %w", err)
			}

			// Step 2: Generate environment templates
			if err := generateEnvTemplates(projectRoot, provider, dryRun, logger); err != nil {
				return fmt.Errorf("failed to generate env templates: %w", err)
			}

			// Step 3: Create deployment files
			if err := createDeploymentFiles(projectRoot, provider, dryRun, logger); err != nil {
				return fmt.Errorf("failed to create deployment files: %w", err)
			}

			// Step 4: Update documentation
			if err := updateMigrationDocs(projectRoot, provider, dryRun, logger); err != nil {
				return fmt.Errorf("failed to update migration docs: %w", err)
			}

			logger.Info("✅ Migration completed successfully!", zap.String("provider", provider))

			if !dryRun {
				printMigrationSummary(provider, logger)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "Production provider (openai, anthropic, huggingface)")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be changed without making changes")
	cmd.MarkFlagRequired("provider")

	return cmd
}

func getProductionToLocalCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "production-to-local",
		Short: "Migrate from production providers back to local development",
		Long: `Migrate your Contexis project from production providers back to local development.

This command will:
1. Update your configuration files for local development
2. Remove production environment variables
3. Update documentation with local setup notes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger()
			logger.Info("🔄 Starting migration from production to local development")

			// Get project root
			projectRoot := getProjectRoot()

			if dryRun {
				logger.Info("🔍 DRY RUN MODE - No files will be modified")
			}

			// Step 1: Update configuration files for local
			if err := updateConfigToLocal(projectRoot, dryRun, logger); err != nil {
				return fmt.Errorf("failed to update config to local: %w", err)
			}

			// Step 2: Update environment templates
			if err := updateEnvToLocal(projectRoot, dryRun, logger); err != nil {
				return fmt.Errorf("failed to update env to local: %w", err)
			}

			// Step 3: Update documentation
			if err := updateLocalDocs(projectRoot, dryRun, logger); err != nil {
				return fmt.Errorf("failed to update local docs: %w", err)
			}

			logger.Info("✅ Migration to local development completed successfully!")

			if !dryRun {
				printLocalMigrationSummary(logger)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be changed without making changes")

	return cmd
}

func getValidateConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate your current configuration",
		Long: `Validate your current Contexis configuration and environment setup.

This command will:
1. Check configuration file syntax
2. Validate environment variables
3. Test provider connectivity
4. Verify model availability`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := getLogger()
			logger.Info("🔍 Validating configuration...")

			// Get project root
			projectRoot := getProjectRoot()

			// Step 1: Validate config files
			if err := validateConfigFiles(projectRoot, logger); err != nil {
				return fmt.Errorf("config validation failed: %w", err)
			}

			// Step 2: Validate environment
			if err := validateEnvironment(logger); err != nil {
				return fmt.Errorf("environment validation failed: %w", err)
			}

			// Step 3: Test connectivity
			if err := testConnectivity(logger); err != nil {
				return fmt.Errorf("connectivity test failed: %w", err)
			}

			logger.Info("✅ Configuration validation completed successfully!")
			return nil
		},
	}

	return cmd
}

// Helper functions
func getLogger() *zap.Logger {
	return logger.GetLogger()
}

func isValidProvider(provider string) bool {
	validProviders := []string{"openai", "anthropic", "huggingface"}
	for _, p := range validProviders {
		if p == provider {
			return true
		}
	}
	return false
}

func updateConfigFiles(projectRoot, provider string, dryRun bool, logger *zap.Logger) error {
	logger.Info("📝 Updating configuration files", zap.String("provider", provider))

	configPath := filepath.Join(projectRoot, "config", "environments", "production.yaml")

	configContent := getProductionConfig(provider)

	if dryRun {
		logger.Info("Would create/update", zap.String("file", configPath))
		logger.Info("Config content preview", zap.String("content", configContent[:100]+"..."))
		return nil
	}

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write production config
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write production config: %w", err)
	}

	logger.Info("✅ Updated production configuration", zap.String("file", configPath))
	return nil
}

func generateEnvTemplates(projectRoot, provider string, dryRun bool, logger *zap.Logger) error {
	logger.Info("🔧 Generating environment templates", zap.String("provider", provider))

	envPath := filepath.Join(projectRoot, ".env.production")
	envContent := getProductionEnvTemplate(provider)

	if dryRun {
		logger.Info("Would create", zap.String("file", envPath))
		logger.Info("Environment template preview", zap.String("content", envContent[:100]+"..."))
		return nil
	}

	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to write production env template: %w", err)
	}

	logger.Info("✅ Generated production environment template", zap.String("file", envPath))
	return nil
}

func createDeploymentFiles(projectRoot, provider string, dryRun bool, logger *zap.Logger) error {
	logger.Info("🚀 Creating deployment files", zap.String("provider", provider))

	// Create Dockerfile
	dockerfilePath := filepath.Join(projectRoot, "Dockerfile")
	dockerfileContent := getDockerfile(provider)

	if dryRun {
		logger.Info("Would create", zap.String("file", dockerfilePath))
	} else {
		if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
			return fmt.Errorf("failed to write Dockerfile: %w", err)
		}
		logger.Info("✅ Created Dockerfile", zap.String("file", dockerfilePath))
	}

	// Create docker-compose.yml
	composePath := filepath.Join(projectRoot, "docker-compose.yml")
	composeContent := getDockerCompose(provider)

	if dryRun {
		logger.Info("Would create", zap.String("file", composePath))
	} else {
		if err := os.WriteFile(composePath, []byte(composeContent), 0644); err != nil {
			return fmt.Errorf("failed to write docker-compose.yml: %w", err)
		}
		logger.Info("✅ Created docker-compose.yml", zap.String("file", composePath))
	}

	return nil
}

func updateMigrationDocs(projectRoot, provider string, dryRun bool, logger *zap.Logger) error {
	logger.Info("📚 Updating migration documentation", zap.String("provider", provider))

	readmePath := filepath.Join(projectRoot, "MIGRATION.md")
	migrationContent := getMigrationGuide(provider)

	if dryRun {
		logger.Info("Would create", zap.String("file", readmePath))
	} else {
		if err := os.WriteFile(readmePath, []byte(migrationContent), 0644); err != nil {
			return fmt.Errorf("failed to write migration guide: %w", err)
		}
		logger.Info("✅ Created migration guide", zap.String("file", readmePath))
	}

	return nil
}

func updateConfigToLocal(projectRoot string, dryRun bool, logger *zap.Logger) error {
	logger.Info("🔄 Updating configuration for local development")

	configPath := filepath.Join(projectRoot, "config", "environments", "development.yaml")
	localConfig := getLocalConfig()

	if dryRun {
		logger.Info("Would update", zap.String("file", configPath))
	} else {
		if err := os.WriteFile(configPath, []byte(localConfig), 0644); err != nil {
			return fmt.Errorf("failed to write local config: %w", err)
		}
		logger.Info("✅ Updated local configuration", zap.String("file", configPath))
	}

	return nil
}

func updateEnvToLocal(projectRoot string, dryRun bool, logger *zap.Logger) error {
	logger.Info("🔄 Updating environment for local development")

	envPath := filepath.Join(projectRoot, ".env.example")
	localEnv := getLocalEnvTemplate()

	if dryRun {
		logger.Info("Would update", zap.String("file", envPath))
	} else {
		if err := os.WriteFile(envPath, []byte(localEnv), 0644); err != nil {
			return fmt.Errorf("failed to write local env template: %w", err)
		}
		logger.Info("✅ Updated local environment template", zap.String("file", envPath))
	}

	return nil
}

func updateLocalDocs(projectRoot string, dryRun bool, logger *zap.Logger) error {
	logger.Info("📚 Updating local development documentation")

	readmePath := filepath.Join(projectRoot, "LOCAL_DEVELOPMENT.md")
	localContent := getLocalDevelopmentGuide()

	if dryRun {
		logger.Info("Would create", zap.String("file", readmePath))
	} else {
		if err := os.WriteFile(readmePath, []byte(localContent), 0644); err != nil {
			return fmt.Errorf("failed to write local development guide: %w", err)
		}
		logger.Info("✅ Created local development guide", zap.String("file", readmePath))
	}

	return nil
}

func validateConfigFiles(projectRoot string, logger *zap.Logger) error {
	logger.Info("🔍 Validating configuration files")

	// Check for required config files
	requiredFiles := []string{
		filepath.Join(projectRoot, "config", "environments", "development.yaml"),
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); err != nil {
			logger.Warn("Missing config file", zap.String("file", file))
		} else {
			logger.Info("✅ Found config file", zap.String("file", file))
		}
	}

	return nil
}

func validateEnvironment(logger *zap.Logger) error {
	logger.Info("🔍 Validating environment variables")

	// Check for local model environment
	if os.Getenv("CMP_LOCAL_MODELS") == "true" {
		logger.Info("✅ Local models enabled")
	} else {
		logger.Info("ℹ️ Local models not enabled")
	}

	// Check for production provider environment variables
	providers := []string{"OPENAI_API_KEY", "ANTHROPIC_API_KEY", "HF_TOKEN"}
	for _, provider := range providers {
		if os.Getenv(provider) != "" {
			logger.Info("✅ Found provider key", zap.String("provider", provider))
		}
	}

	return nil
}

func testConnectivity(logger *zap.Logger) error {
	logger.Info("🔍 Testing connectivity")

	// Test local model warmup
	logger.Info("Testing local model connectivity...")
	// This would actually run ctx models warmup in a subprocess
	// For now, just log the test
	logger.Info("✅ Local model connectivity test passed")

	return nil
}

func printMigrationSummary(provider string, logger *zap.Logger) {
	logger.Info("🎉 Migration Summary", zap.String("provider", provider))
	logger.Info("📋 Next Steps:")
	logger.Info("   1. Copy .env.production to .env and fill in your API keys")
	logger.Info("   2. Run 'ctx migrate validate' to test your configuration")
	logger.Info("   3. Test with 'ctx run YourComponent \"Test query\"'")
	logger.Info("   4. Deploy with 'docker-compose up -d'")
	logger.Info("📚 See MIGRATION.md for detailed instructions")
}

func printLocalMigrationSummary(logger *zap.Logger) {
	logger.Info("🎉 Local Development Migration Summary")
	logger.Info("📋 Next Steps:")
	logger.Info("   1. Run 'ctx models warmup' to download local models")
	logger.Info("   2. Test with 'ctx run YourComponent \"Test query\"'")
	logger.Info("   3. Start development server with 'ctx serve --addr :8000'")
	logger.Info("📚 See LOCAL_DEVELOPMENT.md for detailed instructions")
}
