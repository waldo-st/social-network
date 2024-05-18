package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"social/internal/core/service"
	"strings"
)

type router struct {
	routes       map[string]map[string]http.HandlerFunc
	tokenService service.TokenService
}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewContext(rw http.ResponseWriter, rq *http.Request) *Context {
	return &Context{rw, rq}
}

// add the tokenService here
func New(tokService service.TokenService) *router {
	return &router{
		tokenService: tokService,
		routes:       make(map[string]map[string]http.HandlerFunc),
	}
}

func (c *Context) SendResponse(data interface{}, code int) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	if err := json.NewEncoder(c.ResponseWriter).Encode(data); err != nil {
		log.Println("❌ Error encoding json response")
	}
}

func (c *Context) GetContextValue(ctx context.Context, key string) (any, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, errors.New("bad request")
	}
	return value, nil
}

func (c *Context) HandleError(message string, code int) {
	log.Println("❌ ERROR: ", message)
	response := map[string]string{"message": message}
	c.SendResponse(response, code)
}

func (c *Context) BindJson(data interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

func (r *router) HandleFunc(method string, pattern string, handler http.HandlerFunc) {
	methods, ok := r.routes[pattern]
	if !ok {
		methods = make(map[string]http.HandlerFunc)
		r.routes[pattern] = methods
	}
	methods[method] = handler
}

func matchPattern(pattern string, path string) bool {
	// Split pattern and path into segments
	patternSegments := strings.Split(pattern, "/")
	pathSegments := strings.Split(path, "/")

	// If number of segments do not match, pattern does not match path
	if len(patternSegments) != len(pathSegments) {
		return false
	}

	// Check each segment for match
	for i, segment := range patternSegments {
		if segment != pathSegments[i] && !strings.HasPrefix(segment, "{") {
			return false
		}
	}

	return true
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL.Path)
	methods, ok := r.routes[req.URL.Path]
	if !ok {
		fmt.Println("Not okay", methods)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	handler, ok := methods[req.Method]
	if !ok {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	handler(w, req)
	// for pattern, methods := range r.routes {
	// 	if matchPattern(pattern, req.URL.Path) {
	// 		if handler, ok := methods[req.Method]; ok {
	// 			handler(w, req)
	// 			return
	// 		}
	// 	}
	// }
	// http.NotFound(w, req)
}
