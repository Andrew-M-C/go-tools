package jsonconv
import()


// ==== SetXxx ====
func (this *JsonValue) SetString(s string, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewString(s), first, keys...)
}

func (this *JsonValue) SetBoolean(b bool, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewBool(b), first, keys...)
}

func (this *JsonValue) SetBool(b bool, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewBool(b), first, keys...)
}

func (this *JsonValue) SetNull(first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewNull(), first, keys...)
}

func (this *JsonValue) SetInt64(i int64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewInt64(i), first, keys...)
}

func (this *JsonValue) SetUint64(i uint64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewUint64(i), first, keys...)
}

func (this *JsonValue) SetInt32(i int32, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewInt32(i), first, keys...)
}

func (this *JsonValue) SetUint32(i uint32, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewUint32(i), first, keys...)
}

func (this *JsonValue) SetInt(i int, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewInt(i), first, keys...)
}

func (this *JsonValue) SetUint(i uint, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewUint(i), first, keys...)
}

func (this *JsonValue) SetFloat(f float64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Set(NewFloat(f), first, keys...)
}


// ==== AppendXxx ====
func (this *JsonValue) AppendString(s string, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewString(s), keys...)
}

func (this *JsonValue) AppendBoolean(b bool, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewBool(b), keys...)
}

func (this *JsonValue) AppendBool(b bool, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewBool(b), keys...)
}

func (this *JsonValue) AppendNull(keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewNull(), keys...)
}

func (this *JsonValue) AppendInt64(i int64, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewInt64(i), keys...)
}

func (this *JsonValue) AppendUint64(i uint64, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewUint64(i), keys...)
}

func (this *JsonValue) AppendInt32(i int32, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewInt32(i), keys...)
}

func (this *JsonValue) AppendUint32(i uint32, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewUint32(i), keys...)
}

func (this *JsonValue) AppendInt(i int, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewInt(i), keys...)
}

func (this *JsonValue) AppendUint(i uint, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewUint(i), keys...)
}

func (this *JsonValue) AppendFloat(f float64, keys ...interface{}) (*JsonValue, error) {
	return this.Append(NewFloat(f), keys...)
}


// ==== InsertXxx ====
func (this *JsonValue) InsertString(s string, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewString(s), first, keys...)
}

func (this *JsonValue) InsertBoolean(b bool, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewBool(b), first, keys...)
}

func (this *JsonValue) InsertBool(b bool, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewBool(b), first, keys...)
}

func (this *JsonValue) InsertNull(first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewNull(), first, keys...)
}

func (this *JsonValue) InsertInt64(i int64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewInt64(i), first, keys...)
}

func (this *JsonValue) InsertUint64(i uint64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewUint64(i), first, keys...)
}

func (this *JsonValue) InsertInt32(i int32, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewInt32(i), first, keys...)
}

func (this *JsonValue) InsertUint32(i uint32, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewUint32(i), first, keys...)
}

func (this *JsonValue) InsertInt(i int, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewInt(i), first, keys...)
}

func (this *JsonValue) InsertUint(i uint, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewUint(i), first, keys...)
}

func (this *JsonValue) InsertFloat(f float64, first interface{}, keys ...interface{}) (*JsonValue, error) {
	return this.Insert(NewFloat(f), first, keys...)
}


// ==== NewXxx ====
func NewString(s string) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = String
	obj.stringValue = s
	return obj
}

func NewInt64(i int64) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Number
	obj.intValue = i
	obj.floatValue = float64(i)
	obj.uintValue = uint64(i)
	if i < 0 {
		obj.mustSigned = true
	} else {
		obj.mustUnsigned = true
	}
	return obj
}

func NewUint64(i uint64) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Number
	obj.intValue = int64(i)
	obj.floatValue = float64(i)
	obj.uintValue = i
	if 0 != i & 0x1000000000000000 {
		obj.mustUnsigned = true
	}
	return obj
}

func NewInt32(i int32) *JsonValue {
	return NewInt64(int64(i))
}

func NewUint32(i uint32) *JsonValue {
	return NewUint64(uint64(i))
}

func NewInt(i int) *JsonValue {
	return NewInt64(int64(i))
}

func NewUint(i uint) *JsonValue {
	return NewUint64(uint64(i))
}

func NewFloat(f float64) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Number
	obj.intValue = int64(f)
	obj.floatValue = f
	obj.mustFloat = true
	obj.mustSigned = true
	return obj
}

func NewBool(b bool) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Boolean
	obj.boolValue = b
	return obj
}

func NewBoolean(b bool) *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Boolean
	obj.boolValue = b
	return obj
}

func NewNull() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Null
	return obj
}

func NewObject() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Object
	obj.objChildren = make(map[string]*JsonValue)
	return obj
}

func NewArray() *JsonValue {
	obj := new(JsonValue)
	obj.valueType = Array
	obj.arrChildren = make([]*JsonValue, 0, 0)
	return obj
}
