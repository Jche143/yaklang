package ssa

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func IsObject(v Value) bool {
	_, ok := v.(*Make)
	return ok
}

func ToObject(v Value) *Make {
	o, _ := v.(*Make)
	return o
}

func IsField(v Value) bool {
	_, ok := v.(*Field)
	return ok
}

func ToField(v Value) *Field {
	if o, ok := v.(*Field); ok {
		return o
	} else {
		return nil
	}
}

func GetFields(u User) []*Field {
	f := make([]*Field, 0)
	for _, v := range u.GetValues() {
		if field, ok := v.(*Field); ok {
			if field.Obj == u {
				f = append(f, field)
			}
		}
	}
	return f
}

func NewMake(parentI User, typ Type, low, high, step, Len, Cap Value, block *BasicBlock) *Make {
	i := &Make{
		anInstruction: newAnInstuction(block),
		anNode:        NewNode(),
		parentI:       parentI,
		low:           low,
		high:          high,
		step:          step,
		Len:           Len,
		Cap:           Cap,
	}
	i.SetType(typ)
	//TODO: add this variable
	fixupUseChain(i)
	return i
}

func NewUpdate(address User, v Value, block *BasicBlock) *Update {
	s := &Update{
		anInstruction: newAnInstuction(block),
		Value:         v,
		Address:       address,
	}
	s.AddValue(v)
	s.AddUser(address)
	fixupUseChain(s)
	return s
}

func NewFieldOnly(key Value, obj User, block *BasicBlock) *Field {
	f := &Field{
		anInstruction: newAnInstuction(block),
		anNode:        NewNode(),
		Update:        make([]Value, 0),
		Key:           key,
		Obj:           obj,
	}
	f.AddValue(key)
	f.AddUser(obj)
	return f
}

func (b *FunctionBuilder) CreateInterfaceWithVs(keys []Value, vs []Value) *Make {
	hasKey := true
	if len(keys) == 0 {
		hasKey = false
	}
	lValueLen := NewConst(len(vs))
	ityp := NewObjectType()
	itf := b.EmitMakeWithoutType(lValueLen, lValueLen)
	for i, rv := range vs {
		var key Value
		if hasKey {
			key = keys[i]
		} else {
			key = NewConst(i)
		}
		field := b.EmitFieldMust(itf, key)
		field.SetType(rv.GetType())
		b.EmitUpdate(field, rv)
		ityp.AddField(key, rv.GetType())
	}
	ityp.Finish()
	ityp.Len = len(vs)
	itf.SetType(ityp)
	return itf
}

// --------------- `f.symbol` hanlder, read && write
func (b *FunctionBuilder) getFieldWithCreate(i User, key Value, create bool) Value {
	var ftyp Type
	var isMethod bool
	if I, ok := i.(*Make); ok {
		if I.buildField != nil {
			if v := I.buildField(key.String()); v != nil {
				return v
			}
		}
	}
	if c, ok := key.(*Const); ok && c.IsString() {
		if v := i.GetType().GetMethod(c.VarString()); v != nil {
			isMethod = true
			ftyp = v
		}
	}

	if t := i.GetType(); t.GetTypeKind() == ObjectTypeKind {
		if it, ok := t.(*ObjectType); ok {
			if t, _ := it.GetField(key); t != nil {
				ftyp = t
				isMethod = false
			}
		}
	}
	fields := GetFields(i)
	if index := slices.IndexFunc(fields, func(v *Field) bool {
		return v.Key == key
	}); index != -1 {
		return fields[index]
	}

	if parent := b.parent; parent != nil {
		// find in parent
		if field := parent.builder.ReadField(key.String()); field != nil {
			return field
		}
	}

	if create {
		field := NewFieldOnly(key, i, b.CurrentBlock)
		if ftyp != nil {
			field.SetType(ftyp)
		}
		field.isMethod = isMethod
		b.emit(field)
		fixupUseChain(field)
		return field
	} else {
		return nil
	}
}

func (b *FunctionBuilder) NewCaptureField(text string) *Field {
	f := NewFieldOnly(NewConst(text), b.symbol, b.CurrentBlock)
	f.OutCapture = true
	return f
}

func (b *FunctionBuilder) GetField(i User, key Value, create bool) *Field {
	if field, ok := b.getFieldWithCreate(i, key, create).(*Field); ok {
		return field
	}
	return nil

}
func (b *FunctionBuilder) ReadField(key string) *Field {
	return b.GetField(b.symbol, NewConst(key), false)
}
func (b *FunctionBuilder) NewField(key string) *Field {
	return b.GetField(b.symbol, NewConst(key), true)
}

func (b *FunctionBuilder) writeField(key string, v Value) {
	field := b.NewField(key)
	if field == nil {
		panic(fmt.Sprintf("writeField: %s not found", key))
	}
	if field.GetLastValue() != v {
		b.EmitUpdate(field, v)
	}
}
