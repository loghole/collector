package models

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
)

func TestParam_ToSql(t *testing.T) {
	tests := []struct {
		name          string
		input         *domain.QueryParam
		expectedQuery string
		expectedArgs  []interface{}
		wantErr       bool
		expectedErr   string
	}{
		{
			name: "column1",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "=",
			},
			expectedQuery: "key=?",
			expectedArgs:  []interface{}{"value"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column2",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "!=",
			},
			expectedQuery: "key!=?",
			expectedArgs:  []interface{}{"value"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column3",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "LIKE",
			},
			expectedQuery: "key LIKE ?",
			expectedArgs:  []interface{}{"%value%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column4",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{Item: "value"},
				Operator: "NOT LIKE",
			},
			expectedQuery: "key NOT LIKE ?",
			expectedArgs:  []interface{}{"%value%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column5",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "LIKE",
			},
			expectedQuery: "key LIKE ? AND key LIKE ? AND key LIKE ?",
			expectedArgs:  []interface{}{"%value1%", "%value2%", "%value3%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column6",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "NOT LIKE",
			},
			expectedQuery: "key NOT LIKE ? AND key NOT LIKE ? AND key NOT LIKE ?",
			expectedArgs:  []interface{}{"%value1%", "%value2%", "%value3%"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column7",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "=",
			},
			expectedQuery: "key IN (?,?,?)",
			expectedArgs:  []interface{}{"value1", "value2", "value3"},
			wantErr:       false,
			expectedErr:   "",
		},
		{
			name: "column8",
			input: &domain.QueryParam{
				Type:     "column",
				Key:      "key",
				Value:    domain.ParamValue{List: []string{"value1", "value2", "value3"}},
				Operator: "!=",
			},
			expectedQuery: "key NOT IN (?,?,?)",
			expectedArgs:  []interface{}{"value1", "value2", "value3"},
			wantErr:       false,
			expectedErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, args, err := ColumnParamFromDomain(tt.input).ToSql()
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
