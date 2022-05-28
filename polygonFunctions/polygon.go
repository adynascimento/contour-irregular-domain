package polygonFunctions

import (
	"encoding/json"
	"heatmap/database"
	"log"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
)

type Blocks struct {
	Block_id      int              `json:"block_id"`
	ClientBQ_id   int              `json:"client_id"`
	Block_parent  int              `json:"block_parent"`
	Block_name    string           `json:"block_name"`
	Block_bounds  geojson.Geometry `json:"bounds"`
	Block_abrv    string           `json:"abvr"`
	Centroid      geojson.Geometry `json:"centroid"`
	Centroid_text string           `json:"centroid_text"`
	Properties    []Property       `json:"properties"`
	LeafParent    bool             `json:"leafParent"`
	Date          time.Time        `json:"date"`
	Data          DataDTO          `json:"data"`
}

type DataDTO struct {
	WindSpeed        float64 `json:"windSpeed"`
	SolarIrradiation float64 `json:"solarIrradiation"`
	Temperature      float64 `json:"temperature"`
	Rain             float64 `json:"rain"`
	RelativeHumidity float64 `json:"relativeHumidity"`
}

type Property struct {
	Index int64 // Property index (1 to 10)
	Value []string
	Date  time.Time
}

type KnownPoints struct {
	Lon    []float64
	Lat    []float64
	Values []float64
}

func PointsInMultiPolygon(multiPolygon orb.MultiPolygon, point orb.Point) bool {
	return planar.MultiPolygonContains(multiPolygon, point)
}

func GetBoxBounds(multiPolygon orb.MultiPolygon) orb.Ring {
	return multiPolygon.Bound().ToPolygon()[0][:4]
}

func GetMultiPolygonRedis(leafParent string) (orb.MultiPolygon, KnownPoints) {
	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	keys, err := redisClient.Keys("*:" + leafParent + ":home").Result()
	if err != nil {
		log.Println(err.Error())
	}

	multiPolygon := orb.MultiPolygon{}
	knownPoints := KnownPoints{}

	if len(keys) > 0 {
		for _, key := range keys {
			value, err := redisClient.Get(key).Result()
			if err != nil {
				log.Println(err.Error())
			}
			block := Blocks{}
			json.Unmarshal([]byte(value), &block)
			polygon := block.Block_bounds.Coordinates.(orb.Polygon)
			multiPolygon = append(multiPolygon, polygon)
			knownPoints.Lon = append(knownPoints.Lon, block.Centroid.Coordinates.(orb.Point)[0])
			knownPoints.Lat = append(knownPoints.Lat, block.Centroid.Coordinates.(orb.Point)[1])
			knownPoints.Values = append(knownPoints.Values, block.Data.Temperature)
		}
	}
	return multiPolygon, knownPoints
}
