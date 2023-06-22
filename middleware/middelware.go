package middleware

import (
	"context"
	"errors"
	"net/http"

	"blumer-ms-refers/graph"
	"blumer-ms-refers/model"
)

const CurrentData = "currentData"

type DataCtx struct {
	UserC model.UserCtx
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHeaderValue := r.Header.Get("user")
		roleHeaderValue := r.Header.Get("role")
		ipHeaderValue := r.Header.Get("X-Forward-Ip")
		if userHeaderValue == " " || roleHeaderValue == " " || ipHeaderValue == " " {
			next.ServeHTTP(w, r)
			return
		}
		userCtx := model.UserCtx{UserID: userHeaderValue, Role: roleHeaderValue, Ip: ipHeaderValue}
		dataCtx := DataCtx{UserC: userCtx}
		ctx := context.WithValue(r.Context(), CurrentData, dataCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCurrentUserFromCTX(ctx context.Context) (*model.UserCtx, *graph.ErrorExtensionParams) {
	errNoUserInContext := errors.New("no data in context")
	if ctx.Value(CurrentData) == nil {
		return nil, &graph.ErrorExtensionParams{AppError: graph.NotFound, Reason: errNoUserInContext.Error()}
	}
	currentData, ok := ctx.Value(CurrentData).(DataCtx)
	if !ok || currentData.UserC.UserID == "" {
		return nil, &graph.ErrorExtensionParams{AppError: graph.NotFound, Reason: errNoUserInContext.Error()}
	}
	return &currentData.UserC, nil
}
