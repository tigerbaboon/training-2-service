package base

import (
	ci18n "app/config/i18n"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spf13/viper"
)

// Regexp definitions
var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
var wordBarrierRegex = regexp.MustCompile(`([a-z_0-9])([A-Z])`)

type conventionalMarshallerFromPascal struct {
	Value any
}

func (c conventionalMarshallerFromPascal) MarshalJSON() ([]byte, error) {
	marshalled, err := json.Marshal(c.Value)
	if err != nil {
		return nil, err
	}
	naming := viper.GetString("HTTP_JSON_NAMING")

	val := reflect.TypeOf(c.Value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		field, ok := val.FieldByName("json")
		if ok {
			if field.Tag.Get("naming") != "" {
				naming = field.Tag.Get("naming")
			}
		}
	}

	var converted []byte
	switch naming {
	case "snake_case":

		// https://gist.github.com/Rican7/39a3dc10c1499384ca91
		converted = keyMatchRegex.ReplaceAllFunc(
			marshalled,
			func(match []byte) []byte {
				return bytes.ToLower(wordBarrierRegex.ReplaceAll(
					match,
					[]byte(`${1}_${2}`),
				))
			},
		)
	case "camel_case":
		// https://gist.github.com/piersy/b9934790a8892db1a603820c0c23e4a7
		converted = keyMatchRegex.ReplaceAllFunc(
			marshalled,
			func(match []byte) []byte {
				// Empty keys are valid JSON, only lowercase if we do not have an
				// empty key.
				if len(match) > 2 {
					// check uppercase
					j := 0
					for i, v := range match {
						if unicode.IsUpper(rune(v)) {
							j = i + 1
							break
						}
					}
					// Decode first rune after the double quotes
					for i := 1; i <= j; i++ {
						r, width := utf8.DecodeRune(match[i:])
						r = unicode.ToLower(r)
						utf8.EncodeRune(match[i:width+i], r)
					}
				}
				return match
			},
		)
	case "pascal_case":
		return marshalled, nil
	default:
		return nil, err
	}

	// return converted, nil
	return converted, nil
}

func defaultJSON(ctx *gin.Context, code int, msgID string, data any, paginate *ResponsePaginate, params ...map[string]string) error {
	localizer := i18n.NewLocalizer(ci18n.Bundle, ctx.GetHeader("Accept-Language"))
	var param any
	if len(params) > 0 {
		param = params[0]
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: msgID, TemplateData: param})
	if err != nil || msg == "" {
		ctx.JSON(code, Response[any]{
			ResponseStatus: &ResponseStatus{
				Message: msgID,
				Code:    fmt.Sprintf("%d", code),
			},
			Data:     data,
			Paginate: paginate,
		})
		return nil
	}
	ctx.JSON(code, conventionalMarshallerFromPascal{Response[any]{
		ResponseStatus: &ResponseStatus{
			Message: msg,
			Code:    msgID,
		},
		Data:     data,
		Paginate: paginate,
	}})
	return nil
}

func JSON(ctx *gin.Context, code int, data any) error {
	ctx.JSON(code, conventionalMarshallerFromPascal{data})
	return nil
}

// Success 200 success
func Success(ctx *gin.Context, data any) error {
	return defaultJSON(ctx, http.StatusOK, "200", data, nil)
}

// Paginate 200 success
func Paginate(ctx *gin.Context, data any, page *ResponsePaginate) error {
	return defaultJSON(ctx, http.StatusOK, "200", data, page)
}

// BadRequest 400 other and external error
func BadRequest(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusBadRequest, message, data, nil)
}

// Unauthorized 401 un authentication
func Unauthorized(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusUnauthorized, message, data, nil)
}

// Forbidden 403 No permission
func Forbidden(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusForbidden, message, data, nil)
}

// ValidateFailed 412 Validate error
func ValidateFailed(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusPreconditionFailed, message, data, nil)
}

// InternalServerError 500 internal error
func InternalServerError(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusInternalServerError, message, data, nil)
}

// NotImplemented 501 not implemented
func NotImplemented(ctx *gin.Context, message string, data any, params ...map[string]string) error {
	return defaultJSON(ctx, http.StatusNotImplemented, message, data, nil)
}
