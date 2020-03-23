package cache

// NotFound is a constant to signal that an item is not in the cache.
const NotFound = internalError("not found")

// leveraging the same technique as
// https://github.com/go-redis/redis/blob/e269de20cfd1ccf59cbba825cd4e6b60df5cab3a/internal/proto/reader.go#L19-L27.
type internalError string

func (e internalError) Error() string { return string(e) }
