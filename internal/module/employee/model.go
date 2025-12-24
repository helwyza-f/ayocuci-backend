package employee

// UserOutlet adalah skema tabel database
type UserOutlet struct {
    ID       uint   `gorm:"primaryKey"`
    UserID   uint   `gorm:"index;uniqueIndex:idx_user_outlet"`
    OutletID uint   `gorm:"index;uniqueIndex:idx_user_outlet"`
    Role     string 
    Active   bool   `gorm:"default:true"`
}

// Employee adalah DTO untuk respon ke Flutter (Hasil Join)
type Employee struct {
    ID     uint   `json:"id"`
    Email  string `json:"email"` // Diambil dari tabel users
    Role   string `json:"role"`  // Diambil dari tabel user_outlets
    Active bool   `json:"active"`
}