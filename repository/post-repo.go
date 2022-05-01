package repository

import "github.com/egnptr/rest-api/entity"

type PostRepository interface {
	Save(post *entity.Post) (*entity.Post, error)
	FindAll() ([]entity.Post, error)
	Delete(id int64) error
	Get(id int64) (*entity.Post, error)
}
