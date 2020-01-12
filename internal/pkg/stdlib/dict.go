package stdlib

type Dict interface {
	Has(key string) bool
	Get(key string) interface{}
	Put(key string, value interface{})
	Delete(key string)
	Plus(key string, value interface{}) Dict
	Minus(key string) Dict
	Values() map[string]interface{}
	Clone() Dict
}

type Elem struct {
	Key   string
	Value interface{}
}

type dict struct {
	values map[string]interface{}
}

func NewDict(elem ...Elem) Dict {
	values := make(map[string]interface{})
	for _, p := range elem {
		values[p.Key] = p.Value
	}
	return &dict{values: values}
}

func Pair(key string, value interface{}) Elem {
	return Elem{key, value}
}

func (d *dict) Has(key string) bool {
	_, b := d.values[key]
	return b
}

func (d *dict) Get(key string) interface{} {
	return d.values[key]
}

func (d *dict) Put(key string, value interface{}) {
	d.values[key] = value
}

func (d *dict) Delete(key string) {
	delete(d.values, key)
}

func (d *dict) Plus(key string, value interface{}) Dict {
	x := d.Clone()
	x.Put(key, value)
	return x
}

func (d *dict) Minus(key string) Dict {
	x := d.Clone()
	x.Delete(key)
	return x
}

func (d *dict) Values() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range d.values {
		m[k] = v
	}
	return m
}

func (d *dict) Clone() Dict {
	x := NewDict()
	for k, v := range d.values {
		x.Put(k, v)
	}
	return x
}
