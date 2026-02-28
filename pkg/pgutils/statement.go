package pgutils

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type builder func(reflect.Value, int) (string, []any, error)

var cache sync.Map

func getBuilder(t reflect.Type) builder {
	if cached, ok := cache.Load(t); ok {
		return cached.(builder)
	}

	numField := t.NumField()

	type fieldMeta struct {
		index int
		tag   string
	}

	fields := make([]fieldMeta, 0, numField)

	for i := range numField {
		sf := t.Field(i)

		if sf.Type.Kind() != reflect.Pointer {
			continue
		}

		tag := sf.Tag.Get("db")
		if tag == "" {
			continue
		}

		fields = append(fields, fieldMeta{
			index: i,
			tag:   tag,
		})
	}

	var builder builder = func(v reflect.Value, startId int) (string, []any, error) {
		var sb strings.Builder
		argv := make([]any, 0, len(fields))

		first := true

		for _, f := range fields {
			field := v.Field(f.index)
			if field.IsNil() {
				continue
			}

			if !first {
				sb.WriteString(", ")
			}

			first = false

			sb.WriteString(f.tag)
			sb.WriteString(" = $")
			sb.WriteString(strconv.Itoa(startId))

			argv = append(argv, field.Elem().Interface())
			startId++
		}

		return sb.String(), argv, nil
	}

	cache.Store(t, builder)
	return builder
}

func CreateSetStatement(payload any, startId int) (string, []any, error) {
	if startId < 1 {
		startId = 1
	}

	v := reflect.ValueOf(payload)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil, ErrInvalidType
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", nil, ErrInvalidType
	}

	builder := getBuilder(v.Type())
	return builder(v, startId)
}
