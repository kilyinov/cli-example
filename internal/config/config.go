package config

import "github.com/spf13/viper"

type Config struct {
	JIRA JIRAConfig `mapstructure:"jira"`
	Git  GitConfig  `mapstructure:"git"`
}

type JIRAConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Project string `mapstructure:"project"`
	Email   string `mapstructure:"email"`
	Token   string `mapstructure:"token"`
}

type GitConfig struct {
	DefaultBranch string `mapstructure:"default_branch"`
	BranchPrefix  string `mapstructure:"branch_prefix"`
	Remote        string `mapstructure:"remote"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Git: GitConfig{
			DefaultBranch: "main",
			BranchPrefix:  "feature",
			Remote:        "origin",
		},
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
