package audit_log

import (
	"encoding/json"
)

// Service bertanggung jawab membentuk audit log
// dan menyimpannya via repository
type Service interface {
	Log(
		clientID uint,
		outletID uint,
		userID uint,
		action string,
		entity string,
		entityID uint,
		oldData any,
		newData any,
	) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// ========================
// CREATE AUDIT LOG
// ========================
func (s *service) Log(
	clientID uint,
	outletID uint,
	userID uint,
	action string,
	entity string,
	entityID uint,
	oldData any,
	newData any,
) error {

	var oldJSON string
	var newJSON string

	if oldData != nil {
		if b, err := json.Marshal(oldData); err == nil {
			oldJSON = string(b)
		}
	}

	if newData != nil {
		if b, err := json.Marshal(newData); err == nil {
			newJSON = string(b)
		}
	}

	log := &AuditLog{
		ClientID: clientID,
		OutletID: outletID,
		UserID:   userID,
		Action:   action,
		Entity:   entity,
		EntityID: entityID,
		OldData:  oldJSON,
		NewData:  newJSON,
	}

	return s.repo.Create(log)
}
