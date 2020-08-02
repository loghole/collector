package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
)

func TestParamJson_ToSql(t *testing.T) {
	tests := []struct {
		name          string
		input         *domain.QueryParam
		expectedQuery string
		expectedArgs  []interface{}
		wantErr       bool
		expectedErr   string
	}{
		{
			name: "json1",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "=",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] = ?)",
			expectedArgs:  []interface{}{"key", "key", "value"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json2",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "!=",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] != ?)",
			expectedArgs:  []interface{}{"key", "key", "value"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json3",
			input: &domain.QueryParam{
				Type:     "key",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "=",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] IN (?,?,?))",
			expectedArgs:  []interface{}{"key", "key", "value1", "value2", "value3"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json4",
			input: &domain.QueryParam{
				Type:     "key",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "!=",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] NOT IN (?,?,?))",
			expectedArgs:  []interface{}{"key", "key", "value1", "value2", "value3"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json5",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "LIKE",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] LIKE ?)",
			expectedArgs:  []interface{}{"key", "key", "%value%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json6",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "NOT LIKE",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] NOT LIKE ?)",
			expectedArgs:  []interface{}{"key", "key", "%value%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json7",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "LIKE",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] LIKE ?) AND (has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] LIKE ?) AND (has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] LIKE ?)",
			expectedArgs:  []interface{}{"key", "key", "%value1%", "key", "key", "%value2%", "key", "key", "%value3%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json8",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "NOT LIKE",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] NOT LIKE ?) AND (has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] NOT LIKE ?) AND (has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] NOT LIKE ?)",
			expectedArgs:  []interface{}{"key", "key", "%value1%", "key", "key", "%value2%", "key", "key", "%value3%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json9",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value1"},
				Operator: "<=",
			},
			expectedQuery: "(has(params_string.keys, ?) AND params_string.values[indexOf(params_string.keys, ?)] <= ?)",
			expectedArgs:  []interface{}{"key", "key", "value1"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "json10",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "10"},
				Operator: "<=",
			},
			expectedQuery: "(has(params_float.keys, ?) AND params_float.values[indexOf(params_float.keys, ?)] <= ?)",
			expectedArgs:  []interface{}{"key", "key", (float64)(10)},
			wantErr:       false,
			expectedErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, err := JSONParamFromDomain(tt.input).ToSql()
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}

			assert.Equal(t, tt.expectedQuery, query)
			assert.Equal(t, tt.expectedArgs, args)

			if tt.wantErr {
				assert.EqualError(t, err, tt.expectedErr)
			}
		})
	}
}
