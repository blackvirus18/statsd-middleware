package statsd_middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	statsdv2 "gopkg.in/alexcesaro/statsd.v2"
)

var statsD *statsdv2.Client

type StatsDConfig struct {
	appName string
	host    string
	port    string
	enabled bool
}

func InitStatsDConfig(appName string, host string, port string, enabled bool) StatsDConfig {
	return StatsDConfig{
		appName: appName,
		host:    host,
		port:    port,
		enabled: enabled,
	}
}

func InitiateStatsDMetrics(config StatsDConfig) error {
	statsDConfig := config

	if statsDConfig.enabled {
		address := fmt.Sprintf("%s:%s", statsDConfig.host, statsDConfig.port)
		var err error
		statsD, err = statsdv2.New(statsdv2.Address(address), statsdv2.Prefix(statsDConfig.appName))

		if err != nil {
			return err
		}
	}

	return nil
}

func CloseStatsDClient() {
	if statsD != nil {
		statsD.Close()
	}
}

func IncrementInStatsD(key string) bool {
	if statsD == nil {
		return false
	}
	statsD.Increment(key)
	return true
}

func DecrementInStatsD(key string) bool {
	if statsD == nil {
		return false
	}
	statsD.Count(key, -1)
	return true
}

func TimingKeyInStatsD(key string, value interface{}) bool {
	if statsD == nil {
		return false
	}
	statsD.Timing(key, value)
	return true

}

func GaugeKeyInStatsD(key string, value interface{}) bool {
	if statsD == nil {
		return false
	}
	statsD.Gauge(key, value)
	return true

}

func CountKeyInStatsD(key string, value interface{}) bool {
	if statsD == nil {
		return false
	}
	statsD.Count(key, value)
	return true
}

func TimingInStatsD() *statsdv2.Timing {
	if statsD == nil {
		return nil
	}
	t := statsD.NewTiming()
	return &t
}

func SendInStatsD(key string, t *statsdv2.Timing) bool {
	if statsD == nil {
		return false
	}
	t.Send(key)
	return true
}

func GetKeyStructure(url string) string {
	baseKey := "go.response"
	basePath := strings.Split(url, "/GF")[0]
	keyBasePath := strings.Replace(basePath, "/", ".", len(basePath))
	key := baseKey + keyBasePath
	return key
}

func StatsDMiddlewareLogger() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		path := r.URL.Path
		for _, v := range vars {
			path = strings.Replace(path, v, "", len(path))
		}
		key := GetKeyStructure(r.URL.Path)
		t := TimingInStatsD()
		noOfGoRoutine := runtime.NumGoroutine()
		next(rw, r)
		SendInStatsD(key+".time", t)
		IncrementInStatsD(key + ".calls")
		GaugeKeyInStatsD(key+".goroutines", noOfGoRoutine)
	})
}
func StartHystrixStats(statsDConfig StatsDConfig, hystrixPrefix string) {
	if !statsDConfig.enabled {
		return
	}
	address := fmt.Sprintf("%s:%s", statsDConfig.host, statsDConfig.port)
	c, _ := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: address,
		Prefix:     hystrixPrefix,
	})
	metricCollector.Registry.Register(c.NewStatsdCollector)
}
