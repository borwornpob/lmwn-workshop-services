package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

var (
	config   Config
	mapOrder = map[string][]Order{
		"user-a": {
			newOrder("order-A-1", "user-a", []Item{
				newItem("itemA", 10, 2),
				newItem("itemB", 20, 1),
			}),
		},
		"user-b": {
			newOrder("order-B-1", "user-b", []Item{
				newItem("itemA", 10, 2),
				newItem("itemB", 20, 1),
			}),
			newOrder("order-B-2", "user-b", []Item{
				newItem("itemG", 100, 1),
				newItem("itemY", 50, 1),
				newItem("itemZ", 50, 1),
			}),
		},
		"user-c": {
			newOrder("order-C-1", "user-c", []Item{
				newItem("itemA", 10, 2),
			}),
		},
	}
)

func main() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	config = cfg

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/readiness", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/orders/:user-id", getOrderByUserID)

	r.Run(fmt.Sprintf(":%s", config.Port))
}

type Order struct {
	ID     string
	UserID string
	Items  []Item
}

type Item struct {
	Name       string
	Price      float64
	TotalPrice float64
	Qty        int
}

func newOrder(id string, userID string, items []Item) Order {
	return Order{
		ID:     id,
		UserID: userID,
		Items:  items,
	}
}

func newItem(name string, price float64, qty int) Item {
	return Item{
		Name:       name,
		Price:      price,
		Qty:        qty,
		TotalPrice: price * float64(qty),
	}
}

func getOrderByUserID(c *gin.Context) {
	userID := c.Param("user-id")

	userOrders, ok := mapOrder[userID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": userOrders,
	})
}

type Config struct {
	Port string `envconfig:"PORT" default:"5010"`
}

func getConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
