package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/gofunct/common/pkg/logger/zap"
	"github.com/gofunct/mamba/runtime/config"
)

// Run generator
func Run(cfg *config.Config) {
	if cfg.Storage.Config.Name == "" {
		cfg.Storage.Config.Name = cfg.Project.Name
	}
	if cfg.Storage.MySQL {
		cfg.Storage.Config.Driver = config.StorageMySQL
	}
	if cfg.Storage.Postgres {
		cfg.Storage.Config.Driver = config.StoragePostgres
	}
	zap.LogF("Base templates", copyTemplates(
		path.Join(cfg.Directories.Templates, config.Base),
		cfg.Directories.Service,
	))
	if cfg.API.Enabled {
		zap.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.API, config.Base),
			cfg.Directories.Service,
		))
		if cfg.API.Gateway {
			zap.LogF("Gateway templates for API", copyTemplates(
				path.Join(cfg.Directories.Templates, config.API, config.APIGateway),
				cfg.Directories.Service,
			))
		}
	}
	if cfg.Storage.Enabled {
		zap.LogF("Storage base templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Storage, config.Base),
			cfg.Directories.Service,
		))
		if cfg.Storage.Postgres {
			zap.LogF("Storage templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}
		if cfg.Storage.MySQL {
			zap.LogF("Storage templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Storage, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}
	if cfg.API.Enabled && cfg.Storage.Enabled && cfg.Project.Contract {
		zap.LogF("Contract example templates", copyTemplates(
			path.Join(cfg.Directories.Templates, config.Contract, config.Base),
			cfg.Directories.Service,
		))
		if cfg.Storage.Postgres {
			zap.LogF("Contract templates for postgres", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Contract, config.StoragePostgres),
				cfg.Directories.Service,
			))
		}
		if cfg.Storage.MySQL {
			zap.LogF("Contract templates for mysql", copyTemplates(
				path.Join(cfg.Directories.Templates, config.Contract, config.StorageMySQL),
				cfg.Directories.Service,
			))
		}
	}
	zap.LogF("Render templates", render(cfg))
	zap.LogF("Could not change directory", os.Chdir(cfg.Directories.Service))
	if cfg.API.Enabled && cfg.Storage.Enabled && cfg.Project.Contract {
		log.Println("Prepare contracts:")
		zap.LogF("Generate contracts", Exec("make", "contracts"))
	}
	log.Println("Initialize vendors:")
	zap.LogF("Init dep", Exec("dep", "init", "-skip-tools"))
	zap.LogF("Tests", Exec("make", "check-all"))

	if cfg.Project.GitInit {
		log.Println("Initialize Git repository:")
		zap.LogF("Init git", Exec("git", "init"))
		zap.LogF("Add repo files", Exec("git", "add", "--all"))
		zap.LogF("Initial commit", Exec("git", "commit", "-m", "'Initial commit'"))
	}
	fmt.Printf("New repository was created, use command 'cd %s'", cfg.Directories.Service)
}

// Exec runs the commands
func Exec(command ...string) error {
	execCmd := exec.Command(command[0], command[1:]...) // nolint: gosec
	execCmd.Stderr = os.Stderr
	execCmd.Stdout = os.Stdout
	execCmd.Stdin = os.Stdin
	return execCmd.Run()
}
