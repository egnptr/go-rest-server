package cache

import "github.com/egnptr/rest-api/entity"

type PostCache interface {
	Set(key string, value *entity.Post)
	Get(key string) *entity.Post
}
