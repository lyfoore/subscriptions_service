package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lyfoore/subscriptions_service/internal/domain"
	"github.com/lyfoore/subscriptions_service/internal/service"
	"log"
	"strconv"
	"time"
)

// SubscriptionDTO represents subscription data transfer object
// swagger:model SubscriptionDTO
type SubscriptionDTO struct {
	ID          uint64 `json:"id"`
	UserID      string `json:"user_id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

func (s *SubscriptionDTO) convertToDomain() (*domain.Subscription, error) {
	sub := &domain.Subscription{}

	if s.UserID != "" {
		userID, err := uuid.Parse(s.UserID)
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %v", err)
		}
		sub.UserID = userID
	}

	if s.StartDate != "" {
		startDate, err := time.Parse("01-2006", s.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format, expected MM-YYYY: %v", err)
		}
		sub.StartDate = startDate
	}

	if s.EndDate != "" {
		endDate, err := time.Parse("01-2006", s.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY: %v", err)
		}
		sub.EndDate = endDate
	}

	if s.ServiceName != "" {
		sub.ServiceName = s.ServiceName
	}

	if s.Price != 0 {
		sub.Price = s.Price
	}

	return sub, nil
}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetSubscription godoc
// @Summary Get a subscription by ID
// @Description Get a single subscription by its ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while parsing id": err.Error(),
		})
		log.Println("error while parsing id:", err)
		return
	}
	sub, err := h.service.GetSubscription(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while getting subscription": err.Error(),
		})
		log.Println("error while getting subscription:", err)
		return
	}
	ctx.JSON(200, sub)
}

// GetSubscriptionsList godoc
// @Summary Get all subscriptions
// @Description Get a list of all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} domain.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) GetSubscriptionsList(ctx *gin.Context) {
	list, err := h.service.GetSubscriptionsList()
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while getting subscriptions list": err.Error(),
		})
		log.Println("error while getting subscriptions list:", err)
		return
	}
	ctx.JSON(200, list)
}

// GetSubscriptionsAggregate godoc
// @Summary Get aggregated subscriptions data
// @Description Get sum of subscriptions filtered by date range, user and/or service
// @Tags subscriptions
// @Produce json
// @Param from query string false "Start date in MM-YYYY format"
// @Param to query string false "End date in MM-YYYY format"
// @Param user query string false "User ID"
// @Param service query string false "Service name"
// @Success 200 {object} map[string]int "sum field contains the aggregated value"
// @Failure 400 {object} map[string]string
// @Router /subscriptions/aggregate [get]
func (h *Handler) GetSubscriptionsAggregate(ctx *gin.Context) {
	dateFrom := ctx.Query("from")
	dateTo := ctx.Query("to")
	userID := ctx.Query("user")
	serviceName := ctx.Query("service")
	if dateFrom == "" && dateTo == "" {
		ctx.JSON(400, gin.H{
			"error while getting subscriptions aggregate": "no params",
		})
		log.Println("error while getting subscriptions aggregate: no params")
		return
	}

	subsDTO := SubscriptionDTO{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   dateFrom,
		EndDate:     dateTo,
	}

	subscription, err := subsDTO.convertToDomain()
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while converting subscription": err.Error(),
		})
		log.Println("error while converting subscription:", err)
		return
	}

	sum, err := h.service.GetSubscriptionsAggregate(subscription)
	//sum, err := h.service.GetSubscriptionsAggregate(subscription.StartDate.Format("2006-01-02"),
	//	subscription.EndDate.Format("2006-01-02"), subscription.UserID, subscription.ServiceName) //dateFrom, dateTo, userID, serviceName)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while getting sum of subscriptions": err.Error(),
		})
		log.Println("error while getting sum of subscriptions:", err)
		return
	}
	ctx.JSON(200, gin.H{
		"sum": sum,
	})
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription with the provided data
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body SubscriptionDTO true "Subscription data"
// @Success 200 {object} map[string]uint64 "id field contains the ID of created subscription"
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(ctx *gin.Context) {
	subscriptionDTO := SubscriptionDTO{}
	if err := ctx.ShouldBindJSON(&subscriptionDTO); err != nil {
		ctx.JSON(400, gin.H{
			"error while binding subscription": err.Error(),
		})
		log.Println("error while binding subscription:", err)
		return
	}

	subscription, err := subscriptionDTO.convertToDomain()
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while converting subscription": err.Error(),
		})
		log.Println("error while converting subscription:", err)
		return
	}

	id, err := h.service.CreateSubscription(subscription)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while creating subscription": err.Error(),
		})
		log.Println("error while creating subscription:", err)
		return
	}
	ctx.JSON(200, gin.H{
		"id": id,
	})
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update an existing subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body SubscriptionDTO true "Subscription data"
// @Success 200 {object} map[string]uint64 "id field contains the ID of updated subscription"
// @Failure 400 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while parsing id": err.Error(),
		})
		log.Println("error while parsing id:", err)
		return
	}
	var subscriptionDTO SubscriptionDTO
	if err := ctx.ShouldBindJSON(&subscriptionDTO); err != nil {
		ctx.JSON(400, gin.H{
			"error while binding subscription": err.Error(),
		})
		log.Println("error while binding subscription:", err)
		return
	}

	subscription, err := subscriptionDTO.convertToDomain()
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while converting subscription": err.Error(),
		})
		log.Println("error while converting subscription:", err)
		return
	}

	err = h.service.UpdateSubscription(id, subscription)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while updating subscription": err.Error(),
		})
		log.Println("error while updating subscription:", err)
		return
	}
	ctx.JSON(200, gin.H{
		"id": id,
	})
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]uint64 "id field contains the ID of deleted subscription"
// @Failure 400 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while parsing id": err.Error(),
		})
		log.Println("error while parsing id:", err)
		return
	}
	err = h.service.DeleteSubscription(id)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error while deleting subscription": err.Error(),
		})
		log.Println("error while deleting subscription:", err)
		return
	}
	ctx.JSON(200, gin.H{
		"id": id,
	})
}
