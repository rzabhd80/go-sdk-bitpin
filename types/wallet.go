package types

// Wallet represents a user's wallet in the system, detailing the balance,
// frozen funds, and associated service for a specific asset.
type Wallet struct {
	// Id is the unique identifier for the wallet. It is used to reference and
	// manage specific wallets in the system.
	Id int `json:"id"`

	// Asset represents the asset type in the wallet, such as "BTC", "ETH", or "USDT".
	Asset string `json:"asset"`

	// Balance is the total available balance for the asset in the wallet. It is
	// stored as a string to maintain precision for fractional amounts.
	Balance string `json:"balance"`

	// Frozen represents the amount of the asset that is frozen and not currently
	// available for trading or withdrawal. It is stored as a string for precision.
	Frozen string `json:"frozen"`

	// Service indicates the service or platform associated with the wallet,
	// such as "spot", "margin", or "futures".
	Service string `json:"service"`
}

// GetWalletParams represents the parameters used to fetch a list of wallets.
// It includes optional filters for narrowing down the results.
type GetWalletParams struct {
	// Assets is a list of asset symbols to filter the wallets, such as ["BTC", "ETH"].
	// This field is optional and allows fetching wallets for specific assets.
	Assets []string `json:"assets,omitempty"`

	// Service specifies the service or platform to filter the wallets, such as
	// "spot", "margin", or "futures". This field is optional.
	Service string `json:"service,omitempty"`

	// Offset is the starting index for paginated results. This field is optional
	// and used for pagination.
	Offset int `json:"offset,omitempty"`

	// Limit specifies the maximum number of wallets to return in the response.
	// This field is optional and used for pagination.
	Limit int `json:"limit,omitempty"`
}

// Wallets represents a collection of Wallet objects.
// This type is used to manage and process multiple wallets, such as retrieving
// wallet details for all assets, performing batch operations, or analyzing balances.
type Wallets []Wallet
