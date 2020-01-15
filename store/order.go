package store

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/jackyczj/July/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/jackyczj/July/cache"
)

/*
	OrderNo 	订单号
	Seller  	卖家
	Buyer		卖家
	Payment 	价格
	PaymentType 支付方式
	ShippingTo	邮寄地址
	IsClose		订单是否已关闭
	CreateAt 	订单创建时间

*/
type Order struct {
	OrderNo      string    `json:"OrderNo" bson:"OrderNo"`
	Seller       string    `json:"seller" bson:"seller,omitempty"`
	Buyer        string    `json:"buyer" bson:"buyer"`
	Payment      int32     `json:"payment" bson:"payment" `
	PaymentType  int       `json:"payment_type" bson:"payment_type,omitempty" `
	ShippingTo   Address   `json:"shipping_to" bson:"shipping_to,omitempty"`
	Item         []Product `json:"item" bson:"item,omitempty"`
	CreateTime   time.Time `json:"create_time" bson:"create_time"`
	IsClose      bool      `json:"is_close"`
	EndTime      time.Time `json:"end_time"`
	SendTime     time.Time `json:"send_time"`
	sync.RWMutex `bson:"-"`
}

var err error

//订单创建
func (o *Order) Create() error {
	o.Lock()
	defer o.Unlock()
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err = Client.db.Collection("Order").InsertOne(ctx, o)
	if err != nil {
		return err
	}
	cache.SetCc("Order."+o.OrderNo, o, 1*time.Hour)
	return nil
}

//订单删除
func (o *Order) Delete() error {
	o.Lock()
	defer o.Unlock()
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	m, err := utils.StructToBson(o)
	if err != nil {
		return err
	}
	result := Client.db.Collection("Order").FindOneAndDelete(ctx, m)
	if result.Err() != nil {
		return result.Err()
	}
	cache.DelCc("Order." + o.OrderNo)
	return nil
}

//订单修改
func (o *Order) Update(filed string, value interface{}) error {
	o.Lock()
	defer o.Unlock()
	op := options.FindOneAndUpdate()
	op.SetProjection(bson.D{{Key: "OrderNo", Value: o.OrderNo}})
	update := bson.D{{Key: "$set", Value: bson.D{{Key: filed, Value: value}}}}
	result := Client.db.Collection("Order").FindOneAndUpdate(context.TODO(), o, update)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}