package types

// Currency represents a cryptocurrency or fiat currency with its attributes.
// This struct is typically used to model currencies in trading systems or
// exchanges.
type Currency struct {
	// Currency is the unique identifier for the currency, typically represented
	// by a standardized code such as "USD" for US Dollar or "BTC" for Bitcoin.
	Currency string `json:"currency"`

	// Name provides the full name of the currency, such as "United States Dollar"
	// or "Bitcoin".
	Name string `json:"name"`

	// Tradable indicates whether the currency can be actively traded on the
	// platform. If true, the currency is available for trading; otherwise, it is not.
	Tradable bool `json:"tradable"`

	// Precision specifies the level of precision supported for this currency,
	// usually represented as the number of decimal places allowed in transactions.
	// For example, a precision of "8" for Bitcoin allows values like 0.00000001 BTC.
	Precision string `json:"precision"`
}

// Market represents a trading market on an exchange, characterized by its base
// and quote assets, trading precision, and other attributes.
type Market struct {
	// Symbol is the unique identifier for the market, often represented as a
	// combination of the base and quote assets, such as "BTCUSDT".
	Symbol string `json:"symbol"`

	// Name provides a human-readable name for the market, such as
	// "Bitcoin/US Dollar".
	Name string `json:"name"`

	// Base represents the base asset of the market. The base asset is the
	// currency being bought or sold, such as "BTC" in the "BTC/USDT" market.
	Base string `json:"base"`

	// Quote represents the quote asset of the market. The quote asset is the
	// currency used to price the base asset, such as "USDT" in the "BTC/USDT"
	// market.
	Quote string `json:"quote"`

	// Tradable indicates whether the market is currently active and available
	// for trading. If true, the market can be traded; otherwise, it is not.
	Tradable bool `json:"tradable"`

	// PricePrecision specifies the number of decimal places allowed for the
	// price of the base asset in this market. For example, a precision of 2
	// means prices like 123.45 are valid.
	PricePrecision int `json:"price_precision"`

	// BaseAmountPrecision defines the number of decimal places allowed for the
	// amount of the base asset in transactions. For instance, a precision of 8
	// allows values like 0.12345678 BTC.
	BaseAmountPrecision int `json:"base_amount_precision"`

	// QuoteAmountPrecision defines the number of decimal places allowed for the
	// amount of the quote asset in transactions. For example, a precision of 2
	// allows values like 123.45 USDT.
	QuoteAmountPrecision int `json:"quote_amount_precision"`
}

// Ticker represents real-time market data for a specific trading symbol,
// including its current price, daily price changes, and other related statistics.
type Ticker struct {
	// Symbol is the unique identifier for the trading pair or market, such as
	// "BTCUSDT".
	Symbol string `json:"symbol"`

	// Price represents the current market price of the symbol as a string to
	// maintain precision in cases where high precision is required.
	Price string `json:"price"`

	// DailyChangePrice indicates the price change over the past 24 hours. It is
	// represented as a float64 to allow accurate computations and comparisons.
	DailyChangePrice float64 `json:"daily_change_price"`

	// Low represents the lowest price for the symbol in the past 24 hours, stored
	// as a string for precision.
	Low string `json:"low"`

	// High represents the highest price for the symbol in the past 24 hours,
	// stored as a string for precision.
	High string `json:"high"`

	// Timestamp provides the Unix timestamp (in seconds) when this ticker data
	// was last updated. This allows synchronization with real-time data feeds.
	Timestamp float64 `json:"timestamp"`
}

// OrderBook represents the state of an order book for a specific trading market,
// including the current asks (sell orders) and bids (buy orders).
type OrderBook struct {
	// Asks is a list of sell orders currently available in the market.
	// Each entry is represented as a slice of strings, where:
	// - The first element is the price level of the ask (as a string for precision).
	// - The second element is the quantity available at that price level (as a string for precision).
	// For example, an ask entry could look like ["40000.00", "0.5"], indicating an
	// ask price of 40000.00 and a quantity of 0.5 units of the base asset.
	Asks [][]string `json:"asks"`

	// Bids is a list of buy orders currently available in the market.
	// Each entry is represented as a slice of strings, where:
	// - The first element is the price level of the bid (as a string for precision).
	// - The second element is the quantity available at that price level (as a string for precision).
	// For example, a bid entry could look like ["39000.00", "1.0"], indicating a
	// bid price of 39000.00 and a quantity of 1.0 units of the base asset.
	Bids [][]string `json:"bids"`
}

// Trade represents an individual trade in a trading market, detailing the
// price, amounts, and direction of the trade.
type Trade struct {
	// Id is the unique identifier for the trade, often used to track or reference
	// specific trades in a trading platform or exchange.
	Id string `json:"id"`

	// Price represents the price at which the trade was executed. It is stored
	// as a string to maintain precision for high-value or fractional trades.
	Price string `json:"price"`

	// BaseAmount specifies the amount of the base currency involved in the trade.
	// For example, in a BTC/USDT market, this would represent the amount of BTC
	// traded. It is stored as a string to ensure precision.
	BaseAmount string `json:"base_amount"`

	// QuoteAmount specifies the amount of the quote currency involved in the trade.
	// For example, in a BTC/USDT market, this would represent the equivalent
	// amount of USDT for the trade. It is stored as a string for precision.
	QuoteAmount string `json:"quote_amount"`

	// Side indicates the direction of the trade, either "buy" or "sell", from the
	// perspective of the taker (the trader who initiated the market order).
	Side string `json:"side"`
}

// Currencies represents a collection of Currency objects.
// This type is used to manage and process multiple currencies, such as
// retrieving lists of supported currencies or performing batch operations.
type Currencies []Currency

// Markets represents a collection of Market objects.
// This type is used to manage and process multiple markets, such as
// retrieving available trading pairs or performing market-wide analyses.
type Markets []Market

// Tickers represents a collection of Ticker objects.
// This type is used to handle multiple real-time market data updates, such as
// monitoring price changes or calculating aggregate statistics.
type Tickers []Ticker

// Trades represents a collection of Trade objects.
// This type is used to manage multiple trade records, such as storing trade
// histories, calculating statistics, or analyzing trading activity.
type Trades []Trade
