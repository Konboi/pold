package pold

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	conf Config
}

type Blog struct {
	Title  string
	Author string
	URL    string
}

type View struct {
	Blog    *Blog
	Post    *Post
	Posts   Posts
	Test    string
	Content template.HTML
}

var (
	blog           *Blog
	jetSet         = jet.NewHTMLSet("./templates")
	root, _        = os.Getwd() // todo set config
	topPostNum     = 10         // TODO: set config
	archivePostNum = 9999
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
	router.GET("/archive", ArchiveHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.conf.Port), router))
}

func (s *Server) BlogInfo() (*Blog, error) {
	blog := &Blog{
		Title: s.conf.Title,
		URL:   s.conf.URL,
	}

	if blog.Title == "" {
		return blog, fmt.Errorf("blog title is empty")
	}

	return blog, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(topPostNum)

	if err != nil {
		log.Fatal(err)
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
	}

	jt, err := jetSet.GetTemplate("index.html")
	if err != nil {
		log.Println(err.Error())
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
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

	jt, err := jetSet.GetTemplate("post.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
	}
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(archivePostNum)

	if err != nil {
		log.Fatal(err)
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
	}

	jt, err := jetSet.GetTemplate("archive.html")
	if err != nil {
		log.Fatal(err)
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
	}
}

func notFound(w http.ResponseWriter) {
	code := http.StatusNotFound

	http.Error(w, http.StatusText(code), code)
}
