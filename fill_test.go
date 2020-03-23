package feign

import (
	"testing"
	"time"

	"github.com/deliveroo/assert-go"
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
	}
	assert.ErrorContains(t, Fill(p), "not a pointer value")
	assert.Must(t, Fill(&p))
	assert.True(t, p.Alpha != 0)
	assert.True(t, p.Beta != "")
	assert.NotNil(t, p.Charlie)
	assert.NotNil(t, *p.Delta)
	assert.NotNil(t, p.Echo.Foxtrot)
	assert.NotNil(t, p.Echo.Golf != 0)
	assert.True(t, len(p.Echo.Hotel) > 0)
	assert.NotNil(t, p.Echo.Indigo)
	assert.NotNil(t, p.Echo.Indigo.Juliet)
	assert.False(t, p.Echo.Indigo.Kilo.IsZero())
	assert.True(t, len(p.Echo.Indigo.Lima) > 0)
	assert.True(t, len(p.Echo.Indigo.Mike) > 0)
	for k := range p.Echo.Indigo.Mike {
		assert.True(t, len(p.Echo.Indigo.Mike[k]) > 0)
	}
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
	assert.Must(t, Fill(&p, func(path string) (bool, interface{}) {
		switch path {
		case ".ID":
			return true, "id"
		case
			".Timestamps.Created",
			".Timestamps.Updated":
			return true, now
		default:
			return false, nil
		}
	}))
	assert.Equal(t, p.ID, "id")
	assert.Equal(t, p.Timestamps.Updated, now)
	assert.Equal(t, p.Timestamps.Created, now)
}
