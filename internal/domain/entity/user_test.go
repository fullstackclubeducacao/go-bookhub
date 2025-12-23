package entity

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name         string
		userName     string
		email        string
		passwordHash string
		wantErr      error
	}{
		{
			name:         "valid user",
			userName:     "John Doe",
			email:        "john@example.com",
			passwordHash: "hashedpassword123",
			wantErr:      nil,
		},
		{
			name:         "name too short",
			userName:     "Jo",
			email:        "john@example.com",
			passwordHash: "hashedpassword123",
			wantErr:      ErrInvalidUserName,
		},
		{
			name:         "name too long",
			userName:     string(make([]byte, 101)),
			email:        "john@example.com",
			passwordHash: "hashedpassword123",
			wantErr:      ErrInvalidUserName,
		},
		{
			name:         "invalid email",
			userName:     "John Doe",
			email:        "invalid-email",
			passwordHash: "hashedpassword123",
			wantErr:      ErrInvalidUserEmail,
		},
		{
			name:         "password too short",
			userName:     "John Doe",
			email:        "john@example.com",
			passwordHash: "12345",
			wantErr:      ErrInvalidUserPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.userName, tt.email, tt.passwordHash)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}

			if user.Name != tt.userName {
				t.Errorf("NewUser() name = %v, want %v", user.Name, tt.userName)
			}
			if user.Email != tt.email {
				t.Errorf("NewUser() email = %v, want %v", user.Email, tt.email)
			}
			if !user.Active {
				t.Error("NewUser() user should be active by default")
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	user, _ := NewUser("John Doe", "john@example.com", "hashedpassword123")

	tests := []struct {
		name      string
		newName   string
		newEmail  string
		wantErr   error
		checkName string
	}{
		{
			name:      "update name only",
			newName:   "Jane Doe",
			newEmail:  "",
			wantErr:   nil,
			checkName: "Jane Doe",
		},
		{
			name:     "update email only",
			newName:  "",
			newEmail: "jane@example.com",
			wantErr:  nil,
		},
		{
			name:     "invalid name",
			newName:  "Jo",
			newEmail: "",
			wantErr:  ErrInvalidUserName,
		},
		{
			name:     "invalid email",
			newName:  "",
			newEmail: "invalid",
			wantErr:  ErrInvalidUserEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUser := *user
			err := testUser.Update(tt.newName, tt.newEmail)

			if err != tt.wantErr {
				t.Errorf("User.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_Disable(t *testing.T) {
	t.Run("disable active user", func(t *testing.T) {
		user, _ := NewUser("John Doe", "john@example.com", "hashedpassword123")

		err := user.Disable()
		if err != nil {
			t.Errorf("User.Disable() unexpected error = %v", err)
		}

		if user.Active {
			t.Error("User.Disable() user should be inactive")
		}
	})

	t.Run("disable already disabled user", func(t *testing.T) {
		user, _ := NewUser("John Doe", "john@example.com", "hashedpassword123")
		user.Active = false

		err := user.Disable()
		if err != ErrUserDisabled {
			t.Errorf("User.Disable() error = %v, wantErr %v", err, ErrUserDisabled)
		}
	})
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user.name@domain.org", true},
		{"user+tag@domain.co.uk", true},
		{"invalid", false},
		{"@domain.com", false},
		{"user@", false},
		{"user@.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.valid {
				t.Errorf("isValidEmail(%q) = %v, want %v", tt.email, got, tt.valid)
			}
		})
	}
}
