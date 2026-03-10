package postgres

// SettingModel es el modelo de base de datos para configuraciones globales
type SettingModel struct {
	Key   string `gorm:"type:varchar(100);primaryKey"`
	Value string `gorm:"type:varchar(500);not null"`
}

func (SettingModel) TableName() string {
	return "settings"
}
