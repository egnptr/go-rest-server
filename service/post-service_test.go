package service

import (
	"testing"

	"github.com/egnptr/rest-api/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) Save(post *entity.Post) (*entity.Post, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.(*entity.Post), args.Error(1)
}

func (mock *MockRepository) FindAll() ([]entity.Post, error) {
	args := mock.Called()
	result := args.Get(0)
	return result.([]entity.Post), args.Error(1)
}

func (mock *MockRepository) Delete(id int64) error {
	args := mock.Called(id)
	return args.Error(1)
}

func (mock *MockRepository) Get(id int64) (*entity.Post, error) {
	args := mock.Called(id)
	result := args.Get(0)
	return result.(*entity.Post), args.Error(1)
}

func TestValidateEmptyPost(t *testing.T) {
	testService := NewPostService(nil)

	err := testService.Validate(nil)

	assert.NotNil(t, err)
	assert.Equal(t, "The post is empty", err.Error())
}

func TestValidateEmptyTitle(t *testing.T) {
	post := entity.Post{Id: 1, Title: "", Text: "Text"}
	testService := NewPostService(nil)
	err := testService.Validate(&post)

	assert.NotNil(t, err)
	assert.Equal(t, "The post title is empty", err.Error())
}

func TestCreatePost(t *testing.T) {
	mockRepo := new(MockRepository)

	post := entity.Post{Title: "Title", Text: "Text"}
	mockRepo.On("Save").Return(&post, nil)

	testService := NewPostService(mockRepo)

	result, err := testService.Create(&post)

	mockRepo.AssertExpectations(t)

	assert.NotNil(t, result.Id)
	assert.Equal(t, "Title", result.Title)
	assert.Equal(t, "Text", result.Text)
	assert.Nil(t, err)
}

func TestFindAll(t *testing.T) {
	mockRepo := new(MockRepository)

	var identifier int64 = 1

	post := entity.Post{Id: identifier, Title: "Title", Text: "Text"}
	// Setup expectation
	mockRepo.On("FindAll").Return([]entity.Post{post}, nil)

	testService := NewPostService(mockRepo)

	result, _ := testService.FindAll()

	// Mock assertion
	mockRepo.AssertExpectations(t)

	// Data assertion
	assert.Equal(t, identifier, result[0].Id)
	assert.Equal(t, "Title", result[0].Title)
	assert.Equal(t, "Text", result[0].Text)
}
