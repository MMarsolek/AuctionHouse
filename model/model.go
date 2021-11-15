package model

type User struct {
	Username       string
	DisplayName    string
	HashedPassword string
	Permission     PermissionLevel
}

type AuctionItem struct {
	Name        string
	ImageRef    string
	Description string
}

type AuctionBid struct {
	BidAmount int
	Bidder    *User
	Item      *AuctionItem
}
