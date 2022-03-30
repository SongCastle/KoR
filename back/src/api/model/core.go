package model

import (
	"errors"
	"reflect"

	"github.com/SongCastle/KoR/volume/db"
	"github.com/jinzhu/gorm"
)

var (
	invalidQueryExecution error = errors.New("NotPersisted")
	notPersisted error = errors.New("NotPersisted")
)

type queryFunc = func(*gorm.DB) *gorm.DB

// TODO: 同一のトランザクション内でコネクションを共有したい
func executeQueries(queries ...queryFunc) error {
	return db.Connect(func(d *gorm.DB) error {
		switch len := len(queries); len {
			case 0:
				return invalidQueryExecution
			default:
				for _, query := range queries[:len - 1] {
					d = query(d)
				}
				return queries[len - 1](d).Error
		}
	})
}

// TODO: カラムの存在判定について
func SelectColumns(cols ...string) queryFunc {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(cols)
	}
}

func bindParamsToObject(params interface{}, obj interface{}) {
	vp := reflect.ValueOf(params)
	ivp := reflect.Indirect(vp)
	// 引数が nil ポインタである、または非参照型である (Elem 取得不可)
	if !ivp.IsValid() || vp == ivp {
		return
	}
	vo := reflect.ValueOf(obj)
	ivo := reflect.Indirect(vo)
	if !ivo.IsValid() || vo == ivo {
		return
	}
	rp, ro := vp.Elem(), vo.Elem()
	rpt := rp.Type()
	for i := 0; i < rpt.NumField(); i++ {
		// params のフィールド名を取得
		fn := rpt.Field(i).Name
		if v := rp.FieldByName(fn); !v.IsZero() {
			// obj に同じフィールドがあるか確認
			if v2 := ro.FieldByName(fn); v2 != (reflect.Value{}) {
				if v.Kind() == v2.Kind() {
					// 同じ型のフィールドが存在する場合、値をセットする
					v2.Set(v)
				} else {
					// 型が違う場合、ポインタの参照先を確認する
					if iv := reflect.Indirect(v); iv.IsValid() {
						// 参照先の型が同じ型である場合、値をセットする
						if iv.Kind() == v2.Kind() {
							v2.Set(iv)
						}
					}
				}
			}
		}
	}
}
