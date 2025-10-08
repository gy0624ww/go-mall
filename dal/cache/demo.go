package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-study-lab/go-mall/common/enum"
	"github.com/go-study-lab/go-mall/common/logger"
	"github.com/go-study-lab/go-mall/logic/do"
)

type DummyDemoOrder struct {
	OrderNo string `redis:"orderNo"`
	UserId  int64  `redis:"userId"`
}

func SetDemoOrderStruct(ctx context.Context, demoOrder *do.DemoOrder) error {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	data := struct {
		OrderNo string `redis:"orderNo"`
		UserId  int64  `redis:"userId"`
	}{
		UserId:  demoOrder.UserId,
		OrderNo: demoOrder.OrderNo,
	}
	_, err := Redis().HSet(ctx, redisKey, data).Result()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return err
	}
	return nil
}

func GetDemoOrderStruct(ctx context.Context, orderNo string) (*DummyDemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	data := new(DummyDemoOrder)
	err := Redis().HGetAll(ctx, redisKey).Scan(&data)
	Redis().Get(ctx, redisKey).String()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return nil, err
	}
	logger.Info(ctx, "scan data from redis", "data", &data)
	return data, nil
}

func SetDemoOrder(ctx context.Context, demoOrder *do.DemoOrder) error {
	jsonDataBytes, _ := json.Marshal(demoOrder)
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, demoOrder.OrderNo)
	_, err := Redis().Set(ctx, redisKey, jsonDataBytes, 0).Result()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return err
	}
	return nil
}

func GetDemoOrder(ctx context.Context, orderNo string) (*do.DemoOrder, error) {
	redisKey := fmt.Sprintf(enum.REDIS_KEY_DEMO_ORDER_DETAIL, orderNo)
	jsonBytes, err := Redis().Get(ctx, redisKey).Bytes()
	if err != nil {
		logger.Error(ctx, "redis error", "err", err)
		return nil, err
	}
	data := new(do.DemoOrder)
	json.Unmarshal(jsonBytes, &data)
	return data, nil
}
