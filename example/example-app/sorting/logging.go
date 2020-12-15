package sorting

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"time"
)

type loggingService struct {
	logger  log.Logger
	service SortingService
}

func NewLoggingService(l log.Logger, s SortingService) SortingService {
	return &loggingService{
		logger:  l,
		service: s,
	}
}

func (s *loggingService) BubbleSort(numbers []int) {
	start := time.Now()

	s.service.BubbleSort(numbers)

	logger := log.With(s.logger,
		"method", "BubbleSort",
		"duration", time.Since(start),
	)

	level.Info(logger).Log()
}
