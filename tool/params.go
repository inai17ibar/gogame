package tool

type DatabaseInfo struct {
	url      string
	port     string
	username string
	password string
	db_name  string
}

func (url DatabaseInfo) GetDBUrl() string {
	return "root:root@tcp(0.0.0.0:3306)/gogame_db?parseTime=true&loc=Asia%2FTokyo"
}
