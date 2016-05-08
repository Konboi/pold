package pold

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	conf Config
}

type Blog struct {
	Title  string
	Author string
}

type View struct {
	Blog    *Blog
	Post    *Post
	Posts   Posts
	Test    string
	Content template.HTML
}

var (
	blog       *Blog
	tmpl       = template.Must(template.New("tmpl").ParseGlob("templates/*.html"))
	root, _    = os.Getwd() // todo set config
	topPostNum = 10         // TODO: set config
)

func NewServer(conf Config) (server *Server) {
	return &Server{
		conf: conf,
	}
}

func (s *Server) Run() {
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

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.conf.Port), router))
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
	posts, err := PublishedPosts()

	if err != nil {
		log.Fatal(err)
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
	}

	if err := tmpl.ExecuteTemplate(w, "index", view); err != nil {
		log.Fatal(err)
	}

}

func PostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	path := strings.Trim(ps.ByName("path"), ".html")
	postFilePath := path

	post, err := NewPost(postFilePath)

	if err != nil {
		log.Println(err.Error())
		notFound(w)
		return
	}

	view := &View{
		Test: path,
		Post: post,
		Blog: blog,
	}

	if err := tmpl.ExecuteTemplate(w, "post", view); err != nil {
		log.Fatal(err)
	}
}

func notFound(w http.ResponseWriter) {
	code := http.StatusNotFound

	http.Error(w, http.StatusText(code), code)
}
