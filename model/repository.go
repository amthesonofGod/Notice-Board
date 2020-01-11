package model

import "github.com/amthesonofGod/Notice-Board/entity"

// UserRepository specifies user database operations
type UserRepository interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	StoreSession(session *entity.UserSession) (*entity.UserSession, []error)
	Session(uuid string) (*entity.UserSession, []error)
	DeleteSession(uuid string) (*entity.UserSession, []error)
}

// CompanyRepository specifies company database operations
type CompanyRepository interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company ,[]error)
	StoreSession(session *entity.CompanySession) (*entity.CompanySession, []error)
	Session(uuid string) (*entity.CompanySession, []error)
	DeleteSession(uuid string) (*entity.CompanySession, []error)
}
