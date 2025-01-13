package middleware

import (
	"app/app/modules"
	activitylogsent "app/app/modules/activitylogs/ent"
	"app/app/modules/base"
	helper "app/helper/helpuser"
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LocalOrigin  = "LC-Origin"
	LocalCountry = "LC-COUNTRY"
	LocalCFRay   = "LC-CP-RAY"
	LocalIP      = "LC-IP"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type LogResponseInfo struct {
	Method    string
	Path      string
	IP        string
	UserAgent string
	Header    any
	Query     any
	Request   string
	Response  string
}

func NewLogResponse() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newLog := new(LogResponseInfo)
		newLog.Method = ctx.Request.Method
		newLog.UserAgent = ctx.Request.UserAgent()
		newLog.Path = ctx.FullPath()
		newLog.IP = ctx.ClientIP()
		newLog.Header = ctx.Request.Header
		newLog.Query = ctx.Request.URL.Query()

		// Set Header value
		ctx.Set(LocalOrigin, GetHeader(ctx, `Origin`))
		ctx.Set(LocalCountry, GetHeader(ctx, `CF-IPCountry`))
		ctx.Set(LocalCFRay, GetHeader(ctx, `CF-RAY`))
		ctx.Set(LocalIP, GetHeader(ctx, `CF-Connecting-IP`))

		// GET Data Body
		body, err := io.ReadAll(ctx.Request.Body)
		if errors.Is(err, io.EOF) {

		} else if err != nil {
			base.InternalServerError(ctx, err.Error(), nil)
			ctx.Abort()
			return
		} else {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			newLog.Request = string(body)
		}

		// Set struct Resposne
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw

		// Next Process
		ctx.Next()

		// Get Response Body
		resBody := blw.body.String()
		newLog.Response = string(resBody)

		statusCode := ctx.Writer.Status()
		// Check 404 and redirect to 403
		if statusCode == 404 {
			base.InternalServerError(ctx, "404 Not Found", nil)
			ctx.Abort()
			return
		}

		// Save to Activity Log
		token := ctx.GetHeader("Authorization")
		userId, userType, err := helper.GetUserByToken(ctx, token)
		if err != nil {
			userId = "unknown"
			userType = "unknown"
		}

		modules := modules.Get()

		logTemp := activitylogsent.ActivityLogs{
			Section:       newLog.Path,
			EventType:     newLog.Method,
			StatusCode:    fmt.Sprint(statusCode),
			Detail:        fmt.Sprint(newLog.Query),
			Request:       newLog.Request,
			Responses:     newLog.Response,
			IpAddress:     newLog.IP,
			UserAgent:     newLog.UserAgent,
			CreatedBy:     userId,
			CreatedByType: userType,
			CreatedAt:     time.Now().Unix(),
		}

		_, err = modules.Acticitylogs.Svc.CreateLogs(ctx, logTemp)
		if err != nil {
			base.InternalServerError(ctx, err.Error(), nil)
			ctx.Abort()
			return
		}
	}
}

func GetHeader(ctx *gin.Context, key string) string {
	val, ok := ctx.Get(LocalIP)
	if !ok {
		return `not-found`
	}
	return val.(string)
}
