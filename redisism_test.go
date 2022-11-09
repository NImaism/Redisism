package redisism

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSampleSingleEntry(t *testing.T) {
	cache := InitCache()
	connection := cache.Connect("client 1")
	connection.Set("Nima Gf", "Liana")
	assert.EqualValues(t, "Liana", connection.Get("Nima Gf"))
}

func TestSampleMultiEntry(t *testing.T) {
	cache := InitCache()
	connection := cache.Connect("client 1")
	myMap := map[string]string{
		"Nima":   "Gholami",
		"Liana":  "Rafie",
		"Sobhan": "Sistani",
	}
	for k, v := range myMap {
		connection.Set(k, v)
	}

	for k, v := range myMap {
		actual := connection.Get(k)
		assert.EqualValues(t, v, actual)
	}
}

func TestSampleMultiConn(t *testing.T) {
	cache := InitCache()
	connection1 := cache.Connect("client 1")
	connection2 := cache.Connect("client 2")

	connection1.Set("my key", "my value")

	actual := connection1.Get("my key")
	assert.EqualValues(t, "my value", actual)

	actual = connection2.Get("my key")
	assert.EqualValues(t, "", actual)
}

func TestSampleMultiCache(t *testing.T) {
	cache1 := InitCache()
	connection1 := cache1.Connect("client 1 on cache 1")

	cache2 := InitCache()
	connection2 := cache2.Connect("client 2 on cache 2")

	connection1.Set("a (1)", "b (1)")
	connection2.Set("a (2)", "b (2)")

	actual := connection1.Get("a (1)")
	assert.EqualValues(t, "b (1)", actual)

	actual = connection2.Get("a (1)")
	assert.EqualValues(t, "", actual)

	actual = connection2.Get("a (2)")
	assert.EqualValues(t, "b (2)", actual)

	actual = connection1.Get("a (2)")
	assert.EqualValues(t, "", actual)
}
