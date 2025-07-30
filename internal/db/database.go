package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/lyfoore/subscriptions_service/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = fmt.Errorf("error while connecting to database: %v", err)
		return nil, err
	}

	if err := db.AutoMigrate(&domain.Subscription{}); err != nil {
		err = fmt.Errorf("error while migrating database: %v", err)
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) GetSubscription(id uint64) (*domain.Subscription, error) {
	var subscription domain.Subscription
	err := p.DB.First(&subscription, id).Error
	if err != nil {
		err = fmt.Errorf("error while getting subscription from database: %v", err)
		return nil, err
	}
	return &subscription, nil
}

func (p *PostgresDB) GetSubscriptionsList() ([]domain.Subscription, error) {
	var subscriptions []domain.Subscription
	err := p.DB.Find(&subscriptions).Error
	if err != nil {
		err = fmt.Errorf("error while getting subscriptions list: %v", err)
		return nil, err
	}
	return subscriptions, nil
}

func (p *PostgresDB) GetSubscriptionsAggregate(dateFrom, dateTo string, userID uuid.UUID, serviceName string) (uint64, error) {
	var sum uint64
	query := `
        SELECT COALESCE(SUM(price), 0) AS sum
		FROM subscriptions 
		WHERE start_date >= ? 
    	AND (end_date <= ? OR end_date IS NULL)
    	AND (user_id = ? OR ? = '00000000-0000-0000-0000-000000000000')
    	AND (service_name = ? OR ? = '')`

	err := p.DB.Raw(query,
		dateFrom, dateTo,
		userID, userID,
		serviceName, serviceName).Scan(&sum).Error

	log.Println(dateFrom, dateTo, userID, serviceName)

	if err != nil {
		err = fmt.Errorf("error while getting sum of subscriptions: %v", err)
		return 0, err
	}
	return sum, nil
}

func (p *PostgresDB) CreateSubscription(subscription *domain.Subscription) (uint64, error) {
	err := p.DB.Create(subscription).Error
	if err != nil {
		err = fmt.Errorf("error while creating subscription: %v", err)
		return 0, err
	}
	return subscription.ID, nil
}

func (p *PostgresDB) UpdateSubscription(id uint64, subscription *domain.Subscription) error {
	subscription.ID = id
	err := p.DB.Model(&domain.Subscription{}).Where("id = ?", id).Updates(subscription).Error
	if err != nil {
		err = fmt.Errorf("error while updating subscription: %v", err)
		return err
	}
	return nil
}

func (p *PostgresDB) DeleteSubscription(id uint64) error {
	err := p.DB.Delete(&domain.Subscription{}, id).Error
	if err != nil {
		err = fmt.Errorf("error while deleting subscription: %v", err)
		return err
	}
	return nil
}
