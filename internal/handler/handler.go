package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ZelenyMK/tea-urls/internal/dto"
	"github.com/ZelenyMK/tea-urls/internal/shortener"
	"github.com/ZelenyMK/tea-urls/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"github.com/redis/go-redis/v9"
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
	redis      *redis.Client
	currentURL string
	nextID     int64
}

func NewLinkHandler(URL string, conn *pgx.Conn) *LinkHandler {
	return &LinkHandler{
		db:         conn,
		currentURL: strings.TrimRight(URL, "/"),
		nextID:     1,
		redis:      storage.NewRedisClient()}
}

func (h *LinkHandler) Create(c *gin.Context) {
	var r dto.Request
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid URL",
		})
		return
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "This URL scheme is not supported",
		})
		return
	}

	shortenedURL, alias, err := shortener.ShortenURL(r.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create shortened URL",
		})
		return
	}

	h.mu.Lock()
	_, err = h.db.Exec(context.Background(), "INSERT INTO links (ID, original_url, shorten_url, alias) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
		strconv.Itoa(int(h.nextID)), r.URL, shortenedURL, alias)
	h.nextID++
	h.mu.Unlock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not insert into db",
		})
		return
	}

	err = h.redis.Set(
		c.Request.Context(),
		alias,
		r.URL,
		1*time.Hour,
	).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not set to redis",
		})
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusCreated, `
		<a href="%s" target="_blank" rel="noopener">%s</a>
	`, shortenedURL, shortenedURL)
}

func (h *LinkHandler) Redirect(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "alias is required"})
		return
	}

	h.mu.RLock()
	redisResult, err := h.redis.Get(c.Request.Context(), alias).Result()
	h.mu.RUnlock()

	if err == nil {
		c.Redirect(http.StatusFound, redisResult)
		return
	} else if err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "redis error",
		})
	}

	h.mu.RLock()
	row := h.db.QueryRow(context.Background(), "SELECT original_url FROM links WHERE alias = $1", alias)
	h.mu.RUnlock()

	var result string
	err = row.Scan(&result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QueryRow failed"})
		return
	}
	if result == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "got empty string"})
		return
	}

	err = h.redis.Set(
		c.Request.Context(),
		alias,
		result,
		1*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not write cache to redis"})
	}

	c.Redirect(http.StatusFound, result)
}
