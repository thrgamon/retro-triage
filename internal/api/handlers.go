package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/thrgamon/retro-triage/internal/analyser"
	"github.com/thrgamon/retro-triage/internal/config"
	"github.com/thrgamon/retro-triage/internal/db"
	"github.com/thrgamon/retro-triage/internal/domain"
)

type HandlerConfig struct {
	Queries  *db.Queries
	Analyser *analyser.Analyser
	Cfg      config.Config
}

type Handler struct {
	queries  *db.Queries
	analyser *analyser.Analyser
	cfg      config.Config
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{
		queries:  cfg.Queries,
		analyser: cfg.Analyser,
		cfg:      cfg.Cfg,
	}
}

// Routes registers all HTTP routes on the given router group.
func (h *Handler) Routes(rg *gin.RouterGroup) {
	rg.GET("/health", h.Health)

	retros := rg.Group("/retros")
	{
		retros.POST("", h.CreateRetro)
		retros.GET("", h.ListRetros)
		retros.GET("/:id", h.GetRetro)
		retros.POST("/:id/cards", h.CreateCard)
		retros.DELETE("/:id/cards/:cardId", h.DeleteCard)
		retros.POST("/:id/analyse", h.Analyse)
		retros.GET("/:id/analysis", h.GetAnalysis)
	}
}

// Health godoc
// @Summary Health check
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateRetro godoc
// @Summary Create a new retro
// @Tags retros
// @Accept json
// @Produce json
// @Param body body domain.CreateRetroRequest true "Retro details"
// @Success 201 {object} domain.RetroListItem
// @Failure 400 {object} map[string]string
// @Router /api/retros [post]
func (h *Handler) CreateRetro(c *gin.Context) {
	var req domain.CreateRetroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	retro, err := h.queries.CreateRetro(c.Request.Context(), req.Name)
	if err != nil {
		slog.Error("creating retro", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create retro"})
		return
	}

	c.JSON(http.StatusCreated, domain.RetroListItem{
		ID:        retro.ID,
		Name:      retro.Name,
		CreatedAt: retro.CreatedAt,
	})
}

// ListRetros godoc
// @Summary List recent retros
// @Tags retros
// @Produce json
// @Success 200 {array} domain.RetroListItem
// @Router /api/retros [get]
func (h *Handler) ListRetros(c *gin.Context) {
	retros, err := h.queries.ListRetros(c.Request.Context())
	if err != nil {
		slog.Error("listing retros", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list retros"})
		return
	}

	items := make([]domain.RetroListItem, len(retros))
	for i, r := range retros {
		items[i] = domain.RetroListItem{
			ID:        r.ID,
			Name:      r.Name,
			CreatedAt: r.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, items)
}

// GetRetro godoc
// @Summary Get a retro with its cards
// @Tags retros
// @Produce json
// @Param id path string true "Retro ID"
// @Success 200 {object} domain.RetroResponse
// @Failure 404 {object} map[string]string
// @Router /api/retros/{id} [get]
func (h *Handler) GetRetro(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retro id"})
		return
	}

	retro, err := h.queries.GetRetro(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "retro not found"})
		return
	}

	cards, err := h.queries.GetCardsByRetroID(c.Request.Context(), id)
	if err != nil {
		slog.Error("getting cards", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cards"})
		return
	}

	cardResponses := make([]domain.CardResponse, len(cards))
	for i, card := range cards {
		cardResponses[i] = domain.CardResponse{
			ID:         card.ID,
			ColumnType: domain.ColumnType(card.ColumnType),
			Content:    card.Content,
			AuthorName: card.AuthorName,
			CreatedAt:  card.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, domain.RetroResponse{
		ID:        retro.ID,
		Name:      retro.Name,
		Cards:     cardResponses,
		CreatedAt: retro.CreatedAt,
	})
}

// CreateCard godoc
// @Summary Add a card to a retro
// @Tags retros
// @Accept json
// @Produce json
// @Param id path string true "Retro ID"
// @Param body body domain.CreateCardRequest true "Card details"
// @Success 201 {object} domain.CardResponse
// @Failure 400 {object} map[string]string
// @Router /api/retros/{id}/cards [post]
func (h *Handler) CreateCard(c *gin.Context) {
	retroID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retro id"})
		return
	}

	var req domain.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !domain.ValidColumnTypes[req.ColumnType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column_type"})
		return
	}

	card, err := h.queries.CreateCard(c.Request.Context(), db.CreateCardParams{
		RetroID:    retroID,
		ColumnType: string(req.ColumnType),
		Content:    req.Content,
		AuthorName: req.AuthorName,
	})
	if err != nil {
		slog.Error("creating card", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create card"})
		return
	}

	c.JSON(http.StatusCreated, domain.CardResponse{
		ID:         card.ID,
		ColumnType: domain.ColumnType(card.ColumnType),
		Content:    card.Content,
		AuthorName: card.AuthorName,
		CreatedAt:  card.CreatedAt,
	})
}

// DeleteCard godoc
// @Summary Delete a card from a retro
// @Tags retros
// @Param id path string true "Retro ID"
// @Param cardId path string true "Card ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Router /api/retros/{id}/cards/{cardId} [delete]
func (h *Handler) DeleteCard(c *gin.Context) {
	retroID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retro id"})
		return
	}

	cardID, err := uuid.Parse(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card id"})
		return
	}

	err = h.queries.DeleteCard(c.Request.Context(), db.DeleteCardParams{
		ID:      cardID,
		RetroID: retroID,
	})
	if err != nil {
		slog.Error("deleting card", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete card"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Analyse godoc
// @Summary Trigger LLM analysis of a retro
// @Tags retros
// @Produce json
// @Param id path string true "Retro ID"
// @Success 200 {object} domain.AnalysisResponse
// @Failure 400 {object} map[string]string
// @Router /api/retros/{id}/analyse [post]
func (h *Handler) Analyse(c *gin.Context) {
	retroID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retro id"})
		return
	}

	cards, err := h.queries.GetCardsByRetroID(c.Request.Context(), retroID)
	if err != nil {
		slog.Error("getting cards for analysis", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cards"})
		return
	}

	if len(cards) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no cards to analyse"})
		return
	}

	cardResponses := make([]domain.CardResponse, len(cards))
	for i, card := range cards {
		cardResponses[i] = domain.CardResponse{
			ID:         card.ID,
			ColumnType: domain.ColumnType(card.ColumnType),
			Content:    card.Content,
			AuthorName: card.AuthorName,
			CreatedAt:  card.CreatedAt,
		}
	}

	result, err := h.analyser.Analyse(c.Request.Context(), cardResponses)
	if err != nil {
		slog.Error("analysing retro", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "analysis failed"})
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		slog.Error("marshalling analysis result", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store analysis"})
		return
	}

	analysis, err := h.queries.UpsertAnalysis(c.Request.Context(), db.UpsertAnalysisParams{
		RetroID: retroID,
		Result:  resultJSON,
	})
	if err != nil {
		slog.Error("storing analysis", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store analysis"})
		return
	}

	c.JSON(http.StatusOK, domain.AnalysisResponse{
		ID:        analysis.ID,
		RetroID:   analysis.RetroID,
		Result:    analysis.Result,
		CreatedAt: analysis.CreatedAt,
	})
}

// GetAnalysis godoc
// @Summary Get the analysis for a retro
// @Tags retros
// @Produce json
// @Param id path string true "Retro ID"
// @Success 200 {object} domain.AnalysisResponse
// @Failure 404 {object} map[string]string
// @Router /api/retros/{id}/analysis [get]
func (h *Handler) GetAnalysis(c *gin.Context) {
	retroID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid retro id"})
		return
	}

	analysis, err := h.queries.GetAnalysisByRetroID(c.Request.Context(), retroID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no analysis found"})
		return
	}

	c.JSON(http.StatusOK, domain.AnalysisResponse{
		ID:        analysis.ID,
		RetroID:   analysis.RetroID,
		Result:    analysis.Result,
		CreatedAt: analysis.CreatedAt,
	})
}
