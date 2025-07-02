package database

func NullClause(condition string) string {
	var approvedClause string
	if condition == "0" {
		approvedClause = " IS NULL"
	} else if condition == "1" {
		approvedClause = " IS NOT NULL"
	}
	return approvedClause
}

func NullOrClause(val string) string {
	var approvedClause string
	if val == "0" {
		approvedClause = " IS NULL"
	} else {
		approvedClause = val
	}
	return approvedClause
}
