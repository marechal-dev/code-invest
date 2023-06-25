package transformers

import (
	"github.com/marechal-dev/code-invest/stock/internal/domain/market/dtos"
	"github.com/marechal-dev/code-invest/stock/internal/domain/market/entities"
)

func TransformInput(input dtos.TradeInput) *entities.Order {
	asset := entities.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entities.NewInvestor(input.InvestorID)
	order := entities.NewOrder(input.OrderID, investor, asset, input.Shares, input.Price, input.OrderType)
	if input.CurrentShares > 0 {
		assetPosition := entities.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entities.Order) *dtos.OrderOutput {
	output := &dtos.OrderOutput{
		OrderID:    order.ID,
		InvestorID: order.Investor.ID,
		AssetID:    order.Asset.ID,
		OrderType:  order.OrderType,
		Status:     order.Status,
		Partial:    order.PendingShares,
		Shares:     order.Shares,
	}

	var transactionsOutput []*dtos.TransactionOutput
	for _, transaction := range order.Transactions {
		transactionOutput := &dtos.TransactionOutput{
			TransactionID: transaction.ID,
			BuyerID:       transaction.BuyingOrder.Investor.ID,
			SellerID:      transaction.SellingOrder.Investor.ID,
			AssetID:       transaction.SellingOrder.Asset.ID,
			Price:         transaction.Price,
			Shares:        transaction.SellingOrder.Shares - transaction.SellingOrder.PendingShares,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput)
	}
	output.TransactionsOutput = transactionsOutput

	return output
}
