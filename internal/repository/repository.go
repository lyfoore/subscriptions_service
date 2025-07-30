package repository

import (
	"github.com/google/uuid"
	"github.com/lyfoore/subscriptions_service/internal/domain"
)

type DB interface {
	GetSubscription(id uint64) (*domain.Subscription, error)
	GetSubscriptionsList() ([]domain.Subscription, error)
	GetSubscriptionsAggregate(dateFrom, dateTo string, userID uuid.UUID, serviceName string) (uint64, error)
	CreateSubscription(subscription *domain.Subscription) (uint64, error)
	UpdateSubscription(id uint64, subscription *domain.Subscription) error
	DeleteSubscription(id uint64) error
}
