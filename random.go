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

var defaultTimeFn = func() time.Time {
	return time.Time{}.Add(time.Duration(random.Int63())).UTC()
}

var timeFn = defaultTimeFn

// TimeFn uses the provided function to generate random time values globally.
// This is often relevant when interacting with databases, which may store lower
// resolution timestamps than possible with time.Time, or when time values need
// to be within a certain date range.
func TimeFn(fn func() time.Time) {
	if fn == nil {
		timeFn = defaultTimeFn
	} else {
		timeFn = fn
	}
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
	stringBoundary = boundary{start: 1, end: 32}

	// intBoundary is the possible range of randomly generated integers.
	intBoundary = boundary{start: 0, end: 65536}

	// sliceAndMapBoundary is the possible lengths of randomly generated slices
	// and maps.
	sliceAndMapBoundary = boundary{start: 1, end: 8}
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
