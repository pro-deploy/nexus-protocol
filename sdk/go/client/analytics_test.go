package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func TestLogEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1AnalyticsEvents {
			t.Errorf("Expected path %s, got %s", PathAPIV1AnalyticsEvents, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"data": {
				"event_id": "event-123",
				"message": "Event logged",
				"timestamp": "2025-01-18T10:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.LogEventRequest{
		EventType: "user_action",
		UserID:    "user-123",
	}

	resp, err := client.LogEvent(ctx, req)
	if err != nil {
		t.Fatalf("LogEvent failed: %v", err)
	}

	if resp.EventID != "event-123" {
		t.Errorf("Expected event_id 'event-123', got %s", resp.EventID)
	}
}

func TestGetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1AnalyticsEvents {
			t.Errorf("Expected path %s, got %s", PathAPIV1AnalyticsEvents, r.URL.Path)
		}

		// Проверяем query параметры
		if r.URL.Query().Get("event_type") != "user_action" {
			t.Errorf("Expected event_type 'user_action', got %s", r.URL.Query().Get("event_type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"events": [
					{
						"id": "event-1",
						"event_type": "user_action",
						"user_id": "user-123",
						"timestamp": "2025-01-18T10:00:00Z"
					}
				],
				"total": 1,
				"limit": 10,
				"offset": 0
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.GetEventsRequest{
		EventType: "user_action",
		Limit:     10,
		Offset:    0,
	}

	resp, err := client.GetEvents(ctx, req)
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if resp.Total != 1 {
		t.Errorf("Expected total 1, got %d", resp.Total)
	}

	if len(resp.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(resp.Events))
	}
}

func TestGetStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != PathAPIV1AnalyticsStats {
			t.Errorf("Expected path %s, got %s", PathAPIV1AnalyticsStats, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"period_days": 7,
				"total_events": 100,
				"total_users": 10,
				"active_users": 5,
				"events_today": 20
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{BaseURL: server.URL})
	ctx := context.Background()

	req := &types.GetStatsRequest{
		Days: 7,
	}

	stats, err := client.GetStats(ctx, req)
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	if stats.PeriodDays != 7 {
		t.Errorf("Expected period_days 7, got %d", stats.PeriodDays)
	}

	if stats.TotalEvents != 100 {
		t.Errorf("Expected total_events 100, got %d", stats.TotalEvents)
	}
}

