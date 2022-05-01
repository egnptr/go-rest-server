package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/egnptr/rest-api/cache"
	"github.com/egnptr/rest-api/entity"
	"github.com/egnptr/rest-api/errors"
	"github.com/egnptr/rest-api/service"
	"github.com/gorilla/mux"
)

var (
	postService service.PostService
	postCache   cache.PostCache
)

type controller struct{}

type PostController interface {
	GetAllPost(w http.ResponseWriter, r *http.Request)
	AddPost(w http.ResponseWriter, r *http.Request)
	GetPost(w http.ResponseWriter, r *http.Request)
	DeletePost(w http.ResponseWriter, r *http.Request)
}

func NewPostController(service service.PostService, cache cache.PostCache) PostController {
	postService = service
	postCache = cache
	return &controller{}
}

func (*controller) GetAllPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	posts, err := postService.FindAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error getting the post"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (*controller) AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var post entity.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error unmarshaling the request"})
		return
	}

	if err1 := postService.Validate(&post); err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: err1.Error()})
		return
	}

	res, err2 := postService.Create(&post)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error saving the post"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (*controller) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := postService.Get(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: err.Error()})
		return
	}

	if err1 := postService.Delete(int64(id)); err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Error deleting the post"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (*controller) GetPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	postID := strings.Split(r.URL.Path, "/")[2]
	id, err1 := strconv.ParseInt(postID, 10, 64)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errors.ServiceError{Message: "Post ID cannot be coverted to int64"})
		return
	}

	var post *entity.Post = postCache.Get(postID)
	if post == nil {
		post, err := postService.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errors.ServiceError{Message: "This post doesn't exist"})
			return
		}

		postCache.Set(postID, post)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}
