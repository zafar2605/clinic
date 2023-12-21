package handler

import (
	"context"
	"errors"
	"market_system/config"
	"market_system/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @Summary MakePay
// @Description Pay client's Sale.
// @Tags Pay
// @Accept json
// @Produce json
// @Param sale_id query string ture "sale_id"
// @Param money query float64 true "pay_money"
// @Success 200 {object} models.Coming "Payed"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Sale not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /make_pay [put]
func (h *Handler) MakePay(c *gin.Context) {

	var (
		incrementId = c.Query("sale_id")
		money       = (cast.ToFloat64(c.Query("money")))
		Id          string
		ClientID    string
		BranchID    string
		IncrementID string
		TotalPrice  float64
	)

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	saleList, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{Limit: 10000})
	for _, v := range saleList.Sales {
		if v.IncrementID == incrementId {
			if v.TotalPrice/2 < money {
				if err != nil {
					handleResponse(c, http.StatusInternalServerError, errors.New("not enough money"))
					return
				}
				Id = v.Id
				ClientID = v.ClientID
				BranchID = v.BranchID
				IncrementID = v.IncrementID
				TotalPrice = v.TotalPrice
			}
		}
	}
	_, err = h.strg.Sale().Update(ctx, &models.UpdateSale{
		Id:          Id,
		ClientID:    ClientID,
		BranchID:    BranchID,
		IncrementID: IncrementID,
		TotalPrice:  TotalPrice,
		Paid:        money,
		Debd:        TotalPrice - money,
	})
	if err != nil {
		handleResponse(c, 500, err.Error())
		return
	}

	handleResponse(c, 202, "successful payment")
}

// @Summary Registration
// @Description Get clients.
// @Tags Registration
// @Accept json
// @Produce json
// @Param from query string true "From day"
// @Param to query string true "To day"
// @Success 200 {object} models.GetListClientResponse "Client details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Client not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /registration [get]
func (h *Handler) Registration(c *gin.Context) {

	var (
		from    = c.Query("from")
		to      = c.Query("to")
		clients = models.GetListClientResponse{}
	)
	fromm, err := time.Parse("2006-01-02", from)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	too, err := time.Parse("2006-01-02", to)
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	clientList, err := h.strg.Client().GetList(ctx, &models.GetListClientRequest{Limit: 1000000})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range clientList.Clients {
		created_at, err := time.Parse("2006-01-02", v.CreatedAt[:10])
		if err != nil {
			handleResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		if too.After(created_at) && fromm.Before(created_at) {
			if len(created_at.String()) > 0 {
				clients.Clients = append(clients.Clients, v)
			}
		}
	}

	handleResponse(c, http.StatusOK, clients)
}

// @Summary Branch
// @Description Get Branch sales.
// @Tags Branch Sales
// @Accept json
// @Produce json
// @Param branch_id query string true "Branch Id"
// @Success 200 {object} models.Doc "Branch details"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 404 {object} ErrorResponse "Branch not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /branch_doc [get]
func (h *Handler) BranchDoc(c *gin.Context) {

	var (
		branchId = c.Query("branch_id")
		resp     = models.Doc{}
		sale_id  = ""
	)

	ctx, cancel := context.WithTimeout(context.Background(), config.CtxTimeout)
	defer cancel()

	branch, err := h.strg.Branch().GetByID(ctx, &models.BranchPrimaryKey{Id: branchId})
	if err != nil {
		handleResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	resp.BranchID = branchId
	resp.BranchName = branch.Name

	salesList, err := h.strg.Sale().GetList(ctx, &models.GetListSaleRequest{Limit: 1000000})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	salesProductList, err := h.strg.SaleProduct().GetList(ctx, &models.GetListSaleProductRequest{Limit: 1000000})
	if err != nil {
		handleResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, v := range salesList.Sales {
		if v.BranchID == branchId {
			resp.TotalSalePrice += v.TotalPrice
			sale_id = v.Id
		}
	}
	for _, v := range salesProductList.SaleProducts {
		if v.SaleID == sale_id {
			resp.TotalSaleQuantity += v.Quantity
		}
	}

	handleResponse(c, http.StatusOK, resp)
}
