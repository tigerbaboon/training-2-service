package config

import (
	"app/config"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ConfigService struct {
	conf *config.Config
}

func newService() *ConfigService {
	godotenv.Load()
	conf := configWithDefault(&config.App)
	config.Init(conf)
	return &ConfigService{
		conf,
	}
}

func (s *ConfigService) App() *config.Config {
	return s.conf
}

func (s *ConfigService) Database() *config.Database {
	return &s.conf.Database
}

func configWithDefault(confDefault *config.Config) *config.Config {
	rConfig := reflect.ValueOf(confDefault).Elem()
	t := rConfig.Type()
	for i := 0; i < t.NumField(); i++ {
		key := stringToAllCapsCase(t.Field(i).Name)
		switch t.Field(i).Type.Kind() {
		case reflect.Struct:
			rConfig.Field(i).Set(configStruct(key, rConfig.Field(i)))
		default:
			defaultValue := rConfig.Field(i).Interface()
			rConfig.Field(i).Set(reflect.ValueOf(conf(key, defaultValue)))
		}
	}
	return confDefault
}

func configStruct(prefix string, v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Struct:
		configStructForStruct(prefix, v)
	case reflect.Map:
		configStructForMap(prefix, v)
	case reflect.Ptr:
		configStructForPtr(prefix, v)
	default:
		configStructForDefault(prefix, v)
	}
	return v
}

func configStructForStruct(prefix string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fieldName := v.Type().Field(i).Name
		pf := fieldName[0]
		if pf >= 'A' && pf <= 'Z' {
			key := fmt.Sprintf("%s_%s", prefix, stringToAllCapsCase(fieldName))
			v.Field(i).Set(configStruct(key, v.Field(i)))
		}
	}
}

func configStructForMap(prefix string, v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	mapKeyMap := map[string]bool{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, prefix) {
			env := strings.Split(env, "=")[0]
			rmPrefix := strings.ReplaceAll(env, prefix+"_", "")
			mapKey := strings.ToLower(strings.Split(rmPrefix, "_")[0])
			mapKeyMap[mapKey] = true
		}
	}

	mapKeys := []string{}
	for k := range mapKeyMap {
		mapKeys = append(mapKeys, k)
	}

	for _, mapKey := range mapKeys {
		key := fmt.Sprintf("%s_%s", prefix, strings.ToUpper(mapKey))
		kv := v.MapIndex(reflect.ValueOf(mapKey))
		if kv.Kind() == reflect.Invalid {
			kv = reflect.New(v.Type().Elem().Elem())
		}
		v.SetMapIndex(reflect.ValueOf(mapKey), configStruct(key, kv))
	}
}

func configStructForPtr(prefix string, v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	v.Elem().Set(configStruct(prefix, v.Elem()))
}

func configStructForDefault(prefix string, v reflect.Value) {
	defaultValue := v.Interface()
	v.Set(reflect.ValueOf(conf(prefix, defaultValue)))
}

func conf[T any](key string, fallback T) T {
	viper.Set(key, fallback)
	if value, ok := os.LookupEnv(key); ok {
		viper.Set(key, value)
	}
	viper.UnmarshalKey(key, &fallback)
	return fallback
}

func stringToAllCapsCase(str string) string {
	allCapsBuilder := strings.Builder{}
	defer allCapsBuilder.Reset()
	allCapsBuilder.WriteByte(str[0])
	for _, c := range str[1:] {
		if c >= 'A' && c <= 'Z' {
			allCapsBuilder.WriteString("_")
		}
		allCapsBuilder.WriteRune(c)
	}
	return strings.ToUpper(allCapsBuilder.String())
}
