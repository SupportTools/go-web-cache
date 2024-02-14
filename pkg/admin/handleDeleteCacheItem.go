package admin

import "net/http"

func handleDeleteCacheItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Item key not specified", http.StatusBadRequest)
		return
	}

	cacheManager.Delete(key) // Assuming Delete doesn't return an error; adjust if your implementation differs

	w.WriteHeader(http.StatusOK)
}
