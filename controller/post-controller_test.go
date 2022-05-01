package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/egnptr/rest-api/cache"
	"github.com/egnptr/rest-api/entity"
	"github.com/egnptr/rest-api/repository"
	"github.com/egnptr/rest-api/service"
)

const (
	Id    int64  = 123
	TITLE string = "title 1"
	TEXT  string = "text 1"
)

var (
	postRepo       repository.PostRepository = repository.NewSQLiteRepository()
	postSrv        service.PostService       = service.NewPostService(postRepo)
	postCacheSrv   cache.PostCache           = cache.NewRedisCache("localhost:6379", 0, 10)
	postController PostController            = NewPostController(postSrv, postCacheSrv)
)

func setup() {
	var post entity.Post = entity.Post{
		Id:    Id,
		Title: TITLE,
		Text:  TEXT,
	}
	postRepo.Save(&post)
}

func tearDown(postId int64) {
	// var post entity.Post = entity.Post{
	// 	Id: postId,
	// }
	postRepo.Delete(postId)
}

func TestAddPost(t *testing.T) {
	// Create new HTTP request
	var jsonStr = []byte(`{"title":"` + TITLE + `","text":"` + TEXT + `"}`)
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))

	// Assing HTTP Request handler Function (controller function)
	handler := http.HandlerFunc(postController.AddPost)
	// Record the HTTP Response
	response := httptest.NewRecorder()
	// Dispatch the HTTP Request
	handler.ServeHTTP(response, req)

	// Assert HTTP status
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode HTTP response
	var post entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&post)

	// Assert HTTP response
	assert.NotNil(t, post.Id)
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, TEXT, post.Text)

	// Cleanup database
	tearDown(post.Id)
}

func TestGetPosts(t *testing.T) {
	// Insert new post
	setup()

	// Create new HTTP request
	req, _ := http.NewRequest("GET", "/posts", nil)

	// Assing HTTP Request handler Function (controller function)
	handler := http.HandlerFunc(postController.GetPost)
	// Record the HTTP Response
	response := httptest.NewRecorder()
	// Dispatch the HTTP Request
	handler.ServeHTTP(response, req)

	// Assert HTTP status
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode HTTP response
	var posts []entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&posts)

	// Assert HTTP response
	assert.Equal(t, Id, posts[0].Id)
	assert.Equal(t, TITLE, posts[0].Title)
	assert.Equal(t, TEXT, posts[0].Text)

	// Cleanup database
	tearDown(Id)
}

func TestGetPostByID(t *testing.T) {
	// Insert new post
	setup()

	// Create new HTTP request
	req, _ := http.NewRequest("GET", "/posts/123", nil)

	// Assing HTTP Request handler Function (controller function)
	handler := http.HandlerFunc(postController.GetPost)
	// Record the HTTP Response
	response := httptest.NewRecorder()
	// Dispatch the HTTP Request
	handler.ServeHTTP(response, req)

	// Assert HTTP status
	status := response.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Decode HTTP response
	var post entity.Post
	json.NewDecoder(io.Reader(response.Body)).Decode(&post)

	// Assert HTTP response
	assert.Equal(t, Id, post.Id)
	assert.Equal(t, TITLE, post.Title)
	assert.Equal(t, TEXT, post.Text)

	// Cleanup database
	tearDown(Id)
}
