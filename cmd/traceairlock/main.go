package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 1. Указываем мишень, куда шлюз будет перенаправлять чистый трафик (наш мок)
	targetURL, err := url.Parse("http://localhost:9090")
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	// 2. Инициализируем стандартный Reverse Proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 3. Главный хендлер-перехватчик
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[TraceAirlock] Перехват трафика: %s %s", r.Method, r.URL.Path)

		// 1. Извлекаем заголовок X-Session-ID
		sessionID := r.Header.Get("X-Session-ID")

		// 2. Fail-Closed: нет заголовка — жестко рубим соединение
		if sessionID == "" {
			log.Println("[TraceAirlock] БЛОКИРОВКА: Отсутствует X-Session-ID")
			http.Error(w, "403 Forbidden: Missing Session ID", http.StatusForbidden)
			return
		}

		// 3. Если заголовок есть — стримим трафик на целевой бэкенд
		log.Printf("[TraceAirlock] Сессия %s пропущена", sessionID)
		proxy.ServeHTTP(w, r)
	})

	// Шлюз слушает внешний мир на порту 8080
	log.Println("TraceAirlock L7 Gateway starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Gateway crashed: %v", err)
	}
}
