package domain

// Setting representa un parámetro de configuración global del sistema
type Setting struct {
	Key   string
	Value string
}

// Claves de configuración conocidas
const (
	SettingPointsPerGallon    = "points_per_gallon"
	SettingMinGallons         = "min_gallons"
	SettingMinRedeemPoints    = "min_redeem_points"
	SettingPointsExpiryMonths = "points_expiry_months"
)
