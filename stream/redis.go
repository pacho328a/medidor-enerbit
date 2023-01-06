package stream

import (
	"log"
	"medidor_enerbit/models"
	utils "medidor_enerbit/utils"

	"github.com/go-redis/redis"
)

var host string
var port string
var stream string

func init() {
	host = utils.GetEnvVar("REDIS_HOST")
	port = utils.GetEnvVar("REDIS_PORT")
	stream = "Medidor"
}

func GetRedis() *redis.Client {
	log.Println("Create client")

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	return client
}

func SendStreamMedidor(medidor models.Medidor, client *redis.Client) error {
	log.Println("Publishing event to Redis")

	err := client.XAdd(&redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"ID": string(medidor.ID),
		},
	}).Err()

	return err
}
