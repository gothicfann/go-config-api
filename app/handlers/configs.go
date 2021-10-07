package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gothicfann/go-config-api/app/data"
)

// Configs handler
type Configs struct {
	l *log.Logger
}

// handleError is simple wrapper for error cases used for *Configs handler, where m is message and st is response statusCode.
// It returns json implementation of the error message and logs error if InternalServerError happened.
func (c *Configs) handleError(rw http.ResponseWriter, m string, st int) {
	type response struct {
		Message string `json:"message"`
	}

	resp := response{
		Message: m,
	}

	bs, _ := json.Marshal(resp)
	http.Error(rw, string(bs), st)
	if st == http.StatusInternalServerError {
		c.l.Println(m)
	}
}

// NewConfigs creates new *Configs handler.
func NewConfigs(l *log.Logger) *Configs {
	return &Configs{l}
}

// Health endpoint
func (c *Configs) Health(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(`{"message": "ok"}`))
}

// GetConfigs returns all configs in valid JSON format.
func (c *Configs) GetConfigs(rw http.ResponseWriter, r *http.Request) {
	cl := data.GetConfigs()

	err := cl.ToJSON(rw)
	if err != nil {
		c.handleError(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// GetConfigs returns single config in valid JSON format.
// It searches by "name" path attribute.
func (c *Configs) GetConfig(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	conf, err := data.GetConfig(vars["name"])
	if err != nil {
		c.handleError(rw, err.Error(), http.StatusNotFound)
		return
	}
	err = conf.ToJSON(rw)
	if err != nil {
		c.handleError(rw, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

// AddConfig is used for adds new config to database.
// Gives response with HTTP 409 Conflict if config is already present in database.
func (c *Configs) AddConfig(rw http.ResponseWriter, r *http.Request) {
	conf := r.Context().Value(KeyConfig{}).(data.Config)

	err := data.AddConfig(&conf)
	if err != nil {
		c.handleError(rw, err.Error(), http.StatusConflict)
		return
	}
}

// DeleteConfig deletes config from database by 'name' path attribute.
// It always executes 200 OK even if config is not present in db.
func (c *Configs) DeleteConfig(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data.DeleteConfig(vars["name"])
}

// PutConfig updates config if it exists, otherwise response is 404 Not Found.
func (c *Configs) PutConfig(rw http.ResponseWriter, r *http.Request) {
	conf := r.Context().Value(KeyConfig{}).(data.Config)

	err := data.PutConfig(&conf)
	if err != nil {
		c.handleError(rw, err.Error(), http.StatusNotFound)
		return
	}
}

// PatchConfig patches config in database if it exists, otherwise response will be 404 Not Found.
func (c *Configs) PatchConfig(rw http.ResponseWriter, r *http.Request) {
	conf := r.Context().Value(KeyConfig{}).(data.Config)

	err := data.PatchConfig(&conf)
	if err != nil {
		c.handleError(rw, err.Error(), http.StatusNotFound)
		return
	}
}

// QueryConfigs searches for configs which satisfy provided query parameter.
// Returns empty json object if nothing was found.
func (c *Configs) QueryConfigs(rw http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.RawQuery, "=") {
		q := strings.Split(r.URL.RawQuery, "=")
		k, v := q[0], q[1]
		cl, err := data.QueryConfig(k, v)
		if err != nil {
			c.handleError(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		err = cl.ToJSON(rw)
		if err != nil {
			c.handleError(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	c.handleError(rw, "Wrong query string", http.StatusBadRequest)
}

type KeyConfig struct{}

// MiddlewareValidateConfig is middleware function to validate requests payload before processing it.
func (c *Configs) MiddlewareValidateConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		// Try to decode payload
		conf := data.Config{}
		err := conf.FromJSON(r.Body)
		if err != nil {
			c.handleError(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// For Post method its required to provide name and metadata.
		if r.Method == http.MethodPost {
			if conf.Name == "" || conf.Metadata == nil {
				c.handleError(rw, "Config name or metadata not specified", http.StatusBadRequest)
				return
			}
		}

		// For Put and Patch methods check if metadata is not nil.
		// Also if name is provided it should not be empty and equal to path parameter.
		if r.Method == http.MethodPut || r.Method == http.MethodPatch {
			vars := mux.Vars(r)
			if conf.Metadata == nil {
				c.handleError(rw, "Config metadata not specified", http.StatusBadRequest)
				return
			}
			if len(conf.Name) > 0 && conf.Name != vars["name"] {
				c.handleError(rw, "URI name and config name are different", http.StatusBadRequest)
				return
			}
			conf.Name = vars["name"]
		}

		ctx := context.WithValue(r.Context(), KeyConfig{}, conf)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
