package openapi

// Component references for the shared responses declared once in the
// generated document and pointed at from every operation that can return
// them. They are named because they are INTERNAL cross-references: a typo
// produces a document that still serializes and still validates as JSON,
// but whose $ref resolves to nothing -- a failure that surfaces in a
// client generator long after the change that caused it. The most-used of
// them appears at 25 call sites.
const (
	RespRefBadRequest          = "#/components/responses/BadRequest"
	RespRefForbidden           = "#/components/responses/Forbidden"
	RespRefFound               = "#/components/responses/Found"
	RespRefInternalServerError = "#/components/responses/InternalServerError"
	RespRefNotFound            = "#/components/responses/NotFound"
	RespRefSeeOther            = "#/components/responses/SeeOther"
	RespRefUnauthorized        = "#/components/responses/Unauthorized"
)
