package poller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"tezos-delegation-service/config"
	"tezos-delegation-service/internal/db"
	"tezos-delegation-service/internal/parser"
	"tezos-delegation-service/logger"

	"github.com/patrickmn/go-cache"
)

// Poller is responsible for polling the Tezos API and handling delegations.
type Poller struct {
	cfg        *config.Config // Configuration for the Poller
	querier    db.Querier     // Database querier for interacting with the database
	httpClient *http.Client   // HTTP client for making requests
	cache      *cache.Cache   // Cache to store the last processed delegation ID
	zeroLogger *logger.Logger // Logger for logging messages
}

// TempDelegation represents the structure of a delegation received from the API.
type TempDelegation struct {
	ID        json.RawMessage `json:"id"`
	Delegator json.RawMessage `json:"delegator"`
	Timestamp time.Time       `json:"timestamp"`
	Amount    json.RawMessage `json:"amount"`
	Level     json.RawMessage `json:"level"`
}

// NewPoller creates a new Poller instance with the given configuration and dependencies.
func NewPoller(cfg *config.Config, dbQuerier db.Querier, zeroLog *logger.Logger) *Poller {
	return &Poller{
		cfg:        cfg,
		querier:    dbQuerier,
		httpClient: &http.Client{},
		cache:      cache.New(5*time.Minute, 10*time.Minute),
		zeroLogger: zeroLog,
	}
}

// StartPolling starts the polling process, which runs at intervals defined in the configuration.
func (p *Poller) StartPolling(ctx context.Context) {
	duration, err := time.ParseDuration(p.cfg.PollingInterval)
	if err != nil {
		p.zeroLogger.Log.Error().Msgf("Invalid polling interval duration: %v", err)
		return
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.FetchDelegations(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// FetchDelegations fetches delegations from the API, processes them, and stores them in the database.
func (p *Poller) FetchDelegations(ctx context.Context) {
	resp, err := p.httpClient.Get(p.cfg.TezosAPIURL)
	if err != nil {
		p.zeroLogger.Log.Error().Msgf("Failed to fetch delegations: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		p.zeroLogger.Log.Info().Msgf("Unexpected status code: %d", resp.StatusCode)
		return
	}

	var rawDelegations []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawDelegations); err != nil {
		p.zeroLogger.Log.Error().Msgf("Failed to decode delegations: %v", err)
		return
	}

	if len(rawDelegations) == 0 {
		return
	}

	// Extract the ID of the last delegation in the response
	lastDelegationID := rawDelegations[len(rawDelegations)-1]["id"]
	lastDelegationIDStr, err := json.Marshal(lastDelegationID)
	if err != nil {
		p.zeroLogger.Log.Error().Msgf("Failed to marshal last delegation ID: %v", err)
		return
	}

	// Check if the last delegation ID is different from the cached one
	cachedID, found := p.cache.Get("lastDelegationID")
	if found && string(cachedID.([]byte)) == string(lastDelegationIDStr) {
		p.zeroLogger.Log.Info().Msg("No new delegations found. Skipping database update.")
		return
	}

	// Parse and store new delegations
	var delegations []db.Delegation
	for _, raw := range rawDelegations {
		delegation, err := parser.ParseDelegationParameters(raw)
		if err != nil {
			p.zeroLogger.Log.Error().Msgf("Failed to parse delegation: %v", err)
			return
		}
		delegations = append(delegations, *delegation)
	}

	// Insert the delegations into the database
	err = p.querier.InsertDelegationsBatch(ctx, delegations)
	if err != nil {
		p.zeroLogger.Log.Error().Msgf("Failed to insert delegations batch: %v", err)
		return
	}

	// Update the cache with the new last delegation ID
	p.cache.Set("lastDelegationID", lastDelegationIDStr, cache.DefaultExpiration)
}
