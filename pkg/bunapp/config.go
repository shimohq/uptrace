package bunapp

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ReadConfig(configFile, service string) (*AppConfig, error) {
	configFile, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	cfg := new(AppConfig)
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	cfg.Filepath = configFile
	cfg.Service = service

	if len(cfg.Users) == 0 {
		return nil, fmt.Errorf("config must contain at least one user")
	}
	if len(cfg.Projects) == 0 {
		return nil, fmt.Errorf("config must contain at least one project")
	}

	httpHost, httpPort, err := net.SplitHostPort(cfg.Listen.HTTP)
	if err != nil {
		return nil, fmt.Errorf("can't parse option listen.http addr: %w", err)
	}
	if httpHost == "" {
		httpHost = cfg.Site.Host
	}
	cfg.Listen.HTTPHost = httpHost
	cfg.Listen.HTTPPort = httpPort

	grpcHost, grpcPort, err := net.SplitHostPort(cfg.Listen.GRPC)
	if err != nil {
		return nil, fmt.Errorf("can't parse option listen.grpc addr: %w", err)
	}
	if grpcHost == "" {
		grpcHost = cfg.Site.Host
	}
	cfg.Listen.GRPCHost = grpcHost
	cfg.Listen.GRPCPort = grpcPort

	return cfg, nil
}

type AppConfig struct {
	Filepath string `yaml:"-"`
	Service  string `yaml:"service"`

	Debug     bool   `yaml:"debug"`
	SecretKey string `yaml:"secret_key"`

	Site struct {
		Scheme string `yaml:"scheme"`
		Host   string `yaml:"host"`
	} `yaml:"site"`

	Listen struct {
		HTTP string `yaml:"http"`
		GRPC string `yaml:"grpc"`

		HTTPHost string `yaml:"-"`
		HTTPPort string `yaml:"-"`

		GRPCHost string `yaml:"-"`
		GRPCPort string `yaml:"-"`
	} `yaml:"listen"`

	DB BunConfig `yaml:"db"`
	CH CHConfig  `yaml:"ch"`

	Retention struct {
		TTL string `yaml:"ttl"`
	} `yaml:"retention"`

	Users    []User    `yaml:"users"`
	Projects []Project `yaml:"projects"`

	CHSelectLimits struct {
		SampleRows     int64 `yaml:"sample_rows"`
		MaxRowsToRead  int64 `yaml:"max_rows_to_read"`
		MaxBytesToRead int64 `yaml:"max_bytes_to_read"`
	} `yaml:"ch_select_limits"`
}

type User struct {
	ID       uint64 `yaml:"id" json:"id"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"-"`
}

type Project struct {
	ID    uint32 `yaml:"id" json:"id"`
	Name  string `yaml:"name" json:"name"`
	Token string `yaml:"token" json:"token"`
}

func (c *AppConfig) SiteAddr() string {
	return fmt.Sprintf("%s://%s:%s/", c.Site.Scheme, c.Listen.HTTPHost, c.Listen.HTTPPort)
}

func (c *AppConfig) GRPCEndpoint(project *Project) string {
	return fmt.Sprintf("%s://%s:%s", c.Site.Scheme, c.Listen.GRPCHost, c.Listen.GRPCPort)
}

func (c *AppConfig) HTTPEndpoint(project *Project) string {
	return fmt.Sprintf("%s://%s:%s", c.Site.Scheme, c.Listen.HTTPHost, c.Listen.HTTPPort)
}

func (c *AppConfig) GRPCDsn(project *Project) string {
	return fmt.Sprintf("%s://%s@%s:%s/%d",
		c.Site.Scheme, project.Token, c.Listen.GRPCHost, c.Listen.GRPCPort, project.ID)
}

func (c *AppConfig) HTTPDsn(project *Project) string {
	return fmt.Sprintf("%s://%s@%s:%s/%d",
		c.Site.Scheme, project.Token, c.Listen.HTTPHost, c.Listen.HTTPPort, project.ID)
}

type BunConfig struct {
	DSN string `yaml:"dsn"`
}

type CHConfig struct {
	DSN string `yaml:"dsn"`
}
