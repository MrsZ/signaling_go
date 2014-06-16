package summoners

import (
	"math/rand"
	"time"
)

var randomSource = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

func NewName(n int) (name string) {
	if n > namesVariationWidth {
		panic("Out of my data sample!")
	}
	index := randInt(n*NamesRange, (n+1)*NamesRange)
	return ChampionNames[index]
}

func randInt(min, max int) int {
	if min == max {
		return randInt(min, min+1)
	}
	if min > max {
		return randInt(max, min)
	}
	n := randomSource.Int() % (max - min)
	return n + min
}
