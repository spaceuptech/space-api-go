package mgo

import (
	"github.com/spaceuptech/space-api-go/api/utils"
)

// GenerateFind generates a mongo db find clause from the provided condition
func GenerateFind(condition utils.M) utils.M {
	m := utils.M{}
	switch condition["type"].(string) {
	case "and":
		conds := condition["conds"].([]utils.M)
		for _, c := range conds {
			t := GenerateFind(c)
			for k, v := range t {
				m[k] = v
			}
		}

	case "or":
		conds := condition["conds"].([]utils.M)
		t := []utils.M{}
		for _, c := range conds {
			t = append(t, GenerateFind(c))
		}
		m["$or"] = t

	case "cond":
		f1 := condition["f1"].(string)
		eval := condition["eval"].(string)
		f2 := condition["f2"]

		switch eval {
		case "==":
			m[f1] = map[string]interface{}{"$eq": f2}
		case "!=":
			m[f1] = map[string]interface{}{"$ne": f2}
		case ">":
			m[f1] = map[string]interface{}{"$gt": f2}
		case "<":
			m[f1] = map[string]interface{}{"$lt": f2}
		case ">=":
			m[f1] = map[string]interface{}{"$gte": f2}
		case "<=":
			m[f1] = map[string]interface{}{"$lte": f2}
		case "in":
			m[f1] = map[string]interface{}{"$in": f2}
		case "notIn":
			m[f1] = map[string]interface{}{"$nin": f2}
		}
	}

	return m
}
