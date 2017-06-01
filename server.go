package pold

import (
	"bytes"
	"encoding/json"
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
	Title  string `json:"title"`
	Author string `json:"author"`
	URL    string `json:"url"`
}

type View struct {
	Blog    *Blog
	Post    *Post
	Posts   Posts
	Content template.HTML
	Tag     string
}

type API struct {
	Results Posts `json:"results"`
}

var (
	blog           *Blog
	jetSet         = jet.NewHTMLSet("./templates")
	root, _        = os.Getwd() // todo set config
	topPostNum     = 5          // TODO: set config
	archivePostNum = 9999
	atomFeedNum    = 10
)

func NewServer(conf Config) (server *Server) {
	return &Server{
		conf: conf,
	}
}

func (s *Server) Run() {
	b, err := s.BlogInfo()

	if err != nil {
		log.Fatal("error get blog info", err.Error())
	}
	blog = b

	fmt.Printf("pold server start 0.0.0.0:%d \n", s.conf.Port)

	router := httprouter.New()
	router.GET("/", IndexHandler)
	router.GET("/post/*path", PostHandler)
	router.GET("/tag/*tag", TagHandler)
	router.GET("/archive", ArchiveHandler)
	router.GET("/atom.xml", AtomHandler)
	router.GET("/api/index", IndexAPIHandler)
	router.GET("/api/post/*path", PostAPIHandler)
	router.GET("/api/tag/*tag", TagAPIHandler)
	router.GET("/api/archive", ArchiveAPIHandler)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.conf.Port), router))
}

func (s *Server) BlogInfo() (*Blog, error) {
	blog := &Blog{
		Title:  s.conf.Title,
		URL:    s.conf.URL,
		Author: s.conf.Author,
	}

	if blog.Title == "" {
		return blog, fmt.Errorf("blog title is empty")
	}

	return blog, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(topPostNum)

	if err != nil {
		log.Println("error get published posts", err.Error())
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
	postFilePath := strings.Replace(ps.ByName("path"), ".html", "", -1)

	post, err := NewPost(postFilePath)

	if err != nil {
		log.Println(err.Error())
		notFound(w)
		return
	}

	view := &View{
		Post: post,
		Blog: blog,
	}

	jt, err := jetSet.GetTemplate("post.html")
	if err != nil {
		log.Println("error get post template", err.Error())
		return
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
	}
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(archivePostNum)

	if err != nil {
		log.Println("error get published archive posts", err.Error())
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
	}

	jt, err := jetSet.GetTemplate("archive.html")
	if err != nil {
		log.Println("error get archive template", err.Error())
		return
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
		return
	}
}

func TagHandler(w http.ResponseWriter, r *http.Request, tag httprouter.Params) {
	tagName := strings.Replace(tag.ByName("tag"), "/", "", -1)
	tagName = strings.Replace(tagName, ".html", "", -1)

	posts, err := PublishedPostsByTagName(tagName)
	if err != nil {
		log.Println("get published post by tag name", err.Error())
		return
	}
	if len(posts) == 0 {
		notFound(w)
		return
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
		Tag:   tagName,
	}

	jt, err := jetSet.GetTemplate("tag.html")
	if err != nil {
		log.Println(err.Error())
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
	}
}

func AtomHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(atomFeedNum)

	if err != nil {
		log.Println("error get published posts for atom feed", err.Error())
		return
	}

	view := &View{
		Blog:  blog,
		Posts: posts,
	}

	w.Header().Set("Content-Type", "application/atom+xml")
	jt, err := jetSet.GetTemplate("atom.xml")
	if err != nil {
		log.Println("error get template", err.Error())
		return
	}
	if err := jt.Execute(w, nil, view); err != nil {
		log.Println(err.Error())
	}
}

func IndexAPIHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(topPostNum)

	if err != nil {
		log.Println("error get published posts", err.Error())
		notFound(w)
		return
	}

	result := &API{Results: posts}
	if err := writeJSON(w, result); err != nil {
		log.Println("error write json", err.Error())
	}
	return
}

func PostAPIHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postFilePath := strings.Replace(ps.ByName("path"), ".html", "", -1)
	post, err := NewPost(postFilePath)
	if err != nil {
		log.Println("error get new post", err.Error())
		notFound(w)
		return
	}

	result := &API{Results: Posts{post}}
	if err := writeJSON(w, result); err != nil {
		log.Println("error write json", err.Error())
		notFound(w)
		return
	}

	return
}

func TagAPIHandler(w http.ResponseWriter, r *http.Request, tag httprouter.Params) {
	tagName := strings.Replace(tag.ByName("tag"), "/", "", -1)
	tagName = strings.Replace(tagName, ".html", "", -1)

	posts, err := PublishedPostsByTagName(tagName)
	if err != nil {
		log.Println("not tag post", err.Error())
		notFound(w)
		return
	}
	if len(posts) == 0 {
		notFound(w)
		return
	}

	result := &API{Results: posts}
	if err := writeJSON(w, result); err != nil {
		log.Println("error write json", err.Error())
	}
	return
}

func ArchiveAPIHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	posts, err := PublishedPosts(archivePostNum)
	if err != nil {
		log.Println("error get publied post", err.Error())
		return
	}
	result := &API{Results: posts}
	if err := writeJSON(w, result); err != nil {
		log.Println("error write json", err.Error())
	}
	return
}

func writeJSON(w http.ResponseWriter, result interface{}) error {
	body := &bytes.Buffer{}
	if err := json.NewEncoder(body).Encode(result); err != nil {
		log.Println("error json marshal", err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body.Bytes())
	return nil
}

func notFound(w http.ResponseWriter) {
	code := http.StatusNotFound

	http.Error(w, http.StatusText(code), code)
}
