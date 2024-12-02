package types

import "time"

// OrderStatus represents the status and details of an order in a trading system.
// It provides comprehensive information about the order's lifecycle, including
// its creation, execution, and closure.
type OrderStatus struct {
	// Id is the unique identifier for the order. It is used to reference and track
	// specific orders within the trading system.
	Id int `json:"id"`

	// Symbol is the trading pair associated with the order, such as "BTC_USDT".
	Symbol string `json:"symbol"`

	// Type indicates the type of the order, such as "limit" or "market".
	Type string `json:"type"`

	// Side specifies the direction of the order, either "buy" or "sell".
	Side string `json:"side"`

	// BaseAmount represents the amount of the base currency involved in the order.
	// For example, in a BTC_USDT market, this would represent the amount of BTC.
	BaseAmount string `json:"base_amount"`

	// QuoteAmount represents the amount of the quote currency involved in the order.
	// For example, in a BTC_USDT market, this would represent the equivalent amount
	// of USDT for the order.
	QuoteAmount string `json:"quote_amount"`

	// Price is the price at which the order is placed. It is stored as a string
	// to maintain precision for fractional values.
	Price string `json:"price"`

	// StopPrice is the stop price for stop orders. This field is relevant for
	// orders like stop-limit or stop-market orders.
	StopPrice string `json:"stop_price"`

	// OcoTargetPrice is the target price for One-Cancels-the-Other (OCO) orders,
	// used to specify the secondary order's trigger price.
	OcoTargetPrice string `json:"oco_target_price"`

	// Identifier is a unique client-provided identifier for the order, often used
	// for custom tracking or reconciliation.
	Identifier string `json:"identifier"`

	// State indicates the current state of the order, such as "open", "closed",
	// "cancelled", or "pending".
	State string `json:"state"`

	// CreatedAt is the timestamp when the order was created. It is represented as
	// a time.Time object for accurate time management.
	CreatedAt time.Time `json:"created_at"`

	// ClosedAt is the timestamp when the order was closed, if applicable. It is
	// represented as a string to handle cases where the timestamp may not be
	// available or formatted differently.
	ClosedAt string `json:"closed_at"`

	// DealedBaseAmount specifies the amount of the base currency that has been
	// filled (executed) for the order.
	DealedBaseAmount string `json:"dealed_base_amount"`

	// DealedQuoteAmount specifies the amount of the quote currency that has been
	// filled (executed) for the order.
	DealedQuoteAmount string `json:"dealed_quote_amount"`

	// ReqToCancel indicates whether a request to cancel the order has been made.
	// If true, the order is in the process of being cancelled.
	ReqToCancel bool `json:"req_to_cancel"`

	// Commission represents the fee charged for executing the order. It is stored
	// as a string to maintain precision.
	Commission string `json:"commission"`
}

// CreateOrderParams represents the parameters required to create a new order in
// a trading system. It includes details about the trading pair, order type, side,
// and optional attributes for advanced order functionalities.
type CreateOrderParams struct {
	// Symbol is the trading pair for the order, such as "BTCUSDT".
	Symbol string `json:"symbol"`

	// Type specifies the type of order, such as "limit" or "market".
	Type string `json:"type"`

	// Side indicates whether the order is a "buy" or "sell".
	Side string `json:"side"`

	// BaseAmount specifies the amount of the base currency for the order. It is
	// optional and required for certain order types.
	BaseAmount string `json:"base_amount,omitempty"`

	// QuoteAmount specifies the amount of the quote currency for the order. It is
	// optional and required for certain order types.
	QuoteAmount string `json:"quote_amount,omitempty"`

	// Price is the price at which the order is placed. It is optional and required
	// for limit orders.
	Price string `json:"price,omitempty"`

	// StopPrice is the trigger price for stop orders. This field is optional and
	// relevant for stop-limit or stop-market orders.
	StopPrice string `json:"stop_price,omitempty"`

	// OcoTargetPrice is the target price for One-Cancels-the-Other (OCO) orders.
	// This field is optional and used for advanced order strategies.
	OcoTargetPrice string `json:"oco_target_price,omitempty"`

	// Identifier is an optional unique identifier for the order, often used for
	// client-side tracking or reconciliation.
	Identifier string `json:"identifier,omitempty"`
}

// GetOrdersHistoryParams represents the parameters used to fetch a historical
// list of orders. It includes optional filters for narrowing down the results.
type GetOrdersHistoryParams struct {
	// Symbol is the trading pair for the orders, such as "BTCUSDT". This field
	// is optional and can be used to filter orders by trading pair.
	Symbol string `json:"symbol,omitempty"`

	// Side specifies whether to fetch "buy" or "sell" orders. This field is
	// optional and used for filtering.
	Side string `json:"side,omitempty"`

	// State indicates the state of the orders, such as "open", "closed", or
	// "cancelled". This field is optional.
	State string `json:"state,omitempty"`

	// Type specifies the type of the orders, such as "limit" or "market". This
	// field is optional and used for filtering.
	Type string `json:"type,omitempty"`

	// Identifier is an optional unique identifier for filtering orders.
	Identifier string `json:"identifier,omitempty"`

	// Start specifies the start date-time for fetching orders, formatted as a
	// string. This field is optional.
	Start string `json:"start,omitempty"`

	// End specifies the end date-time for fetching orders, formatted as a string.
	// This field is optional.
	End string `json:"end,omitempty"`

	// IdsIn is a comma-separated string of order IDs to fetch. This field is
	// optional and used to specify a list of specific orders.
	IdsIn string `json:"ids_in,omitempty"`

	// IdentifiersIn is a comma-separated string of order identifiers to fetch.
	// This field is optional and used to specify a list of specific orders.
	IdentifiersIn string `json:"identifiers_in,omitempty"`

	// Offset is the starting index for paginated results. This field is optional
	// and used for pagination.
	Offset int `json:"offset,omitempty"`

	// Limit specifies the maximum number of orders to return in the response.
	// This field is optional and used for pagination.
	Limit int `json:"limit,omitempty"`
}

// UserTrade represents an executed trade performed by a user in a trading system.
// It provides detailed information about the trade, including amounts, price, fees,
// and related metadata.
type UserTrade struct {
	// Id is the unique identifier for the trade. It is used to reference specific
	// trades in the system.
	Id int `json:"id"`

	// Symbol is the trading pair for the trade, such as "BTCUSDT".
	Symbol string `json:"symbol"`

	// BaseAmount represents the amount of the base currency involved in the trade.
	// For example, in a BTC/USDT market, this would represent the amount of BTC
	// traded.
	BaseAmount string `json:"base_amount"`

	// QuoteAmount represents the amount of the quote currency involved in the trade.
	// For example, in a BTC/USDT market, this would represent the equivalent amount
	// of USDT traded.
	QuoteAmount string `json:"quote_amount"`

	// Price is the price at which the trade was executed. It is stored as a string
	// to maintain precision.
	Price string `json:"price"`

	// CreatedAt is the timestamp when the trade was executed. It is represented
	// as a time.Time object for accurate time tracking.
	CreatedAt time.Time `json:"created_at"`

	// Commission represents the fee charged for executing the trade. It is stored
	// as a string to maintain precision.
	Commission string `json:"commission"`

	// Side indicates whether the trade was a "buy" or "sell" from the user's
	// perspective.
	Side string `json:"side"`

	// CommissionCurrency specifies the currency in which the commission was charged.
	// For example, "BTC" or "USDT".
	CommissionCurrency string `json:"commission_currency"`

	// OrderId is the unique identifier of the order associated with this trade.
	// This field links the trade to its originating order.
	OrderId int `json:"order_id"`

	// Identifier is an optional client-provided identifier for the trade, often
	// used for tracking or reconciliation.
	Identifier string `json:"identifier"`
}

// GetUserTradesParams represents the parameters used to fetch a list of user trades.
// It includes optional filters for narrowing down the results.
type GetUserTradesParams struct {
	// Symbol is the trading pair for the trades, such as "BTCUSDT". This field is
	// optional and can be used to filter trades by trading pair.
	Symbol string `json:"symbol,omitempty"`

	// Side specifies whether to fetch "buy" or "sell" trades. This field is
	// optional and used for filtering.
	Side string `json:"side,omitempty"`

	// Offset is the starting index for paginated results. This field is optional
	// and used for pagination.
	Offset int `json:"offset,omitempty"`

	// Limit specifies the maximum number of trades to return in the response.
	// This field is optional and used for pagination.
	Limit int `json:"limit,omitempty"`
}

// OrderStatuses represents a collection of OrderStatus objects.
// This type is used to manage and process multiple order statuses, such as
// retrieving a batch of orders or analyzing the state of multiple orders.
type OrderStatuses []OrderStatus

// UserTrades represents a collection of UserTrade objects.
// This type is used to handle multiple user trades, such as retrieving trade
// history, calculating trade statistics, or analyzing trading activity.
type UserTrades []UserTrade
