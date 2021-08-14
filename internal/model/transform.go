package model

import (
	"github.com/dustyrat/go-webapp/pkg/model"
)

// TransformFromDTO dto model -> database model
func TransformFromDTO(dto model.Document) Document {
	return Document{
		ID: dto.ID,
		Audit: Audit{
			CreatedBy: &User{
				FirstName: dto.Audit.CreatedBy.FirstName,
				LastName:  dto.Audit.CreatedBy.LastName,
				Username:  dto.Audit.CreatedBy.Username,
			},
			CreatedTs: dto.Audit.CreatedTs,
			UpdatedBy: User{
				FirstName: dto.Audit.UpdatedBy.FirstName,
				LastName:  dto.Audit.UpdatedBy.LastName,
				Username:  dto.Audit.UpdatedBy.Username,
			},
			UpdatedTs: dto.Audit.UpdatedTs,
			Version:   dto.Audit.Version,
		},
	}
}

// TransformToDTO database model -> dto model
func TransformToDTO(detail Document) model.Document {
	return model.Document{
		ID: detail.ID,
		Audit: model.Audit{
			CreatedBy: func() model.User {
				if detail.Audit.CreatedBy != nil {
					return model.User{
						FirstName: detail.Audit.CreatedBy.FirstName,
						LastName:  detail.Audit.CreatedBy.LastName,
						Username:  detail.Audit.CreatedBy.Username,
					}
				}
				return model.User{}
			}(),
			CreatedTs: detail.Audit.CreatedTs,
			UpdatedBy: model.User{
				FirstName: detail.Audit.UpdatedBy.FirstName,
				LastName:  detail.Audit.UpdatedBy.LastName,
				Username:  detail.Audit.UpdatedBy.Username,
			},
			UpdatedTs: detail.Audit.UpdatedTs,
			Version:   detail.Audit.Version,
		},
	}
}
