package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/sixojke/test-service/internal/domain"
)

func (h *Handler) createGood(c *gin.Context) {
	var input domain.ItemCreateInp
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusBadRequest,
			Message: "errors.good.badRequest",
		})

		return
	}

	projectId, err := processIntParam(c, "projectId")
	if err != nil {
		return
	}

	input.ProjectId = projectId

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "errors.good.unprocessableEntity",
		})

		return
	}

	good, err := h.services.Goods.Create(input)
	if err != nil {
		log.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "errors.good.internalServerError",
		})

		return
	}

	c.JSON(http.StatusOK, good)
}

func (h *Handler) updateGood(c *gin.Context) {
	var input domain.ItemUpdateInp
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusBadRequest,
			Message: "errors.good.badRequest",
		})

		return
	}

	id, err := processIntParam(c, "id")
	if err != nil {
		return
	}

	input.Id = id

	projectId, err := processIntParam(c, "projectId")
	if err != nil {
		return
	}

	input.ProjectId = projectId

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "errors.good.unprocessableEntity",
		})

		return
	}

	good, err := h.services.Goods.Update(input)
	if err != nil {
		log.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "errors.good.internalServerError",
		})

		return
	}

	if good.Removed {
		newErrorResponse(c, http.StatusNotFound, errorResponse{
			Code:    3,
			Message: "errors.good.notFound",
			Details: map[string]interface{}{},
		})

		return
	}

	c.JSON(http.StatusOK, good)
}

func (h *Handler) deleteGood(c *gin.Context) {
	var input domain.ItemDeleteInp
	id, err := processIntParam(c, "id")
	if err != nil {
		return
	}

	input.Id = id

	projectId, err := processIntParam(c, "projectId")
	if err != nil {
		return
	}

	input.ProjectId = projectId

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "errors.good.unprocessableEntity",
		})

		return
	}

	good, err := h.services.Goods.Delete(input)
	if err != nil {
		if strings.Contains(err.Error(), "item is already removed") {
			newErrorResponse(c, http.StatusInternalServerError, errorResponse{
				Code:    http.StatusNotFound,
				Message: "errors.good.notFound",
			})
		} else {
			log.Error(err)
			newErrorResponse(c, http.StatusInternalServerError, errorResponse{
				Code:    http.StatusInternalServerError,
				Message: "errors.good.internalServerError",
			})
		}

		return
	}

	c.JSON(http.StatusOK, good)
}

func (h *Handler) reprioritiizeGood(c *gin.Context) {
	var input domain.ItemReprioritiizeInp
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusBadRequest,
			Message: "errors.good.badRequest",
		})

		return
	}

	id, err := processIntParam(c, "id")
	if err != nil {
		return
	}

	input.Id = id

	projectId, err := processIntParam(c, "projectId")
	if err != nil {
		return
	}

	input.ProjectId = projectId

	if err := input.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: "errors.good.unprocessableEntity",
		})

		return
	}

	goods, err := h.services.Goods.Reprioritiize(input)
	if err != nil {
		log.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "errors.good.internalServerError",
		})
	}

	c.JSON(http.StatusOK, goods)
}

func (h *Handler) getListGoods(c *gin.Context) {
	limit, err := processIntParam(c, "limit")
	if err != nil {
		return
	}

	offset, err := processIntParam(c, "offset")
	if err != nil {
		return
	}

	good, err := h.services.Goods.GetList(limit, offset)
	if err != nil {
		log.Error(err)
		newErrorResponse(c, http.StatusInternalServerError, errorResponse{
			Code:    http.StatusInternalServerError,
			Message: "errors.good.internalServerError",
		})

		return
	}

	c.JSON(http.StatusOK, good)
}

func processIntParam(c *gin.Context, paramName string) (int, error) {
	num, err := strconv.Atoi(c.Query(paramName))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, errorResponse{
			Code:    http.StatusBadRequest,
			Message: "errors.good.badRequest",
		})

		return 0, err
	}

	return num, nil
}
