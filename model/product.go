package model

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"trade/utils"
)

type Product struct {
	ID          primitive.ObjectID  `json:"id"`
	Name        string              `json:"name"`        // 商品名称
	Title       string              `json:"title"`       // 商品标题
	Description string              `json:"description"` // 商品描述
	Price       float32             `json:"price"`       // 商品价格
	Amount      int                 `json:"amount"`      // 商品数量
	Sales       int                 `json:"sales"`       // 商品销量=
	Created     timestamp.Timestamp `json:"created"`     // 创建时间'
	Updated     timestamp.Timestamp `json:"updated"`     // 更新时间'
}

const (
	KMongoProductCollection = "Product"
)

// ChangeProductAmount 修改产品总量
func ChangeProductAmount(ctx context.Context, productID string, value int) error {
	product, err := RetrieveProduct(productID)
	filter := bson.M{
		"_id": product.ID,
	}
	if value < 0 && product.Amount+value <= 0 {
		return fmt.Errorf("run out of products, ID :%s", product.ID.Hex())
	}
	update := bson.M{"Inc": bson.M{
		"Amount": value,
	}}
	_, err = utils.GlobalDatabase.Collection(KMongoProductCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveProduct 查询指定产品相关信息
func RetrieveProduct(productID string) (Product, error) {
	var product Product
	oID, oErr := primitive.ObjectIDFromHex(productID)
	if oErr != nil {
		return Product{}, oErr
	}
	filter := bson.M{
		"_id": oID,
	}
	res := utils.GlobalDatabase.Collection(KMongoProductCollection).FindOne(context.TODO(), filter)
	err := res.Decode(&product)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return Product{}, nil
		}
		return Product{}, err
	}
	return product, nil
}

// RetrieveAllProduct 查询所用产品相关信息
func RetrieveAllProduct() ([]*Product, error) {
	var products []*Product
	filter := bson.M{
		"amount": bson.M{"$gte": 0},
	}
	res, err := utils.GlobalDatabase.Collection(KMongoProductCollection).Find(context.TODO(), filter)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer res.Close(context.TODO())
	err = res.All(context.TODO(), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
