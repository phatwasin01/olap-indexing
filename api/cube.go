package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phatwasin01/olap-indexing/util"
)

type AddCubeMetaDataRequest struct {
	ID   int  `json:"id"`
	Cube Cube `json:"cube"`
}
type Cube struct {
	Name       string      `json:"name"`
	Dimensions []Dimension `json:"dimensions"`
	Measure    string      `json:"measure"`
}
type Dimension struct {
	Name                 string                 `json:"name"`
	DimensionLevelNumber []DimensionLevelNumber `json:"dimension_level_number"`
	DimensionLevelString []DimensionLevelString `json:"dimension_level_string"`
}

// ex. price min = 0, max = 100000, step =100
// ex. year min = 2010, max = 2020, step = 1
type DimensionLevelNumber struct {
	Name        string `json:"name"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	SettingMin  int    `json:"setting_min"`
	SettingMax  int    `json:"setting_max"`
	SettingStep int    `json:"setting_step"`
}
type DimensionLevelString struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

const (
	SettingDimensionLevelNumberFix   string = "fix"
	SettingDimensionLevelNumberRange string = "range"
	SettingDimensionLevelString      string = "string"
)

func (server *Server) AddCubeMetaData(ctx *gin.Context) {
	var req AddCubeMetaDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := fmt.Sprintf("measure:%s", req.Cube.Measure)
	err := server.redis.SAdd(ctx.Request.Context(), key, req.ID).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, dimension := range req.Cube.Dimensions {
		for _, levelNum := range dimension.DimensionLevelNumber {
			// Create a key for the numeric dimension range
			rangeMin, err := util.FindRange(levelNum.SettingMin, levelNum.SettingMax, levelNum.SettingStep, levelNum.Min)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			rangeMax, err := util.FindRange(levelNum.SettingMin, levelNum.SettingMax, levelNum.SettingStep, levelNum.Max)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			fmt.Println(rangeMin, rangeMax)
			if levelNum.SettingStep == 1 {
				key := fmt.Sprintf("setting:dimension:%s:level:%s", dimension.Name, levelNum.Name)
				err := server.redis.Set(ctx.Request.Context(), key, SettingDimensionLevelNumberFix, 0).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
				key = fmt.Sprintf("setting:dimension:%s:level:%s:fix", dimension.Name, levelNum.Name)
				err = server.redis.HSet(ctx.Request.Context(), key, "min", levelNum.SettingMin, "max", levelNum.SettingMax, "step", levelNum.SettingStep).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
			} else {
				key := fmt.Sprintf("setting:dimension:%s:level:%s", dimension.Name, levelNum.Name)
				err := server.redis.Set(ctx.Request.Context(), key, SettingDimensionLevelNumberRange, 0).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
				key = fmt.Sprintf("setting:dimension:%s:level:%s:range", dimension.Name, levelNum.Name)
				err = server.redis.HSet(ctx.Request.Context(), key, "min", levelNum.SettingMin, "max", levelNum.SettingMax, "step", levelNum.SettingStep).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
			}
			for i := rangeMin; i <= rangeMax; i += levelNum.SettingStep {
				var key string
				if levelNum.SettingStep == 1 {
					key = fmt.Sprintf("dimension:%s:level:%s:fix:%d", dimension.Name, levelNum.Name, i)
				} else {
					key = fmt.Sprintf("dimension:%s:level:%s:range:%d-%d", dimension.Name, levelNum.Name, i, i+levelNum.SettingStep-1)
				}
				err := server.redis.SAdd(ctx.Request.Context(), key, req.ID).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
			}
			keyLevel := fmt.Sprintf("dimension:%s:level:%s", dimension.Name, levelNum.Name)
			err = server.redis.SAdd(ctx.Request.Context(), keyLevel, req.ID).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			fmt.Println(keyLevel)
		}

		for _, levelStr := range dimension.DimensionLevelString {
			key := fmt.Sprintf("setting:dimension:%s:level:%s", dimension.Name, levelStr.Name)
			err := server.redis.Set(ctx.Request.Context(), key, SettingDimensionLevelString, 0).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			fmt.Println(key)
			for _, value := range levelStr.Value {
				key := fmt.Sprintf("dimension:%s:level:%s:value:%s", dimension.Name, levelStr.Name, value)
				err := server.redis.SAdd(ctx.Request.Context(), key, req.ID).Err()
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				fmt.Println(key)
			}
			keyLevel := fmt.Sprintf("dimension:%s:level:%s", dimension.Name, levelStr.Name)
			err = server.redis.SAdd(ctx.Request.Context(), keyLevel, req.ID).Err()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			fmt.Println(keyLevel)

		}
		// Storing Cube dimensions
		key := fmt.Sprintf("cube:%d", req.ID)
		err := server.redis.SAdd(ctx.Request.Context(), key, dimension.Name).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		fmt.Println(key)
		keyDimension := fmt.Sprintf("dimension:%s", dimension.Name)
		err = server.redis.SAdd(ctx.Request.Context(), keyDimension, req.ID).Err()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		fmt.Println(keyDimension)

	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Cube metadata added successfully"})
}

type RetreiveCubeByDimensionsRequest struct {
	Dimensions []RetreiveCubeByDimensions `json:"dimensions"`
	Measure    string                     `json:"measure"`
}

type RetreiveCubeByDimensions struct {
	Name                                string                                `json:"name"`
	RetreiveCubeByDimensionsLevelNumber []RetreiveCubeByDimensionsLevelNumber `json:"numbers"`
	RetreiveCubeByDimensionsLevelString []RetreiveCubeByDimensionsLevelString `json:"strings"`
}
type RetreiveCubeByDimensionsLevelNumber struct {
	Name  string `json:"name"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
	Exact int    `json:"exact"`
}
type RetreiveCubeByDimensionsLevelString struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

func (server *Server) RetreiveCubeByDimensions(ctx *gin.Context) {
	var req RetreiveCubeByDimensionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a unique temporary key for storing the intersection result
	tempKey := fmt.Sprintf("temp:intersection:%d", time.Now().UnixNano())

	// Start a Redis pipeline for efficiency
	pipe := server.redis.TxPipeline()
	var dimensionIntersections []string

	// Process each dimension and add its corresponding Redis keys to the intersection
	for _, dim := range req.Dimensions {
		var keys []string

		// Process numeric dimensions
		for _, num := range dim.RetreiveCubeByDimensionsLevelNumber {
			key := fmt.Sprintf("dimension:%s:level:%s", dim.Name, num.Name)
			keys = append(keys, key)
			levelTypeKey := fmt.Sprintf("setting:dimension:%s:level:%s", dim.Name, num.Name)
			levelTypeValue, err := server.redis.Get(ctx.Request.Context(), levelTypeKey).Result()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if levelTypeValue == SettingDimensionLevelNumberFix {
				if num.Exact != -1 {
					key := fmt.Sprintf("dimension:%s:level:%s:fix:%d", dim.Name, num.Name, num.Exact)
					keys = append(keys, key)
					continue
				}
				key := fmt.Sprintf("setting:dimension:%s:level:%s:fix", dim.Name, num.Name)
				result, err := server.redis.HMGet(ctx, key, "min", "max").Result()
				if err != nil {
					fmt.Println("Error fetching data:", err)
					return
				}
				// Parse the result into integers
				min, errMin := parseInt(result[0])
				max, errMax := parseInt(result[1])
				if errMin != nil || errMax != nil {
					fmt.Println("Error parsing values:", errMin, errMax)
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing values"})
					return
				}
				if (min <= num.Min && num.Min <= max) && (min <= num.Max && num.Max <= max) {
					for i := num.Min; i <= num.Max; i++ {
						key := fmt.Sprintf("dimension:%s:level:%s:fix:%d", dim.Name, num.Name, i)
						keys = append(keys, key)
					}
				}
			} else if levelTypeValue == SettingDimensionLevelNumberRange {
				key := fmt.Sprintf("setting:dimension:%s:level:%s:range", dim.Name, num.Name)
				result, err := server.redis.HMGet(ctx, key, "min", "max", "step").Result()
				if err != nil {
					fmt.Println("Error fetching data:", err)
					return
				}
				// Parse the result into integers
				min, errMin := parseInt(result[0])
				max, errMax := parseInt(result[1])
				step, errStep := parseInt(result[2])
				if errMin != nil || errMax != nil || errStep != nil {
					fmt.Println("Error parsing values:", errMin, errMax, errStep)
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing values"})
					return
				}
				minRange, err := util.FindRange(min, max, step, num.Min)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				maxRange, err := util.FindRange(min, max, step, num.Max)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				for i := minRange; i <= maxRange; i += step {
					key := fmt.Sprintf("dimension:%s:level:%s:range:%d-%d", dim.Name, num.Name, i, i+step-1)
					keys = append(keys, key)
				}
			}

		}

		// Process string dimensions
		for _, str := range dim.RetreiveCubeByDimensionsLevelString {
			key := fmt.Sprintf("dimension:%s:level:%s", dim.Name, str.Name)
			keys = append(keys, key)
			for _, val := range str.Value {
				key := fmt.Sprintf("dimension:%s:level:%s:value:%s", dim.Name, str.Name, val)
				keys = append(keys, key)
			}
		}
		key := fmt.Sprintf("dimension:%s", dim.Name)
		keys = append(keys, key)

		// Perform an intersection for each dimension and store the result in a unique key
		if len(keys) > 0 {
			dimTempKey := fmt.Sprintf("temp:dimension:%s:%d", dim.Name, time.Now().UnixNano())
			dimensionIntersections = append(dimensionIntersections, dimTempKey)
			pipe.SInterStore(ctx.Request.Context(), dimTempKey, keys...)
		}
		fmt.Println("keys", keys)
	}
	measureKey := fmt.Sprintf("measure:%s", req.Measure)
	finalTempKey := fmt.Sprintf("temp:finalIntersection:%d", time.Now().UnixNano())
	fmt.Println("dimensionIntersections", dimensionIntersections)
	if len(dimensionIntersections) > 0 {
		pipe.SInterStore(ctx.Request.Context(), finalTempKey, append(dimensionIntersections, measureKey)...)
	}

	// Execute the pipeline
	if _, err := pipe.Exec(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the intersection result
	// result, err := server.redis.SMembers(ctx.Request.Context(), tempKey).Result()
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// Retrieve the final intersection result
	finalResult, err := server.redis.SMembers(ctx.Request.Context(), finalTempKey).Result()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Clean up all temporary keys, including the final intersection key and all dimension keys
	dimensionIntersections = append(dimensionIntersections, finalTempKey)
	for _, tempKey := range dimensionIntersections {
		server.redis.Del(ctx.Request.Context(), tempKey)
	}

	// Clean up the temporary key
	pipe.Del(ctx.Request.Context(), tempKey)

	// Respond with the result
	ctx.JSON(http.StatusOK, gin.H{"cube_ids": finalResult})
}

// Helper function to parse interface{} to int
func parseInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	case nil:
		return 0, fmt.Errorf("value is nil")
	default:
		return 0, fmt.Errorf("invalid type")
	}
}
