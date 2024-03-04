package domain

import (
	"fmt"
	"time"
)

type Item struct {
	Id          int       `json:"id" db:"id"`
	ProjectId   int       `json:"projectId" db:"project_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Priority    int       `json:"priority" db:"priority"`
	Removed     bool      `json:"removed" db:"removed"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

type ItemCreateInp struct {
	ProjectId int    `db:"project_id"`
	Name      string `db:"name"`
}

func (i *ItemCreateInp) Validate() error {
	if i.ProjectId == 0 || i.Name == "" {
		return fmt.Errorf("invalid data")
	}

	return nil
}

type ItemUpdateInp struct {
	Id          int     `db:"id"`
	ProjectId   int     `db:"project_id"`
	Name        string  `db:"name"`
	Description *string `db:"description,omitempty"`
}

func (i *ItemUpdateInp) Validate() error {
	if i.Id == 0 || i.ProjectId == 0 || i.Name == "" {
		return fmt.Errorf("invalid data")
	}

	return nil
}

type ItemDeleteInp struct {
	Id        int `db:"id"`
	ProjectId int `db:"project_id"`
}

func (i *ItemDeleteInp) Validate() error {
	if i.Id == 0 || i.ProjectId == 0 {
		return fmt.Errorf("invalid data")
	}

	return nil
}

type ItemDeleteOut struct {
	Id        int  `json:"id" db:"id"`
	ProjectId int  `json:"project_id" db:"project_id"`
	Removed   bool `json:"removed" db:"removed"`
}

type ItemReprioritiizeInp struct {
	Id        int `db:"id"`
	ProjectId int `db:"project_id"`
	Priority  int `json:"newPriority" db:"priority"`
}

func (i *ItemReprioritiizeInp) Validate() error {
	if i.Id == 0 || i.ProjectId == 0 || i.Priority == 0 {
		return fmt.Errorf("invalid data")
	}

	return nil
}

type ItemReprioritiizeOut struct {
	Id       int `db:"id"`
	Priority int `json:"Priority"`
}

type List struct {
	Total   int     `json:"total"`
	Removed int     `json:"removed"`
	Limit   int     `json:"limit"`
	Offset  int     `json:"offset"`
	Goods   []*Item `json:"goods"`
}
