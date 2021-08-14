package model

import (
	"testing"

	"github.com/dustyrat/go-webapp/internal/utils"
	"github.com/dustyrat/go-webapp/pkg/model"

	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func Test_TransformFromDTO(t *testing.T) {
	type args struct {
		dto model.Document
	}
	tests := []struct {
		name string
		args args
		want Document
	}{
		{
			name: "test",
			args: args{
				dto: model.Document{
					ID: utils.PPrimitiveObjectID("000000000000000000000001"),
					Audit: model.Audit{
						CreatedBy: model.User{
							FirstName: "John",
							LastName:  "Doe",
							Username:  "john.doe",
						},
						CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						UpdatedBy: model.User{
							FirstName: "Jane",
							LastName:  "Doe",
							Username:  "jane.doe",
						},
						UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						Version:   5,
					},
				},
			},
			want: Document{
				ID: utils.PPrimitiveObjectID("000000000000000000000001"),
				Audit: Audit{
					CreatedBy: &User{
						FirstName: "John",
						LastName:  "Doe",
						Username:  "john.doe",
					},
					CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					UpdatedBy: User{
						FirstName: "Jane",
						LastName:  "Doe",
						Username:  "jane.doe",
					},
					UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					Version:   5,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := TransformFromDTO(test.args.dto); !cmp.Equal(got, test.want) {
				t.Errorf("TransformFromDTO() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}

func Test_TransformToDTO(t *testing.T) {
	type args struct {
		detail Document
	}
	tests := []struct {
		name string
		args args
		want model.Document
	}{
		{
			name: "test",
			args: args{
				detail: Document{
					ID: utils.PPrimitiveObjectID("000000000000000000000001"),
					Audit: Audit{
						CreatedBy: &User{
							FirstName: "John",
							LastName:  "Doe",
							Username:  "john.doe",
						},
						CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						UpdatedBy: User{
							FirstName: "Jane",
							LastName:  "Doe",
							Username:  "jane.doe",
						},
						UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						Version:   5,
					},
				},
			},
			want: model.Document{
				ID: utils.PPrimitiveObjectID("000000000000000000000001"),
				Audit: model.Audit{
					CreatedBy: model.User{
						FirstName: "John",
						LastName:  "Doe",
						Username:  "john.doe",
					},
					CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					UpdatedBy: model.User{
						FirstName: "Jane",
						LastName:  "Doe",
						Username:  "jane.doe",
					},
					UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					Version:   5,
				},
			},
		}, {
			name: "nil CreatedBy",
			args: args{
				detail: Document{
					ID: utils.PPrimitiveObjectID("000000000000000000000001"),
					Audit: Audit{
						CreatedBy: nil,
						CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						UpdatedBy: User{
							FirstName: "Jane",
							LastName:  "Doe",
							Username:  "jane.doe",
						},
						UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
						Version:   5,
					},
				},
			},
			want: model.Document{
				ID: utils.PPrimitiveObjectID("000000000000000000000001"),
				Audit: model.Audit{
					CreatedBy: model.User{},
					CreatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					UpdatedBy: model.User{
						FirstName: "Jane",
						LastName:  "Doe",
						Username:  "jane.doe",
					},
					UpdatedTs: utils.ISODate("2006-01-02T15:04:05.000Z"),
					Version:   5,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := TransformToDTO(test.args.detail); !cmp.Equal(got, test.want) {
				t.Errorf("TransformToDTO() mismatch (-want +got):\n%s", cmp.Diff(test.want, got))
			}
		})
	}
}
