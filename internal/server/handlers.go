package server

import (
	"encoding/json"
	"fmt"
	"gitlab.com/ignitionrobotics/billing/credits/pkg/api"
	"io"
	"net/http"
)

// GetBalance is an HTTP handler to call the api.CreditsV1's GetBalance method.
func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	var in api.GetBalanceRequest
	if err := s.readBodyJSON(w, r, &in); err != nil {
		return
	}

	out, err := s.credits.GetBalance(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeResponse(w, &out)
}

// IncreaseCredits is an HTTP handler to call the api.CreditsV1's IncreaseCredits method.
func (s *Server) IncreaseCredits(w http.ResponseWriter, r *http.Request) {
	var in api.IncreaseCreditsRequest
	if err := s.readBodyJSON(w, r, &in); err != nil {
		return
	}

	out, err := s.credits.IncreaseCredits(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeResponse(w, &out)
}

// DecreaseCredits is an HTTP handler to call the api.CreditsV1's DecreaseCredits method.
func (s *Server) DecreaseCredits(w http.ResponseWriter, r *http.Request) {
	var in api.DecreaseCreditsRequest
	if err := s.readBodyJSON(w, r, &in); err != nil {
		return
	}

	out, err := s.credits.DecreaseCredits(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeResponse(w, &out)
}

// ConvertCurrency is an HTTP handler to call the api.CreditsV1's ConvertCurrency method.
func (s *Server) ConvertCurrency(w http.ResponseWriter, r *http.Request) {
	var in api.ConvertCurrencyRequest
	if err := s.readBodyJSON(w, r, &in); err != nil {
		return
	}

	out, err := s.credits.ConvertCurrency(r.Context(), in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.writeResponse(w, &out)
}

func (s *Server) writeResponse(w http.ResponseWriter, out interface{}) {
	body, err := json.Marshal(out)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s - %s", http.StatusText(http.StatusInternalServerError), "Failed to write JSON body"), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s - %s", http.StatusText(http.StatusInternalServerError), "Failed to write body"), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) readBodyJSON(w http.ResponseWriter, r *http.Request, in interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s - %s", http.StatusText(http.StatusBadRequest), "Failed to read body"), http.StatusBadRequest)
		return err
	}

	if err = json.Unmarshal(body, &in); err != nil {
		http.Error(w, fmt.Sprintf("%s - %s", http.StatusText(http.StatusInternalServerError), "Failed to read JSON body"), http.StatusInternalServerError)
		return err
	}

	return nil
}
