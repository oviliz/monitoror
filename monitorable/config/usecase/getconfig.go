package usecase

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/structs"
	"github.com/monitoror/monitoror/monitorable/config/models"
)

var unknownFieldRegex *regexp.Regexp
var fieldTypeMismatchRegex *regexp.Regexp
var invalidEscapedCharacterRegex *regexp.Regexp

func init() {
	// Based on: https://github.com/golang/go/blob/release-branch.go1.14/src/encoding/json/decode.go#L755
	unknownFieldRegex = regexp.MustCompile(`json: unknown field "(.*)"`)

	// Based on: https://github.com/golang/go/blob/go1.14/src/encoding/json/scanner.go#L343
	fieldTypeMismatchRegex = regexp.MustCompile(`json: cannot unmarshal .+ into Go struct field (.+) of type (.+)`)

	// Based on: https://github.com/golang/go/blob/go1.14/src/encoding/json/scanner.go#L343
	invalidEscapedCharacterRegex = regexp.MustCompile(`'(.*)' in string escape code`)
}

// GetConfig and set default value for Config from repository
func (cu *configUsecase) GetConfig(params *models.ConfigParams) *models.ConfigBag {
	configBag := &models.ConfigBag{}

	var err error
	if params.URL != "" {
		configBag.Config, err = cu.repository.GetConfigFromURL(params.URL)
	} else if params.Path != "" {
		configBag.Config, err = cu.repository.GetConfigFromPath(params.Path)
	}

	if err == nil {
		return configBag
	}

	switch e := err.(type) {
	case *models.ConfigFileNotFoundError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorConfigNotFound,
			Message: e.Error(),
			Data:    models.ConfigErrorData{Value: e.PathOrURL},
		})
	case *models.ConfigVersionFormatError:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnsupportedVersion,
			Message: e.Error(),
			Data: models.ConfigErrorData{
				FieldName: "version",
				Value:     e.WrongVersion,
				Expected:  fmt.Sprintf(`%q >= version >= %q`, MinimalVersion, CurrentVersion),
			},
		})
	case *models.ConfigUnmarshalError:
		// Check if error is "json: unknown field"
		if unknownFieldRegex.MatchString(err.Error()) {
			subMatch := unknownFieldRegex.FindAllStringSubmatch(err.Error(), 1)

			var field = ""
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				field = subMatch[0][1]
			}

			var configField = structs.Fields(models.Config{})
			var tileConfigFields = structs.Fields(models.TileConfig{})
			var expectedFields = append(configField, tileConfigFields...)
			var expectedFieldNames []string

			for _, expectedField := range expectedFields {
				var jsonTag = expectedField.Tag("json")
				if jsonTag != "" && jsonTag != "-" {
					var expectedFieldName = strings.Replace(jsonTag, ",omitempty", "", 1)
					expectedFieldNames = append(expectedFieldNames, expectedFieldName)
				}
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnknownField,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: e.RawConfig,
					Expected:      strings.Join(expectedFieldNames, ", "),
				},
			})
		} else if fieldTypeMismatchRegex.MatchString(err.Error()) {
			subMatch := fieldTypeMismatchRegex.FindAllStringSubmatch(err.Error(), 1)

			var field = ""
			var expectedType = ""
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				fieldParts := strings.Split(subMatch[0][1], ".")
				field = fieldParts[len(fieldParts)-1]
				expectedType = subMatch[0][2]
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorFieldTypeMismatch,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					FieldName:     field,
					ConfigExtract: e.RawConfig,
					Expected:      expectedType,
				},
			})
		} else if invalidEscapedCharacterRegex.MatchString(err.Error()) {
			subMatch := invalidEscapedCharacterRegex.FindAllStringSubmatch(err.Error(), 1)

			var invalidEscapedCharacter = ""
			if len(subMatch) > 0 && len(subMatch[0]) > 1 {
				invalidEscapedCharacter = subMatch[0][1]
			}

			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorInvalidEscapedCharacter,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					ConfigExtract:          e.RawConfig,
					ConfigExtractHighlight: fmt.Sprintf(`\%s`, invalidEscapedCharacter),
				},
			})
		} else {
			configBag.AddErrors(models.ConfigError{
				ID:      models.ConfigErrorUnableToParseConfig,
				Message: e.Error(),
				Data: models.ConfigErrorData{
					ConfigExtract: e.RawConfig,
				},
			})
		}
	default:
		configBag.AddErrors(models.ConfigError{
			ID:      models.ConfigErrorUnexpectedError,
			Message: err.Error(),
		})
	}

	return configBag
}
