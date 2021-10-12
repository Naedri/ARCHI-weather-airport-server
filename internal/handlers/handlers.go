package handlers

import (
	"fmt"
	"meteo_des_aeroports/internal/utils"
)

func GetValueOfDataTypeWithRange(iata string, start string, end string, dataType string) (string, error) {
	listProbes, err := utils.HGetAll(fmt.Sprintf("%s:probes:%s", iata, dataType))

	if err != nil {
		return err.Error(), err
	}

	var result string
	result = `{`

	for index, probeId := range listProbes {
		if index%2 == 0 {

			result += fmt.Sprintf("\"%s\":[", probeId)
			listData, err := utils.ZRANGEBYSCORE(fmt.Sprintf("%s:probe:%s:%s", iata, dataType, probeId), start, end)

			if err != nil {
				return err.Error(), err
			}

			for _, data := range listData {
				result += data + ","
			}
			if len(listData) > 0 {
				result = result[:len(result)-1]
			}

			result += "],"
		}
	}
	if len(listProbes) > 0 {
		result = result[:len(result)-1]
	}
	result += "}"

	return result, nil
}

// func GetAverageValueOfTheDay(iata string) (interface{}, error) {

// }
