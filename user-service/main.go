package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

var (
	config   Config
	mapUsers = map[string]User{
		"user-a": {
			ID:   "user-a",
			Name: "A",
		},
		"user-b": {
			ID:   "user-b",
			Name: "B",
		},
		"user-c": {
			ID:   "user-c",
			Name: "C",
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

	r.GET("/users", getUsers)
	r.GET("/user/:user-id", getUserByID)

	r.Run(fmt.Sprintf(":%s", config.Port))
}

type User struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Orders []Order `json:"orders,omitempty"`
}

func getUsers(c *gin.Context) {
	users := make([]User, 0, len(mapUsers))
	for _, u := range mapUsers {
		users = append(users, u)
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func getUserByID(c *gin.Context) {
	userID := c.Param("user-id")

	user, ok := mapUsers[userID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "user not found",
		})
		return
	}

	orders, err := getOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	user.Orders = orders

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
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

func getOrdersByUserID(userID string) ([]Order, error) {
	requestURL := fmt.Sprintf("%s/orders/%s", config.OrderServiceEndpoint, userID)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return []Order{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Order{}, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []Order{}, err
	}

	var getOrdersResponse GetOrdersResponse
	if err := json.Unmarshal(resBody, &getOrdersResponse); err != nil {
		return []Order{}, err
	}

	return getOrdersResponse.Orders, nil
}

type Config struct {
	Port                 string `envconfig:"PORT" default:"5011"`
	OrderServiceEndpoint string `envconfig:"ORDER_SERVICE_ENDPOINT" default:"http://localhost:5010"`
}

func getConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
