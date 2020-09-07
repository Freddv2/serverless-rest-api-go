package portfolio

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) FindById(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	id := ctx.Param("id")
	portfolio, err := h.service.FindById(ctx, tenantId, id)
	if errors.Is(err, ErrNotFound) {
		ctx.String(404, "Not Found")
	} else if err != nil {
		ctx.JSON(500, err)
	}

	ctx.JSON(200, portfolio)
}

func (h *Handler) Search(ctx *gin.Context) {
	var nbOfReturnedElements int
	if nb := ctx.Query("nbOfReturnedElement"); nb != "" {
		nbOfReturnedElements, _ = strconv.Atoi(nb)
	}
	var ids = make([]string, 0)
	if i := ctx.Query("ids"); i != "" {
		ids = strings.Split(ctx.Query("ids"), ",")
	}
	searchContext := SearchContext{
		TenantId:             ctx.Param("tenantId"),
		Name:                 ctx.Query("name"),
		NextPageCursor:       ctx.Query("nextPageCursor"),
		NbOfReturnedElements: nbOfReturnedElements,
		Ids:                  ids,
	}
	portfolios, err := h.service.Search(ctx, searchContext)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, portfolios)
}

func (h *Handler) Create(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")

	var portfolio Portfolio
	err := json.NewDecoder(ctx.Request.Body).Decode(&portfolio)
	if err != nil {
		ctx.JSON(500, err)
	}

	id, err := h.service.Create(ctx, tenantId, portfolio)
	if errors.Is(err, ErrAlreadyExists) {
		ctx.String(409, "Portfolio already exists")
	} else if err != nil {
		ctx.JSON(500, err)
	} else {
		ctx.String(201, id)
	}
}

func (h *Handler) Update(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")

	var portfolio Portfolio
	err := json.NewDecoder(ctx.Request.Body).Decode(&portfolio)
	if err != nil {
		ctx.JSON(500, err)
	}

	err = h.service.Update(ctx, tenantId, portfolio)
	if errors.Is(err, ErrNotFound) {
		ctx.String(404, "Not Found")
	} else if err != nil {
		ctx.JSON(500, err)
	}

	ctx.String(200, "")
}

func (h *Handler) Delete(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	id := ctx.Param("id")

	err := h.service.Delete(ctx, tenantId, id)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.String(200, "")
}
