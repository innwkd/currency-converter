package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"

	"github.com/yddmat/currency-converter/app/types"
)

const pattern = "cur:%s"

type RedisStorage struct {
	conn *redis.Client
}

func NewRedisStorage(conn *redis.Client) *RedisStorage {
	return &RedisStorage{conn: conn}
}

func (re *RedisStorage) Set(pair types.CurrencyPair, rate types.CurrencyRate, duration time.Duration) (types.CurrencyRate, error) {
	jdata, err := json.Marshal(rate)
	if err != nil {
		return types.CurrencyRate{}, errors.Wrapf(err, "can't marshall rate to json")
	}

	key := fmt.Sprintf(pattern, pair)

	pipe := re.conn.Pipeline()
	pipe.SetNX(key, jdata, duration)
	get := pipe.Get(key)
	if _, err := pipe.Exec(); err != nil {
		return types.CurrencyRate{}, errors.Wrapf(err, "can't exec pipe in redis")
	}

	// All errors would appear in exec func above
	rateValue, _ := get.Result()
	rate = types.CurrencyRate{}
	if err := json.Unmarshal([]byte(rateValue), &rate); err != nil {
		return types.CurrencyRate{}, errors.Wrapf(err, "can't unmarshal rate to struct")
	}

	return rate, nil
}

func (re *RedisStorage) Get(pair types.CurrencyPair) (types.CurrencyRate, error) {
	rate := types.CurrencyRate{}
	value, err := re.conn.Get(fmt.Sprintf(pattern, pair.String())).Result()
	if err != nil {
		if err == redis.Nil {
			return rate, ErrNotExists
		}

		return rate, err
	}

	if err := json.Unmarshal([]byte(value), &rate); err != nil {
		return rate, errors.Wrapf(err, "can't unmarshal rate to struct")
	}

	return rate, nil
}

func (re *RedisStorage) GetAll() ([]types.CurrencyRate, error) {
	keys := make([]string, 0)
	iter := re.conn.Scan(0, fmt.Sprintf(pattern, "*"), 100).Iterator()
	for iter.Next() {
		if iter.Err() != nil {
			return nil, errors.Wrap(iter.Err(), "can't scan keys in redis")
		}
		keys = append(keys, iter.Val())
	}

	if len(keys) == 0 {
		return []types.CurrencyRate{}, nil
	}

	values, err := re.conn.MGet(keys...).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "can't exec batch get in redis")
	}

	rates := make([]types.CurrencyRate, 0)
	for _, value := range values {
		rate := types.CurrencyRate{}
		if err := json.Unmarshal([]byte(value.(string)), &rate); err != nil {
			return nil, errors.Wrapf(err, "can't unmarshal value from redis to struct")
		}

		rates = append(rates, rate)
	}

	return rates, nil
}
