package routes

import (
	"net/http"
	"rentx/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := order.Create(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create order"})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func getOrderByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order, err := models.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func listOrders(c *gin.Context) {
	orders, err := models.ListOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func deleteOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	order := models.Order{Id: id}
	if err := order.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order and items deleted"})
}
