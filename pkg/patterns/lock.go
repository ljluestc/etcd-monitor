// Package patterns provides common etcd patterns and implementations
package patterns

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

// DistributedLock provides distributed mutual exclusion using etcd
type DistributedLock struct {
	client    *clientv3.Client
	session   *concurrency.Session
	mutex     *concurrency.Mutex
	lockKey   string
	ttl       int
	logger    *zap.Logger
	mu        sync.Mutex
	isLocked  bool
}

// LockConfig holds configuration for distributed lock
type LockConfig struct {
	LockKey   string        // Key to lock on
	TTL       int           // Time-to-live in seconds
	Timeout   time.Duration // Lock acquisition timeout
}

// NewDistributedLock creates a new distributed lock
func NewDistributedLock(client *clientv3.Client, config LockConfig, logger *zap.Logger) (*DistributedLock, error) {
	if config.TTL == 0 {
		config.TTL = 60 // Default 60 seconds
	}

	// Create a session with TTL
	session, err := concurrency.NewSession(client, concurrency.WithTTL(config.TTL))
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Create a mutex on the lock key
	mutex := concurrency.NewMutex(session, config.LockKey)

	return &DistributedLock{
		client:   client,
		session:  session,
		mutex:    mutex,
		lockKey:  config.LockKey,
		ttl:      config.TTL,
		logger:   logger,
		isLocked: false,
	}, nil
}

// Lock acquires the distributed lock
func (dl *DistributedLock) Lock(ctx context.Context) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if dl.isLocked {
		return errors.New("already locked")
	}

	dl.logger.Info("Attempting to acquire lock", zap.String("key", dl.lockKey))

	if err := dl.mutex.Lock(ctx); err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	dl.isLocked = true
	dl.logger.Info("Lock acquired", zap.String("key", dl.lockKey))

	return nil
}

// TryLock attempts to acquire the lock without blocking
func (dl *DistributedLock) TryLock(ctx context.Context) (bool, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if dl.isLocked {
		return false, errors.New("already locked")
	}

	// Try to acquire the lock with immediate timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err := dl.mutex.Lock(timeoutCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return false, nil // Lock not available
		}
		return false, fmt.Errorf("failed to try lock: %w", err)
	}

	dl.isLocked = true
	dl.logger.Info("Lock acquired (try)", zap.String("key", dl.lockKey))

	return true, nil
}

// Unlock releases the distributed lock
func (dl *DistributedLock) Unlock(ctx context.Context) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if !dl.isLocked {
		return errors.New("not locked")
	}

	if err := dl.mutex.Unlock(ctx); err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	dl.isLocked = false
	dl.logger.Info("Lock released", zap.String("key", dl.lockKey))

	return nil
}

// IsLocked returns whether the lock is currently held
func (dl *DistributedLock) IsLocked() bool {
	dl.mu.Lock()
	defer dl.mu.Unlock()
	return dl.isLocked
}

// Close closes the session and releases the lock
func (dl *DistributedLock) Close() error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if dl.isLocked {
		if err := dl.mutex.Unlock(context.Background()); err != nil {
			dl.logger.Error("Failed to unlock during close", zap.Error(err))
		}
		dl.isLocked = false
	}

	return dl.session.Close()
}

// WithLock executes a function while holding the lock
func (dl *DistributedLock) WithLock(ctx context.Context, fn func() error) error {
	if err := dl.Lock(ctx); err != nil {
		return err
	}
	defer dl.Unlock(ctx)

	return fn()
}

// LockManager manages multiple distributed locks
type LockManager struct {
	client *clientv3.Client
	locks  map[string]*DistributedLock
	mu     sync.RWMutex
	logger *zap.Logger
}

// NewLockManager creates a new lock manager
func NewLockManager(client *clientv3.Client, logger *zap.Logger) *LockManager {
	return &LockManager{
		client: client,
		locks:  make(map[string]*DistributedLock),
		logger: logger,
	}
}

// AcquireLock acquires a lock with the given key
func (lm *LockManager) AcquireLock(ctx context.Context, key string, ttl int) (*DistributedLock, error) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	// Check if lock already exists
	if lock, exists := lm.locks[key]; exists {
		if lock.IsLocked() {
			return nil, fmt.Errorf("lock %s already acquired", key)
		}
	}

	// Create new lock
	lock, err := NewDistributedLock(lm.client, LockConfig{
		LockKey: key,
		TTL:     ttl,
	}, lm.logger)
	if err != nil {
		return nil, err
	}

	// Acquire the lock
	if err := lock.Lock(ctx); err != nil {
		lock.Close()
		return nil, err
	}

	lm.locks[key] = lock
	return lock, nil
}

// ReleaseLock releases a lock with the given key
func (lm *LockManager) ReleaseLock(ctx context.Context, key string) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lock, exists := lm.locks[key]
	if !exists {
		return fmt.Errorf("lock %s not found", key)
	}

	if err := lock.Unlock(ctx); err != nil {
		return err
	}

	if err := lock.Close(); err != nil {
		lm.logger.Error("Failed to close lock", zap.Error(err))
	}

	delete(lm.locks, key)
	return nil
}

// ReleaseAll releases all locks
func (lm *LockManager) ReleaseAll(ctx context.Context) error {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	var lastErr error
	for key, lock := range lm.locks {
		if err := lock.Unlock(ctx); err != nil {
			lm.logger.Error("Failed to unlock", zap.String("key", key), zap.Error(err))
			lastErr = err
		}
		if err := lock.Close(); err != nil {
			lm.logger.Error("Failed to close lock", zap.String("key", key), zap.Error(err))
		}
	}

	lm.locks = make(map[string]*DistributedLock)
	return lastErr
}

// ReadWriteLock provides distributed read-write lock
type ReadWriteLock struct {
	client      *clientv3.Client
	keyPrefix   string
	writeKey    string
	readKey     string
	ttl         int
	logger      *zap.Logger
	writeLock   *DistributedLock
	mu          sync.Mutex
	readerCount int
}

// NewReadWriteLock creates a new distributed read-write lock
func NewReadWriteLock(client *clientv3.Client, keyPrefix string, ttl int, logger *zap.Logger) (*ReadWriteLock, error) {
	return &ReadWriteLock{
		client:      client,
		keyPrefix:   keyPrefix,
		writeKey:    keyPrefix + "/write",
		readKey:     keyPrefix + "/read",
		ttl:         ttl,
		logger:      logger,
		readerCount: 0,
	}, nil
}

// RLock acquires a read lock
func (rwl *ReadWriteLock) RLock(ctx context.Context) error {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	// First reader acquires the write lock to block writers
	if rwl.readerCount == 0 {
		lock, err := NewDistributedLock(rwl.client, LockConfig{
			LockKey: rwl.writeKey,
			TTL:     rwl.ttl,
		}, rwl.logger)
		if err != nil {
			return err
		}
		if err := lock.Lock(ctx); err != nil {
			return err
		}
		rwl.writeLock = lock
	}

	rwl.readerCount++
	rwl.logger.Debug("Read lock acquired", zap.Int("reader_count", rwl.readerCount))
	return nil
}

// RUnlock releases a read lock
func (rwl *ReadWriteLock) RUnlock(ctx context.Context) error {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	if rwl.readerCount == 0 {
		return errors.New("not read locked")
	}

	rwl.readerCount--

	// Last reader releases the write lock
	if rwl.readerCount == 0 && rwl.writeLock != nil {
		if err := rwl.writeLock.Unlock(ctx); err != nil {
			return err
		}
		rwl.writeLock.Close()
		rwl.writeLock = nil
	}

	rwl.logger.Debug("Read lock released", zap.Int("reader_count", rwl.readerCount))
	return nil
}

// Lock acquires a write lock
func (rwl *ReadWriteLock) Lock(ctx context.Context) error {
	lock, err := NewDistributedLock(rwl.client, LockConfig{
		LockKey: rwl.writeKey,
		TTL:     rwl.ttl,
	}, rwl.logger)
	if err != nil {
		return err
	}

	if err := lock.Lock(ctx); err != nil {
		return err
	}

	rwl.mu.Lock()
	rwl.writeLock = lock
	rwl.mu.Unlock()

	rwl.logger.Debug("Write lock acquired")
	return nil
}

// Unlock releases a write lock
func (rwl *ReadWriteLock) Unlock(ctx context.Context) error {
	rwl.mu.Lock()
	defer rwl.mu.Unlock()

	if rwl.writeLock == nil {
		return errors.New("not write locked")
	}

	if err := rwl.writeLock.Unlock(ctx); err != nil {
		return err
	}

	rwl.writeLock.Close()
	rwl.writeLock = nil

	rwl.logger.Debug("Write lock released")
	return nil
}
