package sqlgenerator

import "github.com/Masterminds/squirrel"

var PGsql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
