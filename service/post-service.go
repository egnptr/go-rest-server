package service

import (
	"errors"
	"math/rand"

	"github.com/egnptr/rest-api/entity"
	"github.com/egnptr/rest-api/repository"
)

type PostService interface {
	Validate(post *entity.Post) error
	FindAll() ([]entity.Post, error)
	Create(post *entity.Post) (*entity.Post, error)
	Delete(id int64) error
	Get(id int64) (*entity.Post, error)
}

type service struct{}

var (
	repo repository.PostRepository
)

func NewPostService(repository repository.PostRepository) PostService {
	repo = repository
	return &service{}
}

func (s *service) Validate(post *entity.Post) error {
	if post == nil {
		err := errors.New("The post is empty")
		return err
	}
	if post.Title == "" {
		err := errors.New("The post title is empty")
		return err
	}
	return nil
}

func (s *service) Create(post *entity.Post) (*entity.Post, error) {
	post.Id = rand.Int63()
	return repo.Save(post)
}

func (s *service) FindAll() ([]entity.Post, error) {
	return repo.FindAll()
}

func (s *service) Delete(id int64) error {
	return repo.Delete(id)
}

func (s *service) Get(id int64) (*entity.Post, error) {
	return repo.Get(id)
}
