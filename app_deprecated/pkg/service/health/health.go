package health

var _ HealthService = (*Service)(nil)

// Service implements the HealthChecker interface for a health checking service.
type Service struct{}

// NewService creates a new instance of a Service.
func NewService() *Service {
	return &Service{}
}

// CheckHealth performs a health check and returns a CheckHealthResponse.
// It returns a more descriptive status compatible with typical monitoring setups in AWS environments.
func (s *Service) CheckHealth() CheckHealthResponse {
	// Example check (this should be replaced with actual health checks)
	healthy := true // Assume the service is healthy after checks
	if healthy {
		return CheckHealthResponse{
			Healthy: healthy,
			Status:  "Service is healthy",
		}
	} else {
		return CheckHealthResponse{
			Healthy: healthy,
			Status:  "Service is unhealthy",
		}
	}
}
