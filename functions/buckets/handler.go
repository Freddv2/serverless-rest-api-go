package buckets

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) FindById(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	bucketId := ctx.Param("bucketId")
	bucket, err := h.service.FindById(ctx, tenantId, bucketId)
	if errors.Is(err, ErrNotFound) {
		ctx.String(404, "Not Found")
	} else if err != nil {
		ctx.JSON(500, err)
	}

	ctx.JSON(200, bucket)
}

func (h *handler) Search(ctx *gin.Context) {
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
	buckets, err := h.service.Search(ctx, searchContext)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.JSON(200, buckets)
}

func (h *handler) Create(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")

	var bucket Bucket
	err := json.NewDecoder(ctx.Request.Body).Decode(&bucket)
	if err != nil {
		ctx.JSON(500, err)
	}

	id, err := h.service.Create(ctx, tenantId, bucket)
	if errors.Is(err, ErrAlreadyExists) {
		ctx.String(409, "Already Exists")
	} else if err != nil {
		ctx.JSON(500, err)
	}

	ctx.String(201, id)
}

func (h *handler) Update(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")

	var bucket Bucket
	err := json.NewDecoder(ctx.Request.Body).Decode(&bucket)
	if err != nil {
		ctx.JSON(500, err)
	}

	err = h.service.Update(ctx, tenantId, bucket)
	if errors.Is(err, ErrNotFound) {
		ctx.String(404, "Not Found")
	} else if err != nil {
		ctx.JSON(500, err)
	}

	ctx.String(200, "")
}

func (h *handler) Delete(ctx *gin.Context) {
	tenantId := ctx.Param("tenantId")
	bucketId := ctx.Param("bucketId")

	err := h.service.Delete(ctx, tenantId, bucketId)
	if err != nil {
		ctx.JSON(500, err)
	}
	ctx.String(200, "")
}
