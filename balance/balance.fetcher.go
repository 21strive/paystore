package balance

type Fetcher struct {
}

func (bf *Fetcher) FetchByUUID(uuid string) {}

func (bf *Fetcher) IsBlankByUUID(uuid string) bool {
	return false
}

func (bf *Fetcher) FetchByExternalID(externalID string) {}

func (bf *Fetcher) IsBlankByExternalID(externalID string) bool {
	return false
}

func (bf *Fetcher) FetchPartial(lastRandId []string) ([]Account, error) {
	return nil, nil
}

func NewFetcher() *Fetcher {
	return &Fetcher{}
}
