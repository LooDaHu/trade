package model

import (
	"context"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"trade/utils"
)

type Order struct {
	ID          primitive.ObjectID  `json:"id"`           // 订单id
	ProductItem string              `json:"product_item"` // 商品项
	TotalPrice  float32             `json:"total_price"`  // 合计
	Status      string              `json:"status"`       // 订单状态
	AddressID   string              `json:"address_id"`   // 地址id
	UserID      string              `json:"user_id"`      // 用户id
	Nickname    string              `json:"nick_name"`    // 用户昵称,
	Created     timestamp.Timestamp `json:"created"`      // 创建时间
	Updated     timestamp.Timestamp `json:"updated"`      // 更新时间
}

const (
	KMongoOrderCollection = "Order"
)

// AddOrder 新增订单
func AddOrder(ctx context.Context, order *Order) error {
	_, err := utils.GlobalDatabase.Collection(KMongoOrderCollection).InsertOne(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveOrder 查询订单
func RetrieveOrder(id string) ([]*Order, error) {
	var orders []*Order
	oID, oErr := primitive.ObjectIDFromHex(id)
	if oErr != nil {
		return nil, oErr
	}
	filter := bson.M{
		"_id": oID,
	}
	res, err := utils.GlobalDatabase.Collection(KMongoProductCollection).Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer res.Close(context.TODO())
	err = res.All(context.TODO(), &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
