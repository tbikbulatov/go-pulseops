package http

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	uc "github.com/tbikbulatov/go-pulseops/internal/alert/usecase/ingestalert"
	"github.com/tbikbulatov/go-pulseops/internal/platform/metrics"
	"github.com/tbikbulatov/go-pulseops/internal/platform/validation"
)

type Handler struct {
	validator *validation.Validator
	usecase   uc.IngestAlertUsecase
	metrics   metrics.AlertMetrics
}

func NewAlertHandler(v *validation.Validator, uc uc.IngestAlertUsecase, m metrics.AlertMetrics) *Handler {
	return &Handler{v, uc, m}
}

func (h *Handler) IngestAlert(c *echo.Context) error {
	var req IngestAlertRequest
	if err := c.Bind(&req); err != nil {
		h.metrics.IngestTotal.WithLabelValues(metrics.ResultValidationError).Inc()
		return c.JSON(http.StatusBadRequest, map[string]string{
			"code":    "invalid_json",
			"message": "invalid request body",
		})
	}

	if errs := h.validator.Validate(req); errs != nil {
		var validationErrs validator.ValidationErrors
		if !errors.As(errs, &validationErrs) {
			h.metrics.IngestTotal.WithLabelValues(metrics.ResultValidationError).Inc()
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}

		errmap := make(map[string]string, len(validationErrs))
		for _, f := range validationErrs {
			errmap[f.Field()] = f.Error()
		}
		h.metrics.IngestTotal.WithLabelValues(metrics.ResultValidationError).Inc()
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{"errors": errmap})
	}

	req.IntegrationKey = c.Param("integration_key")

	res, err := h.usecase.Handle(c.Request().Context(), req.toCommand())
	if err != nil {
		h.metrics.IngestTotal.WithLabelValues(metrics.ResultUsecaseError).Inc()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"code":    "internal_error",
			"message": "failed to ingest alert",
		})
	}

	h.metrics.IngestTotal.WithLabelValues(metrics.ResultSuccess).Inc()
	return c.JSON(http.StatusAccepted, map[string]string{
		"alert_id": res.AlertID,
		"status":   res.Status,
	})
}
