package user

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Jhon-2801/courses-domain/domain"
	"github.com/Jhon-2801/lib-response/response"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *domain.User) error
	GetAll(ctx context.Context, filters Filters, offsite, limit int) ([]domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(ctx context.Context, filters Filters) (int, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repo{
		db: db,
	}
}

func (repo *repo) Create(ctx context.Context, user *domain.User) error {

	if err := repo.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (repo *repo) GetAll(ctx context.Context, filters Filters, offsite, limit int) ([]domain.User, error) {
	var u []domain.User
	tx := repo.db.WithContext(ctx).Model(&u)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offsite)
	result := tx.Order("created_at desc").Find(&u)
	if result.Error != nil {
		log.Fatal(result.Error)
		return nil, result.Error
	}
	return u, nil
}

func (repo *repo) Get(ctx context.Context, id string) (*domain.User, error) {
	user := domain.User{ID: id}
	if err := repo.db.WithContext(ctx).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *repo) Delete(ctx context.Context, id string) error {
	user := domain.User{ID: id}

	result := repo.db.WithContext(ctx).Delete(&user)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return response.NotFound("user doesn't exists")
	}
	return nil
}

func (repo *repo) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
	values := make(map[string]interface{})
	if firstName != nil {
		values["first_name"] = *firstName
	}
	if lastName != nil {
		values["last_name"] = *lastName
	}
	if email != nil {
		values["email"] = *email
	}
	if phone != nil {
		values["phone"] = *phone
	}

	result := repo.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return response.NotFound("user doesn't exists")
	}

	return nil
}
func (repo *repo) Count(ctx context.Context, filters Filters) (int, error) {
	var count int64
	tx := repo.db.WithContext(ctx).Model(domain.User{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, nil
	}
	return int(count), nil
}
func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(first_name) like ?", filters.FirstName)
	}
	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(last_name) like ?", filters.LastName)
	}
	return tx

}
