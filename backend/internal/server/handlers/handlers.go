package handlers

import (
	"backend/internal/models"
	"backend/internal/repository/postgres"
	"backend/internal/repository/redis"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	Repo  *postgres.Storage
	Cache *redis.Storage
}

func (h *Handler) Main(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("<h1>Crypto Metrics Server</h1><p>Status: Running</p><a href='/swagger/index.html'>View API Docs</a>"))
}

// DumpStat godoc
// @Summary      Сохранить метрику
// @Description  Принимает JSON с данными о валюте и сохраняет в базу
// @Tags         metrics
// @Accept       json
// @Produce      json
// @Param        stat  body      models.Stat  true  "Данные метрики"
// @Success      201   {string}  string       "Created"
// @Failure      400   {string}  string       "Bad Request"
// @Failure      500   {string}  string       "Internal Server Error"
// @Router       /update [post]
func (h *Handler) DumpStat(w http.ResponseWriter, r *http.Request) {
	var s models.Stat
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		log.Printf("JSON Decode Error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// for some reason agent sends incorrect data
	if s.Timedump.IsZero() {
		s.Timedump = time.Now()
	}

	// checkin agent logs
	log.Printf("Saving to DB: 	Exchange=%s, 	Base=%s, 	Quote=%s, 	AskPrice=%f, 	BidPrice=%f, 	Time=%v",
		s.Source, s.Base, s.Quote, s.AskPrice, s.BidPrice, s.Timedump)

	if err := h.Repo.Save(s); err != nil {
		log.Printf("DATABASE SAVE ERROR: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := h.Cache.Save(s); err != nil {
		log.Printf("CACHE SAVE ERROR: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetStat godoc
// @Summary      Get all crypto stats
// @Description  Returns a list of all collected cryptocurrency statistics from the database to the stat router
// @Tags         stats
// @Produce      json
// @Success      200  {array}   models.Stat
// @Failure      500  {object}  string
// @Router       /stat [get]
func (h *Handler) GetStat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stats, err := h.Cache.GetStat()
	if err != nil {
		log.Printf("error getting data from db: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
