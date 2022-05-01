package main

import (
	"fmt"
	"net/http"

	"github.com/egnptr/rest-api/cache"
	"github.com/egnptr/rest-api/controller"
	router "github.com/egnptr/rest-api/http"
	"github.com/egnptr/rest-api/repository"
	"github.com/egnptr/rest-api/service"
)

var (
	postRepo       repository.PostRepository = repository.NewSQLiteRepository()
	postService    service.PostService       = service.NewPostService(postRepo)
	postCache      cache.PostCache           = cache.NewRedisCache("localhost:6379", 1, 10)
	postController controller.PostController = controller.NewPostController(postService, postCache)
	httpRouter     router.Router             = router.NewMuxRouter()
)

func main() {
	const port string = ":8080"
	httpRouter.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
	httpRouter.GET("/posts", postController.GetAllPost)
	httpRouter.POST("/posts", postController.AddPost)
	httpRouter.GET("/posts/{id}", postController.GetPost)
	httpRouter.DELETE("/posts/{id}", postController.DeletePost)
	httpRouter.SERVE(port)
}
