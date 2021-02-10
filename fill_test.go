package feign

import (
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
)

type Romeo interface {
	Sierra() string
	Tango() int
}

type (
	Baker   map[int]map[Fisher][]Canteen
	Canteen int
	Fisher  uint
)

func TestFill(t *testing.T) {
	var p struct {
		Alpha   int
		Beta    string
		Charlie *string
		Delta   **string
		Echo    struct {
			Foxtrot *bool
			Golf    uint16
			Hotel   []string
			Indigo  *struct {
				Juliet []float32
				Kilo   time.Time
				Lima   map[int]string
				Mike   map[string][]int
			}
		}
		November interface{}
		Oscar    interface{}
		Papa     interface{}
		Quebec   Romeo
		Uniform  map[interface{}]int
		Victor   []interface{}
		Whiskey  map[string]interface{}
		Yankee   chan int
		Zulu     func(string) bool

		Actor Baker
		Diver map[string]Canteen
		Eagle map[Canteen]float64
	}
	assert.ErrorContains(t, Fill(p), "not a pointer value")
	assert.Must(t, Fill(&p, func(path string) (interface{}, bool) {
		switch path {
		case ".Charlie":
			return nil, true
		case ".Oscar":
			return nil, false
		case ".Papa":
			return 1, true
		case ".Zulu":
			return func(s string) bool {
				return s == ""
			}, true
		default:
			return nil, false
		}
	}))
	assert.True(t, p.Zulu(""))
	assert.True(t, p.Alpha != 0)
	assert.True(t, p.Beta != "")
	assert.Nil(t, p.Charlie)
	assert.NotNil(t, *p.Delta)
	assert.NotNil(t, p.Echo.Foxtrot)
	assert.NotNil(t, p.Echo.Golf != 0)
	assert.True(t, len(p.Echo.Hotel) > 0)
	assert.NotNil(t, p.Echo.Indigo)
	assert.NotNil(t, p.Echo.Indigo.Juliet)
	assert.False(t, p.Echo.Indigo.Kilo.IsZero())
	assert.True(t, len(p.Echo.Indigo.Lima) > 0)
	assert.True(t, len(p.Echo.Indigo.Mike) > 0)
	assert.Nil(t, p.November)
	for k := range p.Echo.Indigo.Mike {
		assert.True(t, len(p.Echo.Indigo.Mike[k]) > 0)
	}
}

func TestUnhandledTypes(t *testing.T) {
	var p interface{}
	assert.NotNil(t, Fill(&p))

	var q chan error
	assert.NotNil(t, Fill(&q))

	var r func(bool) int32
	assert.NotNil(t, Fill(&r))
}

func TestFillers(t *testing.T) {
	var p struct {
		ID         string
		Name       string
		Timestamps struct {
			Created time.Time
			Updated time.Time
		}
	}
	now := time.Now()
	assert.Must(t, Fill(&p, func(path string) (interface{}, bool) {
		switch path {
		case ".ID":
			return "id", true
		case
			".Timestamps.Created",
			".Timestamps.Updated":
			return now, true
		default:
			return nil, false
		}
	}))
	assert.Equal(t, p.ID, "id")
	assert.Equal(t, p.Timestamps.Updated, now)
	assert.Equal(t, p.Timestamps.Created, now)
}
