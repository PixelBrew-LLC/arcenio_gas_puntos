package postgres

// RoleModel es el modelo de base de datos para roles
type RoleModel struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null"`
}

func (RoleModel) TableName() string {
	return "roles"
}
