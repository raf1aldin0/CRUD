package cbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

// ✅ Breaker global yang dapat digunakan di seluruh aplikasi
var Breaker *gobreaker.CircuitBreaker

// ✅ Fungsi untuk membuat breaker dengan setting default
func NewDefaultBreaker(name string) *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3, // Saat half-open, hanya 3 request yang boleh masuk
		Interval:    60 * time.Second, // Reset hitungan sukses/gagal setiap 60 detik
		Timeout:     10 * time.Second, // Waktu tunggu sebelum mencoba kembali setelah breaker open
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Jika terjadi 5 kegagalan berturut-turut, breaker akan masuk mode "open"
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log atau observasi saat terjadi perubahan status breaker
			logStateChange(name, from, to)
		},
	}
	return gobreaker.NewCircuitBreaker(settings)
}

// ✅ Helper logging perubahan status breaker (optional)
func logStateChange(name string, from, to gobreaker.State) {
	stateToStr := map[gobreaker.State]string{
		gobreaker.StateClosed:   "CLOSED",
		gobreaker.StateOpen:     "OPEN",
		gobreaker.StateHalfOpen: "HALF-OPEN",
	}
	println("⚡ Circuit Breaker [" + name + "] berubah dari " + stateToStr[from] + " ke " + stateToStr[to])
}
