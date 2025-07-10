package http

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/InVisionApp/go-health/v2"
	"github.com/qq-mercantil/qq-framework-basic-golang/cache"
	"github.com/qq-mercantil/qq-framework-basic-kafka/kafka"
	"github.com/qq-mercantil/qq-framework-db-golang/db"
	kafkaGo "github.com/segmentio/kafka-go"
)

type Checker interface {
	Status() (interface{}, error)
}

type HealthChecks struct {
	Health        *health.Health
	KafkaCheck    Checker
	PostgresCheck Checker
	RedisCheck    Checker
}

type KafkaCheck struct {
	Config kafka.IKafkaProvider
}

type PostgresCheck struct {
	db     *db.GormDatabase
	config db.IDbProvider
}

type RedisCheck struct {
	client cache.ICacheClient
	config cache.ICacheProvider
}

func (k *KafkaCheck) Status() (details interface{}, errorKafka error) {
	if len(k.Config.GetBrokers()) == 0 {
		return nil, fmt.Errorf("nenhum broker configurado")
	}

	details = map[string]interface{}{
		"brokers": k.Config.GetBrokers(),
	}

	var failedBrokers []string
	var failedControllerBroker []string

	for _, broker := range k.Config.GetBrokers() {
		conn, err := kafkaGo.DialContext(context.Background(), "tcp", broker)
		if err != nil {
			fmt.Println(err)
			failedBrokers = append(failedBrokers, broker)
			continue
		}
		defer conn.Close()

		_, err = conn.Controller()
		if err != nil {
			failedBrokers = append(failedBrokers, broker)
			failedControllerBroker = append(failedControllerBroker, broker)
			continue
		}
	}

	errorMessages := []string{}

	if len(failedBrokers) > 0 {
		errorMessages = append(errorMessages, fmt.Sprintf("falha ao conectar aos brokers: %s", strings.Join(failedBrokers, ", ")))
	}

	if len(failedControllerBroker) > 0 {
		errorMessages = append(errorMessages, fmt.Sprintf("falha ao obter o controller dos brokers: %s", strings.Join(failedControllerBroker, ", ")))
	}

	if len(errorMessages) > 0 {
		errorKafka = errors.New(strings.Join(errorMessages, "\n"))
	}

	return details, errorKafka
}

func (p *PostgresCheck) Status() (interface{}, error) {
	sqlDB, err := p.db.DB.DB()

	details := map[string]interface{}{
		"db":   p.config.GetName(),
		"host": p.config.GetHost(),
		"port": p.config.GetPort(),
	}

	if err != nil {
		return details, fmt.Errorf("falha ao obter a conexão do banco de dados: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return details, fmt.Errorf("falha ao fazer ping no PostgreSQL: %v", err)
	}

	return details, nil
}

func (r *RedisCheck) Status() (interface{}, error) {
	err := r.client.Ping(context.Background())

	details := map[string]interface{}{
		"host": r.config.GetHost(),
		"port": r.config.GetPort(),
		"db":   r.config.GetDB(),
	}

	if err != nil {
		return details, fmt.Errorf("falha ao conectar no redis: %v", err)
	}

	return details, nil
}

func RegisterHealthChecks(kafkaCfg kafka.IKafkaProvider, dbCfg db.IDbProvider, cacheProvider cache.ICacheProvider, db *db.GormDatabase, cache cache.ICacheClient) (*HealthChecks, error) {
	h := health.New()

	kafkaCheck := &KafkaCheck{Config: kafkaCfg}
	postgresCheck := &PostgresCheck{db: db, config: dbCfg}
	redisCheck := &RedisCheck{client: cache, config: cacheProvider}

	err := h.AddChecks([]*health.Config{
		{
			Name:    "kafka",
			Checker: kafkaCheck,
			Fatal:   false,
		},
		{
			Name:    "postgres",
			Checker: postgresCheck,
			Fatal:   false,
		},
		{
			Name:    "redis",
			Checker: redisCheck,
			Fatal:   false,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao adicionar health checks: %v", err)
	}

	return &HealthChecks{
		Health:        h,
		KafkaCheck:    kafkaCheck,
		PostgresCheck: postgresCheck,
		RedisCheck:    redisCheck,
	}, nil
}
