package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_GetValue(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	st.Nil(err)

	c := Default()

	// error on get
	_, ok := GetValue("")
	st.False(ok)

	_, ok = c.GetValue("notExist")
	st.False(ok)
	_, ok = c.GetValue("name.sub")
	st.False(ok)
	st.Error(c.Error())

	_, ok = c.GetValue("map1.key", false)
	st.False(ok)
	st.False(Exists("map1.key", false))

	val, ok := GetValue("map1.notExist")
	st.Nil(val)
	st.False(ok)
	st.False(Exists("map1.notExist"))

	val, ok = GetValue("notExist.sub")
	st.False(ok)
	st.Nil(val)
	st.False(Exists("notExist.sub"))

	val, ok = c.GetValue("arr1.100")
	st.Nil(val)
	st.False(ok)
	st.False(Exists("arr1.100"))

	val, ok = c.GetValue("arr1.notExist")
	st.Nil(val)
	st.False(ok)
	st.False(Exists("arr1.notExist"))

	ClearAll()
}

func TestGet(t *testing.T) {
	st := assert.New(t)

	ClearAll()
	err := LoadStrings(JSON, jsonStr)
	st.Nil(err)

	// fmt.Printf("%#v\n", Data())
	c := Default()

	st.False(c.IsEmpty())
	st.True(Exists("age"))
	st.True(Exists("map1.key"))
	st.True(Exists("arr1.1"))
	st.False(Exists("arr1.1", false))
	st.False(Exists("not-exist.sub"))
	st.False(Exists(""))
	st.False(Exists("not-exist"))

	// get value
	val := Get("age")
	st.Equal(float64(123), val)
	st.Equal("float64", fmt.Sprintf("%T", val))

	val = Get("not-exist")
	st.Nil(val)

	val = Get("name")
	st.Equal("app", val)

	// get string array
	arr := Strings("notExist")
	st.Empty(arr)

	arr = Strings("map1")
	st.Empty(arr)

	arr = Strings("arr1")
	st.Equal(`[]string{"val", "val1", "val2"}`, fmt.Sprintf("%#v", arr))

	val = String("arr1.1")
	st.Equal("val1", val)

	err = LoadStrings(JSON, `{
"iArr": [12, 34, 36],
"iMap": {"k1": 12, "k2": 34, "k3": 36}
}`)
	st.Nil(err)

	// Ints: get int arr
	iarr := Ints("name")
	st.False(Exists("name.1"))
	st.Empty(iarr)

	iarr = Ints("notExist")
	st.Empty(iarr)

	iarr = Ints("iArr")
	st.Equal(`[]int{12, 34, 36}`, fmt.Sprintf("%#v", iarr))

	iv := Int("iArr.1")
	st.Equal(34, iv)

	iv = Int("iArr.100")
	st.Equal(0, iv)

	// IntMap: get int map
	imp := IntMap("name")
	st.Empty(imp)
	imp = IntMap("notExist")
	st.Empty(imp)

	imp = IntMap("iMap")
	st.NotEmpty(imp)

	iv = Int("iMap.k2")
	st.Equal(34, iv)
	st.True(Exists("iMap.k2"))

	iv = Int("iMap.notExist")
	st.Equal(0, iv)
	st.False(Exists("iMap.notExist"))

	// set a intMap
	err = Set("intMap0", map[string]int{"a": 1, "b": 2})
	st.Nil(err)

	imp = IntMap("intMap0")
	st.NotEmpty(imp)
	st.Equal(1, imp["a"])
	st.Equal(2, Get("intMap0.b"))
	st.True(Exists("intMap0.a"))
	st.False(Exists("intMap0.c"))

	// StringMap: get string map
	smp := StringMap("map1")
	st.Equal("val1", smp["key1"])

	// like load from yaml content
	// c = New("test")
	err = c.LoadData(map[string]interface{}{
		"newIArr":    []int{2, 3},
		"newSArr":    []string{"a", "b"},
		"newIArr1":   []interface{}{12, 23},
		"newIArr2":   []interface{}{12, "abc"},
		"invalidMap": map[string]int{"k": 1},
		"yMap": map[interface{}]interface{}{
			"k0": "v0",
			"k1": 23,
		},
		"yMap1": map[interface{}]interface{}{
			"k":  "v",
			"k1": 23,
			"k2": []interface{}{23, 45},
		},
		"yMap10": map[string]interface{}{
			"k":  "v",
			"k1": 23,
			"k2": []interface{}{23, 45},
		},
		"yMap2": map[interface{}]interface{}{
			"k":  2,
			"k1": 23,
		},
		"yArr": []interface{}{23, 45, "val", map[string]interface{}{"k4": "v4"}},
	})
	st.Nil(err)

	iarr = Ints("newIArr")
	st.Equal("[2 3]", fmt.Sprintf("%v", iarr))

	iarr = Ints("newIArr1")
	st.Equal("[12 23]", fmt.Sprintf("%v", iarr))
	iarr = Ints("newIArr2")
	st.Empty(iarr)

	iv = Int("newIArr.1")
	st.True(Exists("newIArr.1"))
	st.Equal(3, iv)

	iv = Int("newIArr.200")
	st.False(Exists("newIArr.200"))
	st.Equal(0, iv)

	// invalid intMap
	imp = IntMap("yMap1")
	st.Empty(imp)

	imp = IntMap("yMap10")
	st.Empty(imp)

	imp = IntMap("yMap2")
	st.Equal(2, imp["k"])

	val = String("newSArr.1")
	st.True(Exists("newSArr.1"))
	st.Equal("b", val)

	val = String("newSArr.100")
	st.False(Exists("newSArr.100"))
	st.Equal("", val)

	smp = StringMap("invalidMap")
	st.Nil(smp)

	smp = StringMap("yMap.notExist")
	st.Nil(smp)

	smp = StringMap("yMap")
	st.True(Exists("yMap.k0"))
	st.False(Exists("yMap.k100"))
	st.Equal("v0", smp["k0"])

	iarr = Ints("yMap1.k2")
	st.Equal("[23 45]", fmt.Sprintf("%v", iarr))
}

func TestInt(t *testing.T) {
	st := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	st.True(Exists("age"))

	iv := Int("age")
	st.Equal(123, iv)

	iv = Int("name")
	st.Equal(0, iv)

	iv = Int("notExist", 34)
	st.Equal(34, iv)

	c := Default()
	iv = c.Int("age")
	st.Equal(123, iv)
	iv = c.Int("notExist")
	st.Equal(0, iv)

	ClearAll()
}

func TestInt64(t *testing.T) {
	st := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	// get int64
	iv64 := Int64("age")
	st.Equal(int64(123), iv64)

	iv64 = Int64("name")
	st.Equal(iv64, int64(0))

	iv64 = Int64("age", 34)
	st.Equal(int64(123), iv64)
	iv64 = Int64("notExist", 34)
	st.Equal(int64(34), iv64)

	c := Default()
	iv64 = c.Int64("age")
	st.Equal(int64(123), iv64)
	iv64 = c.Int64("notExist")
	st.Equal(int64(0), iv64)

	ClearAll()
}

func TestFloat(t *testing.T) {
	st := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)
	c := Default()

	// get float
	err := c.Set("flVal", 23.45)
	st.Nil(err)
	flt := c.Float("flVal")
	st.Equal(23.45, flt)

	flt = Float("name")
	st.Equal(float64(0), flt)

	flt = c.Float("notExists")
	st.Equal(float64(0), flt)

	flt = c.Float("notExists", 10)
	st.Equal(float64(10), flt)

	flt = Float("flVal", 0)
	st.Equal(23.45, flt)

	ClearAll()
}

func TestString(t *testing.T) {
	st := assert.New(t)
	ClearAll()
	_ = LoadStrings(JSON, jsonStr)

	// get string
	val := String("arr1")
	st.Equal("[val val1 val2]", val)

	str := String("notExists")
	st.Equal("", str)

	str = String("notExists", "defVal")
	st.Equal("defVal", str)

	c := Default()
	str = c.String("name")
	st.Equal("app", str)
	str = c.String("notExist")
	st.Equal("", str)

	ClearAll()
}

func TestBool(t *testing.T) {
	st := assert.New(t)
	ClearAll()
	_ = LoadSources(JSON, []byte(jsonStr))

	// get bool
	val := Get("debug")
	st.Equal(true, val)

	bv := Bool("debug")
	st.Equal(true, bv)

	bv = Bool("age")
	st.Equal(false, bv)

	bv = Bool("debug", false)
	st.Equal(true, bv)

	bv = Bool("notExist", false)
	st.Equal(false, bv)

	c := Default()
	bv = c.Bool("debug")
	st.True(bv)
	bv = c.Bool("notExist")
	st.False(bv)

	ClearAll()
}

func TestConfig_MapStructure(t *testing.T) {
	st := assert.New(t)

	cfg := Default()
	cfg.ClearAll()

	err := cfg.LoadStrings(JSON, `{
"age": 28,
"name": "inhere",
"sports": ["pingPong", "跑步"]
}`)

	st.Nil(err)

	user := &struct {
		Age    int
		Name   string
		Sports []string
	}{}
	// map all
	err = MapStruct("", user)
	st.Nil(err)

	st.Equal(28, user.Age)
	st.Equal("inhere", user.Name)
	st.Equal("pingPong", user.Sports[0])

	// map some
	err = cfg.LoadStrings(JSON, `{
"sec": {
	"key": "val",
	"age": 120,
	"tags": [12, 34]
}
}`)
	st.Nil(err)

	some := struct {
		Age  int
		Kye  string
		Tags []int
	}{}
	err = cfg.MapStruct("sec", &some)
	st.Nil(err)
	st.Equal(120, some.Age)
	st.Equal(12, some.Tags[0])

	cfg.ClearAll()
}
