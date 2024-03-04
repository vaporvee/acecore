package main

import "encoding/json"

type ModalJsonField struct {
	Label       string `json:"label"`
	IsParagraph bool   `json:"is_paragraph"`
	Value       string `json:"value"`
	Required    bool   `json:"required"`
	MinLength   int    `json:"min_length"`
	MaxLength   int    `json:"max_length"`
}

type ModalJson struct {
	FormType string           `json:"form_type"`
	Title    string           `json:"title"`
	Form     []ModalJsonField `json:"form"`
}

func jsonStringShowModal(jsonString string, id string) {
	var modal ModalJson
	json.Unmarshal([]byte(jsonString), &modal)

}
