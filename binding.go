package MyGin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"github.com/go-playground/validator/v10"
)

const (
	CONTENT_TYPE = "Content-Type"

	APPLICATION_JSON = "application/json"
	APPLICATION_FORM = "application/x-www-form-urlencoded"
	APPLICATION_XML = "application/xml"
)
var once sync.Once

//判断标签
func Validate(obj interface{}) error {
	validate:=validator.New()
	once.Do(func() {
		validate.SetTagName("binding")
	})
	return validate.Struct(obj)
}

func Bind(obj interface{},c *Context) error {

	values :=c.Querys()
	Binddata(obj,values,"")



	//body
	content_type := c.Request.Header.Get(CONTENT_TYPE)
	switch content_type {
	case APPLICATION_JSON:
		return BindJSON(obj,c)
	case APPLICATION_FORM:
		params, err := c.ParseForm()
		if err != nil {
			return err
		}
		if err = Binddata(obj, params, "form"); err != nil {
			return err
		}

	case APPLICATION_XML:
		return Bindxml(obj,c)
	default:
		return errors.New("contenttype error")
	}
	return nil
}

func Binddata(ptr interface{},data map[string][]string,tag string) error {
	if ptr == nil || len(data) == 0 {
		return nil
	}
	typ := reflect.TypeOf(ptr).Elem() //type
	val := reflect.ValueOf(ptr).Elem() //value

	// Map
	if typ.Kind() == reflect.Map {
		for k, v := range data {
			val.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}
		return nil
	}

	//!struct
	if typ.Kind() != reflect.Struct {
		return errors.New("binding element must be a struct")
	}

	//struct
	for i:=0;i<typ.NumField();i++ {
		typefield :=typ.Field(i)
		structField :=val.Field(i)
		if !structField.CanSet() {
			continue
		}
		structFieldKind := structField.Kind()
		inputname:=typefield.Tag.Get(tag)
		if inputname == "" {
			inputname = typefield.Name
			if structFieldKind == reflect.Struct {
				if err:= Binddata(structField.Addr().Interface(), data, tag) ;err!=nil{
					return err
				}
				continue
			}
		}
		inputValue, exists := data[inputname]
		if !exists {

			for k, v := range data {
				if strings.EqualFold(k, inputname) {
					inputValue = v
					exists = true
					break
				}
			}
		}

		if !exists {
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for j := 0; j < numElems; j++ {
				if err := setWithType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else if err := setWithType(typefield.Type.Kind(), inputValue[0], structField); err != nil {
			return err

		}




	}

	return nil
}

func setWithType(valueKind reflect.Kind,val string,structField reflect.Value)error  {
	switch valueKind {
	case reflect.Ptr:
		return setWithType(structField.Elem().Kind(), val, structField.Elem())
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}


func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}


func BindJSON(obj interface{},c *Context) error {
	decoder:=json.NewDecoder(c.Request.Body)
	if err:=decoder.Decode(obj);err!=nil{

		return err
	}
	if err:= Validate(obj);err!=nil{

		return err
	}
	return nil
}


func Bindxml(obj interface{},c *Context) error {
	decoder:=xml.NewDecoder(c.Request.Body)
	if err:=decoder.Decode(obj);err!=nil{
		return err
	}
	if err:= Validate(obj);err!=nil{
		return err
	}
	return nil
}
