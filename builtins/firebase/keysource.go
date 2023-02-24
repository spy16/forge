package firebase

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pquerna/cachecontrol/cacheobject"

	"github.com/spy16/forge/core/errors"
)

const syncURL = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

type KeySource struct {
	mu       sync.RWMutex
	keys     map[string]string
	nextSync time.Time
}

func (ks *KeySource) Find(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	if err := ks.syncIfNeeded(); err != nil {
		return nil, err
	}

	ks.mu.RLock()
	defer ks.mu.RUnlock()
	keyStr, found := ks.keys[kid]
	if !found {
		return nil, errors.NotFound
	}

	return jwt.ParseRSAPublicKeyFromPEM([]byte(keyStr))
}

func (ks *KeySource) syncIfNeeded() error {
	ks.mu.RLock()
	syncNeeded := ks.nextSync.Before(time.Now())
	ks.mu.RUnlock()

	if !syncNeeded {
		return nil
	}

	ks.mu.Lock()
	defer ks.mu.Unlock()

	resp, err := http.Get(syncURL)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != 200 {
		return errors.InternalIssue.
			CausedBy(fmt.Errorf("unexpected status: %s", resp.Status))
	}

	dirs, err := cacheobject.ParseResponseCacheControl(resp.Header.Get("Cache-Control"))
	if err != nil {
		return err
	}
	ks.nextSync = time.Now().Add(time.Duration(dirs.MaxAge) * time.Second)
	return json.NewDecoder(resp.Body).Decode(&ks.keys)
}
