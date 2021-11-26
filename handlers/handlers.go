package handlers

import (
	GardenClient "code.cloudfoundry.org/garden/client"
	"net/http"

	"code.cloudfoundry.org/executor"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/locket/metrics/helpers"
	"code.cloudfoundry.org/rep"
	"code.cloudfoundry.org/rep/auctioncellrep"
	"code.cloudfoundry.org/rep/evacuation/evacuation_context"
	"github.com/tedsuo/rata"
)

func New(
	localCellClient auctioncellrep.AuctionCellClient,
	localMetricCollector MetricCollector,
// Added for PaaS-TA
	gardenClient GardenClient.Client,
	executorClient executor.Client,
	evacuatable evacuation_context.Evacuatable,
	requestMetrics helpers.RequestMetrics,
	logger lager.Logger,
	secure bool,
) rata.Handlers {

	handlers := rata.Handlers{}
	if secure {
		stateHandler := newStateHandler(localCellClient, requestMetrics)
		containerMetricsHandler := newContainerMetricsHandler(localMetricCollector, requestMetrics)
		performHandler := newPerformHandler(localCellClient, requestMetrics)
		resetHandler := newResetHandler(localCellClient, requestMetrics)
		stopLrpHandler := NewStopLRPInstanceHandler(executorClient, requestMetrics)
		cancelTaskHandler := newCancelTaskHandler(executorClient, requestMetrics)

		handlers[rep.StateRoute] = logWrap(stateHandler.ServeHTTP, logger)
		handlers[rep.ContainerMetricsRoute] = logWrap(containerMetricsHandler.ServeHTTP, logger)
		handlers[rep.PerformRoute] = logWrap(performHandler.ServeHTTP, logger)
		handlers[rep.SimResetRoute] = logWrap(resetHandler.ServeHTTP, logger)

		handlers[rep.StopLRPInstanceRoute] = logWrap(stopLrpHandler.ServeHTTP, logger)
		handlers[rep.CancelTaskRoute] = logWrap(cancelTaskHandler.ServeHTTP, logger)
	} else {
		pingHandler := newPingHandler(requestMetrics)
		evacuationHandler := newEvacuationHandler(evacuatable, requestMetrics)

		handlers[rep.PingRoute] = logWrap(pingHandler.ServeHTTP, logger)
		handlers[rep.EvacuateRoute] = logWrap(evacuationHandler.ServeHTTP, logger)

		// Added for PaaS-TA
		containerHandler := NewContainerListHandler(logger, executorClient, gardenClient)
		handlers[rep.ContainerListRoute] = logWrap(containerHandler.ServeHTTP, logger)
	}

	return handlers
}

// this isn't being used in the Rep anymore. It is used in tests that run a
// fake cell. Without this function those tests will have to replicate the code
// below. Those places are auctioneer fake_cell_test.go and rep's
// handlers_suite_test.go
func NewLegacy(
	localCellClient auctioncellrep.AuctionCellClient,
	localMetricCollector MetricCollector,

// Added for PaaS-TA
	gardenClient GardenClient.Client,

	executorClient executor.Client,
	evacuatable evacuation_context.Evacuatable,
	requestMetrics helpers.RequestMetrics,
	logger lager.Logger,
) rata.Handlers {
	// Added for PaaS-TA
	//insecureHandlers := New(localCellClient, localMetricCollector, executorClient, evacuatable, requestMetrics, logger, false)
	//secureHandlers := New(localCellClient, localMetricCollector, executorClient, evacuatable, requestMetrics, logger, true)
	insecureHandlers := New(localCellClient, localMetricCollector, gardenClient, executorClient, evacuatable, requestMetrics, logger, false)
	secureHandlers := New(localCellClient, localMetricCollector, gardenClient, executorClient, evacuatable, requestMetrics, logger, true)
	for name, handler := range secureHandlers {
		insecureHandlers[name] = handler
	}
	return insecureHandlers
}

func logWrap(loggable func(http.ResponseWriter, *http.Request, lager.Logger), logger lager.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestLog := logger.Session("request", lager.Data{
			"method":  r.Method,
			"request": r.URL.String(),
		})

		defer requestLog.Debug("done")
		requestLog.Debug("serving")

		loggable(w, r, requestLog)
	}
}