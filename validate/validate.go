// Package validate provides common validations
package validate

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cfg "github.com/arutselvan15/estore-common/config"
)

var (
	checkFreezeEnabled = checkFreezeEnabledImpl
)

// PatchOperation patch operations
type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

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

// AdmissionRequired check admission required
func AdmissionRequired(ignoredNamespaces []string, admissionAnnotationKey string, metadata *metav1.ObjectMeta) (bool, string) {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredNamespaces {
		if metadata.Namespace == namespace {
			return false, fmt.Sprintf("skip validation for %s for its in special namespace: %s", metadata.Name, metadata.Namespace)
		}
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	switch strings.ToLower(annotations[admissionAnnotationKey]) {
	case "n", "no", "false", "off":
		return false, fmt.Sprintf("skip validation for %s for its special annotation %s: %s", metadata.Name,
			admissionAnnotationKey, annotations[admissionAnnotationKey])
	}

	return true, ""
}

// CreatePatchAnnotations create annotations patch
func CreatePatchAnnotations(availableAnnotations, addAnnotations map[string]string) (patch []PatchOperation) {
	if availableAnnotations == nil {
		patch = append(patch, PatchOperation{
			Op:    "add",
			Path:  "/metadata/annotations",
			Value: addAnnotations,
		})

		return patch
	}

	for key, value := range addAnnotations {
		op := "add"
		if _, ok := availableAnnotations[key]; ok {
			op = "replace"
		}

		patch = append(patch, PatchOperation{
			Op:    op,
			Path:  "/metadata/annotations/" + key,
			Value: value,
		})
	}

	return patch
}

// CreatePatchLabels create labels patch
func CreatePatchLabels(availableLabels, addLabels map[string]string) (patch []PatchOperation) {
	if availableLabels == nil {
		patch = append(patch, PatchOperation{
			Op:    "add",
			Path:  "/metadata/labels",
			Value: addLabels,
		})

		return patch
	}

	for key, value := range addLabels {
		op := "add"
		if _, ok := availableLabels[key]; ok {
			op = "replace"
		}

		patch = append(patch, PatchOperation{
			Op:    op,
			Path:  "/metadata/labels/" + key,
			Value: value,
		})
	}

	return patch
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
		st, err := time.Parse(cfg.TimeLayout, startTime)
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
				et, err := time.Parse(cfg.TimeLayout, endTime)
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
