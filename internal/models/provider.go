package models

type Provider struct {
	ProviderName string `json:"provider_name"`
}

// ProviderInfo is an entry of TMDb's provider catalogue (/watch/providers/movie).
type ProviderInfo struct {
	ID              int    `json:"provider_id"`
	Name            string `json:"provider_name"`
	DisplayPriority int    `json:"display_priority"`
}

type ProviderListResponse struct {
	Results []ProviderInfo `json:"results"`
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