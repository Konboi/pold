package pold

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	conf Config
}

type Blog struct {
	Title  string
	Author string
}

type Post struct {
	Title   string
	Content string
}

type Posts []Post

type View struct {
	Blog  *Blog
	Post  *Post
	Posts *Posts
}

var (
	blog *Blog
	tmpl = template.Must(template.New("tmpl").ParseGlob("templates/*.html"))
)

func NewServer(conf Config) (server *Server) {
	return &Server{
		conf: conf,
	}
}

func (s *Server) Run() error {
	b, err := s.BlogInfo()
	blog = b

	if err != nil {
		log.Println("error")
		log.Fatalf(err.Error())
	}

	fmt.Printf("pold server start 0.0.0.0:%d \n", s.conf.Port)

	router := httprouter.New()
	router.GET("/", IndexHandler)
	router.GET("/post/*path", PostHandler)

	return http.ListenAndServe(":"+strconv.Itoa(s.conf.Port), router)
}

func (s *Server) BlogInfo() (*Blog, error) {
	blog := &Blog{
		Title: s.conf.Title,
	}

	if blog.Title == "" {
		return blog, fmt.Errorf("blog title is empty")
	}

	return blog, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	view := &View{
		Blog: blog,
	}

	if err := tmpl.ExecuteTemplate(w, "index", view); err != nil {
		log.Fatal(err)
	}

}

func PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, ps.ByName("path"))
}
