package cgroups

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

func Value(fieldName string) (string, error) {
	path := fmt.Sprintf("%s/%s", CgroupsRoot, fieldName)
	if res, err := os.ReadFile(path); err != nil {
		logrus.Errorf("failed to read value of %s: %s", fieldName, err.Error())
		return "", err
	} else {
		return string(res), nil
	}
}

func SetValue(fieldName, value string) error {
	path := fmt.Sprintf("%s/%s", CgroupsRoot, fieldName)
	if err := os.WriteFile(path, []byte(value), 0644); err != nil {
		logrus.Errorf("failed to set %s value %s: %s\n", fieldName, value, err.Error())
		return err
	}
	return nil
}

func Remove(fieldName string) error {
	path := fmt.Sprintf("%s/%s", CgroupsRoot, fieldName)
	if err := os.Remove(path); err != nil {
		logrus.Errorf("failed to remove %s: %s\n", fieldName, err.Error())
		return err
	}
	return nil
}

func SetField[T any](prefixName string, controllers *T) error {
	controllersVal := reflect.ValueOf(controllers).Elem()
	controllersType := controllersVal.Type()
	for i := range controllersVal.NumField() {
		fieldVal := controllersVal.Field(i)
		if fieldVal.IsZero() {
			continue
		}
		fieldType := controllersType.Field(i)
		fieldName := fmt.Sprintf("%s.%s", prefixName, strings.ToLower(SplitWords(fieldType.Name)))
		if err := SetValue(fieldName, string(fieldVal.Interface().(string))); err != nil {
			logrus.Errorf("failed to set %s value %s\n", fieldName, fieldVal.Interface().(string))
			return err
		}
	}
	return nil
}

func SplitWords(x string) string {
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	return re.ReplaceAllString(x, "${1}.${2}")
}
