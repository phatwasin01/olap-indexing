package api

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) MockupDataForTesting(ctx *gin.Context) {
	//random number
	rand := rand.Intn(10000)

	// Mock data based on AddCubeMetaDataRequest format
	mockData := []AddCubeMetaDataRequest{
		{
			ID: rand,
			Cube: Cube{
				Name: "MarketAnalysisCube",
				Dimensions: []Dimension{
					{
						Name: "Date",
						DimensionLevelNumber: []DimensionLevelNumber{
							{
								Name:        "Year",
								Min:         2001,
								Max:         2012,
								SettingMin:  1901,
								SettingMax:  2500,
								SettingStep: 1,
							},
							// ... add other dimensions similarly
						},
						DimensionLevelString: []DimensionLevelString{
							// ... add string dimensions similarly
						},
					},
					// ... add other cubes similarly
				},
			},
		},
		// ... add more mock data if needed
	}

	// Iterate over mock data and add to Redis
	for _, data := range mockData {
		// You can call your existing AddCubeMetaData logic here
		// For example, server.AddCubeMetaDataLogic(data)
		// Or write the logic to add the data directly to Redis
		// For simplicity, I am just printing the data
		fmt.Printf("Adding mock data to Redis: %+v\n", data)
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Mock data added successfully"})
}

// You would also need to create a function that mimics the logic of AddCubeMetaData
// and populates Redis based on the provided data.
