package polygonFunctions

import (
	"encoding/json"
	"heatmap/database"
	"log"
	"testing"

	"github.com/paulmach/orb"
)

func TestGetLimitBounds(t *testing.T) {
	t.Run("Testing two separate squares", func(t *testing.T) {
		square1 := orb.Polygon{{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}}}
		square2 := orb.Polygon{{{5, 5}, {5, 8}, {9, 9}, {9, 2}, {5, 5}}}
		multiPolygon := orb.MultiPolygon{square1, square2}
		got := GetBoxBounds(multiPolygon)
		want := orb.Ring{
			{0, 0},
			{9, 0},
			{9, 9},
			{0, 9},
		}
		if !got.Equal(want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
func TestPointsInPolygon(t *testing.T) {
	cases := []struct {
		Description  string
		Point        orb.Point
		MultiPolygon orb.MultiPolygon
		want         bool
	}{
		{"Verify if point is inside square",
			orb.Point{0.5, 0.5},
			orb.MultiPolygon{
				orb.Polygon{
					{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}},
				},
			},

			true,
		},
		{"Verify if point is outside square",
			orb.Point{3, 3},
			orb.MultiPolygon{
				orb.Polygon{
					{{0, 0}, {0, 2}, {2, 2}, {2, 0}, {0, 0}},
				},
			},
			false,
		},
		{"Verify if point is inside complex polygon",
			orb.Point{-56.4798333, -14.623567},
			orb.MultiPolygon{
				getRealPolygon(),
			},
			true,
		},
		{"Verify if point is outside complex polygon",
			orb.Point{-57.4798333, -15.623567},
			orb.MultiPolygon{
				getRealPolygon(),
			},
			false,
		},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got := PointsInMultiPolygon(test.MultiPolygon, test.Point)
			if got != test.want {
				t.Errorf("got %t, want %t", got, test.want)
			}
		})
	}
}

func getRealPolygon() orb.Polygon {
	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	value, err := redisClient.Get("4546:4532").Result()
	if err != nil {
		log.Println(err.Error())
	}
	block := Blocks{}
	json.Unmarshal([]byte(value), &block)
	return block.Block_bounds.Coordinates.(orb.Polygon)

}
