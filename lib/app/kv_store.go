package app

import "fmt"

type KVStore struct {
	kvs map[string]string
}

func (sm *KVStore) Print() {
	fmt.Println(sm.kvs)
}

func NewKVStore() *KVStore {
	return &KVStore{
		kvs: make(map[string]string),
	}
}

func (sm *KVStore) GetAll() map[string]string {
	return sm.kvs
}

func (sm *KVStore) KeyAvailable(key string) bool {
	_, ok := sm.kvs[key]
	return ok
}

// Layanan get digunakan untuk mendapatkan sebuah nilai dari key yang diberikan.
// Kembalikan string kosong jika key belum ada
func (sm *KVStore) Get(key string) (string, error) {
	if value, ok := sm.kvs[key]; ok {
		return value, nil
	}
	return "", nil
}

// Layanan set digunakan untuk menetapkan nilai dengan key yang diberikan.
// Jika key sudah ada, overwrite nilai yang lama.
func (sm *KVStore) Set(key string, value string) (string, error) {
	if key == "" {
		return "Key is empty", fmt.Errorf("Key cannot be empty")
	}

	sm.kvs[key] = value
	return "OK", nil
}

// Layanan strln digunakan untuk mendapatkan panjang value dari key yang diberikan
func (sm *KVStore) Strln(key string) (int, error) {
	if value, ok := sm.kvs[key]; ok {
		return len(value), nil
	} else {
		return -1, nil
	}
}

// Layanan del digunakan untuk menghapus entry dari key yang diberikan.
// Mengembalikan nilai yang dihapus. Kembalikan string kosong jika key belum ada.
func (sm *KVStore) Delete(key string) (string, error) {
	if value, ok := sm.kvs[key]; ok {
		delete(sm.kvs, key)
		return value, nil
	} else {
		return "Key not found", fmt.Errorf("Key not found")
	}
}

// Layanan append digunakan untuk nilai dengan key yang diberikan.
// Jika key belum ada, buat key dengan nilai string kosong sebelum melakukan append.
func (sm *KVStore) Append(key string, newValue string) (string, error) {
	if key == "" {
		return "Key is empty", fmt.Errorf("Key cannot be empty")
	} else {
		if value, ok := sm.kvs[key]; ok {
			sm.kvs[key] = value + newValue
		} else {
			sm.kvs[key] = newValue
		}
		return "OK", nil
	}
}

// func main() {
// 	sm := KVStore{
// 		Map: map[string]string{
// 			"state1": "value1",
// 			"state2": "value2",
// 			"state3": "value3",
// 		},
// 	}

// 	// Test : Get
// 	fmt.Println(sm.Get("state1"))
// 	fmt.Println(sm.Get("state4"))

// 	// Test : Set
// 	res, err := sm.Set("state4", "value4")
// 	if err != nil {
// 		fmt.Println(res, ":", err)
// 	} else {
// 		fmt.Println(res)
// 	}

// 	res2, err2 := sm.Set("", "value4")
// 	if err2 != nil {
// 		fmt.Println(res2, ":", err2)
// 	} else {
// 		fmt.Println(res2)
// 	}

// 	fmt.Println(sm.Get("state4"))

// 	// Test : Strln
// 	fmt.Println(sm.Strln("state1"))
// 	fmt.Println(sm.Strln("state5"))

// 	// Test : Strln
// 	res3, err3 := sm.Append("state4", "valueAppend4")
// 	if err3 != nil {
// 		fmt.Println(res3, ":", err3)
// 	} else {
// 		fmt.Println(res3)
// 	}

// 	res5, err5 := sm.Append("state5", "valueAppend4")
// 	if err3 != nil {
// 		fmt.Println(res5, ":", err5)
// 	} else {
// 		fmt.Println(res5)
// 	}

// 	fmt.Println(sm.Get("state5"))
// 	fmt.Println(sm.GetAll())

// 	// Test : Delete
// 	fmt.Println(sm.Delete("state1"))
// 	fmt.Println(sm.Get("state1"))
// 	fmt.Println(sm.GetAll())
// }
