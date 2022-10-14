package domain

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type Event struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Identifier  string             `bson:"identifier,omitempty" json:"identifier"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ListID      primitive.ObjectID `bson:"list_id" json:"list_id"`
	Date        time.Time          `bson:"date" json:"date"`
	Active      bool               `bson:"active" json:"active"`
}

func (e Event) FormatDate() string {
	month := strconv.Itoa(int(e.Date.Month()))
	if e.Date.Month() < 10 {
		month = "0" + month
	}

	day := strconv.Itoa(e.Date.Day())
	if e.Date.Day() < 10 {
		day = "0" + day
	}

	hour := strconv.Itoa(e.Date.Hour())
	if e.Date.Hour() < 10 {
		hour = "0" + hour
	}

	minute := strconv.Itoa(e.Date.Minute())
	if e.Date.Minute() < 10 {
		minute = "0" + minute
	}

	return fmt.Sprintf("%s.%s  o  %s:%s", day, month, hour, minute)
}

func (e Event) Format(list List) string {
	return fmt.Sprintf("Назва: %s\nДата: %s\nКількість місць: %d\nОпис: %s", e.Name, e.FormatDate(), list.Capacity, e.Description)
}

func (e Event) Preview(list List) string {
	return fmt.Sprintf("*%s*\n\nКоли: *%s*\n\n%s\n\nКількість місць: %d\nВільних місць: %d", e.Name, e.FormatDate(), e.Description, list.Capacity, list.Capacity-len(list.List))
}

func (e Event) PreviewForProgram(list List) string {
	return fmt.Sprintf("*%s*\n%s\nВільних місць: %d\n\n", e.FormatDate(), e.Name, list.Capacity-len(list.List))
}
