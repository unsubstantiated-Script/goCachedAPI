package cache

import (
	"context"
	"fmt"
	"goCachedAPI/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ProductCache struct {
	client        *redis.Client
	productTTL    time.Duration
	recentListKey string
	recentLimit   int64
}

func NewProductCache(
	client *redis.Client,
	productTTL time.Duration,
	recentListKey string,
	recentLimit int64,
) *ProductCache {
	return &ProductCache{
		client:        client,
		productTTL:    productTTL,
		recentListKey: recentListKey,
		recentLimit:   recentLimit,
	}
}

func productKey(id uint) string {
	return fmt.Sprintf("product:%d", id)
}

func (c *ProductCache) GetProduct(ctx context.Context, id uint) (models.Product, bool, error) {
	key := productKey(id)

	res, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return models.Product{}, false, err
	}

	if len(res) == 0 {
		return models.Product{}, false, nil
	}

	price, err := strconv.Atoi(res["price"])
	if err != nil {
		return models.Product{}, false, err
	}

	product := models.Product{
		ID:    id,
		Name:  res["name"],
		Price: price,
	}

	return product, true, nil
}

func (c *ProductCache) SetProduct(ctx context.Context, product models.Product) error {
	key := productKey(product.ID)

	if err := c.client.HSet(ctx, key, map[string]any{
		"name":  product.Name,
		"price": product.Price,
	}).Err(); err != nil {
		return err
	}

	return c.client.Expire(ctx, key, c.productTTL).Err()
}

func (c *ProductCache) DeleteProduct(ctx context.Context, id uint) error {
	return c.client.Del(ctx, productKey(id)).Err()
}

func (c *ProductCache) AddRecentProduct(ctx context.Context, id uint) error {
	if err := c.client.LPush(ctx, c.recentListKey, id).Err(); err != nil {
		return err
	}

	return c.client.LTrim(ctx, c.recentListKey, 0, c.recentLimit-1).Err()
}

func (c *ProductCache) GetRecentProductIDs(ctx context.Context) ([]uint, error) {
	rawIDs, err := c.client.LRange(ctx, c.recentListKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]uint, len(rawIDs))

	for _, raw := range rawIDs {
		n, err := strconv.Atoi(raw)
		if err != nil {
			continue
		}
		ids = append(ids, uint(n))
	}

	return ids, nil
}

func (c *ProductCache) UpdateProductWithTransaction(
	ctx context.Context,
	id uint,
	name string,
	price int,
) error {
	key := productKey(id)

	_, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, map[string]any{
			"name":  name,
			"price": price,
		})
		pipe.Expire(ctx, key, c.productTTL)
		return nil
	})

	return err
}
