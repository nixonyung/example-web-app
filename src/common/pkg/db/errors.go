package db

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// (ref.) [PostgREST - Errors](https://postgrest.org/en/stable/references/errors.html)
var pgErrorToHTTPStatus = map[string]int{
	// 08* Connection Exception
	"08000": http.StatusServiceUnavailable, // 503
	"08003": http.StatusServiceUnavailable, // 503
	"08006": http.StatusServiceUnavailable, // 503
	"08001": http.StatusServiceUnavailable, // 503
	"08004": http.StatusServiceUnavailable, // 503
	"08007": http.StatusServiceUnavailable, // 503
	"08P01": http.StatusServiceUnavailable, // 503

	// 09* Triggered Action Exception
	"09000": http.StatusInternalServerError, // 500

	// 0L* Invalid Grantor
	"0L000": http.StatusForbidden, // 403
	"0LP01": http.StatusForbidden, // 403

	// 0P* Invalid Role Specification
	"0P000": http.StatusForbidden, // 403

	// 23503 foreign_key_violation
	"23503": http.StatusConflict, // 409

	// 23505 unique_violation
	"23505": http.StatusConflict, // 409

	// 25006 read_only_sql_transaction
	"25006": http.StatusMethodNotAllowed, // 405

	// 25* Invalid Transaction State
	"25000": http.StatusInternalServerError, // 500
	"25001": http.StatusInternalServerError, // 500
	"25002": http.StatusInternalServerError, // 500
	"25008": http.StatusInternalServerError, // 500
	"25003": http.StatusInternalServerError, // 500
	"25004": http.StatusInternalServerError, // 500
	"25005": http.StatusInternalServerError, // 500
	"25007": http.StatusInternalServerError, // 500
	"25P01": http.StatusInternalServerError, // 500
	"25P02": http.StatusInternalServerError, // 500
	"25P03": http.StatusInternalServerError, // 500

	// 28* Invalid Authorization Specification
	"28000": http.StatusForbidden, // 403
	"28P01": http.StatusForbidden, // 403

	// 2D* Invalid Transaction Termination
	"2D000": http.StatusInternalServerError, // 500

	// 38* External Routine Exception
	"38000": http.StatusInternalServerError, // 500
	"38001": http.StatusInternalServerError, // 500
	"38002": http.StatusInternalServerError, // 500
	"38003": http.StatusInternalServerError, // 500
	"38004": http.StatusInternalServerError, // 500

	// 39* External Routine Invocation Exception
	"39000": http.StatusInternalServerError, // 500
	"39001": http.StatusInternalServerError, // 500
	"39004": http.StatusInternalServerError, // 500
	"39P01": http.StatusInternalServerError, // 500
	"39P02": http.StatusInternalServerError, // 500
	"39P03": http.StatusInternalServerError, // 500

	// 3B* Savepoint Exception
	"3B000": http.StatusInternalServerError, // 500
	"3B001": http.StatusInternalServerError, // 500

	// 40* Transaction Rollback
	"40000": http.StatusInternalServerError, // 500
	"40002": http.StatusInternalServerError, // 500
	"40001": http.StatusInternalServerError, // 500
	"40003": http.StatusInternalServerError, // 500
	"40P01": http.StatusInternalServerError, // 500

	// 42501 insufficient_privilege
	"42501": http.StatusForbidden, // 403

	// 42P01 undefined_table
	"42P01": http.StatusNotFound, // 404

	// 42883 undefined_function
	"42883": http.StatusNotFound, // 404

	// 53* Insufficient Resources
	"53000": http.StatusServiceUnavailable, // 503
	"53100": http.StatusServiceUnavailable, // 503
	"53200": http.StatusServiceUnavailable, // 503
	"53300": http.StatusServiceUnavailable, // 503
	"53400": http.StatusServiceUnavailable, // 503

	// 54* Program Limit Exceeded
	"54000": http.StatusRequestEntityTooLarge, // 413
	"54001": http.StatusRequestEntityTooLarge, // 413
	"54011": http.StatusRequestEntityTooLarge, // 413
	"54023": http.StatusRequestEntityTooLarge, // 413

	// 55* Object Not In Prerequisite State
	"55000": http.StatusInternalServerError, // 500
	"55006": http.StatusInternalServerError, // 500
	"55P02": http.StatusInternalServerError, // 500
	"55P03": http.StatusInternalServerError, // 500
	"55P04": http.StatusInternalServerError, // 500

	// 57* Operator Intervention
	"57000": http.StatusInternalServerError, // 500
	"57014": http.StatusInternalServerError, // 500
	"57P01": http.StatusInternalServerError, // 500
	"57P02": http.StatusInternalServerError, // 500
	"57P03": http.StatusInternalServerError, // 500
	"57P04": http.StatusInternalServerError, // 500
	"57P05": http.StatusInternalServerError, // 500

	// 58* System Error (errors external to PostgreSQL itself)
	"58000": http.StatusInternalServerError, // 500
	"58030": http.StatusInternalServerError, // 500
	"58P01": http.StatusInternalServerError, // 500
	"58P02": http.StatusInternalServerError, // 500

	// F0* Configuration File Error
	"F0000": http.StatusInternalServerError, // 500
	"F0001": http.StatusInternalServerError, // 500

	// HV* Foreign Data Wrapper Error (SQL/MED)
	"HV000": http.StatusInternalServerError, // 500
	"HV005": http.StatusInternalServerError, // 500
	"HV002": http.StatusInternalServerError, // 500
	"HV010": http.StatusInternalServerError, // 500
	"HV021": http.StatusInternalServerError, // 500
	"HV024": http.StatusInternalServerError, // 500
	"HV007": http.StatusInternalServerError, // 500
	"HV008": http.StatusInternalServerError, // 500
	"HV004": http.StatusInternalServerError, // 500
	"HV006": http.StatusInternalServerError, // 500
	"HV091": http.StatusInternalServerError, // 500
	"HV00B": http.StatusInternalServerError, // 500
	"HV00C": http.StatusInternalServerError, // 500
	"HV00D": http.StatusInternalServerError, // 500
	"HV090": http.StatusInternalServerError, // 500
	"HV00A": http.StatusInternalServerError, // 500
	"HV009": http.StatusInternalServerError, // 500
	"HV014": http.StatusInternalServerError, // 500
	"HV001": http.StatusInternalServerError, // 500
	"HV00P": http.StatusInternalServerError, // 500
	"HV00J": http.StatusInternalServerError, // 500
	"HV00K": http.StatusInternalServerError, // 500
	"HV00Q": http.StatusInternalServerError, // 500
	"HV00R": http.StatusInternalServerError, // 500
	"HV00L": http.StatusInternalServerError, // 500
	"HV00M": http.StatusInternalServerError, // 500
	"HV00N": http.StatusInternalServerError, // 500

	// P0001 raise_exception
	"P0001": http.StatusBadRequest, // 400

	// P0* PL/pgSQL Error
	"P0000": http.StatusInternalServerError, // 500
	"P0002": http.StatusInternalServerError, // 500
	"P0003": http.StatusInternalServerError, // 500
	"P0004": http.StatusInternalServerError, // 500

	// XX* Internal Error
	"XX000": http.StatusInternalServerError, // 500
	"XX001": http.StatusInternalServerError, // 500
	"XX002": http.StatusInternalServerError, // 500

	// other: 400
}

// (ref.) https://github.com/go-gorm/gorm/issues/4135#issuecomment-790602200
func HTTPStatus(err *pgconn.PgError) int {
	if code, ok := pgErrorToHTTPStatus[err.Code]; ok {
		return code
	} else {
		return http.StatusBadRequest
	}
}

// (ref.) [Status codes and their use in gRPC](https://chromium.googlesource.com/external/github.com/grpc/grpc/+/refs/tags/v1.21.4-pre1/doc/statuscodes.md)
var pgErrorToGRPCCode = map[string]codes.Code{
	// 08* Connection Exception
	"08000": codes.Unavailable,
	"08003": codes.Unavailable,
	"08006": codes.Unavailable,
	"08001": codes.Unavailable,
	"08004": codes.Unavailable,
	"08007": codes.Unavailable,
	"08P01": codes.Unavailable,

	// 09* Triggered Action Exception
	"09000": codes.Internal,

	// 0L* Invalid Grantor
	"0L000": codes.PermissionDenied,
	"0LP01": codes.PermissionDenied,

	// 0P* Invalid Role Specification
	"0P000": codes.PermissionDenied,

	// 23503 foreign_key_violation
	"23503": codes.AlreadyExists,

	// 23505 unique_violation
	"23505": codes.AlreadyExists,

	// 25006 read_only_sql_transaction
	"25006": codes.InvalidArgument,

	// 25* Invalid Transaction State
	"25000": codes.Internal,
	"25001": codes.Internal,
	"25002": codes.Internal,
	"25008": codes.Internal,
	"25003": codes.Internal,
	"25004": codes.Internal,
	"25005": codes.Internal,
	"25007": codes.Internal,
	"25P01": codes.Internal,
	"25P02": codes.Internal,
	"25P03": codes.Internal,

	// 28* Invalid Authorization Specification
	"28000": codes.PermissionDenied,
	"28P01": codes.PermissionDenied,

	// 2D* Invalid Transaction Termination
	"2D000": codes.Internal,

	// 38* External Routine Exception
	"38000": codes.Internal,
	"38001": codes.Internal,
	"38002": codes.Internal,
	"38003": codes.Internal,
	"38004": codes.Internal,

	// 39* External Routine Invocation Exception
	"39000": codes.Internal,
	"39001": codes.Internal,
	"39004": codes.Internal,
	"39P01": codes.Internal,
	"39P02": codes.Internal,
	"39P03": codes.Internal,

	// 3B* Savepoint Exception
	"3B000": codes.Internal,
	"3B001": codes.Internal,

	// 40* Transaction Rollback
	"40000": codes.Internal,
	"40002": codes.Internal,
	"40001": codes.Internal,
	"40003": codes.Internal,
	"40P01": codes.Internal,

	// 42501 insufficient_privilege
	"42501": codes.PermissionDenied,

	// 42P01 undefined_table
	"42P01": codes.InvalidArgument,

	// 42883 undefined_function
	"42883": codes.InvalidArgument,

	// 53* Insufficient Resources
	"53000": codes.Unavailable,
	"53100": codes.Unavailable,
	"53200": codes.Unavailable,
	"53300": codes.Unavailable,
	"53400": codes.Unavailable,

	// 54* Program Limit Exceeded
	"54000": codes.ResourceExhausted,
	"54001": codes.ResourceExhausted,
	"54011": codes.ResourceExhausted,
	"54023": codes.ResourceExhausted,

	// 55* Object Not In Prerequisite State
	"55000": codes.Internal,
	"55006": codes.Internal,
	"55P02": codes.Internal,
	"55P03": codes.Internal,
	"55P04": codes.Internal,

	// 57* Operator Intervention
	"57000": codes.Internal,
	"57014": codes.Internal,
	"57P01": codes.Internal,
	"57P02": codes.Internal,
	"57P03": codes.Internal,
	"57P04": codes.Internal,
	"57P05": codes.Internal,

	// 58* System Error (errors external to PostgreSQL itself)
	"58000": codes.Internal,
	"58030": codes.Internal,
	"58P01": codes.Internal,
	"58P02": codes.Internal,

	// F0* Configuration File Error
	"F0000": codes.Internal,
	"F0001": codes.Internal,

	// HV* Foreign Data Wrapper Error (SQL/MED)
	"HV000": codes.Internal,
	"HV005": codes.Internal,
	"HV002": codes.Internal,
	"HV010": codes.Internal,
	"HV021": codes.Internal,
	"HV024": codes.Internal,
	"HV007": codes.Internal,
	"HV008": codes.Internal,
	"HV004": codes.Internal,
	"HV006": codes.Internal,
	"HV091": codes.Internal,
	"HV00B": codes.Internal,
	"HV00C": codes.Internal,
	"HV00D": codes.Internal,
	"HV090": codes.Internal,
	"HV00A": codes.Internal,
	"HV009": codes.Internal,
	"HV014": codes.Internal,
	"HV001": codes.Internal,
	"HV00P": codes.Internal,
	"HV00J": codes.Internal,
	"HV00K": codes.Internal,
	"HV00Q": codes.Internal,
	"HV00R": codes.Internal,
	"HV00L": codes.Internal,
	"HV00M": codes.Internal,
	"HV00N": codes.Internal,

	// P0001 raise_exception
	"P0001": codes.InvalidArgument,

	// P0* PL/pgSQL Error
	"P0000": codes.Internal,
	"P0002": codes.Internal,
	"P0003": codes.Internal,
	"P0004": codes.Internal,

	// XX* Internal Error
	"XX000": codes.Internal,
	"XX001": codes.Internal,
	"XX002": codes.Internal,

	// other: codes.InvalidArgument
}

// (ref.) [How to assert gRPC error codes client side in Go](https://stackoverflow.com/questions/52969205/how-to-assert-grpc-error-codes-client-side-in-go)
func GRPCError(err error) error {
	if err == nil {
		return nil
	} else if e := &(pgconn.PgError{}); errors.As(err, &e) {
		return status.Error(pgErrorToGRPCCode[e.Code], e.Error())
	} else {
		return status.Error(codes.Unknown, e.Error())
	}
}
