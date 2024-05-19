package data

import (
	"fmt"
	"module/model/service"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"go.uber.org/zap"
)

type PostgresClient struct {
	PostgresConfig PostgresConfig
	Logger         *zap.Logger
	client         *pg.DB
}

type PostgresConfig struct {
	Addr     string
	User     string
	Password string
	Database string
}

func (p *PostgresClient) StartPostgres() {
	p.Logger.Info("Starting Postgres client")
	p.client = pg.Connect(&pg.Options{
		Addr:     p.PostgresConfig.Addr,
		User:     p.PostgresConfig.User,
		Password: p.PostgresConfig.Password,
		Database: p.PostgresConfig.Database,
	})
	p.Logger.Info("Postgres client created")
	defer p.client.Close()
}

// createSchema creates database schema for Account
func (p *PostgresClient) CreateSchema() error {
	models := []interface{}{
		(*service.Account)(nil),
	}

	for _, model := range models {
		err := p.client.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	p.Logger.Info("Account table created")
	return nil
}

func (p *PostgresClient) GetUser(account *service.Account, userQuery string) {
	p.client.Model(account).
		Where("name = ?", userQuery).
		Select()
}

func (p *PostgresClient) PutInStore(model interface{}, sendChan chan<- bool) {
	p.Logger.Info("Persisting data in database")
	_, err := p.client.Model(model).Insert()
	if err != nil {
		fmt.Printf("Insert failed: %s", err.Error())
		sendChan <- false
	} else {
		sendChan <- true
	}
}
