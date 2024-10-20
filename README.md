# TimedMap
`TimedMap` is a thread-safe, generic map that automatically removes expired entries. It is designed to be used for caching scenarios where data is only valid for a limited time and needs to be automatically purged after expiration.

## Features
- **Generic**: Supports any key and value types, as long as the key type is comparable.
- **Automatic Expiration**: Entries are automatically removed once they expire after the specified duration.
- **Background Cleanup**: A cleanup process periodically scans and removes expired entries to ensure efficient memory usage.
- **Thread-Safe**: `TimedMap` uses a `sync.RWMutex` to synchronize access, allowing safe concurrent reads and writes.

## API Documentation
The `TimedMap` library provides the following API:

*   `New[K, V]()` - Creates a new `TimedMap` with the default cleanup interval of 1 minute.
*   `NewWithCleanupInterval[K, V](interval time.Duration)` - Creates a new `TimedMap` with the given cleanup interval.
*   `Put(key K, value V, ttl time.Duration)` - Adds a value and its time-to-live duration to the `TimedMap` for the given key.
*   `Get(key K) (V, bool)` - Returns the value associated with the given key and a boolean indicating if the key exists.
*   `Delete(key K)` - Removes the value associated with the given key regardless of its expiration time.
*   `Clear()` - Removes all entries from the `TimedMap`.
*   `Len() int` - Returns the number of entries in the `TimedMap`.

## Example

```go
package main

import (
    "fmt"
    "time"

    "github.com/mxmlkzdh/timedmap"
)

func main() {

    // Initialize a TimedMap with a cleanup interval of 30 seconds
    tm := timedmap.NewWithCleanupInterval[string, string](30 * time.Second)

    // Add entries
    tm.Put("key1", "value1", 10*time.Second)  // This entry will expire in 10 seconds
    tm.Put("key2", "value2", 1*time.Minute)   // This entry will expire in 1 minute

    // Fetch values
    value, ok := tm.Get("key1")
    if ok {
        fmt.Println("key1:", value)
    } else {
        fmt.Println("key1 expired or does not exist")
    }

    // Delete an entry
    tm.Delete("key2")

    // Get the map length
    fmt.Println("Map size:", tm.Len())

    // Wait for expiration
    time.Sleep(11 * time.Second)

    // Fetch an expired entry
    value, ok = tm.Get("key1")
    if ok {
        fmt.Println("key1:", value)
    } else {
        fmt.Println("key1 expired or does not exist")
    }
}
```
## Installation
To use `TimedMap`, install it using `go get`:
```bash
go get github.com/mxmlkzdh/timedmap
```
Then import it in your Go code:
```go
import "github.com/mxmlkzdh/timedmap"
```

## Thread Safety
`TimedMap` is designed to be used concurrently from multiple goroutines. It uses `sync.RWMutex` to provide safe access for both readers and writers.
- **Concurrent Reads**: Multiple goroutines can safely read from the map simultaneously.
- **Concurrent Writes**: Writes (setting or deleting entries) are synchronized, ensuring that the map's state remains consistent.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
