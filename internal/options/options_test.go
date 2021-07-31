package options

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestNew(t *testing.T) {
	type args struct {
		sort  map[string]string
		page  string
		limit string
	}
	tests := []struct {
		name string
		args args
		want Options
	}{
		{
			name: "",
			args: args{
				sort: map[string]string{
					"date": "asc",
				},
				page:  "1",
				limit: "25",
			},
			want: Options{
				sort: map[string]string{
					"date": "asc",
				},
				page:  "1",
				limit: "25",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := New(test.args.sort, test.args.page, test.args.limit); !reflect.DeepEqual(got, test.want) {
				t.Errorf("New() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestOptions_Sort(t *testing.T) {
	type args struct {
		fn convert
	}
	tests := []struct {
		name string
		o    *Options
		args args
		want bson.M
	}{
		{
			name: "createdTs desc",
			o: &Options{
				sort: map[string]string{
					"createdTs": "desc",
				},
			},
			args: args{
				fn: func(fields map[string]string) bson.M {
					sort := bson.M{}
					for key, direction := range fields {
						switch key {
						case "createdTs":
							key = "createdTs"
						default:
							key = ""
						}

						if key != "" {
							switch direction {
							case "desc":
								sort[key] = -1
							default:
								sort[key] = 1
							}
							break
						}
					}
					return sort
				},
			},
			want: map[string]interface{}{
				"createdTs": -1,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.o.Sort(test.args.fn); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Options.Sort() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestOptions_Page(t *testing.T) {
	tests := []struct {
		name string
		o    *Options
		want int
	}{
		{
			name: "Test",
			o: &Options{
				page: "1",
			},
			want: 1,
		}, {
			name: "Test",
			o: &Options{
				page: "",
			},
			want: 1,
		}, {
			name: "Test",
			o: &Options{
				page: "42",
			},
			want: 42,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.o.Page(); got != test.want {
				t.Errorf("Options.Page() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestOptions_Limit(t *testing.T) {
	tests := []struct {
		name string
		o    *Options
		want int
	}{
		{
			name: "Test",
			o: &Options{
				limit: "",
			},
			want: 25,
		}, {
			name: "Test",
			o: &Options{
				limit: "5",
			},
			want: 5,
		}, {
			name: "Test",
			o: &Options{
				limit: "50",
			},
			want: 50,
		}, {
			name: "Test",
			o: &Options{
				limit: "51",
			},
			want: 25,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.o.Limit(); got != test.want {
				t.Errorf("Options.Limit() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestOptions_Skip(t *testing.T) {
	tests := []struct {
		name string
		o    *Options
		want int
	}{
		{
			name: "Test",
			o: &Options{
				page:  "1",
				limit: "25",
			},
			want: 0,
		}, {
			name: "Test",
			o: &Options{
				page:  "2",
				limit: "25",
			},
			want: 25,
		}, {
			name: "Test",
			o: &Options{
				page:  "100",
				limit: "50",
			},
			want: 99 * 50,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.o.Skip(); got != test.want {
				t.Errorf("Options.Skip() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func TestCreatePipeline(t *testing.T) {
	type args struct {
		match bson.M
		skip  int
		limit int
		sort  bson.M
	}
	tests := []struct {
		name string
		args args
		want []bson.M
	}{
		{
			name: "Test",
			args: args{
				match: bson.M{
					"_id": primitive.NilObjectID,
				},
				skip:  25,
				limit: 100,
				sort: bson.M{
					"date": -1,
				},
			},
			want: []bson.M{
				{"$match": bson.M{"_id": primitive.NilObjectID}},
				{"$sort": bson.M{"date": -1}},
				{"$skip": 25},
				{"$limit": 100},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := CreatePipeline(test.args.match, test.args.skip, test.args.limit, test.args.sort); !reflect.DeepEqual(got, test.want) {
				t.Errorf("CreatePipeline() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}
