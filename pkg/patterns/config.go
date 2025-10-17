package patterns

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// ConfigManager manages application configuration in etcd
type ConfigManager struct {
	client    *clientv3.Client
	keyPrefix string
	logger    *zap.Logger
	cache     map[string]interface{}
	watchers  map[string][]ConfigWatcher
	mu        sync.RWMutex
}

// ConfigWatcher is called when configuration changes
type ConfigWatcher func(key string, oldValue, newValue interface{})

// NewConfigManager creates a new configuration manager
func NewConfigManager(client *clientv3.Client, keyPrefix string, logger *zap.Logger) *ConfigManager {
	return &ConfigManager{
		client:    client,
		keyPrefix: keyPrefix,
		logger:    logger,
		cache:     make(map[string]interface{}),
		watchers:  make(map[string][]ConfigWatcher),
	}
}

// Set stores a configuration value
func (cm *ConfigManager) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	fullKey := cm.getFullKey(key)
	_, err = cm.client.Put(ctx, fullKey, string(data))
	if err != nil {
		return fmt.Errorf("failed to set config: %w", err)
	}

	cm.mu.Lock()
	cm.cache[key] = value
	cm.mu.Unlock()

	cm.logger.Debug("Configuration set", zap.String("key", key))
	return nil
}

// Get retrieves a configuration value
func (cm *ConfigManager) Get(ctx context.Context, key string, dest interface{}) error {
	// Check cache first
	cm.mu.RLock()
	if cached, exists := cm.cache[key]; exists {
		cm.mu.RUnlock()
		return cm.copyValue(cached, dest)
	}
	cm.mu.RUnlock()

	// Fetch from etcd
	fullKey := cm.getFullKey(key)
	resp, err := cm.client.Get(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return fmt.Errorf("config key not found: %s", key)
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, dest); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Update cache
	cm.mu.Lock()
	cm.cache[key] = dest
	cm.mu.Unlock()

	return nil
}

// GetString retrieves a string configuration value
func (cm *ConfigManager) GetString(ctx context.Context, key string) (string, error) {
	var value string
	if err := cm.Get(ctx, key, &value); err != nil {
		return "", err
	}
	return value, nil
}

// GetInt retrieves an integer configuration value
func (cm *ConfigManager) GetInt(ctx context.Context, key string) (int, error) {
	var value int
	if err := cm.Get(ctx, key, &value); err != nil {
		return 0, err
	}
	return value, nil
}

// GetBool retrieves a boolean configuration value
func (cm *ConfigManager) GetBool(ctx context.Context, key string) (bool, error) {
	var value bool
	if err := cm.Get(ctx, key, &value); err != nil {
		return false, err
	}
	return value, nil
}

// Delete removes a configuration value
func (cm *ConfigManager) Delete(ctx context.Context, key string) error {
	fullKey := cm.getFullKey(key)
	_, err := cm.client.Delete(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("failed to delete config: %w", err)
	}

	cm.mu.Lock()
	delete(cm.cache, key)
	cm.mu.Unlock()

	cm.logger.Debug("Configuration deleted", zap.String("key", key))
	return nil
}

// GetAll retrieves all configuration values under a prefix
func (cm *ConfigManager) GetAll(ctx context.Context, prefix string) (map[string]interface{}, error) {
	fullPrefix := cm.getFullKey(prefix)
	resp, err := cm.client.Get(ctx, fullPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get configs: %w", err)
	}

	configs := make(map[string]interface{})
	for _, kv := range resp.Kvs {
		var value interface{}
		if err := json.Unmarshal(kv.Value, &value); err != nil {
			cm.logger.Warn("Failed to unmarshal config", zap.String("key", string(kv.Key)), zap.Error(err))
			continue
		}
		// Strip prefix from key
		relativeKey := string(kv.Key)[len(cm.keyPrefix)+1:]
		configs[relativeKey] = value
	}

	return configs, nil
}

// Watch watches for changes to a configuration key
func (cm *ConfigManager) Watch(ctx context.Context, key string, watcher ConfigWatcher) {
	cm.mu.Lock()
	cm.watchers[key] = append(cm.watchers[key], watcher)
	isFirst := len(cm.watchers[key]) == 1
	cm.mu.Unlock()

	// Start watching if this is the first watcher
	if isFirst {
		go cm.watchKey(ctx, key)
	}
}

// watchKey watches for changes to a specific key
func (cm *ConfigManager) watchKey(ctx context.Context, key string) {
	fullKey := cm.getFullKey(key)
	watchChan := cm.client.Watch(ctx, fullKey)

	for {
		select {
		case <-ctx.Done():
			return
		case watchResp := <-watchChan:
			if watchResp.Err() != nil {
				cm.logger.Error("Watch error", zap.Error(watchResp.Err()))
				continue
			}

			for _, event := range watchResp.Events {
				cm.handleConfigEvent(key, event)
			}
		}
	}
}

// handleConfigEvent processes a configuration change event
func (cm *ConfigManager) handleConfigEvent(key string, event *clientv3.Event) {
	var oldValue, newValue interface{}

	// Get old value from cache
	cm.mu.RLock()
	oldValue = cm.cache[key]
	cm.mu.RUnlock()

	// Parse new value
	if event.Type == clientv3.EventTypePut {
		if err := json.Unmarshal(event.Kv.Value, &newValue); err != nil {
			cm.logger.Warn("Failed to unmarshal new value", zap.Error(err))
			return
		}

		// Update cache
		cm.mu.Lock()
		cm.cache[key] = newValue
		cm.mu.Unlock()
	} else if event.Type == clientv3.EventTypeDelete {
		// Remove from cache
		cm.mu.Lock()
		delete(cm.cache, key)
		cm.mu.Unlock()
		newValue = nil
	}

	// Notify watchers
	cm.mu.RLock()
	watchers := cm.watchers[key]
	cm.mu.RUnlock()

	for _, watcher := range watchers {
		go watcher(key, oldValue, newValue)
	}

	cm.logger.Debug("Config changed",
		zap.String("key", key),
		zap.String("type", event.Type.String()))
}

// LoadConfig loads all configuration into cache
func (cm *ConfigManager) LoadConfig(ctx context.Context) error {
	configs, err := cm.GetAll(ctx, "")
	if err != nil {
		return err
	}

	cm.mu.Lock()
	cm.cache = configs
	cm.mu.Unlock()

	cm.logger.Info("Configuration loaded", zap.Int("count", len(configs)))
	return nil
}

// getFullKey returns the full etcd key
func (cm *ConfigManager) getFullKey(key string) string {
	if key == "" {
		return cm.keyPrefix
	}
	return path.Join(cm.keyPrefix, key)
}

// copyValue copies a value to destination
func (cm *ConfigManager) copyValue(src, dest interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// FeatureFlag manages feature flags
type FeatureFlag struct {
	configManager *ConfigManager
	flagPrefix    string
}

// NewFeatureFlag creates a new feature flag manager
func NewFeatureFlag(configManager *ConfigManager) *FeatureFlag {
	return &FeatureFlag{
		configManager: configManager,
		flagPrefix:    "features",
	}
}

// IsEnabled checks if a feature flag is enabled
func (ff *FeatureFlag) IsEnabled(ctx context.Context, featureName string) (bool, error) {
	key := path.Join(ff.flagPrefix, featureName)
	enabled, err := ff.configManager.GetBool(ctx, key)
	if err != nil {
		// Default to false if not found
		return false, nil
	}
	return enabled, nil
}

// Enable enables a feature flag
func (ff *FeatureFlag) Enable(ctx context.Context, featureName string) error {
	key := path.Join(ff.flagPrefix, featureName)
	return ff.configManager.Set(ctx, key, true)
}

// Disable disables a feature flag
func (ff *FeatureFlag) Disable(ctx context.Context, featureName string) error {
	key := path.Join(ff.flagPrefix, featureName)
	return ff.configManager.Set(ctx, key, false)
}

// GetAllFlags returns all feature flags
func (ff *FeatureFlag) GetAllFlags(ctx context.Context) (map[string]bool, error) {
	configs, err := ff.configManager.GetAll(ctx, ff.flagPrefix)
	if err != nil {
		return nil, err
	}

	flags := make(map[string]bool)
	for key, value := range configs {
		if boolVal, ok := value.(bool); ok {
			flags[key] = boolVal
		}
	}

	return flags, nil
}

// ConfigSnapshot represents a point-in-time configuration snapshot
type ConfigSnapshot struct {
	Timestamp int64                  `json:"timestamp"`
	Version   int64                  `json:"version"`
	Config    map[string]interface{} `json:"config"`
}

// SnapshotManager manages configuration snapshots
type SnapshotManager struct {
	configManager *ConfigManager
	snapshotKey   string
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager(configManager *ConfigManager) *SnapshotManager {
	return &SnapshotManager{
		configManager: configManager,
		snapshotKey:   "snapshots",
	}
}

// CreateSnapshot creates a configuration snapshot
func (sm *SnapshotManager) CreateSnapshot(ctx context.Context, version int64) (*ConfigSnapshot, error) {
	configs, err := sm.configManager.GetAll(ctx, "")
	if err != nil {
		return nil, err
	}

	snapshot := &ConfigSnapshot{
		Timestamp: time.Now().Unix(), // Use current timestamp
		Version:   version,
		Config:    configs,
	}

	// Store snapshot
	key := path.Join(sm.snapshotKey, fmt.Sprintf("v%d", version))
	if err := sm.configManager.Set(ctx, key, snapshot); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// RestoreSnapshot restores configuration from a snapshot
func (sm *SnapshotManager) RestoreSnapshot(ctx context.Context, version int64) error {
	key := path.Join(sm.snapshotKey, fmt.Sprintf("v%d", version))

	var snapshot ConfigSnapshot
	if err := sm.configManager.Get(ctx, key, &snapshot); err != nil {
		return err
	}

	// Restore each configuration value
	for key, value := range snapshot.Config {
		if err := sm.configManager.Set(ctx, key, value); err != nil {
			return fmt.Errorf("failed to restore config %s: %w", key, err)
		}
	}

	return nil
}
