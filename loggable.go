package loggable

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"reflect"
	"time"
)

type ChangeLog struct {
	ID           string    `gorm:"type:uuid;primary_key;"`
	CreatedAt    time.Time `sql:"DEFAULT:current_timestamp"`
	ChangedBy    string    `gorm:"index"`
	ChangedWhere string    `gorm:"index"`
	Action       string
	ObjectID     string `gorm:"index"`
	ObjectType   string `gorm:"index"`
	Object       JSONB  `sql:"type:JSONB"`
}

type LoggableInterface interface {
	SetEnabled(v bool)
	Enabled() bool
}

type LoggableModel struct {
	Disabled bool
}

func (m *LoggableModel) SetEnabled(v bool) {
	m.Disabled = !v
}

func (m *LoggableModel) Enabled() bool {
	return !m.Disabled
}

type JSONB []byte

func (j JSONB) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not string")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSONB) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSONB) Equals(j1 JSONB) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

func RecursiveSetLoggableEnabled(v interface{}, value bool) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}

	recursiveSetLoggableEnabledAdv(val, value)
	return nil
}

func recursiveSetLoggableEnabledAdv(val reflect.Value, value bool) {
	switch val.Kind() {
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if val.CanAddr() {
				recursiveSetLoggableEnabledAdv(val.Field(i).Addr(), value)
			} else {
				recursiveSetLoggableEnabledAdv(val.Field(i), value)
			}
		}
	case reflect.Ptr:
		if !val.IsNil() {
			if val.CanInterface() {
				ll, ok := val.Interface().(LoggableInterface)
				if ok {
					ll.SetEnabled(value)
				}
			}

			recursiveSetLoggableEnabledAdv(val.Elem(), value)
		}
	case reflect.Slice:
		for j := 0; j < val.Len(); j++ {
			if val.Index(j).CanAddr() {
				recursiveSetLoggableEnabledAdv(val.Index(j).Addr(), value)
			}
		}
	}
}
