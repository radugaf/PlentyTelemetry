package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/radugaf/PlentyTelemetry/mocks"
	p "github.com/radugaf/PlentyTelemetry/ports"
	"go.uber.org/mock/gomock"
)

func TestLogger_Log(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	tags := map[string]string{
		"service": "test",
		"user":    "123",
	}
	txID := "test-transaction-id"

	mockWriter.EXPECT().Write(gomock.Any()).Do(func(entry p.LogEntry) {
		if entry.Level != p.Info {
			t.Errorf("Expected level INFO, got %s", entry.Level)
		}
	}).Times(1)

	logger.Log(p.Info, "test message", tags, txID)
}

func TestLogger_Info(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	mockWriter.EXPECT().Write(gomock.Any()).Do(func(entry p.LogEntry) {
		if entry.Level != p.Info {
			t.Errorf("Expected level INFO, got %s", entry.Level)
		}
		if entry.Message != "info message" {
			t.Errorf("Expected message 'info message', got '%s'", entry.Message)
		}
	}).Times(1)

	logger.Info("info message")
}

func TestLogger_Debug(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	mockWriter.EXPECT().Write(gomock.Any()).Do(func(entry p.LogEntry) {
		if entry.Level != p.Debug {
			t.Errorf("Expected level DEBUG, got %s", entry.Level)
		}
	}).Times(1)

	logger.Debug("debug message")
}

func TestLogger_Warning(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	mockWriter.EXPECT().Write(gomock.Any()).Do(func(entry p.LogEntry) {
		if entry.Level != p.Warning {
			t.Errorf("Expected level WARNING, got %s", entry.Level)
		}
	}).Times(1)

	logger.Warning("warning message")
}

func TestLogger_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	mockWriter.EXPECT().Write(gomock.Any()).Do(func(entry p.LogEntry) {
		if entry.Level != p.Error {
			t.Errorf("Expected level ERROR, got %s", entry.Level)
		}
	}).Times(1)

	logger.Error("error message")
}

func TestLogger_StartTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockWriter := mocks.NewMockLogWriter(ctrl)
	logger := NewLogger(mockWriter)

	txID1 := logger.StartTransaction()
	txID2 := logger.StartTransaction()

	if _, err := uuid.Parse(txID1); err != nil {
		t.Errorf("Expected valid UUID for transaction ID, got %s: %v", txID1, err)
	}

	if _, err := uuid.Parse(txID2); err != nil {
		t.Errorf("Expected valid UUID for transaction ID, got %s: %v", txID2, err)
	}

	if txID1 == txID2 {
		t.Errorf("Expected unique transaction IDs, but got same: %s", txID1)
	}
}

func TestParseArgs(t *testing.T) {
	inputTags := map[string]string{"env": "test"}
	tags, txID := parseArgs(123, inputTags, "tx-789", 45.6)

	if len(tags) != 1 {
		t.Errorf("Expected 1 tag, got %d", len(tags))
	}
	if tags["env"] != "test" {
		t.Errorf("Expected tag env=test, got %s", tags["env"])
	}
	if txID[0] != "tx-789" {
		t.Errorf("Expected txID tx-789, got %s", txID[0])
	}
}
