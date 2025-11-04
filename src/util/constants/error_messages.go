package constants

type UserMessage string

const (
	// Generic / unknown
	DEFAULT    UserMessage = ""
	UNKNOWN    UserMessage = "An unexpected error occurred."
	UNKNOWN_DB UserMessage = "A database error occurred."

	// CRUD / presence
	NOT_FOUND      UserMessage = "No entry found."
	ALREADY_EXISTS UserMessage = "This record already exists."
	DUPLICATE      UserMessage = "This record already exists."

	// Integrity constraints
	FOREIGN_KEY     UserMessage = "This record is linked to another and cannot be deleted."
	MISSING_FIELD   UserMessage = "A required field was left blank."
	CHECK_VIOLATION UserMessage = "One or more values violate a database rule."
	TOO_LONG        UserMessage = "One or more values are too long."
	INVALID_FORMAT  UserMessage = "Invalid data format."

	// Data or validation
	OUT_OF_RANGE UserMessage = "A value is outside the allowed range."

	// System & availability
	TIMEOUT     UserMessage = "The request took too long to complete. Please try again."
	UNAVAILABLE UserMessage = "The system is temporarily unavailable. Please try again later."
	CONFLICT    UserMessage = "The resource is currently locked or being modified."
)
