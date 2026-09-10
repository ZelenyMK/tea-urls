package handler

import (
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/ZelenyMK/tea-urls/dto"
	"github.com/ZelenyMK/tea-urls/internal/shortener"
	"github.com/gin-gonic/gin"
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortenURL  string
	Alias       string
}

type LinkHandler struct {
	mu         sync.RWMutex
	links      map[string]string
	currentURL string
	nextID     int64
}

func NewLinkHandler(URL string) *LinkHandler {
	return &LinkHandler{
		links:      make(map[string]string),
		currentURL: strings.TrimRight(URL, "/"),
		nextID:     1}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var r dto.Request
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid URL",
		})
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This URL scheme is not supported",
		})
	}

	shortenedURL, alias, err := shortener.ShortenURL(h.currentURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create shortened URL",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		ID:          h.nextID,
		OriginalURL: r.URL,
		ShortenURL:  shortenedURL,
		Alias:       alias,
	})

	h.links[alias] = r.URL // ToDo: replace with PostgreSQL later
	h.nextID++
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias is required"})
		return
	}
	h.mu.RLock()
	originalURL, ok := h.links[alias]
	h.mu.RUnlock()
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this alias is not in a database"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
