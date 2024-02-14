package admin

import (
	"encoding/json"
	"net/http"

	"github.com/supporttools/go-web-cache/pkg/cache"
)

func handleListCacheItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	keys := cacheManager.List()
	items := make([]cache.CacheItem, 0, len(keys))
	for _, key := range keys {
		if item, found := cacheManager.Read(key); found {
			items = append(items, item)
		}
	}

	jsonData, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Failed to serialize cache items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
