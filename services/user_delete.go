package services

// DeleteUser handles user deletion
func (s *UserService) DeleteUser(id uint) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
