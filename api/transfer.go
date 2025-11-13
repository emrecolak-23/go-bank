package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/emrecolak-23/go-bank/db/sqlc"
	"github.com/emrecolak-23/go-bank/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(context *gin.Context) {

	var req createTransferRequest

	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validAccount(context, req.FromAccountID, req.Currency)

	authPayload := context.MustGet(authorizationKey).(*token.Payload)

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("fromaccount doesnt belong to the authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if !valid {
		return
	}

	_, valid = server.validAccount(context, req.ToAccountID, req.Currency)

	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(context, arg)

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(context *gin.Context, accountId int64, currency string) (db.Account, bool) {

	account, err := server.store.GetAccount(context, accountId)

	if err != nil {

		if errors.Is(err, db.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		context.JSON(http.StatusInternalServerError, errorResponse(err))

		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] curreny mismatch: %s vs %s", accountId, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true

}
