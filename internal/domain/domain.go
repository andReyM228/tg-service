package domain

import (
	"time"
)

const (
	RegistrationStepStart    = "start"
	RegistrationStepName     = "name"
	RegistrationStepSurname  = "surname"
	RegistrationStepEmail    = "email"
	RegistrationStepPhone    = "phone"
	RegistrationStepPassword = "password"
	RegistrationStepAddress  = "address"

	LoginStepStart    = "start"
	LoginStepPassword = "password"
)

type Car struct {
	ID        int
	Name      string
	Model     string
	Price     int64
	Image     string
	CreatedAt time.Time
}

type Cars struct {
	Cars []Car
}

type User struct {
	ID             int
	Name           string
	Surname        string
	Phone          string
	Email          string
	Password       string
	ChatID         int64
	AccountAddress string
	Cars           []Car
	CreatedAt      time.Time
}

type CarCharacteristics struct {
	Engine       string `json:"engine,omitempty"`
	Power        string `json:"power,omitempty"`
	Acceleration string `json:"acceleration,omitempty"`
	TopSpeed     string `json:"top_speed,omitempty"`
	FuelEconomy  string `json:"fuel_economy,omitempty"`
	Transmission string `json:"transmission,omitempty"`
	DriveType    string `json:"drive_type,omitempty"`
}

type ProcessingLoginUser struct {
	ChatID int64
	Step   string
}

type ProcessingLoginUsers []ProcessingLoginUser

type ProcessingRegistrationUser struct {
	ChatID int64
	Step   string
	User   User
}

type ProcessingRegistrationUsers []ProcessingRegistrationUser

func (r ProcessingRegistrationUsers) SetName(chatID int64, name string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.Name = name
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) SetSurname(chatID int64, surname string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.Surname = surname
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) SetPhone(chatID int64, phone string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.Phone = phone
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) SetEmail(chatID int64, email string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.Email = email
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) SetPassword(chatID int64, password string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.Password = password
			usr.User.ChatID = chatID
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) SetAddress(chatID int64, address string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.User.AccountAddress = address
			r[i].User = usr.User

			return
		}
	}
}

func (r ProcessingRegistrationUsers) UpdateRegistrationStep(chatID int64, step string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.Step = step
			r[i] = usr

			return
		}
	}
}

func (r ProcessingRegistrationUsers) IfExists(chatID int64) bool {
	for _, usr := range r {
		if usr.ChatID == chatID {
			return true
		}
	}

	return false
}

func (r *ProcessingRegistrationUsers) GetOrCreate(chatID int64) ProcessingRegistrationUser {
	for _, usr := range *r {
		if usr.ChatID == chatID {
			return usr
		}
	}

	newUser := ProcessingRegistrationUser{
		ChatID: chatID,
		Step:   RegistrationStepStart,
	}

	*r = append(*r, newUser)

	return newUser
}

func (r *ProcessingRegistrationUsers) Delete(chatID int64) {
	users := *r

	for i, usr := range users {
		if usr.ChatID == chatID {
			if i >= 0 && i < len(users) {
				*r = append(users[:i], users[i+1:]...)
			}

			return
		}
	}
}

//---------------------------------------------------------------------------------------

func (r ProcessingLoginUsers) UpdateLoginStep(chatID int64, step string) {
	for i, usr := range r {
		if usr.ChatID == chatID {
			usr.Step = step
			r[i] = usr

			return
		}
	}
}

func (r ProcessingLoginUsers) IfExists(chatID int64) bool {
	for _, usr := range r {
		if usr.ChatID == chatID {
			return true
		}
	}

	return false
}

func (r *ProcessingLoginUsers) GetOrCreate(chatID int64) ProcessingLoginUser {
	for _, usr := range *r {
		if usr.ChatID == chatID {
			return usr
		}
	}

	newUser := ProcessingLoginUser{
		ChatID: chatID,
		Step:   RegistrationStepStart,
	}

	*r = append(*r, newUser)

	return newUser
}

func (r *ProcessingLoginUsers) Delete(chatID int64) {
	users := *r

	for i, usr := range users {
		if usr.ChatID == chatID {
			if i >= 0 && i < len(users) {
				*r = append(users[:i], users[i+1:]...)
			}

			return
		}
	}
}

type ProcessingBuyUsers map[int64]int64

func (p ProcessingBuyUsers) Create(chatID, carID int64) {
	p[chatID] = carID
}

func (p ProcessingBuyUsers) Delete(chatID int64) {
	delete(p, chatID)
}

func (p ProcessingBuyUsers) IfExists(chatID int64) bool {
	_, ok := p[chatID]
	if ok {
		return true
	}

	return false
}

func (p ProcessingBuyUsers) GetCarID(chatID int64) (int64, bool) {
	carID, ok := p[chatID]
	if ok {
		return carID, true
	}

	return 0, false
}
