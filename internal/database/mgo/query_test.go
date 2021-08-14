package mgo

import (
	"net/url"
	"testing"
	"time"

	"github.com/dustyrat/go-webapp/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestGetSort(t *testing.T) {
	type args struct {
		fields map[string]string
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "Default",
			args: args{
				fields: map[string]string{
					"createdTs": "NA",
				},
			},
			want: map[string]interface{}{
				"createdTs": 1,
			},
		},
		{
			name: "createdTs acs",
			args: args{
				fields: map[string]string{
					"createdTs": "asc",
				},
			},
			want: map[string]interface{}{
				"createdTs": 1,
			},
		}, {
			name: "updatedTs desc",
			args: args{
				fields: map[string]string{
					"updatedTs": "desc",
				},
			},
			want: map[string]interface{}{
				"updatedTs": -1,
			},
		}, {
			name: "NA",
			args: args{
				fields: map[string]string{
					"NA": "asdfg",
				},
			},
			want: map[string]interface{}{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := GetSort(test.args.fields); !cmp.Equal(test.want, got) {
				t.Errorf("query.go GetSort() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		query url.Values
	}

	type want struct {
		Filter   bson.M
		Errs     []error
		Warnings []error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Multi Value Warning",
			args: args{
				query: map[string][]string{
					"page": {"1", "2"},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs:   []error{},
				Warnings: []error{
					errors.New("param 'page' used 2 times (limit 1);"),
				},
			},
		},
		{
			name: "No Value Warning",
			args: args{
				query: map[string][]string{
					"id": {""},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs:   []error{},
				Warnings: []error{
					errors.New("no value supplied for 'id='[0]"),
				},
			},
		},
		{
			name: "No Key Error",
			args: args{
				query: map[string][]string{
					"": {"value"},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs: []error{
					errors.New("no key supplied for '=value'"),
				},
				Warnings: []error{},
			},
		},
		{
			name: "Invalid Parameter Error",
			args: args{
				query: map[string][]string{
					"key": {"value"},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs: []error{
					errors.New("'key=value' not processable: 'key' is not a valid query parameter;"),
				},
				Warnings: []error{},
			},
		},
		{
			name: "id",
			args: args{
				query: map[string][]string{
					"id": {"000000000000000000000001"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{"_id": bson.M{"$in": []interface{}{utils.PrimitiveObjectID("000000000000000000000001")}}},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "createdOn",
			args: args{
				query: map[string][]string{
					"createdOn": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"createdTs": bson.M{
										"$gte": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date
										}(),
										"$lt": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date.Add(24 * time.Hour)
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "createdOn (Not a Date)",
			args: args{
				query: map[string][]string{
					"createdOn": {"NaN"},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs: []error{
					errors.New("invalid input: 'createdOn=NaN'; date must be in the format of 'YYYY-MM-DD' (createdOn=2006-01-02)"),
				},
				Warnings: []error{},
			},
		},
		{
			name: "createdAfter",
			args: args{
				query: map[string][]string{
					"createdAfter": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"createdTs": bson.M{
										"$gte": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date.Add(24 * time.Hour)
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "createdBefore",
			args: args{
				query: map[string][]string{
					"createdBefore": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"createdTs": bson.M{
										"$lt": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "updatedOn",
			args: args{
				query: map[string][]string{
					"updatedOn": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"updatedTs": bson.M{
										"$gte": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date
										}(),
										"$lt": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date.Add(24 * time.Hour)
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "updatedOn (Not a Date)",
			args: args{
				query: map[string][]string{
					"updatedOn": {"NaN"},
				},
			},
			want: want{
				Filter: map[string]interface{}{},
				Errs: []error{
					errors.New("invalid input: 'updatedOn=NaN'; date must be in the format of 'YYYY-MM-DD' (updatedOn=2006-01-02)"),
				},
				Warnings: []error{},
			},
		},
		{
			name: "updatedAfter",
			args: args{
				query: map[string][]string{
					"updatedAfter": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"updatedTs": bson.M{
										"$gte": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date.Add(24 * time.Hour)
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
		{
			name: "updatedBefore",
			args: args{
				query: map[string][]string{
					"updatedBefore": {"2020-01-01"},
				},
			},
			want: want{
				Filter: map[string]interface{}{
					"$and": []bson.M{
						{
							"$or": []bson.M{
								{
									"updatedTs": bson.M{
										"$lt": func() time.Time {
											date, _ := time.Parse("2006-01-02", "2020-01-01")
											return date
										}(),
									},
								},
							},
						},
					},
				},
				Errs:     []error{},
				Warnings: []error{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filter, errs, warnings := ParseQuery(test.args.query)
			got := want{
				Filter:   filter,
				Errs:     errs,
				Warnings: warnings,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("query.go ParseQuery() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func Test_parseDateRange(t *testing.T) {
	type args struct {
		key  string
		date time.Time
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "AfterDate",
			args: args{
				key: "dateAfter",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$gte": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date.Add(24 * time.Hour)
					}(),
				},
			},
		}, {
			name: "BeforeDate",
			args: args{
				key: "dateBefore",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$lt": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date
					}(),
				},
			},
		}, {
			name: "OnDate",
			args: args{
				key: "dateOn",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$gte": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date
					}(),
					"$lt": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date.Add(24 * time.Hour)
					}(),
				},
			},
		}, {
			name: "Default",
			args: args{
				key: "Default",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"Default": bson.M{
					"$gte": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date
					}(),
					"$lt": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date.Add(24 * time.Hour)
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDateRange(tt.args.key, tt.args.date); !cmp.Equal(tt.want, got) {
				t.Errorf("query.go parseDateRange() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}

func Test_getDateRange(t *testing.T) {
	type args struct {
		op    string
		field string
		date  time.Time
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		{
			name: "AfterDate",
			args: args{
				op:    "After",
				field: "date",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$gte": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date.Add(24 * time.Hour)
					}(),
				},
			},
		}, {
			name: "BeforeDate",
			args: args{
				op:    "Before",
				field: "date",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$lt": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date
					}(),
				},
			},
		}, {
			name: "OnDate",
			args: args{
				op:    "On",
				field: "date",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"date": bson.M{
					"$gte": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date
					}(),
					"$lt": func() time.Time {
						date, _ := time.Parse("2006-01-02", "2020-01-01")
						return date.Add(24 * time.Hour)
					}(),
				},
			},
		}, {
			name: "Default",
			args: args{
				op:    "Default",
				field: "Default",
				date: func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
			want: map[string]interface{}{
				"Default": func() time.Time {
					date, _ := time.Parse("2006-01-02", "2020-01-01")
					return date
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDateRange(tt.args.op, tt.args.field, tt.args.date); !cmp.Equal(tt.want, got) {
				t.Errorf("query.go getDateRange() mismatch (-want +got):\n%s", cmp.Diff(tt.want, got))
			}
		})
	}
}
