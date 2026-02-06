package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Manager struct {
	buckets    map[string]*TokenBucket
	lastAccess map[string]time.Time
	capacity   int
	fillRate   float64
	ttl        time.Duration
	mu         sync.RWMutex
}

func NewManager(capacity int, fillRate float64, ttl time.Duration) *Manager {
	m := &Manager{
		capacity:   capacity,
		fillRate:   fillRate,
		ttl:        ttl,
		buckets:    make(map[string]*TokenBucket),
		lastAccess: make(map[string]time.Time),
	}
	go m.cleanupLoop()
	return m
}

func (m *Manager) GetBucket(key string) *TokenBucket {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.lastAccess[key] = time.Now()

	if bucket, exists := m.buckets[key]; exists {
		return bucket
	}
	bucket := NewTokenBucket(m.capacity, m.fillRate)
	m.buckets[key] = bucket
	return bucket
}

func (m *Manager) ExtractClientKey(r *http.Request, useUserID bool) string {
	if useUserID {
		userID := r.Header.Get("X-User-ID")
		if userID != "" {
			return "user:" + userID
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		ip = r.RemoteAddr
	}

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		ip = strings.TrimSpace(ips[0])
	}
	return "ip:" + ip
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(m.ttl)
	for range ticker.C {
		m.cleanup()
	}
}

func (m *Manager) cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for key, lastSeen := range m.lastAccess {
		if now.Sub(lastSeen) > m.ttl {
			delete(m.buckets, key)
			delete(m.lastAccess, key)
		}
	}
}
