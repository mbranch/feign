package feign

import (
	cryptorand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// Seed uses the provided seed value to initialize the random number generator
// to a deterministic state. Seed should not be called concurrently.
func Seed(seed int64) {
	random = rand.New(&safeSource{
		source: rand.NewSource(seed),
	})
}

func Rand() rand.Source {
	return random
}

// random holds a thread-safe source of random numbers.
var random *rand.Rand

func init() {
	var seed int64
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(math.MaxInt64))
	if err == nil {
		seed = n.Int64()
	} else {
		seed = time.Now().UnixNano()
	}
	random = rand.New(&safeSource{
		source: rand.NewSource(seed),
	})
}

type boundary struct {
	start int
	end   int
}

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "
)

var (
	// stringBoundary is the possible lengths of randomly generated strings.
	stringBoundary = boundary{start: 1, end: 64}

	// intBoundary is the possible range of randomly generated integers.
	intBoundary = boundary{start: 0, end: 65536}

	// sliceAndMapBoundary is the possible lengths of randomly generated slices
	// and maps.
	sliceAndMapBoundary = boundary{start: 1, end: 64}
)

// safeSource holds a thread-safe implementation of rand.Source64.
type safeSource struct {
	source rand.Source
	sync.Mutex
}

func (rs *safeSource) Seed(seed int64) {
	rs.Lock()
	rs.source.Seed(seed)
	rs.Unlock()
}

func (rs *safeSource) Int63() int64 {
	rs.Lock()
	n := rs.source.Int63()
	rs.Unlock()
	return n
}

func randomString() string {
	n := stringBoundary.start + random.Intn(stringBoundary.end-stringBoundary.start)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}

func randomInteger() int {
	return random.Intn(intBoundary.end-intBoundary.start) + intBoundary.start
}

func randomSliceAndMapSize() int {
	return sliceAndMapBoundary.start + random.Intn(sliceAndMapBoundary.end-sliceAndMapBoundary.start)
}
