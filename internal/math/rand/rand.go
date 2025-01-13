package rand

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"sync"
)

var gRand = rand.New(new(cryptoSource))

var _ rand.Source64 = (*cryptoSource)(nil)

type cryptoSource struct {
	lk sync.Mutex
}

func (s *cryptoSource) Seed(_ int64) {}

func (s *cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s *cryptoSource) Uint64() (v uint64) {
	s.lk.Lock()
	defer s.lk.Unlock()
	n, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		panic(err)
	}
	return n.Uint64()
}

// Seed uses the provided seed value to initialize the default Source to a
// deterministic state. If Seed is not called, the generator behaves as
// if seeded by Seed(1). Seed values that have the same remainder when
// divided by 2³¹-1 generate the same pseudo-random sequence.
// Seed, unlike the Rand.Seed method, is safe for concurrent use.
func Seed(seed int64) { gRand.Seed(seed) }

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64
// from the default Source.
func Int63() int64 { return gRand.Int63() }

// Uint32 returns a pseudo-random 32-bit value as a uint32
// from the default Source.
func Uint32() uint32 { return gRand.Uint32() }

// Uint64 returns a pseudo-random 64-bit value as a uint64
// from the default Source.
func Uint64() uint64 { return gRand.Uint64() }

// Int31 returns a non-negative pseudo-random 31-bit integer as an int32
// from the default Source.
func Int31() int32 { return gRand.Int31() }

// Int returns a non-negative pseudo-random int from the default Source.
func Int() int { return gRand.Int() }

// Int63n returns, as an int64, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Int63n(n int64) int64 { return gRand.Int63n(n) }

// Int31n returns, as an int32, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Int31n(n int32) int32 { return gRand.Int31n(n) }

// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
// from the default Source.
// It panics if n <= 0.
func Intn(n int) int { return gRand.Intn(n) }

// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0)
// from the default Source.
func Float64() float64 { return gRand.Float64() }

// Float32 returns, as a float32, a pseudo-random number in [0.0,1.0)
// from the default Source.
func Float32() float32 { return gRand.Float32() }

// Perm returns, as a slice of n ints, a pseudo-random permutation of the integers [0,n)
// from the default Source.
func Perm(n int) []int { return gRand.Perm(n) }

// Shuffle pseudo-randomizes the order of elements using the default Source.
// n is the number of elements. Shuffle panics if n < 0.
// swap swaps the elements with indexes i and j.
func Shuffle(n int, swap func(i, j int)) { gRand.Shuffle(n, swap) }

// Read generates len(p) random bytes from the default Source and
// writes them into p. It always returns len(p) and a nil error.
// Read, unlike the Rand.Read method, is safe for concurrent use.
func Read(p []byte) (n int, err error) { return gRand.Read(p) }

// NormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// standard normal distribution (mean = 0, stddev = 1)
// from the default Source.
// To produce a different normal distribution, callers can
// adjust the output using:
//
// sample = NormFloat64() * desiredStdDev + desiredMean
func NormFloat64() float64 { return gRand.NormFloat64() }

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
// sample = ExpFloat64() / desiredRateParameter
func ExpFloat64() float64 { return gRand.ExpFloat64() }

// ByWeights returns a random index into the array chosen by the relative
func ByWeights(weights []int64) int64 {
	var total int64
	for _, w := range weights {
		total += w
	}
	if total == 0 {
		return -1
	}
	r := gRand.Int63n(total)
	for i, w := range weights {
		if r < w {
			return int64(i)
		}
		r -= w
	}
	return r
}
