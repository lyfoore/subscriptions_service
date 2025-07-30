package service

import (
	"errors"
	"github.com/lyfoore/subscriptions_service/internal/domain"
	"github.com/lyfoore/subscriptions_service/internal/repository"
)

type Service struct {
	db repository.DB
}

func NewService(db repository.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetSubscription(id uint64) (*domain.Subscription, error) {
	return s.db.GetSubscription(id)
}

func (s *Service) GetSubscriptionsList() ([]domain.Subscription, error) {
	return s.db.GetSubscriptionsList()
}

func (s *Service) GetSubscriptionsAggregate(subscription *domain.Subscription) (uint64, error) {
	dateFrom := subscription.StartDate.Format("2006-01-02")
	dateTo := subscription.EndDate.Format("2006-01-02")
	userID := subscription.UserID
	serviceName := subscription.ServiceName
	return s.db.GetSubscriptionsAggregate(dateFrom, dateTo, userID, serviceName)
}

func (s *Service) CreateSubscription(subscription *domain.Subscription) (uint64, error) {
	if subscription.Price < 0 {
		return 0, errors.New("price cannot be negative")
	}
	return s.db.CreateSubscription(subscription)
}

func (s *Service) UpdateSubscription(id uint64, subscription *domain.Subscription) error {
	return s.db.UpdateSubscription(id, subscription)
}

func (s *Service) DeleteSubscription(id uint64) error {
	return s.db.DeleteSubscription(id)
}
