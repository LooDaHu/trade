package business

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"net/http"
	"trade/message"
	"trade/model"
	"trade/utils"
)

// RetrieveAvailableProductList 查询所有可用产品列表
func RetrieveAvailableProductList(c *gin.Context) {
	utils.SugarLogger.Info("IncrBalance Call")
	products, err := model.RetrieveAllProduct()
	if err != nil || products == nil {
		utils.SugarLogger.Error("MongoDB Ops Failed @RetrieveAvailableProductList", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	data, _ := json.Marshal(products)
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"data":    data,
	})

}

// BuyProduct 用户购买产品
func BuyProduct(c *gin.Context) {
	var info model.Session
	req := new(message.BuyReq)
	utils.SugarLogger.Info("BuyProduct Call", req)
	//校验参数
	if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
		utils.SugarLogger.Error("Params Check Failed @BuyProduct", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数校验失败"})
	}
	session := sessions.Default(c)
	session.Get(info)
	// 开启事务
	opts := options.Session().SetDefaultReadConcern(readconcern.Majority())
	sess, err := utils.GlobalMongoClient.StartSession(opts)
	product, err := model.RetrieveProduct(req.ProductID)
	defer sess.EndSession(context.TODO())
	// 开启new session
	err = mongo.WithSession(context.TODO(), sess, func(sessCtx mongo.SessionContext) error {
		if err := sess.StartTransaction(); err != nil {
			utils.SugarLogger.Error("Session Failed @BuyProduct", err)
			return err
		}
		defer func() {
			if err != nil {
				_ = sess.AbortTransaction(context.Background())
			}
		}()
		totalPrince := float32(req.Amount) * product.Price
		err := model.ChangeBalance(sessCtx, info.User, -totalPrince)
		if err != nil {
			utils.SugarLogger.Error("ChangeBalance Failed @BuyProduct", err)
			return err
		}
		err = model.ChangeProductAmount(sessCtx, req.ProductID, -req.Amount)
		if err != nil {
			utils.SugarLogger.Error("ChangeProductAmount Failed @BuyProduct", err)
			return err
		}
		err = model.AddOrder(sessCtx, &model.Order{
			ProductItem: "",
			TotalPrice:  totalPrince,
			Status:      "未交付",
			UserID:      info.User,
			Created:     timestamp.Timestamp{},
			Updated:     timestamp.Timestamp{},
		})
		if err != nil {
			utils.SugarLogger.Error("AddOrder Failed @BuyProduct", err)
			return err
		}
		return sess.CommitTransaction(context.Background())
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
	})
}
