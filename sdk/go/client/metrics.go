package client

import (
	"context"
	"net/http"
	"time"
)

// MetricsCollector собирает метрики для запросов
type MetricsCollector interface {
	// RecordRequest записывает метрику запроса
	RecordRequest(method, path string, statusCode int, duration time.Duration)

	// RecordError записывает метрику ошибки
	RecordError(method, path string, err error)
}

// MetricsInterceptor реализует Interceptor для сбора метрик
type MetricsInterceptor struct {
	collector MetricsCollector
}

// NewMetricsInterceptor создает новый interceptor для метрик
func NewMetricsInterceptor(collector MetricsCollector) *MetricsInterceptor {
	return &MetricsInterceptor{
		collector: collector,
	}
}

// BeforeRequest вызывается перед запросом
func (m *MetricsInterceptor) BeforeRequest(ctx context.Context, req *http.Request) error {
	// Сохраняем время начала запроса в контексте
	req.Header.Set("X-Request-Start-Time", time.Now().Format(time.RFC3339Nano))
	return nil
}

// AfterResponse вызывается после ответа
func (m *MetricsInterceptor) AfterResponse(ctx context.Context, req *http.Request, resp *http.Response) error {
	if m.collector == nil {
		return nil
	}

	// Вычисляем длительность запроса
	startTimeStr := req.Header.Get("X-Request-Start-Time")
	if startTimeStr != "" {
		if startTime, err := time.Parse(time.RFC3339Nano, startTimeStr); err == nil {
			duration := time.Since(startTime)
			m.collector.RecordRequest(req.Method, req.URL.Path, resp.StatusCode, duration)
		}
	}

	// Записываем ошибки
	if resp.StatusCode >= 400 {
		m.collector.RecordError(req.Method, req.URL.Path, nil)
	}

	return nil
}

// SimpleMetricsCollector простая реализация MetricsCollector
type SimpleMetricsCollector struct {
	requests  map[string]int64
	errors    map[string]int64
	durations map[string][]time.Duration
}

// NewSimpleMetricsCollector создает простой коллектор метрик
func NewSimpleMetricsCollector() *SimpleMetricsCollector {
	return &SimpleMetricsCollector{
		requests:  make(map[string]int64),
		errors:    make(map[string]int64),
		durations: make(map[string][]time.Duration),
	}
}

// RecordRequest записывает метрику запроса
func (s *SimpleMetricsCollector) RecordRequest(method, path string, statusCode int, duration time.Duration) {
	key := method + " " + path
	s.requests[key]++

	if s.durations[key] == nil {
		s.durations[key] = make([]time.Duration, 0, 100)
	}
	s.durations[key] = append(s.durations[key], duration)

	// Ограничиваем размер массива
	if len(s.durations[key]) > 1000 {
		s.durations[key] = s.durations[key][len(s.durations[key])-1000:]
	}
}

// RecordError записывает метрику ошибки
func (s *SimpleMetricsCollector) RecordError(method, path string, err error) {
	key := method + " " + path
	s.errors[key]++
}

// GetStats возвращает статистику
func (s *SimpleMetricsCollector) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	requests := make(map[string]int64)
	errors := make(map[string]int64)
	avgDurations := make(map[string]time.Duration)

	for key, count := range s.requests {
		requests[key] = count

		if durations, ok := s.durations[key]; ok && len(durations) > 0 {
			var total time.Duration
			for _, d := range durations {
				total += d
			}
			avgDurations[key] = total / time.Duration(len(durations))
		}
	}

	for key, count := range s.errors {
		errors[key] = count
	}

	stats["requests"] = requests
	stats["errors"] = errors
	stats["avg_durations"] = avgDurations

	return stats
}

// Reset сбрасывает метрики
func (s *SimpleMetricsCollector) Reset() {
	s.requests = make(map[string]int64)
	s.errors = make(map[string]int64)
	s.durations = make(map[string][]time.Duration)
}
