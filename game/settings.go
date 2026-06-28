package game

import u "gollux/utils"

type Setting struct {
	ID          int    `json:"id"`
	Game        string `json:"game"`
	Win         string `json:"win"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type SysSetting struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

func (s *Setting) Update() map[string]interface{} {
	err := DBM.Update(s)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Setting updated!")
	return response
}

func (s *SysSetting) Get() *[]SysSetting {
	var items []SysSetting
	_, err := DBM.Query(&items, `SELECT * from sys_settings`)
	if err != nil {
		panic(err)
	}

	return &items
}
func (s *Setting) Get() map[string]interface{} {
	var items []Setting
	_, err := DBM.Query(&items, `SELECT * from settings`)
	if err != nil {
		panic(err)
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Results")
	response["settings"] = items
	return response
}

func (s *Setting) Query() *Setting {
	var item Setting
	_, err := DBM.Query(&item, `SELECT * from settings where game=?`, s.Game)
	if err != nil {
		panic(err)
	}

	return &item
}
