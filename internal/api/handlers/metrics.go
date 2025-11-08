package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/connect-univyn/connect-server/internal/live"
	"github.com/gin-gonic/gin"
)


type MetricsHandler struct {
	liveService *live.Service
}


func NewMetricsHandler(liveService *live.Service) *MetricsHandler {
	return &MetricsHandler{
		liveService: liveService,
	}
}


func (h *MetricsHandler) HandlePrometheusMetrics(c *gin.Context) {
	if h.liveService == nil {
		c.String(http.StatusServiceUnavailable, "# Live service not available\n")
		return
	}

	
	wsMetrics := h.liveService.GetWebSocketMetrics()

	
	brokerMetrics := h.liveService.GetBrokerMetrics()

	
	var output string
	now := time.Now().UnixMilli()

	
	output += "# HELP websocket_active_connections Current number of active WebSocket connections\n"
	output += "# TYPE websocket_active_connections gauge\n"
	output += fmt.Sprintf("websocket_active_connections %d %d\n", wsMetrics.ActiveConnections, now)

	output += "# HELP websocket_total_connections Total number of WebSocket connections since start\n"
	output += "# TYPE websocket_total_connections counter\n"
	output += fmt.Sprintf("websocket_total_connections %d %d\n", wsMetrics.TotalConnections, now)

	output += "# HELP websocket_connections_rejected Total number of rejected connections due to limits\n"
	output += "# TYPE websocket_connections_rejected counter\n"
	output += fmt.Sprintf("websocket_connections_rejected %d %d\n", wsMetrics.ConnectionsRejected, now)

	output += "# HELP websocket_messages_received Total number of messages received\n"
	output += "# TYPE websocket_messages_received counter\n"
	output += fmt.Sprintf("websocket_messages_received %d %d\n", wsMetrics.MessagesReceived, now)

	output += "# HELP websocket_messages_sent Total number of messages sent\n"
	output += "# TYPE websocket_messages_sent counter\n"
	output += fmt.Sprintf("websocket_messages_sent %d %d\n", wsMetrics.MessagesSent, now)

	output += "# HELP websocket_errors Total number of WebSocket errors\n"
	output += "# TYPE websocket_errors counter\n"
	output += fmt.Sprintf("websocket_errors %d %d\n", wsMetrics.Errors, now)

	output += "# HELP websocket_average_latency_ms Average message latency in milliseconds\n"
	output += "# TYPE websocket_average_latency_ms gauge\n"
	output += fmt.Sprintf("websocket_average_latency_ms %.2f %d\n", wsMetrics.GetAverageLatencyMs(), now)

	output += "# HELP websocket_message_throughput_sec Message throughput per second\n"
	output += "# TYPE websocket_message_throughput_sec gauge\n"
	output += fmt.Sprintf("websocket_message_throughput_sec %.2f %d\n", wsMetrics.GetMessageThroughput(), now)

	
	if brokerMetrics != nil {
		output += "# HELP redis_events_published Total events published to Redis\n"
		output += "# TYPE redis_events_published counter\n"
		output += fmt.Sprintf("redis_events_published %d %d\n", brokerMetrics.GetEventsPublished(), now)

		output += "# HELP redis_events_received Total events received from Redis\n"
		output += "# TYPE redis_events_received counter\n"
		output += fmt.Sprintf("redis_events_received %d %d\n", brokerMetrics.GetEventsReceived(), now)

		output += "# HELP redis_publish_errors Total Redis publish errors\n"
		output += "# TYPE redis_publish_errors counter\n"
		output += fmt.Sprintf("redis_publish_errors %d %d\n", brokerMetrics.GetPublishErrors(), now)

		output += "# HELP redis_reconnect_count Total Redis reconnection attempts\n"
		output += "# TYPE redis_reconnect_count counter\n"
		output += fmt.Sprintf("redis_reconnect_count %d %d\n", brokerMetrics.GetReconnectCount(), now)

		if !brokerMetrics.LastReconnectTime.IsZero() {
			output += "# HELP redis_last_reconnect_timestamp Unix timestamp of last reconnection\n"
			output += "# TYPE redis_last_reconnect_timestamp gauge\n"
			output += fmt.Sprintf("redis_last_reconnect_timestamp %d %d\n", brokerMetrics.LastReconnectTime.Unix(), now)
		}
	}

	c.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	c.String(http.StatusOK, output)
}


func (h *MetricsHandler) HandleJSONMetrics(c *gin.Context) {
	if h.liveService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Live service not available",
		})
		return
	}

	
	wsMetrics := h.liveService.GetWebSocketMetrics()

	
	brokerMetrics := h.liveService.GetBrokerMetrics()

	response := gin.H{
		"websocket": gin.H{
			"active_connections":     wsMetrics.ActiveConnections,
			"total_connections":      wsMetrics.TotalConnections,
			"connections_rejected":   wsMetrics.ConnectionsRejected,
			"messages_received":      wsMetrics.MessagesReceived,
			"messages_sent":          wsMetrics.MessagesSent,
			"errors":                 wsMetrics.Errors,
			"last_error":             wsMetrics.LastError,
			"last_error_time":        wsMetrics.LastErrorTime,
			"average_latency_ms":     wsMetrics.GetAverageLatencyMs(),
			"message_throughput_sec": wsMetrics.GetMessageThroughput(),
		},
	}

	if brokerMetrics != nil {
		response["redis"] = gin.H{
			"events_published":    brokerMetrics.GetEventsPublished(),
			"events_received":     brokerMetrics.GetEventsReceived(),
			"publish_errors":      brokerMetrics.GetPublishErrors(),
			"reconnect_count":     brokerMetrics.GetReconnectCount(),
			"last_reconnect_time": brokerMetrics.LastReconnectTime,
		}
	}

	c.JSON(http.StatusOK, response)
}
