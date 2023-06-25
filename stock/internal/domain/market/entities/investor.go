package entities

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}

type Investor struct {
	ID            string                   `json:"id"`
	AssetPosition []*InvestorAssetPosition `json:"assetPosition"`
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition)
}

func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}

	return nil
}

func (i *Investor) UpdateAssetPosition(assetID string, sharesQuantity int) {
	assetPosition := i.GetAssetPosition(assetID)

	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, sharesQuantity))
	} else {
		assetPosition.Shares += sharesQuantity
	}
}
