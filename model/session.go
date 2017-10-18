package model

type Session struct {
	Id				uint		`db:"id" json:"id"`
	SessionId		string		`db:"session_id" json:"session"`
	VpnIp			string		`db:"vpn_ip" json:"vpn"`
	Lk				string		`db:"lk" json:"lk"`
	Started			int			`db:"started" json:"started"`
}

type SessionOutput struct {
	SessionId		string
	VpnIp			string
	Lk				string
	Started			string
}