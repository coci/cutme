package repositories

import (
	"fmt"
	"log"

	"github.com/apache/cassandra-gocql-driver/v2"
	"github.com/coci/cutme/internal/core/domain"
	"github.com/coci/cutme/internal/infra/config"
)

type LinkRepository struct {
	conn *gocql.Session
	cfg  *config.Config
}

func NewLinkRepository(cfg *config.Config) LinkRepository {

	cluster := gocql.NewCluster(cfg.CassandraCfg.Hosts...)
	cluster.Keyspace = cfg.CassandraCfg.Keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.CassandraCfg.Username,
		Password: cfg.CassandraCfg.Password,
	}

	// Create Cassandra session
	conn, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	return LinkRepository{
		conn: conn,
		cfg:  cfg,
	}
}

func (l LinkRepository) Save(link domain.Link) error {
	if err := l.conn.Query(
		`INSERT INTO ? (code, link) VALUES (?, ?)`,
		l.cfg.CassandraCfg.LinkTableName, link.Code, link.Link,
	).Exec(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (l LinkRepository) FindByCode(code string) domain.Link {
	result := domain.Link{}

	err := l.conn.Query(
		`SELECT link,code FROM ? WHERE code = ? LIMIT 1`,
		l.cfg.CassandraCfg.LinkTableName, code,
	).Consistency(gocql.One).Scan(&result.Link, &result.Code)

	if err != nil {
		return domain.Link{}
	}
	return result
}
