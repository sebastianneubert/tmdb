package models

type Provider struct {
	ProviderName string `json:"provider_name"`
}

type RegionProviders struct {
	Link     string     `json:"link"`
	Flatrate []Provider `json:"flatrate"`
	Rent     []Provider `json:"rent"`
	Buy      []Provider `json:"buy"`
}

type WatchProviderResponse struct {
	ID      int                        `json:"id"`
	Results map[string]RegionProviders `json:"results"`
}