package service

import (
	"../proto"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sort"
)

type User struct {
	UID       string
	Email     string
	Name      string
	LastLogin int64
}

type Location struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type DriverLocation struct {
	UID       string
	TimeStamp int64
	Location  Location
	distance  float64
}

type Distance struct {
	calculated float64
	location   Location
}

type Server struct {
}

func (s Server) PickDriverAndAssign(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	var results []DriverLocation
	filter2 := bson.D{
		{"$geoNear", bson.D{
			{"distanceField", "distance"},
			{"maxDistance", 50000},
			{"spherical", true},
		}},
	}

	curGeoNear, err := DriverLocationCollection.Aggregate(context.Background(), mongo.Pipeline{filter2})
	var geoDriverCoordinates []bson.M
	if err = curGeoNear.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	curGeoNear, err = DriverLocationCollection.Aggregate(context.Background(), mongo.Pipeline{filter2})
	if err = curGeoNear.All(ctx, &geoDriverCoordinates); err != nil {
		log.Fatal(err)
	}

	for index, result := range geoDriverCoordinates {
		results[index].distance = result["distance"].(float64)
		fmt.Println(result["distance"].(float64))
	}

	fmt.Println("Results are ... ")
	fmt.Println(results)

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].distance < results[j].distance
	})

	return &proto.Response{
		DriverID:             results[0],
	}, nil
}

func NewLocation(lng, lat float64) Location {
	return Location{
		Type:        "Point",
		Coordinates: []float64{lng, lat},
	}
}
