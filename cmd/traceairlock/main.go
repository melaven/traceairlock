package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// RequireSessionID — это наша Middleware.
// Она принимает следующий хендлер в цепочке (next) и возвращает новый хендлер.
func RequireSessionID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("X-Session-ID")

		if sessionID == "" {
			log.Println("[TraceAirlock] БЛОКИРОВКА: Отсутствует X-Session-ID")
			http.Error(w, "403 Forbidden: Missing Session ID", http.StatusForbidden)
			return
		}

		log.Printf("[TraceAirlock] Сессия %s пропущена", sessionID)

		// Если проверка пройдена, передаем запрос дальше по цепочке
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Читаем URL бэкенда из окружения, если нет — берем дефолтный
	targetStr := os.Getenv("TARGET_URL")
	if targetStr == "" {
		targetStr = "http://localhost:9090"
	}

	targetURL, err := url.Parse(targetStr)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// Читаем порт шлюза, если нет — ставим 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	protectedHandler := RequireSessionID(proxy)
	http.Handle("/", protectedHandler)

	log.Printf("TraceAirlock L7 Gateway starting on :%s, routing to %s", port, targetStr)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gateway crashed: %v", err)
	}
}
