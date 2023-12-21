package grpcbase

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// (ref.) [Status codes and their use in gRPC](https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md)
var gRPCCodeToHTTPStatus = map[codes.Code]int{
	codes.OK:                 http.StatusOK,                  // 200
	codes.Canceled:           http.StatusBadRequest,          // 400, (ref.) [What is the Correct HTTP Status Code for a Cancelled Request](https://stackoverflow.com/questions/46234679/what-is-the-correct-http-status-code-for-a-cancelled-request)
	codes.Unknown:            http.StatusInternalServerError, // 500
	codes.InvalidArgument:    http.StatusBadRequest,          // 400
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,      // 504
	codes.NotFound:           http.StatusNotFound,            // 404
	codes.AlreadyExists:      http.StatusConflict,            // 409
	codes.PermissionDenied:   http.StatusForbidden,           // 403
	codes.ResourceExhausted:  http.StatusTooManyRequests,     // 429
	codes.FailedPrecondition: http.StatusBadRequest,          // 400
	codes.Aborted:            http.StatusConflict,            // 409
	codes.OutOfRange:         http.StatusBadRequest,          // 400
	codes.Unimplemented:      http.StatusNotImplemented,      // 501
	codes.Internal:           http.StatusInternalServerError, // 500
	codes.Unavailable:        http.StatusServiceUnavailable,  // 503
	codes.DataLoss:           http.StatusInternalServerError, // 500
	codes.Unauthenticated:    http.StatusUnauthorized,        // 401
}

func HTTPStatus(status *status.Status) int {
	return gRPCCodeToHTTPStatus[codes.Code(status.Code())]
}
