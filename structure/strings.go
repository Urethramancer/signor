package structure

// protoType returns the protobuf equivalent for a Go type.
// NOTE: Some assumptions are made. If you specify an int* type instead
// of a uint, it's assumed you expect negative numbers, and the variable-
// length sint types are used in the resulting protobuf definition.
func protoType(s string) string {
	switch s {
	case "bool", "string", "uint32", "uint64":
		return s
	case "int32":
		return "sint32"
	case "int64":
		return "sint64"
	case "float64":
		return "double"
	case "float32":
		return "float"
	case "[]byte":
		return "bytes"
	}

	return ""
}
