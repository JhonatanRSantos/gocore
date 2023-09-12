package pgdb

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrSuccessfulCompletion                            = errors.New("successful_completion")
	ErrWarning                                         = errors.New("warning")
	ErrDynamicResultSetsReturned                       = errors.New("dynamic_result_sets_returned")
	ErrImplicitZeroBitPadding                          = errors.New("implicit_zero_bit_padding")
	ErrNullValueEliminatedInSetFunction                = errors.New("null_value_eliminated_in_set_function")
	ErrPrivilegeNotGranted                             = errors.New("privilege_not_granted")
	ErrPrivilegeNotRevoked                             = errors.New("privilege_not_revoked")
	ErrStringDataRightTruncation                       = errors.New("string_data_right_truncation")
	ErrDeprecatedFeature                               = errors.New("deprecated_feature")
	ErrNoData                                          = errors.New("no_data")
	ErrNoAdditionalDynamicResultSetsReturned           = errors.New("no_additional_dynamic_result_sets_returned")
	ErrSqlStatementNotYetComplete                      = errors.New("sql_statement_not_yet_complete")
	ErrConnectionException                             = errors.New("connection_exception")
	ErrConnectionDoesNotExist                          = errors.New("connection_does_not_exist")
	ErrConnectionFailure                               = errors.New("connection_failure")
	ErrSqlclientUnableToEstablishSqlConnection         = errors.New("sqlclient_unable_to_establish_sqlconnection")
	ErrSqlserverRejectedEstablishmentOfSqlConnection   = errors.New("sqlserver_rejected_establishment_of_sqlconnection")
	ErrTransactionResolutionUnknown                    = errors.New("transaction_resolution_unknown")
	ErrProtocolViolation                               = errors.New("protocol_violation")
	ErrTriggeredActionException                        = errors.New("triggered_action_exception")
	ErrFeatureNotSupported                             = errors.New("feature_not_supported")
	ErrInvalidTransactionInitiation                    = errors.New("invalid_transaction_initiation")
	ErrLocatorException                                = errors.New("locator_exception")
	ErrInvalidLocatorSpecification                     = errors.New("invalid_locator_specification")
	ErrInvalidGrantor                                  = errors.New("invalid_grantor")
	ErrInvalidGrantOperation                           = errors.New("invalid_grant_operation")
	ErrInvalidRoleSpecification                        = errors.New("invalid_role_specification")
	ErrDiagnosticsException                            = errors.New("diagnostics_exception")
	ErrStackedDiagnosticsAccessedWithoutActiveHandler  = errors.New("stacked_diagnostics_accessed_without_active_handler")
	ErrCaseNotFound                                    = errors.New("case_not_found")
	ErrCardinalityViolation                            = errors.New("cardinality_violation")
	ErrDataException                                   = errors.New("data_exception")
	ErrArraySubscript                                  = errors.New("array_subscript_error")
	ErrCharacterNotInRepertoire                        = errors.New("character_not_in_repertoire")
	ErrDatetimeFieldOverflow                           = errors.New("datetime_field_overflow")
	ErrDivisionByZero                                  = errors.New("division_by_zero")
	ErrErrorInAssignment                               = errors.New("error_in_assignment")
	ErrEscapeCharacterConflict                         = errors.New("escape_character_conflict")
	ErrIndicatorOverflow                               = errors.New("indicator_overflow")
	ErrIntervalFieldOverflow                           = errors.New("interval_field_overflow")
	ErrInvalidArgumentForLogarithm                     = errors.New("invalid_argument_for_logarithm")
	ErrInvalidArgumentForNtileFunction                 = errors.New("invalid_argument_for_ntile_function")
	ErrInvalidArgumentForNthValueFunction              = errors.New("invalid_argument_for_nth_value_function")
	ErrInvalidArgumentForPowerFunction                 = errors.New("invalid_argument_for_power_function")
	ErrInvalidArgumentForWidthBucketFunction           = errors.New("invalid_argument_for_width_bucket_function")
	ErrInvalidCharacterValueForCast                    = errors.New("invalid_character_value_for_cast")
	ErrInvalidDatetimeFormat                           = errors.New("invalid_datetime_format")
	ErrInvalidEscapeCharacter                          = errors.New("invalid_escape_character")
	ErrInvalidEscapeOctet                              = errors.New("invalid_escape_octet")
	ErrInvalidEscapeSequence                           = errors.New("invalid_escape_sequence")
	ErrNonstandardUseOfEscapeCharacter                 = errors.New("nonstandard_use_of_escape_character")
	ErrInvalidIndicatorParameterValue                  = errors.New("invalid_indicator_parameter_value")
	ErrInvalidParameterValue                           = errors.New("invalid_parameter_value")
	ErrInvalidRegularExpression                        = errors.New("invalid_regular_expression")
	ErrInvalidRowCountInLimitClause                    = errors.New("invalid_row_count_in_limit_clause")
	ErrInvalidRowCountInResultOffsetClause             = errors.New("invalid_row_count_in_result_offset_clause")
	ErrInvalidTimeZoneDisplacementValue                = errors.New("invalid_time_zone_displacement_value")
	ErrInvalidUseOfEscapeCharacter                     = errors.New("invalid_use_of_escape_character")
	ErrMostSpecificTypeMismatch                        = errors.New("most_specific_type_mismatch")
	ErrNullValueNotAllowed                             = errors.New("null_value_not_allowed")
	ErrNullValueNoIndicatorParameter                   = errors.New("null_value_no_indicator_parameter")
	ErrNumericValueOutOfRange                          = errors.New("numeric_value_out_of_range")
	ErrSequenceGeneratorLimitExceeded                  = errors.New("sequence_generator_limit_exceeded")
	ErrStringDataLengthMismatch                        = errors.New("string_data_length_mismatch")
	ErrSubstring                                       = errors.New("substring_error")
	ErrTrim                                            = errors.New("trim_error")
	ErrUnterminatedCString                             = errors.New("unterminated_c_string")
	ErrZeroLengthCharacterString                       = errors.New("zero_length_character_string")
	ErrFloatingPointException                          = errors.New("floating_point_exception")
	ErrInvalidTextRepresentation                       = errors.New("invalid_text_representation")
	ErrInvalidBinaryRepresentation                     = errors.New("invalid_binary_representation")
	ErrBadCopyFileFormat                               = errors.New("bad_copy_file_format")
	ErrUntranslatableCharacter                         = errors.New("untranslatable_character")
	ErrNotAnXmlDocument                                = errors.New("not_an_xml_document")
	ErrInvalidXmlDocument                              = errors.New("invalid_xml_document")
	ErrInvalidXmlContent                               = errors.New("invalid_xml_content")
	ErrInvalidXmlComment                               = errors.New("invalid_xml_comment")
	ErrInvalidXmlProcessingInstruction                 = errors.New("invalid_xml_processing_instruction")
	ErrIntegrityConstraintViolation                    = errors.New("integrity_constraint_violation")
	ErrRestrictViolation                               = errors.New("restrict_violation")
	ErrNotNullViolation                                = errors.New("not_null_violation")
	ErrForeignKeyViolation                             = errors.New("foreign_key_violation")
	ErrUniqueViolation                                 = errors.New("unique_violation")
	ErrCheckViolation                                  = errors.New("check_violation")
	ErrExclusionViolation                              = errors.New("exclusion_violation")
	ErrInvalidCursorState                              = errors.New("invalid_cursor_state")
	ErrInvalidTransactionState                         = errors.New("invalid_transaction_state")
	ErrActiveSqlTransaction                            = errors.New("active_sql_transaction")
	ErrBranchTransactionAlreadyActive                  = errors.New("branch_transaction_already_active")
	ErrHeldCursorRequiresSameIsolationLevel            = errors.New("held_cursor_requires_same_isolation_level")
	ErrInappropriateAccessModeForBranchTransaction     = errors.New("inappropriate_access_mode_for_branch_transaction")
	ErrInappropriateIsolationLevelForBranchTransaction = errors.New("inappropriate_isolation_level_for_branch_transaction")
	ErrNoActiveSqlTransactionForBranchTransaction      = errors.New("no_active_sql_transaction_for_branch_transaction")
	ErrReadOnlySqlTransaction                          = errors.New("read_only_sql_transaction")
	ErrSchemaAndDataStatementMixingNotSupported        = errors.New("schema_and_data_statement_mixing_not_supported")
	ErrNoActiveSqlTransaction                          = errors.New("no_active_sql_transaction")
	ErrInFailedSqlTransaction                          = errors.New("in_failed_sql_transaction")
	ErrInvalidSqlStatementName                         = errors.New("invalid_sql_statement_name")
	ErrTriggeredDataChangeViolation                    = errors.New("triggered_data_change_violation")
	ErrInvalidAuthorizationSpecification               = errors.New("invalid_authorization_specification")
	ErrInvalidPassword                                 = errors.New("invalid_password")
	ErrDependentPrivilegeDescriptorsStillExist         = errors.New("dependent_privilege_descriptors_still_exist")
	ErrDependentObjectsStillExist                      = errors.New("dependent_objects_still_exist")
	ErrInvalidTransactionTermination                   = errors.New("invalid_transaction_termination")
	ErrSqlRoutineException                             = errors.New("sql_routine_exception")
	ErrFunctionExecutedNoReturnStatement               = errors.New("function_executed_no_return_statement")
	ErrModifyingSqlDataNotPermitted                    = errors.New("modifying_sql_data_not_permitted")
	ErrProhibitedSqlStatementAttempted                 = errors.New("prohibited_sql_statement_attempted")
	ErrReadingSqlDataNotPermitted                      = errors.New("reading_sql_data_not_permitted")
	ErrInvalidCursorName                               = errors.New("invalid_cursor_name")
	ErrExternalRoutineException                        = errors.New("external_routine_exception")
	ErrContainingSqlNotPermitted                       = errors.New("containing_sql_not_permitted")
	ErrExternalRoutineInvocationException              = errors.New("external_routine_invocation_exception")
	ErrInvalidSqlstateReturned                         = errors.New("invalid_sqlstate_returned")
	ErrTriggerProtocolViolated                         = errors.New("trigger_protocol_violated")
	ErrSrfProtocolViolated                             = errors.New("srf_protocol_violated")
	ErrSavepointException                              = errors.New("savepoint_exception")
	ErrInvalidSavepointSpecification                   = errors.New("invalid_savepoint_specification")
	ErrInvalidCatalogName                              = errors.New("invalid_catalog_name")
	ErrInvalidSchemaName                               = errors.New("invalid_schema_name")
	ErrTransactionRollback                             = errors.New("transaction_rollback")
	ErrTransactionIntegrityConstraintViolation         = errors.New("transaction_integrity_constraint_violation")
	ErrSerializationFailure                            = errors.New("serialization_failure")
	ErrStatementCompletionUnknown                      = errors.New("statement_completion_unknown")
	ErrDeadlockDetected                                = errors.New("deadlock_detected")
	ErrSyntaxErrorOrAccessRuleViolation                = errors.New("syntax_error_or_access_rule_violation")
	ErrSyntax                                          = errors.New("syntax_error")
	ErrInsufficientPrivilege                           = errors.New("insufficient_privilege")
	ErrCannotCoerce                                    = errors.New("cannot_coerce")
	ErrGrouping                                        = errors.New("grouping_error")
	ErrWindowing                                       = errors.New("windowing_error")
	ErrInvalidRecursion                                = errors.New("invalid_recursion")
	ErrInvalidForeignKey                               = errors.New("invalid_foreign_key")
	ErrInvalidName                                     = errors.New("invalid_name")
	ErrNameTooLong                                     = errors.New("name_too_long")
	ErrReservedName                                    = errors.New("reserved_name")
	ErrDatatypeMismatch                                = errors.New("datatype_mismatch")
	ErrIndeterminateDatatype                           = errors.New("indeterminate_datatype")
	ErrCollationMismatch                               = errors.New("collation_mismatch")
	ErrIndeterminateCollation                          = errors.New("indeterminate_collation")
	ErrWrongObjectType                                 = errors.New("wrong_object_type")
	ErrUndefinedColumn                                 = errors.New("undefined_column")
	ErrUndefinedFunction                               = errors.New("undefined_function")
	ErrUndefinedTable                                  = errors.New("undefined_table")
	ErrUndefinedParameter                              = errors.New("undefined_parameter")
	ErrUndefinedObject                                 = errors.New("undefined_object")
	ErrDuplicateColumn                                 = errors.New("duplicate_column")
	ErrDuplicateCursor                                 = errors.New("duplicate_cursor")
	ErrDuplicateDatabase                               = errors.New("duplicate_database")
	ErrDuplicateFunction                               = errors.New("duplicate_function")
	ErrDuplicatePreparedStatement                      = errors.New("duplicate_prepared_statement")
	ErrDuplicateSchema                                 = errors.New("duplicate_schema")
	ErrDuplicateTable                                  = errors.New("duplicate_table")
	ErrDuplicateAlias                                  = errors.New("duplicate_alias")
	ErrDuplicateObject                                 = errors.New("duplicate_object")
	ErrAmbiguousColumn                                 = errors.New("ambiguous_column")
	ErrAmbiguousFunction                               = errors.New("ambiguous_function")
	ErrAmbiguousParameter                              = errors.New("ambiguous_parameter")
	ErrAmbiguousAlias                                  = errors.New("ambiguous_alias")
	ErrInvalidColumnReference                          = errors.New("invalid_column_reference")
	ErrInvalidColumnDefinition                         = errors.New("invalid_column_definition")
	ErrInvalidCursorDefinition                         = errors.New("invalid_cursor_definition")
	ErrInvalidDatabaseDefinition                       = errors.New("invalid_database_definition")
	ErrInvalidFunctionDefinition                       = errors.New("invalid_function_definition")
	ErrInvalidPreparedStatementDefinition              = errors.New("invalid_prepared_statement_definition")
	ErrInvalidSchemaDefinition                         = errors.New("invalid_schema_definition")
	ErrInvalidTableDefinition                          = errors.New("invalid_table_definition")
	ErrInvalidObjectDefinition                         = errors.New("invalid_object_definition")
	ErrWithCheckOptionViolation                        = errors.New("with_check_option_violation")
	ErrInsufficientResources                           = errors.New("insufficient_resources")
	ErrDiskFull                                        = errors.New("disk_full")
	ErrOutOfMemory                                     = errors.New("out_of_memory")
	ErrTooManyConnections                              = errors.New("too_many_connections")
	ErrConfigurationLimitExceeded                      = errors.New("configuration_limit_exceeded")
	ErrProgramLimitExceeded                            = errors.New("program_limit_exceeded")
	ErrStatementTooComplex                             = errors.New("statement_too_complex")
	ErrTooManyColumns                                  = errors.New("too_many_columns")
	ErrTooManyArguments                                = errors.New("too_many_arguments")
	ErrObjectNotInPrerequisiteState                    = errors.New("object_not_in_prerequisite_state")
	ErrObjectInUse                                     = errors.New("object_in_use")
	ErrCantChangeRuntimeParam                          = errors.New("cant_change_runtime_param")
	ErrLockNotAvailable                                = errors.New("lock_not_available")
	ErrOperatorIntervention                            = errors.New("operator_intervention")
	ErrQueryCanceled                                   = errors.New("query_canceled")
	ErrAdminShutdown                                   = errors.New("admin_shutdown")
	ErrCrashShutdown                                   = errors.New("crash_shutdown")
	ErrCannotConnectNow                                = errors.New("cannot_connect_now")
	ErrDatabaseDropped                                 = errors.New("database_dropped")
	ErrSystem                                          = errors.New("system_error")
	ErrIO                                              = errors.New("io_error")
	ErrUndefinedFile                                   = errors.New("undefined_file")
	ErrDuplicateFile                                   = errors.New("duplicate_file")
	ErrConfigFile                                      = errors.New("config_file_error")
	ErrLockFileExists                                  = errors.New("lock_file_exists")
	ErrFdw                                             = errors.New("fdw_error")
	ErrFdwColumnNameNotFound                           = errors.New("fdw_column_name_not_found")
	ErrFdwDynamicParameterValueNeeded                  = errors.New("fdw_dynamic_parameter_value_needed")
	ErrFdwFunctionSequence                             = errors.New("fdw_function_sequence_error")
	ErrFdwInconsistentDescriptorInformation            = errors.New("fdw_inconsistent_descriptor_information")
	ErrFdwInvalidAttributeValue                        = errors.New("fdw_invalid_attribute_value")
	ErrFdwInvalidColumnName                            = errors.New("fdw_invalid_column_name")
	ErrFdwInvalidColumnNumber                          = errors.New("fdw_invalid_column_number")
	ErrFdwInvalidDataType                              = errors.New("fdw_invalid_data_type")
	ErrFdwInvalidDataTypeDescriptors                   = errors.New("fdw_invalid_data_type_descriptors")
	ErrFdwInvalidDescriptorFieldIdentifier             = errors.New("fdw_invalid_descriptor_field_identifier")
	ErrFdwInvalidHandle                                = errors.New("fdw_invalid_handle")
	ErrFdwInvalidOptionIndex                           = errors.New("fdw_invalid_option_index")
	ErrFdwInvalidOptionName                            = errors.New("fdw_invalid_option_name")
	ErrFdwInvalidStringLengthOrBufferLength            = errors.New("fdw_invalid_string_length_or_buffer_length")
	ErrFdwInvalidStringFormat                          = errors.New("fdw_invalid_string_format")
	ErrFdwInvalidUseOfNullPointer                      = errors.New("fdw_invalid_use_of_null_pointer")
	ErrFdwTooManyHandles                               = errors.New("fdw_too_many_handles")
	ErrFdwOutOfMemory                                  = errors.New("fdw_out_of_memory")
	ErrFdwNoSchemas                                    = errors.New("fdw_no_schemas")
	ErrFdwOptionNameNotFound                           = errors.New("fdw_option_name_not_found")
	ErrFdwReplyHandle                                  = errors.New("fdw_reply_handle")
	ErrFdwSchemaNotFound                               = errors.New("fdw_schema_not_found")
	ErrFdwTableNotFound                                = errors.New("fdw_table_not_found")
	ErrFdwUnableToCreateExecution                      = errors.New("fdw_unable_to_create_execution")
	ErrFdwUnableToCreateReply                          = errors.New("fdw_unable_to_create_reply")
	ErrFdwUnableToEstablishConnection                  = errors.New("fdw_unable_to_establish_connection")
	ErrPlpgsql                                         = errors.New("plpgsql_error")
	ErrRaiseException                                  = errors.New("raise_exception")
	ErrNoDataFound                                     = errors.New("no_data_found")
	ErrTooManyRows                                     = errors.New("too_many_rows")
	ErrInternal                                        = errors.New("internal_error")
	ErrDataCorrupted                                   = errors.New("data_corrupted")
	ErrIndexCorrupted                                  = errors.New("index_corrupted")
	ErrUnknown                                         = errors.New("unknown error")
)

// http://www.postgresql.org/docs/9.3/static/errcodes-appendix.html
var pgErros = map[string]error{
	// Class 00 - Successful Completion
	"00000": ErrSuccessfulCompletion,
	// Class 01 - Warning
	"01000": ErrWarning,
	"0100C": ErrDynamicResultSetsReturned,
	"01008": ErrImplicitZeroBitPadding,
	"01003": ErrNullValueEliminatedInSetFunction,
	"01007": ErrPrivilegeNotGranted,
	"01006": ErrPrivilegeNotRevoked,
	"01004": ErrStringDataRightTruncation,
	"01P01": ErrDeprecatedFeature,
	// Class 02 - No Data (this is also a warning class per the SQL standard)
	"02000": ErrNoData,
	"02001": ErrNoAdditionalDynamicResultSetsReturned,
	// Class 03 - SQL Statement Not Yet Complete
	"03000": ErrSqlStatementNotYetComplete,
	// Class 08 - Connection Exception
	"08000": ErrConnectionException,
	"08003": ErrConnectionDoesNotExist,
	"08006": ErrConnectionFailure,
	"08001": ErrSqlclientUnableToEstablishSqlConnection,
	"08004": ErrSqlserverRejectedEstablishmentOfSqlConnection,
	"08007": ErrTransactionResolutionUnknown,
	"08P01": ErrProtocolViolation,
	// Class 09 - Triggered Action Exception
	"09000": ErrTriggeredActionException,
	// Class 0A - Feature Not Supported
	"0A000": ErrFeatureNotSupported,
	// Class 0B - Invalid Transaction Initiation
	"0B000": ErrInvalidTransactionInitiation,
	// Class 0F - Locator Exception
	"0F000": ErrLocatorException,
	"0F001": ErrInvalidLocatorSpecification,
	// Class 0L - Invalid Grantor
	"0L000": ErrInvalidGrantor,
	"0LP01": ErrInvalidGrantOperation,
	// Class 0P - Invalid Role Specification
	"0P000": ErrInvalidRoleSpecification,
	// Class 0Z - Diagnostics Exception
	"0Z000": ErrDiagnosticsException,
	"0Z002": ErrStackedDiagnosticsAccessedWithoutActiveHandler,
	// Class 20 - Case Not Found
	"20000": ErrCaseNotFound,
	// Class 21 - Cardinality Violation
	"21000": ErrCardinalityViolation,
	// Class 22 - Data Exception
	"22000": ErrDataException,
	"2202E": ErrArraySubscript,
	"22021": ErrCharacterNotInRepertoire,
	"22008": ErrDatetimeFieldOverflow,
	"22012": ErrDivisionByZero,
	"22005": ErrErrorInAssignment,
	"2200B": ErrEscapeCharacterConflict,
	"22022": ErrIndicatorOverflow,
	"22015": ErrIntervalFieldOverflow,
	"2201E": ErrInvalidArgumentForLogarithm,
	"22014": ErrInvalidArgumentForNtileFunction,
	"22016": ErrInvalidArgumentForNthValueFunction,
	"2201F": ErrInvalidArgumentForPowerFunction,
	"2201G": ErrInvalidArgumentForWidthBucketFunction,
	"22018": ErrInvalidCharacterValueForCast,
	"22007": ErrInvalidDatetimeFormat,
	"22019": ErrInvalidEscapeCharacter,
	"2200D": ErrInvalidEscapeOctet,
	"22025": ErrInvalidEscapeSequence,
	"22P06": ErrNonstandardUseOfEscapeCharacter,
	"22010": ErrInvalidIndicatorParameterValue,
	"22023": ErrInvalidParameterValue,
	"2201B": ErrInvalidRegularExpression,
	"2201W": ErrInvalidRowCountInLimitClause,
	"2201X": ErrInvalidRowCountInResultOffsetClause,
	"22009": ErrInvalidTimeZoneDisplacementValue,
	"2200C": ErrInvalidUseOfEscapeCharacter,
	"2200G": ErrMostSpecificTypeMismatch,
	"22004": ErrNullValueNotAllowed,
	"22002": ErrNullValueNoIndicatorParameter,
	"22003": ErrNumericValueOutOfRange,
	"2200H": ErrSequenceGeneratorLimitExceeded,
	"22026": ErrStringDataLengthMismatch,
	"22001": ErrStringDataRightTruncation,
	"22011": ErrSubstring,
	"22027": ErrTrim,
	"22024": ErrUnterminatedCString,
	"2200F": ErrZeroLengthCharacterString,
	"22P01": ErrFloatingPointException,
	"22P02": ErrInvalidTextRepresentation,
	"22P03": ErrInvalidBinaryRepresentation,
	"22P04": ErrBadCopyFileFormat,
	"22P05": ErrUntranslatableCharacter,
	"2200L": ErrNotAnXmlDocument,
	"2200M": ErrInvalidXmlDocument,
	"2200N": ErrInvalidXmlContent,
	"2200S": ErrInvalidXmlComment,
	"2200T": ErrInvalidXmlProcessingInstruction,
	// Class 23 - Integrity Constraint Violation
	"23000": ErrIntegrityConstraintViolation,
	"23001": ErrRestrictViolation,
	"23502": ErrNotNullViolation,
	"23503": ErrForeignKeyViolation,
	"23505": ErrUniqueViolation,
	"23514": ErrCheckViolation,
	"23P01": ErrExclusionViolation,
	// Class 24 - Invalid Cursor State
	"24000": ErrInvalidCursorState,
	// Class 25 - Invalid Transaction State
	"25000": ErrInvalidTransactionState,
	"25001": ErrActiveSqlTransaction,
	"25002": ErrBranchTransactionAlreadyActive,
	"25008": ErrHeldCursorRequiresSameIsolationLevel,
	"25003": ErrInappropriateAccessModeForBranchTransaction,
	"25004": ErrInappropriateIsolationLevelForBranchTransaction,
	"25005": ErrNoActiveSqlTransactionForBranchTransaction,
	"25006": ErrReadOnlySqlTransaction,
	"25007": ErrSchemaAndDataStatementMixingNotSupported,
	"25P01": ErrNoActiveSqlTransaction,
	"25P02": ErrInFailedSqlTransaction,
	// Class 26 - Invalid SQL Statement Name
	"26000": ErrInvalidSqlStatementName,
	// Class 27 - Triggered Data Change Violation
	"27000": ErrTriggeredDataChangeViolation,
	// Class 28 - Invalid Authorization Specification
	"28000": ErrInvalidAuthorizationSpecification,
	"28P01": ErrInvalidPassword,
	// Class 2B - Dependent Privilege Descriptors Still Exist
	"2B000": ErrDependentPrivilegeDescriptorsStillExist,
	"2BP01": ErrDependentObjectsStillExist,
	// Class 2D - Invalid Transaction Termination
	"2D000": ErrInvalidTransactionTermination,
	// Class 2F - SQL Routine Exception
	"2F000": ErrSqlRoutineException,
	"2F005": ErrFunctionExecutedNoReturnStatement,
	"2F002": ErrModifyingSqlDataNotPermitted,
	"2F003": ErrProhibitedSqlStatementAttempted,
	"2F004": ErrReadingSqlDataNotPermitted,
	// Class 34 - Invalid Cursor Name
	"34000": ErrInvalidCursorName,
	// Class 38 - External Routine Exception
	"38000": ErrExternalRoutineException,
	"38001": ErrContainingSqlNotPermitted,
	"38002": ErrModifyingSqlDataNotPermitted,
	"38003": ErrProhibitedSqlStatementAttempted,
	"38004": ErrReadingSqlDataNotPermitted,
	// Class 39 - External Routine Invocation Exception
	"39000": ErrExternalRoutineInvocationException,
	"39001": ErrInvalidSqlstateReturned,
	"39004": ErrNullValueNotAllowed,
	"39P01": ErrTriggerProtocolViolated,
	"39P02": ErrSrfProtocolViolated,
	// Class 3B - Savepoint Exception
	"3B000": ErrSavepointException,
	"3B001": ErrInvalidSavepointSpecification,
	// Class 3D - Invalid Catalog Name
	"3D000": ErrInvalidCatalogName,
	// Class 3F - Invalid Schema Name
	"3F000": ErrInvalidSchemaName,
	// Class 40 - Transaction Rollback
	"40000": ErrTransactionRollback,
	"40002": ErrTransactionIntegrityConstraintViolation,
	"40001": ErrSerializationFailure,
	"40003": ErrStatementCompletionUnknown,
	"40P01": ErrDeadlockDetected,
	// Class 42 - Syntax Error or Access Rule Violation
	"42000": ErrSyntaxErrorOrAccessRuleViolation,
	"42601": ErrSyntax,
	"42501": ErrInsufficientPrivilege,
	"42846": ErrCannotCoerce,
	"42803": ErrGrouping,
	"42P20": ErrWindowing,
	"42P19": ErrInvalidRecursion,
	"42830": ErrInvalidForeignKey,
	"42602": ErrInvalidName,
	"42622": ErrNameTooLong,
	"42939": ErrReservedName,
	"42804": ErrDatatypeMismatch,
	"42P18": ErrIndeterminateDatatype,
	"42P21": ErrCollationMismatch,
	"42P22": ErrIndeterminateCollation,
	"42809": ErrWrongObjectType,
	"42703": ErrUndefinedColumn,
	"42883": ErrUndefinedFunction,
	"42P01": ErrUndefinedTable,
	"42P02": ErrUndefinedParameter,
	"42704": ErrUndefinedObject,
	"42701": ErrDuplicateColumn,
	"42P03": ErrDuplicateCursor,
	"42P04": ErrDuplicateDatabase,
	"42723": ErrDuplicateFunction,
	"42P05": ErrDuplicatePreparedStatement,
	"42P06": ErrDuplicateSchema,
	"42P07": ErrDuplicateTable,
	"42712": ErrDuplicateAlias,
	"42710": ErrDuplicateObject,
	"42702": ErrAmbiguousColumn,
	"42725": ErrAmbiguousFunction,
	"42P08": ErrAmbiguousParameter,
	"42P09": ErrAmbiguousAlias,
	"42P10": ErrInvalidColumnReference,
	"42611": ErrInvalidColumnDefinition,
	"42P11": ErrInvalidCursorDefinition,
	"42P12": ErrInvalidDatabaseDefinition,
	"42P13": ErrInvalidFunctionDefinition,
	"42P14": ErrInvalidPreparedStatementDefinition,
	"42P15": ErrInvalidSchemaDefinition,
	"42P16": ErrInvalidTableDefinition,
	"42P17": ErrInvalidObjectDefinition,
	// Class 44 - WITH CHECK OPTION Violation
	"44000": ErrWithCheckOptionViolation,
	// Class 53 - Insufficient Resources
	"53000": ErrInsufficientResources,
	"53100": ErrDiskFull,
	"53200": ErrOutOfMemory,
	"53300": ErrTooManyConnections,
	"53400": ErrConfigurationLimitExceeded,
	// Class 54 - Program Limit Exceeded
	"54000": ErrProgramLimitExceeded,
	"54001": ErrStatementTooComplex,
	"54011": ErrTooManyColumns,
	"54023": ErrTooManyArguments,
	// Class 55 - Object Not In Prerequisite State
	"55000": ErrObjectNotInPrerequisiteState,
	"55006": ErrObjectInUse,
	"55P02": ErrCantChangeRuntimeParam,
	"55P03": ErrLockNotAvailable,
	// Class 57 - Operator Intervention
	"57000": ErrOperatorIntervention,
	"57014": ErrQueryCanceled,
	"57P01": ErrAdminShutdown,
	"57P02": ErrCrashShutdown,
	"57P03": ErrCannotConnectNow,
	"57P04": ErrDatabaseDropped,
	// Class 58 - System Error (errors external to PostgreSQL itself)
	"58000": ErrSystem,
	"58030": ErrIO,
	"58P01": ErrUndefinedFile,
	"58P02": ErrDuplicateFile,
	// Class F0 - Configuration File Error
	"F0000": ErrConfigFile,
	"F0001": ErrLockFileExists,
	// Class HV - Foreign Data Wrapper Error (SQL/MED)
	"HV000": ErrFdw,
	"HV005": ErrFdwColumnNameNotFound,
	"HV002": ErrFdwDynamicParameterValueNeeded,
	"HV010": ErrFdwFunctionSequence,
	"HV021": ErrFdwInconsistentDescriptorInformation,
	"HV024": ErrFdwInvalidAttributeValue,
	"HV007": ErrFdwInvalidColumnName,
	"HV008": ErrFdwInvalidColumnNumber,
	"HV004": ErrFdwInvalidDataType,
	"HV006": ErrFdwInvalidDataTypeDescriptors,
	"HV091": ErrFdwInvalidDescriptorFieldIdentifier,
	"HV00B": ErrFdwInvalidHandle,
	"HV00C": ErrFdwInvalidOptionIndex,
	"HV00D": ErrFdwInvalidOptionName,
	"HV090": ErrFdwInvalidStringLengthOrBufferLength,
	"HV00A": ErrFdwInvalidStringFormat,
	"HV009": ErrFdwInvalidUseOfNullPointer,
	"HV014": ErrFdwTooManyHandles,
	"HV001": ErrFdwOutOfMemory,
	"HV00P": ErrFdwNoSchemas,
	"HV00J": ErrFdwOptionNameNotFound,
	"HV00K": ErrFdwReplyHandle,
	"HV00Q": ErrFdwSchemaNotFound,
	"HV00R": ErrFdwTableNotFound,
	"HV00L": ErrFdwUnableToCreateExecution,
	"HV00M": ErrFdwUnableToCreateReply,
	"HV00N": ErrFdwUnableToEstablishConnection,
	// Class P0 - PL/pgSQL Error
	"P0000": ErrPlpgsql,
	"P0001": ErrRaiseException,
	"P0002": ErrNoDataFound,
	"P0003": ErrTooManyRows,
	// Class XX - Internal Error
	"XX000": ErrInternal,
	"XX001": ErrDataCorrupted,
	"XX002": ErrIndexCorrupted,
}

func ParseRowsTo[T any](ctx context.Context, rows pgx.Rows) ([]*T, error) {
	if value, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T]); err != nil {
		if pgErr := getPgError(err); err != nil {
			return nil, pgErr
		}
		return nil, err
	} else {
		return value, err
	}
}

func ParseRowTo[T any](ctx context.Context, rows pgx.Rows) (*T, error) {
	if value, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[T]); err != nil {
		if pgErr := getPgError(err); err != nil {
			return nil, pgErr
		}
		return nil, err
	} else {
		return value, err
	}
}

// getPgError Return the PG error or nil if the error is not recognized
func getPgError(err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if e, ok := pgErros[pgErr.Code]; ok {
			return e
		}
		return ErrUnknown
	}
	return nil
}
