package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
	"github.com/shkiperko0/auth-go-ms/models"
)

type ITokenRepository interface {
	Create(item *models.Token, duration time.Duration) error
	Update(item *models.Token, duration time.Duration) error
	Get(id string, userID uint) (*models.Token, error)
	UpdateAccessToken(userID uint, accessToken string, duration time.Duration) error
	Remove(token *models.Token) error
	RemoveAllOfUser(token *models.Token) error
}

type TokenRepository struct {
	Prefix       string
	Ctx          context.Context
	Client       *redis.Client
	defaultLimit int
}

func (r *TokenRepository) getKey(id string, userID uint) string {
	return fmt.Sprintf("%s:%d:%s", r.Prefix, userID, id)
}

func (r *TokenRepository) Remove(token *models.Token) error {
	return r.Client.Del(r.Ctx, r.getKey(token.ID, token.UserID)).Err()
}

func (r *TokenRepository) RemoveAllOfUser(token *models.Token) error {
	keys, err := r.Client.Keys(r.Ctx, r.getKey("*", token.UserID)).Result()

	if err != nil {
		return err
	}

	return r.Client.Del(r.Ctx, keys...).Err()
}

func (r *TokenRepository) Create(value *models.Token, duration time.Duration) error {
	if len(value.ID) == 0 {
		uuid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		value.ID = uuid.String()
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(r.Ctx, r.getKey(value.ID, value.UserID), payload, duration).Err()
}

func (r *TokenRepository) Update(value *models.Token, duration time.Duration) error {
	return r.Create(value, duration)
}

func (r *TokenRepository) Get(id string, userID uint) (*models.Token, error) {
	var item models.Token
	res := r.Client.Get(r.Ctx, r.getKey(id, userID))
	if res.Err() != nil {
		return nil, res.Err()
	}

	data, err := res.Bytes()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &item)
	if err != nil {
		return nil, err
	}

	return &item, err
}

func (r *TokenRepository) UpdateAccessToken(userID uint, accessToken string, duration time.Duration) error {
	list, err := r.Client.Keys(r.Ctx, r.getKey("*", userID)).Result()
	if err != nil {
		return err
	}

	for _, data := range list {
		var item models.Token

		res := r.Client.Get(r.Ctx, data)
		if res.Err() != nil {
			return res.Err()
		}

		dataBytes, err := res.Bytes()
		if err != nil {
			return err
		}

		err = json.Unmarshal(dataBytes, &item)
		if err != nil {
			return err
		}

		item.AccessToken = accessToken
		err = r.Update(&item, duration)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewTokenRepository(client *redis.Client, defaultLimit int) ITokenRepository {
	return &TokenRepository{
		Prefix:       "tokens",
		Ctx:          context.Background(),
		Client:       client,
		defaultLimit: defaultLimit,
	}
}
