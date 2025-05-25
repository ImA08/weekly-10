package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
	"minitask1.go/internal/utils"
	"minitask1.go/pkg"
)

type OrderHandler struct {
	OrderRepo *repositories.OrderRepository
}

func NewOrderHandler(OrderRepo *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{OrderRepo: OrderRepo}
}

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	claims, _ := ctx.Get("Payload")
	userClaims := claims.(*pkg.Claims)

	var formBody models.OrderStruct
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Printf("Binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	invoiceNumb, err := utils.GenerateNumber()
	if err != nil {
		log.Printf("[DEBUG] generate number")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	cmd, err := o.OrderRepo.SetOrder(ctx.Request.Context(), formBody, userClaims.Id, string(invoiceNumb))

	if err != nil {
		log.Println("[DEBUG] TRANSACTIONS", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}

	if cmd.RowsAffected() == 0 {
		log.Println("Query Gagal, Tidak merubah data di DB")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Data yang diberikan salah",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Order Succes",
	})
	// 1. id user
	// 2. method
	// 3. seat

}

func (o *OrderHandler) FindOrderById(ctx *gin.Context) {
	idStr, ok := ctx.Params.Get("orderId")
	log.Println("[DEBUG] ID DEBUGGING", idStr)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Need param id",
		})
		return
	}
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server Error",
		})
		return
	}

	result, err := o.OrderRepo.GetOrder(ctx.Request.Context(), idInt)

	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "Internal server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Succes",
		"user": result,
	})
}

func (o *OrderHandler) UpdateOrderHandler(Ctx *gin.Context){
	
}