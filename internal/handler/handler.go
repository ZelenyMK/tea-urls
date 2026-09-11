package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/ZelenyMK/tea-urls/dto"
	"github.com/ZelenyMK/tea-urls/internal/shortener"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortenURL  string
	Alias       string
}

type LinkHandler struct {
	mu         sync.RWMutex
	db         *pgx.Conn
	currentURL string
	nextID     int64
}

func NewLinkHandler(URL string, conn *pgx.Conn) *LinkHandler {
	return &LinkHandler{
		db:         conn,
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

	shortenedURL, alias, err := shortener.ShortenURL(r.URL)
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

	h.mu.Lock()
	h.db.Exec(context.Background(), "INSERT INTO links (ID, original_url, shorten_url, alias) VALUES ($1, $2, $3, $4)",
		strconv.Itoa(int(h.nextID)), r.URL, shortenedURL, alias)
	h.nextID++
	h.mu.Unlock()
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias is required"})
		return
	}
	h.mu.RLock()
	row := h.db.QueryRow(context.Background(), "SELECT original_url FROM links WHERE alias = $1", alias)
	h.mu.RUnlock()

	var result string
	err := row.Scan(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QueryRow failed"})
		return
	}
	if result == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "got empty string"})
		return
	}

	c.Redirect(http.StatusFound, result)
}
