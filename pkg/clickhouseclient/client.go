package clickhouseclient

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/ClickHouse/clickhouse-go" // driver
)

const (
	connectTryCount = 3
)

type Config struct {
	Addr         string
	User         string
	Password     string
	Database     string
	ReadTimeout  int
	WriteTimeout int
}

type Client struct {
	db *sqlx.DB
}

func NewClient(config *Config) (client *Client, err error) {
	for i := 0; i < connectTryCount; i++ {
		db, err := sqlx.Connect("clickhouse", connString(config))
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		return &Client{db: db}, nil
	}

	return nil, err
}

func (c *Client) Client() *sqlx.DB {
	return c.db
}

func (c *Client) Close() error {
	return c.db.Close()
}

func connString(config *Config) string {
	return fmt.Sprintf("tcp://%s?username=%s&database=%s&read_timeout=%d&write_timeout=%d",
		config.Addr, config.User, config.Database, config.ReadTimeout, config.WriteTimeout)
}
