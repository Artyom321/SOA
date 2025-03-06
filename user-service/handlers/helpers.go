package handlers

import "social-network/common/models"

func UpdateUserFields(user *models.User, input models.ProfileUpdateRequest) bool {
	updated := false

	if input.Name != nil && *input.Name != user.Name {
		user.Name = *input.Name
		updated = true
	}
	if input.Surname != nil && *input.Surname != user.Surname {
		user.Surname = *input.Surname
		updated = true
	}
	if input.Email != nil && *input.Email != user.Email {
		user.Email = *input.Email
		updated = true
	}
	if input.PhoneNumber != nil && *input.PhoneNumber != user.PhoneNumber {
		user.PhoneNumber = *input.PhoneNumber
		updated = true
	}
	if input.BirthDate != nil && (user.BirthDate == nil || *input.BirthDate != *user.BirthDate) {
		user.BirthDate = input.BirthDate
		updated = true
	}

	return updated
}
