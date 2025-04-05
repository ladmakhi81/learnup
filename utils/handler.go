package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ladmakhi81/learnup/types"
	"log"
	"net/http"
	"time"
)

type Handler func(*gin.Context) (*types.ApiResponse, error)

func JsonHandler(fn Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, err := fn(ctx)
		if err != nil {
			errorHandler(ctx, err)
			return
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func errorHandler(ctx *gin.Context, err error) {
	timestamp := time.Now().Unix()
	traceId := uuid.New().String()

	if serverErr, ok := err.(*types.ServerError); ok {
		internalServerErrorHandler(
			ctx,
			serverErr,
			timestamp,
			traceId,
		)
		return
	}

	if clientErr, ok := err.(*types.ClientError); ok {
		clientErrorHandler(
			ctx,
			clientErr,
			timestamp,
			traceId,
		)
		return
	}

	unknownErrorHandler(
		ctx,
		err,
		timestamp,
		traceId,
	)
}

func generateErrorLogMessage(
	timestamp int64,
	url string,
	location string,
	message string,
	traceId string,
) string {
	return fmt.Sprintf(
		"Timestamp: %v, URL: %s, Location: %s, Message: %s, TraceID: %s",
		timestamp,
		url,
		location,
		message,
		traceId,
	)
}

func internalServerErrorHandler(
	ctx *gin.Context,
	error *types.ServerError,
	timestamp int64,
	traceId string,
) {
	logErr := SaveMessageIntoLog(
		"error",
		generateErrorLogMessage(
			timestamp,
			ctx.Request.RequestURI,
			error.Location,
			error.Message,
			traceId,
		),
	)

	if logErr != nil {
		log.Printf("Unable to store error in log file ( internal server error handler ): %s\n", logErr)
	}

	ctx.JSON(
		http.StatusInternalServerError,
		types.NewApiError(
			http.StatusInternalServerError,
			"Internal Server Error",
			timestamp,
			traceId,
		))
}

func clientErrorHandler(
	ctx *gin.Context,
	clientErr *types.ClientError,
	timestamp int64,
	traceId string,
) {
	var errorMessage any
	if clientErr.Message == "" {
		errorMessage = clientErr.Metadata
	} else {
		errorMessage = clientErr.Message
	}
	ctx.JSON(
		clientErr.StatusCode,
		types.NewApiError(
			clientErr.StatusCode,
			errorMessage,
			timestamp,
			traceId,
		),
	)
}

func unknownErrorHandler(
	ctx *gin.Context,
	err error,
	timestamp int64,
	traceId string,
) {
	logErr := SaveMessageIntoLog(
		"error",
		generateErrorLogMessage(
			timestamp,
			ctx.Request.RequestURI,
			"Unknown",
			err.Error(),
			traceId,
		),
	)

	if logErr != nil {
		log.Printf("Unable to store error in log file ( unknown error handler ): %s\n", logErr)
	}

	ctx.JSON(
		http.StatusInternalServerError,
		types.NewApiError(
			http.StatusInternalServerError,
			"Internal Server Error",
			timestamp,
			traceId,
		),
	)
}
