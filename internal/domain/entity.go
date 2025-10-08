package domain

type Order struct {
	OrderUID          string
	TrackNumber       string
	Entry             string
	Delivery          DeliveryInfo
	Payment           PaymentInfo
	Items             []ItemInfo
	Locale            string
	InternalSignature string
	CustomerID        string
	DeliveryService   string
	ShardKey          string
	SMID              int
	DateCreated       string
	OOFShard          string
}

type DeliveryInfo struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type PaymentInfo struct {
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       int
	PaymentDT    int
	Bank         string
	DeliveryCost int
	GoodsTotal   int
	CustomFee    int
}

type ItemInfo struct {
	OrderUID    string
	ChrtID      int
	TrackNumber string
	Price       int
	RID         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int
	NmID        int
	Brand       string
	Status      int
}
