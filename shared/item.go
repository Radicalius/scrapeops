package shared

import "encoding/json"

type Item struct {
	dict map[string]interface{}
}

func NewItem(data map[string]interface{}) *Item {
	return &Item{
		dict: data,
	}
}

func (i *Item) GetString(key string) (*string, bool) {
	return Get[string](i.dict, key)
}

func (i *Item) GetInt(key string) (*int64, bool) {
	return Get[int64](i.dict, key)
}

func (i *Item) GetFloat(key string) (*float64, bool) {
	return Get[float64](i.dict, key)
}

func (i *Item) GetBool(key string) (*bool, bool) {
	return Get[bool](i.dict, key)
}

func (i *Item) SetString(key string, val string) {
	i.dict[key] = val
}

func (i *Item) SetInt(key string, val int64) {
	i.dict[key] = val
}

func (i *Item) SetFloat(key string, val float64) {
	i.dict[key] = val
}

func (i *Item) SetBool(key string, val bool) {
	i.dict[key] = val
}

func (i *Item) Serialize() ([]byte, error) {
	return json.Marshal(i.dict)
}

func DeserializeItem(data []byte) (*Item, error) {
	var dict map[string]interface{}
	err := json.Unmarshal(data, &dict)
	if err != nil {
		return nil, err
	}

	return NewItem(dict), err
}

func Get[T any](m map[string]interface{}, key string) (*T, bool) {
	val, exists := m[key]
	if !exists {
		return nil, false
	}

	valCast, ok := val.(T)
	if !ok {
		return nil, false
	}

	return &valCast, true
}
