package main

import (
	"fmt"
	"net/http"
	"time"

	"api5back/src/database"
	"api5back/src/server"

	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Contador de requisições HTTP
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Número total de requisições HTTP recebidas.",
		},
		[]string{"method", "endpoint"}, // Labels para diferenciar métricas
	)

	// Histograma para medir o tempo de resposta das requisições
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos.",
			Buckets: prometheus.DefBuckets, // Faixas padrão: {0.005, 0.01, 0.025...}
		},
		[]string{"endpoint"}, // Label para diferenciar endpoints
	)

	// Gauge para monitorar conexões ao banco
	dbConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connections",
			Help: "Número de conexões abertas ao banco de dados.",
		},
		[]string{"database_type"}, // Exemplo: "DB" ou "DW"
	)
)

func init() {
	// Registrar métricas no Prometheus
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbConnections)
}

func main() {
	// Conectar ao banco de dados
	dbClient, err := database.Setup("DB")
	if err != nil {
		panic(fmt.Errorf("failed to setup normalized database: %v", err))
	}
	defer dbClient.Close()

	// Atualizar métrica de conexões ao banco
	dbConnections.WithLabelValues("DB").Set(1)

	dwClient, err := database.Setup("DW")
	if err != nil {
		panic(fmt.Errorf("failed to setup data warehouse: %v", err))
	}
	defer dwClient.Close()

	// Atualizar métrica de conexões ao banco
	dbConnections.WithLabelValues("DW").Set(1)

	// Criar servidor
	srv := server.NewServer(dbClient, dwClient)

	// Configurar métricas padrão do gin-metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(srv)

	// Middleware para medir métricas personalizadas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Incrementar contador de requisições
		httpRequestsTotal.WithLabelValues(r.Method, "/").Inc()

		// Processar requisição
		srv.ServeHTTP(w, r)

		// Medir tempo de execução
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues("/").Observe(duration)
	})

	// Expor endpoint de métricas
	http.Handle("/metrics", promhttp.Handler())

	// Iniciar servidor
	http.ListenAndServe(":8080", nil)
}
