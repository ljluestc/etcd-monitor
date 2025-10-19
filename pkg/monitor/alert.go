package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"sync"
	"time"

	"go.uber.org/zap"
)

// AlertLevel defines the severity of an alert
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// AlertType defines the type of alert
type AlertType string

const (
	AlertTypeClusterHealth      AlertType = "cluster_health"
	AlertTypeLeaderElection     AlertType = "leader_election"
	AlertTypeNetworkPartition   AlertType = "network_partition"
	AlertTypeHighLatency        AlertType = "high_latency"
	AlertTypeHighDiskUsage      AlertType = "high_disk_usage"
	AlertTypeHighProposalQueue  AlertType = "high_proposal_queue"
	AlertTypeEtcdAlarm          AlertType = "etcd_alarm"
)

// Alert represents an alert to be sent
type Alert struct {
	Level     AlertLevel
	Type      AlertType
	Message   string
	Details   map[string]interface{}
	Timestamp time.Time
}

// AlertManager manages alerting and notifications
type AlertManager struct {
	thresholds   AlertThresholds
	logger       *zap.Logger
	mu           sync.RWMutex
	alertHistory []Alert
	maxHistory   int
	channels     []AlertChannel

	// Deduplication
	activeAlerts map[string]time.Time
	dedupWindow  time.Duration
}

// AlertChannel is an interface for sending alerts
type AlertChannel interface {
	Send(alert Alert) error
	Name() string
}

// NewAlertManager creates a new alert manager
func NewAlertManager(thresholds AlertThresholds, logger *zap.Logger) *AlertManager {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	return &AlertManager{
		thresholds:   thresholds,
		logger:       logger,
		alertHistory: make([]Alert, 0),
		maxHistory:   1000,
		channels:     make([]AlertChannel, 0),
		activeAlerts: make(map[string]time.Time),
		dedupWindow:  5 * time.Minute,
	}
}

// AddChannel adds an alert channel
func (am *AlertManager) AddChannel(channel AlertChannel) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.channels = append(am.channels, channel)
	am.logger.Info("Alert channel added", zap.String("channel", channel.Name()))
}

// TriggerAlert triggers an alert
func (am *AlertManager) TriggerAlert(alert Alert) {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Check for deduplication
	alertKey := fmt.Sprintf("%s:%s", alert.Type, alert.Message)
	if lastSent, exists := am.activeAlerts[alertKey]; exists {
		if time.Since(lastSent) < am.dedupWindow {
			am.logger.Debug("Alert deduplicated",
				zap.String("type", string(alert.Type)),
				zap.String("message", alert.Message))
			return
		}
	}

	// Record alert
	am.alertHistory = append(am.alertHistory, alert)
	if len(am.alertHistory) > am.maxHistory {
		am.alertHistory = am.alertHistory[1:]
	}

	// Mark as active
	am.activeAlerts[alertKey] = time.Now()

	// Send to all channels
	am.logger.Info("Triggering alert",
		zap.String("level", string(alert.Level)),
		zap.String("type", string(alert.Type)),
		zap.String("message", alert.Message))

	for _, channel := range am.channels {
		go func(ch AlertChannel) {
			if err := ch.Send(alert); err != nil {
				am.logger.Error("Failed to send alert",
					zap.String("channel", ch.Name()),
					zap.Error(err))
			}
		}(channel)
	}
}

// GetAlertHistory returns the alert history
func (am *AlertManager) GetAlertHistory() []Alert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	history := make([]Alert, len(am.alertHistory))
	copy(history, am.alertHistory)
	return history
}

// ActiveAlert represents an alert that is currently active
type ActiveAlert struct {
	Alert
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}

// GetActiveAlerts returns all currently active alerts
func (am *AlertManager) GetActiveAlerts() []ActiveAlert {
	am.mu.RLock()
	defer am.mu.RUnlock()

	activeAlerts := make([]ActiveAlert, 0)
	now := time.Now()

	for alertKey, lastSeen := range am.activeAlerts {
		// Consider alert active if seen within the last dedup window
		if now.Sub(lastSeen) < am.dedupWindow*2 {
			// Find the alert in history
			for _, alert := range am.alertHistory {
				key := fmt.Sprintf("%s:%s", alert.Type, alert.Message)
				if key == alertKey {
					activeAlerts = append(activeAlerts, ActiveAlert{
						Alert:     alert,
						FirstSeen: alert.Timestamp,
						LastSeen:  lastSeen,
					})
					break
				}
			}
		}
	}

	return activeAlerts
}

// ClearAlert clears an active alert
func (am *AlertManager) ClearAlert(alertType AlertType, message string) {
	am.mu.Lock()
	defer am.mu.Unlock()

	alertKey := fmt.Sprintf("%s:%s", alertType, message)
	delete(am.activeAlerts, alertKey)

	am.logger.Info("Alert cleared",
		zap.String("type", string(alertType)),
		zap.String("message", message))
}

// EmailChannel sends alerts via email
type EmailChannel struct {
	SMTPServer   string
	SMTPPort     int
	From         string
	To           []string
	Username     string
	Password     string
}

func (ec *EmailChannel) Name() string {
	return "email"
}

func (ec *EmailChannel) Send(alert Alert) error {
	// Compose email message
	subject := fmt.Sprintf("[%s] etcd-monitor: %s", alert.Level, alert.Type)
	body := ec.formatEmailBody(alert)

	// Create message
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", ec.From, ec.To[0], subject, body)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", ec.SMTPServer, ec.SMTPPort)

	// Setup authentication if credentials provided
	var auth smtp.Auth
	if ec.Username != "" && ec.Password != "" {
		auth = smtp.PlainAuth("", ec.Username, ec.Password, ec.SMTPServer)
	}

	// Send email
	err := smtp.SendMail(addr, auth, ec.From, ec.To, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (ec *EmailChannel) formatEmailBody(alert Alert) string {
	detailsHTML := ""
	if len(alert.Details) > 0 {
		detailsHTML = "<h3>Details:</h3><ul>"
		for key, value := range alert.Details {
			detailsHTML += fmt.Sprintf("<li><strong>%s:</strong> %v</li>", key, value)
		}
		detailsHTML += "</ul>"
	}

	color := "#28a745" // green
	if alert.Level == AlertLevelWarning {
		color = "#ffc107" // yellow
	} else if alert.Level == AlertLevelCritical {
		color = "#dc3545" // red
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; }
        .alert-box { border-left: 4px solid %s; padding: 15px; margin: 10px 0; background-color: #f8f9fa; }
        .alert-header { font-size: 18px; font-weight: bold; color: %s; }
        .alert-message { font-size: 14px; margin: 10px 0; }
        .alert-meta { font-size: 12px; color: #6c757d; }
    </style>
</head>
<body>
    <div class="alert-box">
        <div class="alert-header">%s Alert: %s</div>
        <div class="alert-message">%s</div>
        <div class="alert-meta">
            <p><strong>Time:</strong> %s</p>
            <p><strong>Level:</strong> %s</p>
            <p><strong>Type:</strong> %s</p>
        </div>
        %s
    </div>
    <hr>
    <p style="font-size: 12px; color: #6c757d;">
        This alert was generated by etcd-monitor.
        <a href="http://your-etcd-monitor-url/alerts">View Dashboard</a>
    </p>
</body>
</html>
`, color, color, alert.Level, alert.Type, alert.Message,
		alert.Timestamp.Format(time.RFC3339), alert.Level, alert.Type, detailsHTML)
}

// SlackChannel sends alerts to Slack
type SlackChannel struct {
	WebhookURL string
	Channel    string
	Username   string
}

func (sc *SlackChannel) Name() string {
	return "slack"
}

func (sc *SlackChannel) Send(alert Alert) error {
	payload := map[string]interface{}{
		"channel":  sc.Channel,
		"username": sc.Username,
		"text":     fmt.Sprintf("[%s] %s", alert.Level, alert.Message),
		"attachments": []map[string]interface{}{
			{
				"color": sc.getColor(alert.Level),
				"fields": []map[string]interface{}{
					{
						"title": "Type",
						"value": string(alert.Type),
						"short": true,
					},
					{
						"title": "Timestamp",
						"value": alert.Timestamp.Format(time.RFC3339),
						"short": true,
					},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(sc.WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send to Slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack returned non-200 status: %d", resp.StatusCode)
	}

	return nil
}

func (sc *SlackChannel) getColor(level AlertLevel) string {
	switch level {
	case AlertLevelCritical:
		return "danger"
	case AlertLevelWarning:
		return "warning"
	default:
		return "good"
	}
}

// PagerDutyChannel sends alerts to PagerDuty
type PagerDutyChannel struct {
	IntegrationKey string
}

func (pdc *PagerDutyChannel) Name() string {
	return "pagerduty"
}

func (pdc *PagerDutyChannel) Send(alert Alert) error {
	payload := map[string]interface{}{
		"routing_key":  pdc.IntegrationKey,
		"event_action": "trigger",
		"payload": map[string]interface{}{
			"summary":   alert.Message,
			"severity":  string(alert.Level),
			"source":    "etcd-monitor",
			"timestamp": alert.Timestamp.Format(time.RFC3339),
			"custom_details": alert.Details,
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post("https://events.pagerduty.com/v2/enqueue",
		"application/json",
		bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to send to PagerDuty: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pagerduty returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}

// WebhookChannel sends alerts to a generic webhook
type WebhookChannel struct {
	URL     string
	Headers map[string]string
}

func (wc *WebhookChannel) Name() string {
	return "webhook"
}

func (wc *WebhookChannel) Send(alert Alert) error {
	payload := map[string]interface{}{
		"level":     string(alert.Level),
		"type":      string(alert.Type),
		"message":   alert.Message,
		"details":   alert.Details,
		"timestamp": alert.Timestamp.Format(time.RFC3339),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", wc.URL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range wc.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}

// ConsoleChannel logs alerts to the console (for testing)
type ConsoleChannel struct {
	logger *zap.Logger
}

func NewConsoleChannel(logger *zap.Logger) *ConsoleChannel {
	return &ConsoleChannel{logger: logger}
}

func (cc *ConsoleChannel) Name() string {
	return "console"
}

func (cc *ConsoleChannel) Send(alert Alert) error {
	cc.logger.Info("ALERT",
		zap.String("level", string(alert.Level)),
		zap.String("type", string(alert.Type)),
		zap.String("message", alert.Message),
		zap.Any("details", alert.Details),
		zap.Time("timestamp", alert.Timestamp))
	return nil
}
