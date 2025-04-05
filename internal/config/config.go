package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Connectors ConnectorsConfig `yaml:"connectors"`
}

// ConnectorsConfig holds configuration for all connectors
type ConnectorsConfig struct {
	Telegram TelegramConfig `yaml:"telegram"`
	RSS      RSSConfig      `yaml:"rss"`
	Reddit   RedditConfig   `yaml:"reddit"`
}

// TelegramConfig holds configuration for Telegram connector
type TelegramConfig struct {
	Enabled    bool             `yaml:"enabled"`
	Channels   []ChannelConfig  `yaml:"channels"`
	Credentials TelegramCredentials `yaml:"credentials"`
}

// TelegramCredentials holds authentication information for Telegram
type TelegramCredentials struct {
	APIID   string `yaml:"api_id"`
	APIHash string `yaml:"api_hash"`
}

// ChannelConfig holds configuration for a Telegram channel
type ChannelConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// RSSConfig holds configuration for RSS connector
type RSSConfig struct {
	Enabled  bool           `yaml:"enabled"`
	Feeds    []FeedConfig   `yaml:"feeds"`
	Settings RSSSettings    `yaml:"settings"`
}

// FeedConfig holds configuration for an RSS feed
type FeedConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// RSSSettings holds settings for RSS connector
type RSSSettings struct {
	Timeout    time.Duration `yaml:"timeout"`
	UserAgent  string        `yaml:"user_agent"`
}

// RedditConfig holds configuration for Reddit connector
type RedditConfig struct {
	Enabled    bool             `yaml:"enabled"`
	Subreddits []SubredditConfig `yaml:"subreddits"`
	Settings   RedditSettings   `yaml:"settings"`
}

// SubredditConfig holds configuration for a subreddit
type SubredditConfig struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

// RedditSettings holds settings for Reddit connector
type RedditSettings struct {
	Timeout      time.Duration `yaml:"timeout"`
	UserAgent    string        `yaml:"user_agent"`
	ClientID     string        `yaml:"client_id"`
	ClientSecret string        `yaml:"client_secret"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Limit        int           `yaml:"limit"`
	Sort         string        `yaml:"sort"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadConnectorsConfig loads just the connectors configuration
func LoadConnectorsConfig(path string) (*ConnectorsConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ConnectorsConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}