package output

import (
	"strconv"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// NewVoidType func
func NewVoidType() *graphql.Scalar {
	return graphql.NewScalar(graphql.ScalarConfig{
		Name: "void",
		ParseValue: func(value interface{}) interface{} {
			return nil
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			return nil
		},
		Serialize: func(value interface{}) interface{} {
			return nil
		},
	})
}

func coerceInt64(value interface{}) interface{} {
	switch value := value.(type) {
	case int64:
		return value
	case *int64:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}

var Int64 = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Int64",
	Description: "Int64",
	Serialize:   coerceInt64,
	ParseValue:  coerceInt64,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return intValue
			}
		}
		return nil
	},
})

var BigInt = graphql.NewScalar(graphql.ScalarConfig{
	Name: "BigInt",
	Description: "The `BigInt` scalar type represents non-fractional signed whole numeric " +
		"values. BigInt represent values correspond to golang Int ",
	Serialize:  coerceBigInt,
	ParseValue: coerceBigInt,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.Atoi(valueAST.Value); err == nil {
				return intValue
			}
		}
		return nil
	},
})

func coerceBigInt(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		if value {
			return 1
		}
		return 0
	case *bool:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case int:
		return value
	case *int:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case int8:
		return int(value)
	case *int8:
		if value == nil {
			return nil
		}
		return int(*value)
	case int16:
		return int(value)
	case *int16:
		if value == nil {
			return nil
		}
		return int(*value)
	case int32:
		return int(value)
	case *int32:
		if value == nil {
			return nil
		}
		return int(*value)
	case int64:
		return int(value)
	case *int64:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case uint:
		return int(value)
	case *uint:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case uint8:
		return int(value)
	case *uint8:
		if value == nil {
			return nil
		}
		return int(*value)
	case uint16:
		return int(value)
	case *uint16:
		if value == nil {
			return nil
		}
		return int(*value)
	case uint32:
		return int(value)
	case *uint32:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case uint64:
		return int(value)
	case *uint64:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case float32:
		return int(value)
	case *float32:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case float64:
		return int(value)
	case *float64:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	case string:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil
		}
		return coerceBigInt(val)
	case *string:
		if value == nil {
			return nil
		}
		return coerceBigInt(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}
