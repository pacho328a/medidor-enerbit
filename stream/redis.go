package stream

import (
	"github.com/go-redis/redis"
	"log"
	"medidor_enerbit/models"
)

func GetRedis() *redis.Client {
	log.Println("Create client")

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func SendStreamMedidor(medidor models.Medidor, client *redis.Client) error {
	log.Println("Publishing event to Redis")

	err := client.XAdd(&redis.XAddArgs{
		Stream: "tickets",
		Values: map[string]interface{}{
			"brand":   string(medidor.Brand),
			"address": string(medidor.Address),
			"serial":  string(medidor.Serial),
		},
	}).Err()

	return err
}
