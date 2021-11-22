package model

// User defines how to identify and name someone who can create bids or manage the system.
type User struct {
	Username       string
	DisplayName    string
	HashedPassword string
	Permission     PermissionLevel
}

// AuctionItem defines the item that is being auctioned off.
type AuctionItem struct {
	Name        string
	ImageRef    string
	Description string
}

// AuctionBid creates the link between the user and the item and how much was being bid.
type AuctionBid struct {
	BidAmount int
	Bidder    *User
	Item      *AuctionItem
}
