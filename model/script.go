package model

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	intSlice []int
	strSlice []string
)

type Script struct {
	ScriptName  string `json:"script_name" gorm:"column:script_name"`
	ScriptIntro string `json:"script_intro" gorm:"column:script_intro"`
	//type json型
	ScriptTag             intSlice `json:"script_tag" gorm:"column:script_tag" type:"JsonArray"`
	ScriptScore           float64  `json:"script_score" gorm:"column:script_score"`
	GroupDuration         int      `json:"group_duration" gorm:"column:group_duration"`
	ScriptCoverUrl        string   `json:"script_cover_url" gorm:"column:script_cover_url"`
	ScriptTextContext     string   `json:"script_text_context" gorm:"column:script_text_context"`
	ScriptPlotScore       float64  `json:"script_plot_score" gorm:"column:script_plot_score"`
	ScriptImageContent    strSlice `json:"script_image_content" gorm:"column:script_image_content"`
	ScriptMalePlayer      int      `json:"script_male_player" gorm:"column:script_male_player"`
	ScriptFemalePlayer    int      `json:"script_female_player" gorm:"column:script_female_player"`
	ScriptDifficultDegree string   `json:"script_difficult_degree" gorm:"column:script_difficult_degree"`
	ScriptPlayerLimit     int      `json:"script_player_limit" gorm:"column:script_player_limit"`
	Uuid                  string   `json:"uuid" gorm:"column:uuid" valid:"no_empty"`
	ScriptComplexScore    float64  `json:"script_complex_score" gorm:"column:script_complex_score"`
	Qid                   string   `json:"qid" gorm:"column:qid"`
	IsDel                 int      `json:"-" gorm:"column:is_del"`
	Tags                  strSlice `json:"script_tags" gorm:"-"`
}

func (m *Script) TableName() string {
	return "scripts"
}

func (p intSlice) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *intSlice) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}

func (p strSlice) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *strSlice) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), &p)
}
