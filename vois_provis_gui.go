// Author: Denis Komnenovic
package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// VoipProvisGUI provides methods for generating a textual form from a model, product, and brand.
type VoipProvisGUI struct {
	Base string
}

// NewVoipProvisGUI initializes the class with a base directory path.
func NewVoipProvisGUI(base string) *VoipProvisGUI {
	return &VoipProvisGUI{Base: base}
}

// GenerateTextualForm generates a textual form based on the provided data.
func (v *VoipProvisGUI) GenerateTextualForm(model, product, brand string, hideArray []string) map[string]map[string]map[string]string {
	templateArray := v.GenerateCompleteArray(model, product, brand)
	html := make(map[string]map[string]map[string]string)

	for category, subs := range templateArray["data"] {
		for subcategories, its := range subs {
			for kitems, items := range its {
				if !stringInSlice(kitems, hideArray) {
					if strings.HasPrefix(kitems, "option|") {
						html[category] = make(map[string]map[string]string)
						html[category][subcategories] = make(map[string]string)
						html[category][subcategories][kitems] = v.Convert2HTML(kitems, items[0])
					}

					if strings.HasPrefix(kitems, "loop|") {
						for loopKey, loopData := range items {
							key := fmt.Sprintf("%s|%s", kitems, loopKey)
							html[category] = make(map[string]map[string]string)
							html[category][subcategories] = make(map[string]string)
							html[category][subcategories][key] = v.Convert2HTML(key, loopData)
						}
					}

					if strings.HasPrefix(kitems, "lineloop|") {
						for loopKey, loopData := range items {
							split := strings.Split(kitems, "|")
							line := split[1]
							key := fmt.Sprintf("lineloop|%s|%s", line, loopKey)
							html[category] = make(map[string]map[string]string)
							html[category][subcategories] = make(map[string]string)
							html[category][subcategories][key] = v.Convert2HTML(key, loopData)
						}
					}

					if strings.HasPrefix(kitems, "break") {
						html[category] = make(map[string]map[string]string)
						html[category][subcategories] = make(map[string]string)
						html[category][subcategories][kitems] = "<br />"
					}
				}
			}
		}
	}

	return html
}

// GenerateCompleteArray generates a structured data array based on the model, product, and brand.
func (v *VoipProvisGUI) GenerateCompleteArray(model, product, brand string) map[string]interface{} {
	data := make(map[string]interface{})
	fdJSON := v.File2JSON(fmt.Sprintf("%s/%s/%s/family_data.json", v.Base, brand, product))
	modelLocation := v.ArraySearchRecursive(model, fdJSON, "model")

	if modelLocation == nil {
		panic("can't find model")
	}

	modelInformation := fdJSON["data"].(map[string]interface{})["model_list"].([]interface{})[modelLocation[2]].(map[string]interface{})

	data["phone_data"] = map[string]interface{}{
		"brand":   brand,
		"product": product,
		"model":   model,
	}
	data["lines"] = modelInformation["lines"]
	files := modelInformation["template_data"].([]interface{})
	files = append([]interface{}{"/../../global_template_data.json"}, files...)
	b := 0

	for _, filesData := range files {
		filePath := fmt.Sprintf("%s/%s/%s/%s", v.Base, brand, product, filesData)
		if fileExists(filePath) {
			templateData := v.File2JSON(filePath)
			templateDataMap := templateData.(map[string]interface{})
			categories := templateDataMap["template_data"].(map[string]interface{})["category"].([]interface{})

			for _, category := range categories {
				categoryMap := category.(map[string]interface{})
				categoryName := categoryMap["name"].(string)
				subcategories := categoryMap["subcategory"].([]interface{})

				for _, subcategory := range subcategories {
					subcategoryMap := subcategory.(map[string]interface{})
					subcategoryName := subcategoryMap["name"].(string)
					items := subcategoryMap["item"].([]interface{})
					itemsFin := make([]interface{}, 0)
					itemsLoop := make(map[string]interface{})

					for _, item := range items {
						itemMap := item.(map[string]interface{})
						switch itemMap["type"] {
						case "loop_line_options":
							for i := 1; i <= int(modelInformation["lines"].(float64)); i++ {
								varNam := fmt.Sprintf("lineloop|line_%d", i)
								itemsLoopMap, exists := itemsLoop[varNam].(map[string]interface{})

								if !exists {
									itemsLoopMap = make(map[string]interface{})
								}

								for _, itemLoop := range itemMap["data"].(map[string]interface{})["item"].([]interface{}) {
									if itemLoop.(map[string]interface{})["type"].(string) != "break" {
										z := strings.TrimPrefix(itemLoop.(map[string]interface{})["variable"].(string), "$")
										itemLoopMap, exists := itemsLoopMap[z].(map[string]interface{})

										if !exists {
											itemLoopMap = make(map[string]interface{})
										}

										itemLoopMap["description"] = strings.Replace(itemLoopMap["description"].(string), "{$count}", fmt.Sprintf("%d", i), -1)
										itemLoopMap["default_value"] = strings.Replace(itemLoopMap["default_value"].(string), "{$count}", fmt.Sprintf("%d", i), -1)
										itemLoopMap["line_loop"] = true
										itemLoopMap["line_count"] = i
										itemsLoopMap[z] = itemLoopMap
									}
								}

								if len(itemsLoopMap) == 0 {
									itemsLoopMap = append(itemsLoopMap, map[string]interface{}{"type": "break"})
								}

								itemsLoop[varNam] = itemsLoopMap
							}

							if len(itemsLoop) > 0 {
								itemsFin = append(itemsFin, itemsLoop)
							}
						case "loop":
							loopStart := int(itemMap["loop_start"].(float64))
							loopEnd := int(itemMap["loop_end"].(float64))

							for i := loopStart; i <= loopEnd; i++ {
								nameSplit := strings.Split(itemMap["data"].(map[string]interface{})["item"].([]interface{})[0].(map[string]interface{})["variable"].(string), "_")
								varNam := fmt.Sprintf("loop|%s_%d", strings.TrimPrefix(nameSplit[0], "$"), i)
								itemsLoopMap, exists := itemsLoop[varNam].(map[string]interface{})

								if !exists {
									itemsLoopMap = make(map[string]interface{})
								}

								for _, itemLoop := range itemMap["data"].(map[string]interface{})["item"].([]interface{}) {
									if itemLoop.(map[string]interface{})["type"].(string) != "break" {
										zTmp := strings.Split(itemLoop.(map[string]interface{})["variable"].(string), "_")
										z := zTmp[1]
										itemLoopMap, exists := itemsLoopMap[z].(map[string]interface{})

										if !exists {
											itemLoopMap = make(map[string]interface{})
										}

										itemLoopMap["description"] = strings.Replace(itemLoopMap["description"].(string), "{$count}", fmt.Sprintf("%d", i), -1)
										itemMap["variable"] = strings.Replace(itemMap["variable"].(string), "_", fmt.Sprintf("_%d_", i), -1)
										itemLoopMap["default_value"] = ""
										itemLoopMap["loop"] = true
										itemLoopMap["loop_count"] = i
										itemsLoopMap[z] = itemLoopMap
									}
								}

								itemsFin = append(itemsFin, itemsLoopMap)
							}
						case "break":
							itemsFin = append(itemsFin, "break")
						default:
							varNam := fmt.Sprintf("option|%s", strings.TrimPrefix(itemMap["variable"].(string), "$"))
							itemsFinMap, exists := itemsFin.(map[string]interface{})

							if !exists {
								itemsFinMap = make(map[string]interface{})
							}

							itemsFinMap[varNam] = append(itemsFinMap[varNam].([]interface{}), item)
							itemsFin = itemsFinMap
						}
					}

					if categoryMap, exists := data[categoryName].(map[string]interface{}); exists {
						if subcategoryMap, exists := categoryMap[subcategoryName].(map[string]interface{}); exists {
							data[categoryName] = categoryMap
							data[categoryName][subcategoryName] = append(data[categoryName][subcategoryName].([]interface{}), itemsFin)
						} else {
							subcategoryMap := make(map[string]interface{})
							subcategoryMap[categoryName] = subcategoryMap
							subcategoryMap[categoryName][subcategoryName] = itemsFin
						}
					} else {
						categoryMap := make(map[string]interface{})
						categoryMap[categoryName] = subcategoryMap
						categoryMap[categoryName][subcategoryName] = itemsFin
					}

					if data[categoryName] != nil {
						if oldC, exists := data[categoryName].(map[string]interface{}); exists {
							if newC, exists := subcategoryMap.(map[string]interface{}); exists {
								data[categoryName] = oldC
								data[categoryName] = newC
							} else {
								data[categoryName] = subcategoryMap
							}
						} else {
							data[categoryName] = subcategoryMap
						}
					}
				}
			}
		}
	}

	return data
}

// ArraySearchRecursive performs a recursive search in an array.
func (v *VoipProvisGUI) ArraySearchRecursive(needle interface{}, haystack map[string]interface{}, needleKey string, strict bool) []string {
	path := make([]string, 0)
	if haystack == nil {
		return path
	}

	if haystackArray, ok := haystack[needleKey].([]interface{}); ok {
		for key, val := range haystackArray {
			valMap := val.(map[string]interface{})
			subPath := v.ArraySearchRecursive(needle, valMap, needleKey, strict)

			if len(subPath) > 0 {
				path = append(path, fmt.Sprintf("%d", key))
				path = append(path, subPath...)
				return path
			}
		}
	}

	if !strict && haystack[needleKey] == needle {
		path = append(path, needleKey)
	}

	return path
}

// Convert2HTML converts structured data into HTML form elements.
func (v *VoipProvisGUI) Convert2HTML(key string, data map[string]interface{}) string {
	htmlReturn := ""

	switch data["type"] {
	case "input":
		value := data["value"].(string)
		if value == "" {
			value = data["default_value"].(string)
		}
		htmlReturn = fmt.Sprintf("%s: <input type='text' name='%s' value='%s'/><br />", data["description"], key, value)
	case "break":
		htmlReturn = "<br/>"
	case "list":
		htmlReturn = fmt.Sprintf("%s: <select name='%s'>", data["description"], key)
		value := data["value"].(string)
		if value == "" {
			value = data["default_value"].(string)
		}
		listData := data["data"].([]interface{})
		for _, list := range listData {
			listMap := list.(map[string]interface{})
			selected := ""
			if value == listMap["value"] {
				selected = "selected"
			}
			htmlReturn += fmt.Sprintf("<option value='%s' %s>%s</option>", listMap["value"], selected, listMap["text"])
		}
		htmlReturn += "</select><br />"
	case "radio":
		htmlReturn = fmt.Sprintf("%s:", data["description"])
		radioData := data["data"].([]interface{})
		for _, list := range radioData {
			listMap := list.(map[string]interface{})
			checked := ""
			if listMap["checked"] {
				checked = "checked"
			}
			htmlReturn += fmt.Sprintf("|<input type='radio' name='%s' value='%s' %s/>%s", key, key, checked, listMap["description"])
		}
		htmlReturn += "<br />"
	case "checkbox":
		value := data["value"].(string)
		if value == "" {
			value = data["default_value"].(string)
		}
		checked := ""
		if value != "" {
			checked = "checked"
		}
		htmlReturn = fmt.Sprintf("%s: <input type='checkbox' name='%s' %s/><br />", data["description"], key, checked)
	}

	return htmlReturn
}

// File2JSON reads a JSON file and returns its contents as a map[string]interface{}.
func (v *VoipProvisGUI) File2JSON(file string) map[string]interface{} {
	data := readJSONFile(file)
	return data
}

func main() {
	// Initialize the class with the base directory path.
	voipProvis := NewVoipProvisGUI("/path/to/base_directory")

	// Generate a textual form, and provide the model, product, brand, and an array of items to hide.
	hideArray := []string{"hide_item1", "hide_item2"}
	formData := voipProvis.GenerateTextualForm("model_name", "product_name", "brand_name", hideArray)

	// Print the generated HTML form.
	fmt.Printf("%+v\n", formData)
}

func readJSONFile(file string) map[string]interface{} {
	// Simulate reading a JSON file and converting it to a map[string]interface{}.
	// In your actual code, you should read the file and parse it into a map.
	return map[string]interface{}{}
}

func fileExists(filename string) bool {
	// Simulate checking if a file exists. Replace with your actual code.
	return true
}

func stringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
