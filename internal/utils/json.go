package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/harry713j/minurly/internal/apperrors"
	"github.com/rs/zerolog"
)

func RespondJSON[T any](w http.ResponseWriter, code int, payload T, logger zerolog.Logger) {
	w.Header().Add("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		logger.Err(err).Msg(fmt.Sprintf("failed to marshalling the json payload %v", payload))

		apiErr := apperrors.NewInternalServerErr()

		RespondError(w, apiErr)
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(data); err != nil {
		logger.Err(err).Msg(fmt.Sprintf("failed to write the response %v", payload))
		apiErr := apperrors.NewInternalServerErr()

		RespondError(w, apiErr)
		return
	}

	logger.Info().Msg("successfully written the response!")
}

func RespondError(w http.ResponseWriter, err error) {
	apiErr, ok := err.(*apperrors.ApiError)

	if !ok {
		apiErr = apperrors.NewInternalServerErr()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)
	json.NewEncoder(w).Encode(apiErr)
}
