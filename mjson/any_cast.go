package mjson

func any2String(v interface{}, def string) string {
	if rv, ok := v.(string); ok {
		return rv
	}
	if rv, ok := v.(*string); ok {
		return *rv
	}
	return def
}

func any2Int(v interface{}, def int) int {
	out := 0
	switch num := v.(type) {

	//
	case uint8:
		out = int(num)
	case uint16:
		out = int(num)
	case uint:
		out = int(num)
	case uint32:
		out = int(num)
	case uint64:
		out = int(num)
	case *uint8:
		out = int(*num)
	case *uint16:
		out = int(*num)
	case *uint:
		out = int(*num)
	case *uint32:
		out = int(*num)
	case *uint64:
		out = int(*num)
	//
	case int8:
		out = int(num)
	case int16:
		out = int(num)
	case int:
		out = num
	case int32:
		out = int(num)
	case int64:
		out = int(num)
	case float32:
		out = int(num)
	case float64:
		out = int(num)
	case *int8:
		out = int(*num)
	case *int16:
		out = int(*num)
	case *int:
		out = *num
	case *int32:
		out = int(*num)
	case *int64:
		out = int(*num)
	case *float32:
		out = int(*num)
	case *float64:
		out = int(*num)
	default:
		return def
	}
	return out
}
