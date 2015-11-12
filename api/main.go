package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/db/testdb"
	"github.com/loopfz/scecret/utils/tonic"
)

var db *gorp.DbMap

func main() {

	tdb, err := testdb.InitTestDB()
	if err != nil {
		panic(err)
	}

	db = tdb

	router := gin.Default()

	// Auth
	router.POST("/register", tonic.Handler(RegisterUser, 201))
	router.POST("/auth", tonic.Handler(Auth, 200))
	router.GET("/me", tonic.Handler(GetMe, 200))

	// Scenarios
	router.POST("/scenario", tonic.Handler(NewScenario, 201))
	router.GET("/scenario", tonic.Handler(ListScenarios, 200))
	router.GET("/scenario/:scenario", tonic.Handler(GetScenario, 200))
	router.PUT("/scenario/:scenario", tonic.Handler(UpdateScenario, 200))
	router.DELETE("/scenario/:scenario", tonic.Handler(DeleteScenario, 204))
	router.GET("/scenario/:scenario/graph", tonic.Handler(GetGraph, 200))

	// Locations
	router.POST("/scenario/:scenario/location", tonic.Handler(NewLocation, 201))
	router.GET("/scenario/:scenario/location", tonic.Handler(ListLocations, 200))
	router.GET("/scenario/:scenario/location/:location", tonic.Handler(GetLocation, 200))
	router.PUT("/scenario/:scenario/location/:location", tonic.Handler(UpdateLocation, 200))
	router.DELETE("/scenario/:scenario/location/:location", tonic.Handler(DeleteLocation, 204))

	// Location cards
	router.POST("/scenario/:scenario/location/:location/card", tonic.Handler(NewLocationCard, 201))
	router.GET("/scenario/:scenario/location/:location/card", tonic.Handler(ListLocationCards, 200))
	router.GET("/scenario/:scenario/location/:location/card/:card", tonic.Handler(GetLocationCard, 200))
	router.PUT("/scenario/:scenario/location/:location/card/:location_card", tonic.Handler(UpdateLocationCard, 200))
	router.DELETE("/scenario/:scenario/location/:location/card/:location_card", tonic.Handler(DeleteLocationCard, 204))

	// Location links
	router.POST("/scenario/:scenario/locationlink", tonic.Handler(NewLocationLink, 201))
	router.GET("/scenario/:scenario/locationlink", tonic.Handler(ListLocationLinks, 200))
	router.GET("/scenario/:scenario/locationlink/:locationlink", tonic.Handler(GetLocationLink, 200))
	router.DELETE("/scenario/:scenario/locationlink/:locationlink", tonic.Handler(DeleteLocationLink, 204))

	// Element links
	router.POST("/scenario/:scenario/elementlink", tonic.Handler(NewElementLink, 201))
	router.GET("/scenario/:scenario/elementlink", tonic.Handler(ListElementLinks, 200))
	router.GET("/scenario/:scenario/elementlink/:elementlink", tonic.Handler(GetElementLink, 200))
	router.DELETE("/scenario/:scenario/elementlink/:elementlink", tonic.Handler(DeleteElementLink, 204))

	// State tokens
	router.GET("/scenario/:scenario/statetoken", tonic.Handler(ListStateTokens, 200))
	router.GET("/scenario/:scenario/statetoken/:statetoken", tonic.Handler(GetStateToken, 200))

	// State token links
	router.POST("/scenario/:scenario/statetokenlink", tonic.Handler(NewStateTokenLink, 201))
	router.GET("/scenario/:scenario/statetokenlink", tonic.Handler(ListStateTokenLinks, 200))
	router.GET("/scenario/:scenario/statetokenlink/:statetokenlink", tonic.Handler(GetStateTokenLink, 200))
	router.DELETE("/scenario/:scenario/statetokenlink/:statetokenlink", tonic.Handler(DeleteStateTokenLink, 204))

	// Stats
	router.POST("/scenario/:scenario/stat", tonic.Handler(NewStat, 201))
	router.GET("/scenario/:scenario/stat", tonic.Handler(ListStats, 200))
	router.GET("/scenario/:scenario/stat/:stat", tonic.Handler(GetStat, 200))
	router.PUT("/scenario/:scenario/stat/:stat", tonic.Handler(UpdateStat, 200))
	router.DELETE("/scenario/:scenario/stat/:stat", tonic.Handler(DeleteStat, 204))

	// Skill tests
	router.POST("/scenario/:scenario/skilltest", tonic.Handler(CreateSkillTest, 201))
	router.GET("/scenario/:scenario/skilltest", tonic.Handler(ListSkillTests, 200))
	router.GET("/scenario/:scenario/skilltest/:skilltest", tonic.Handler(GetSkillTest, 200))
	router.PUT("/scenario/:scenario/skilltest/:skilltest", tonic.Handler(UpdateSkillTest, 200))
	router.DELETE("/scenario/:scenario/skilltest/:skilltest", tonic.Handler(DeleteSkillTest, 204))

	// Icons
	router.POST("/scenario/:scenario/icon", tonic.Handler(NewIcon, 201))
	router.GET("/scenario/:scenario/icon", tonic.Handler(ListIcons, 200))
	router.GET("/scenario/:scenario/icon/:icon", tonic.Handler(GetIcon, 200))
	router.PUT("/scenario/:scenario/icon/:icon", tonic.Handler(UpdateIcon, 200))
	router.DELETE("/scenario/:scenario/icon/:icon", tonic.Handler(DeleteIcon, 204))

	// Elements
	router.POST("/scenario/:scenario/element", tonic.Handler(NewElement, 201))
	router.GET("/scenario/:scenario/element", tonic.Handler(ListElements, 200))
	router.GET("/scenario/:scenario/element/:element", tonic.Handler(GetElement, 200))
	router.PUT("/scenario/:scenario/element/:element", tonic.Handler(UpdateElement, 200))
	router.DELETE("/scenario/:scenario/element/:element", tonic.Handler(DeleteElement, 204))

	// TODO GET /scenario/:scenario/card
	// TODO GET /scenario/:scenario/card/:card
	// TODO PUT /scenario/:scenario/card/:card
	// TODO POST /scenario/:scenario/card/:card/icon
	// TODO GET /scenario/:scenario/card/:card/icon
	// TODO GET /scenario/:scenario/card/:card/icon/:icon
	// TODO PUT /scenario/:scenario/card/:card/icon/:icon
	// TODO DELETE /scenario/:scenario/card/:card/icon/:icon

	// TODO POST /sandbox

	router.Run(":8080")
}
