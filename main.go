package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/gograz/gograz-meetup/meetupcom"
)

//go:embed templates
var rootFS embed.FS

type server struct {
	client    *meetupcom.Client
	urlName   string
	cache     *cache.Cache
	templates *template.Template
}

type attendee struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ThumbLink string `json:"thumbLink"`
	PhotoLink string `json:"photoLink"`
	Guests    int64  `json:"guests"`
}

type rsvps struct {
	EventID string     `json:"-"`
	Yes     []attendee `json:"yes"`
	No      []attendee `json:"no"`
}

func convertRSVPs(in meetupcom.RSVPsResponse) rsvps {
	out := rsvps{
		Yes: make([]attendee, 0, 2),
		No:  make([]attendee, 0, 2),
	}
	for _, item := range in {
		m := attendee{
			ID:        item.Member.ID,
			Name:      item.Member.Name,
			PhotoLink: item.Member.Photo.PhotoLink,
			ThumbLink: item.Member.Photo.ThumbLink,
		}
		if item.Response == "YES" {
			out.Yes = append(out.Yes, m)
		} else if item.Response == "NO" {
			out.No = append(out.No, m)
		}
	}
	return out
}

func (s *server) handleGetRSVPs(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "eventID")
	cacheKey := fmt.Sprintf("rsvps:%s", eventID)
	var rsvps *meetupcom.RSVPsResponse

	cached, found := s.cache.Get(cacheKey)
	if found {
		s.encodeRSVPs(w, r, eventID, cached.(meetupcom.RSVPsResponse))
		return
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	rsvps, err := s.client.GetRSVPs(ctx, eventID, s.urlName)
	if err != nil {
		log.WithError(err).Errorf("Failed to fetch RSVPs for %s", eventID)
		http.Error(w, "Failed to fetch RSVPs from backend", http.StatusInternalServerError)
		return
	}
	s.cache.Set(cacheKey, *rsvps, 0)
	s.encodeRSVPs(w, r, eventID, *rsvps)
}

func (s *server) encodeRSVPs(w http.ResponseWriter, r *http.Request, eventID string, in meetupcom.RSVPsResponse) {
	out := convertRSVPs(in)
	out.EventID = eventID
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("Content-Type", "text/html")
		_ = s.templates.ExecuteTemplate(w, "rsvps.tmpl", out)
		return
	}
	w.Header().Set("Content-Type", "text/json")
	_ = json.NewEncoder(w).Encode(out)
}

func main() {
	var addr string
	var urlName string
	var allowedOrigins []string

	flag.StringVar(&addr, "addr", "127.0.0.1:8080", "Address to listen on")
	flag.StringVar(&urlName, "url-name", "Graz-Open-Source-Meetup", "URL name of the meetup group on meetup.com")
	flag.StringArrayVar(&allowedOrigins, "allowed-origins", []string{"http://localhost:1313", "https://gograz.org"}, "Allowed origin hosts")
	flag.Parse()

	ch := cache.New(5*time.Minute, 10*time.Minute)

	templates, err := template.ParseFS(rootFS, "templates/*.tmpl")
	if err != nil {
		log.WithError(err).Fatal("cannot load templates")
	}

	s := server{
		client:    meetupcom.NewClient(meetupcom.ClientOptions{}),
		urlName:   urlName,
		cache:     ch,
		templates: templates,
	}

	router := chi.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedHeaders:   []string{"HX-Current-URL", "HX-Request", "HX-Target", "HX-Trigger"},
	})
	router.Get("/{eventID}/rsvps", s.handleGetRSVPs)
	log.Infof("Starting HTTPD on %s", addr)
	_ = http.ListenAndServe(addr, c.Handler(router))
}
