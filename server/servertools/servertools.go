package servertools

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
)

// Checks if a file exists and isn't a directory.
// see: https://golangcode.com/check-if-a-file-exists/
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !info.IsDir()
}

func IsEntitled(keys []string, endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Api-Key"] != nil {

			key := r.Header["Api-Key"][0]

			if contains(keys, key) {
				endpoint(w, r)
			} else {
				RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}

		} else {
			RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
	})
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
