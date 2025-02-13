package validators

import (
	errs "errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/buonotti/apisense/v2/errors"
	"github.com/buonotti/apisense/v2/log"
	"github.com/goccy/go-json"
	"github.com/spf13/viper"
)

// ValidatorDefinition is the definition of an external validator
type ValidatorDefinition struct {
	Name  string   // Name is the name of the validator
	Path  string   // Path is the path to the executable
	Args  []string // Args are the arguments to pass to the executable
	Fatal bool     // Fatal controls whether the validator is fatal or not that is if it fails the pipeline should stop
	Slim  bool     // Slim controls how much data the validator gets. Setting this to true reduces the context the validator gets
}

// parse parses the external validators in the config file and returns a slice containing all validators to later use in the pipeline
func parse() ([]ValidatorDefinition, error) {
	object := viper.Get("validation.external_validators")
	if object == nil {
		return []ValidatorDefinition{}, nil
	}

	arr, isArray := object.([]interface{})
	if !isArray {
		return nil, errors.ExternalValidatorParseError.New("cannot parse external validators. Expected []any, got %T", object)
	}

	validators := make([]ValidatorDefinition, len(arr))

	for i, arrayEntry := range arr {
		obj, isStringMap := arrayEntry.(map[string]interface{})
		if !isStringMap {
			return nil, errors.ExternalValidatorParseError.New("cannot parse external validators. Expected map[string]any, got %T", arrayEntry)
		}

		args, err := parseArgs(obj["args"])
		if err != nil {
			return nil, err
		}

		validators[i] = ValidatorDefinition{
			Name:  obj["name"].(string),
			Path:  obj["path"].(string),
			Args:  args,
			Fatal: obj["fatal"].(bool),
		}
	}
	return validators, nil
}

// parseArgs parses the args to pass to the validator
func parseArgs(obj any) ([]string, error) {
	arr, isArray := obj.([]interface{})
	if !isArray {
		return nil, errors.ExternalValidatorParseError.New("cannot parse external validator. expected []interface{}, got %T", obj)
	}

	if len(arr) == 0 {
		return []string{}, nil
	}

	args := make([]string, len(arr))
	for i, elem := range arr {
		if _, isString := elem.(string); !isString {
			return nil, errors.ExternalValidatorParseError.New("cannot parse external validator. expected []string, got []%T", elem)
		}
		args[i] = elem.(string)
	}

	return args, nil
}

// LoadExternalValidators loads the exernal validators from the definitions
func LoadExternalValidators() ([]Validator, error) {
	definitions, err := parse()
	if err != nil {
		return nil, err
	}

	externalValidators := make([]Validator, len(definitions))
	for i, definition := range definitions {
		externalValidators[i] = NewExternalValidator(definition)
	}

	return externalValidators, nil
}

// NewExternalValidator creates a new external validator based on the given definition and returns a validation.Validator
func NewExternalValidator(definition ValidatorDefinition) Validator {
	return externalValidator{
		Definition: definition,
	}
}

// externalValidator is a validator that was defined in the configuration file
type externalValidator struct {
	Definition ValidatorDefinition // Definition is the external.ValidatorDefinition that defines the external validator
}

// Name returns the name of the validator: external.<name> where name is the name
// of the external validator defined in the config
func (v externalValidator) Name() string {
	return "external." + v.Definition.Name
}

// Validate validates an item by serializing it and sending it to the external
// process then returning an error according to the status code of the external
// program
func (v externalValidator) Validate(item ValidationItem) error {
	var err error
	var jsonString []byte
	if v.Definition.Slim {
		jsonString, err = json.Marshal(SlimValidationItem{
			response: item.Response(),
		})
	} else {
		jsonString, err = json.Marshal(ExtendedValidationItem{
			response:   item.Response(),
			definition: item.Definition(),
		})
	}

	jsonString, err = json.Marshal(item)
	outString := &strings.Builder{}
	if err != nil {
		return errors.CannotSerializeItemError.Wrap(err, "cannot serialize item: %s", err)
	}
	cmd := exec.Command(v.Definition.Path, v.Definition.Args...)
	cmd.Stdin = strings.NewReader(string(jsonString))
	cmd.Stdout = outString

	validatorOut := strings.Builder{}
	validatorErr := strings.Builder{}
	cmd.Stdout = &validatorOut
	cmd.Stderr = &validatorErr

	err = cmd.Run()

	log.DaemonLogger().With("validator", fmt.Sprintf("external.%s", v.Definition.Name)).Debugf("validator output: %s", validatorOut.String())
	log.DaemonLogger().With("validator", fmt.Sprintf("external.%s", v.Definition.Name)).Debugf("validator error: %s", validatorErr.String())

	if err != nil {
		var exitErr *exec.ExitError
		if errs.As(err, &exitErr) {
			if exitErr.ExitCode() == 1 {
				return errors.ValidationError.New("validation failed: %s", validatorErr.String())
			} else {
				return errors.ValidationError.New("validation failed: unexpected exit code from external validator: %d", exitErr.ExitCode())
			}
		}
	}
	return nil
}

func (v externalValidator) IsFatal() bool {
	return v.Definition.Fatal
}

func (v externalValidator) IsSlim() bool {
	return v.Definition.Slim
}
