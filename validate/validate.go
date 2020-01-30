// Package validate provides common validations
package validate

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/arutselvan15/estore-common/config"
)

var (
	checkFreezeEnabled = checkFreezeEnabledImpl
)

// FreezeEnabled freeze enable check
func FreezeEnabled(component string) (bool, string, error) {
	enabled, err := checkFreezeEnabled(viper.GetString("app.freeze.startTime"), viper.GetString("app.freeze.endTime"))
	if err != nil {
		return false, "", err
	}

	if enabled {
		components := getFreezeComponents()
		// freeze for all components
		if _, ok := components["all"]; ok {
			return true, viper.GetString("app.freeze.message"), nil
		}

		// freeze for selective components
		if _, ok := components[component]; ok {
			return true, viper.GetString("app.freeze.message"), nil
		}
	}

	return false, "", nil
}

func getFreezeComponents() map[string]bool {
	components := map[string]bool{}
	for _, i := range strings.Split(viper.GetString("app.freeze.components"), ",") {
		components[i] = true
	}

	return components
}

func checkFreezeEnabledImpl(startTime, endTime string) (bool, error) {
	if startTime != "" {
		st, err := time.Parse(config.TimeLayout, startTime)
		if err != nil {
			return false, fmt.Errorf("unable to parse freeze start time %s, %s", startTime, err.Error())
		}

		tn := time.Now()
		tn.Format(time.RFC3339)

		// check for freeze period
		// time now - freeze start time, if positive, time now is greater than freeze start time, possible in freeze period
		if tn.Sub(st) > 0 {
			// check freeze end time to see if still in freeze
			if endTime != "" {
				et, err := time.Parse(config.TimeLayout, endTime)
				if err != nil {
					return false, fmt.Errorf("unable to parse freeze end time %s, %s", endTime, err.Error())
				}

				// time now - freeze end time, if positive, time now is greater than freeze end time, freeze is over
				if tn.Sub(et) > 0 {
					return false, nil
				}

				return true, nil
			}

			// no end time mentioned considering as no end freeze window
			return true, nil
		}

		//  time now is less than freeze start time so freeze period is not started yet
		return false, nil
	}

	return false, nil
}
