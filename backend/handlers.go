package main

import (
	"encoding/csv"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) GetState(c *gin.Context) {
	state := h.Service.GetState()
	resp := StateResponse{
		PositionX: state.Robot.PositionX,
		PositionY: state.Robot.PositionY,
		Holding:   state.Robot.Holding,
		Grid:      state.Grid,
		Won:       h.Service.HasWon(),
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ProcessCommand(c *gin.Context) {
	var req CommandRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var (
		state State
		err   error
	)

	switch req.Action {
	case Move:
		if req.Direction == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing direction for move action"})
			return
		}
		state, err = h.Service.Move(req.Direction)
	case PickUp:
		state, err = h.Service.Pick()
	case Drop:
		state, err = h.Service.Drop()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := StateResponse{
		PositionX: state.Robot.PositionX,
		PositionY: state.Robot.PositionY,
		Holding:   state.Robot.Holding,
		Grid:      state.Grid,
		Won:       h.Service.HasWon(),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ExportHistory(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=history.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	if err := writer.Write([]string{"Timestamp", "Moves"}); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, record := range h.Service.GetHistory() {
		if err := writer.Write([]string{record.Timestamp.Format(time.RFC3339), record.Moves}); err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
	}
}
