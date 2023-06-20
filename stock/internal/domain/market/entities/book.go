package entities

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders              []*Order
	Transactions        []*Transaction
	OrdersChannel       chan *Order
	OrdersChannelOutput chan *Order
	WaitGroup           *sync.WaitGroup
}

func NewBook(orderChannel chan *Order, orderChannelOutput chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Orders:              []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannel:       orderChannel,
		OrdersChannelOutput: orderChannelOutput,
		WaitGroup:           waitGroup,
	}
}

func (b *Book) Trade() {
	buyOrders := NewOrderQueue()
	sellOrders := NewOrderQueue()

	heap.Init(buyOrders)
	heap.Init(sellOrders)

	for order := range b.OrdersChannel {
		switch order.OrderType {
		case "BUY":
			buyOrders.Push(order)

			if sellOrders.Len() > 0 && sellOrders.Orders[0].Price <= order.Price {
				sellOrder := sellOrders.Pop().(*Order)

				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)
					b.AddTransaction(transaction, b.WaitGroup)

					sellOrder.AddTransaction(transaction)
					order.AddTransaction(transaction)

					b.sendToOutputChannel(sellOrder, order)

					if sellOrder.PendingShares > 0 {
						sellOrders.Push(sellOrder)
					}
				}
			}
		case "SELL":
			sellOrders.Push(order)

			if buyOrders.Len() > 0 && buyOrders.Orders[0].Price <= order.Price {
				buyOrder := buyOrders.Pop().(*Order)

				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(buyOrder, order, order.Shares, buyOrder.Price)
					b.AddTransaction(transaction, b.WaitGroup)

					buyOrder.AddTransaction(transaction)
					order.AddTransaction(transaction)

					b.sendToOutputChannel(buyOrder, order)

					if buyOrder.PendingShares > 0 {
						sellOrders.Push(buyOrder)
					}
				}
			}
		}
	}
}

func (b *Book) AddTransaction(transaction *Transaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < minShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellOrderPendingShares(-minShares)

	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.Price)

	transaction.CloseSellingOrder()
	transaction.CloseBuyingOrder()

	b.Transactions = append(b.Transactions, transaction)
}

func (b *Book) sendToOutputChannel(buyOrSellOrder *Order, order *Order) {
	b.OrdersChannelOutput <- buyOrSellOrder
	b.OrdersChannelOutput <- order
}
