package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"trade/utils"
)

type User struct {
	ID       primitive.ObjectID  `bson:"_id" json:"_id"`             // 订单id
	Username string              `bson:"username" json:"username"`   //用户名称
	RealName string              `bson:"real_name" json:"real_name"` // 用户真实名称
	RoleID   int                 `bson:"role_id" json:"role_id"`     //用户角色，1表示普通用户
	Password string              `bson:"password" json:"password"`   //用户密码
	Phone    string              `bson:"phone" json:"phone"`         // 用户电话
	Balance  float32             `bson:"balance" json:"balance"`     // 用户余额
	Status   int                 `bson:"status" json:"status"`       // 用户状态，1表示正常，0表示暂停
	Created  primitive.Timestamp `bson:"created" json:"created"`     // 创建时间
	Updated  primitive.Timestamp `bson:"updated" json:"updated"`     // 更新时间
}

type Session struct {
	User string `json:"user"`
}

const (
	KMongoUserCollection = "User"
	NormalStatus         = 1
)

// CheckUserInfo 检查用户信息
func CheckUserInfo(username string, password string) (bool, string, error) {
	var user User
	filter := bson.M{
		"username": username,
		"password": password,
		"status":   NormalStatus,
	}
	res := utils.GlobalDatabase.Collection(KMongoUserCollection).FindOne(context.TODO(), filter)
	err := res.Decode(&user)
	if err != nil {
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			return false, "", nil
		}
		return false, "", err
	}
	return true, user.ID.Hex(), nil
}

// RetrieveUserInfo 查询用户信息
func RetrieveUserInfo(id string) (User, error) {
	var user User
	oID, oErr := primitive.ObjectIDFromHex(id)
	if oErr != nil {
		return User{}, oErr
	}
	filter := bson.M{
		"_id":    oID,
		"status": NormalStatus,
	}
	res := utils.GlobalDatabase.Collection(KMongoUserCollection).FindOne(context.TODO(), filter)
	err := res.Decode(&user)
	return user, err
}

// ChangeBalance 修改用户余额
func ChangeBalance(ctx context.Context, id string, val float32) error {
	user, err := RetrieveUserInfo(id)
	if val < 0 && user.Balance+val <= 0 {
		return fmt.Errorf("run out of balance, ID :%s", user.ID.Hex())
	}
	filter := bson.M{
		"_id":    user.ID,
		"status": NormalStatus,
	}
	update := bson.M{
		"$inc": bson.M{
			"balance": val,
		},
	}
	res, err := utils.GlobalDatabase.Collection(KMongoUserCollection).
		UpdateOne(ctx, filter, update)
	if err != nil || res.MatchedCount == 0 {
		return err
	}
	return nil
}
