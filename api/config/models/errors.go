package models

import (
	"fmt"
)

// ConfigFileNotFoundError
type ConfigFileNotFoundError struct {
	PathOrURL string
	Err       error
}

func (e *ConfigFileNotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(`CoreConfig not found at: %s, %v`, e.PathOrURL, e.Err.Error())
	}
	return fmt.Sprintf(`CoreConfig not found at: %s`, e.PathOrURL)
}
func (e *ConfigFileNotFoundError) Unwrap() error { return e.Err }

// ConfigVersionFormatError
type ConfigVersionFormatError struct {
	WrongVersion string
}

func (e *ConfigVersionFormatError) Error() string {
	return fmt.Sprintf(`json: cannot unmarshal %s into Go struct field CoreConfig.ConfigVersion of type string and X.y format`, e.WrongVersion)
}

//ConfigUnmarshalError
type ConfigUnmarshalError struct {
	Err       error
	RawConfig string
}

func (e *ConfigUnmarshalError) Error() string {
	return e.Err.Error()
}
func (e *ConfigUnmarshalError) Unwrap() error { return e.Err }
