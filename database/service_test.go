package database_test

import (
	"encoding/json"
	"github.com/whitesmith/brand-digital-box/database"
	"reflect"
	"testing"
)

// TestService_InsertRecord tests inserting a database record.
func TestService_InsertRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	msg := ConfigMessage{}
	msg.Screen.Brightness = 0.5
	msg.Screen.Resolution = "1920x1080"
	msg.Screen.Position = "0x0"
	msg.Screen.Output = "HDMI"

	r, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Save("TestCollection", "TestKey", r)
	if err != nil {
		t.Fatal(err)
	}
}

// TestService_LoadRecord tests retrieving a database record.
func TestService_LoadRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	msg := ConfigMessage{}
	msg.Screen.Brightness = 0.5
	msg.Screen.Resolution = "1920x1080"
	msg.Screen.Position = "0x0"
	msg.Screen.Output = "HDMI"

	r, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	record, err := s.Load("TestCollection", "TestKey")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(r, record) {
		t.Fatal(err)
	}
}

// TestService_IterateCollection tests iterating over a database collection.
func TestService_IterateCollection(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	// insert struct 1.
	test1 := TestStruct{Field1: "01", Field2: 1}
	r1, err := json.Marshal(test1)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test1.Field1, r1)
	if err != nil {
		t.Fatal(err)
	}

	// insert struct 2.
	test2 := TestStruct{Field1: "02", Field2: 2}
	r2, err := json.Marshal(test2)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test2.Field1, r2)
	if err != nil {
		t.Fatal(err)
	}

	// insert struct 3.
	test3 := TestStruct{Field1: "03", Field2: 3}
	r3, err := json.Marshal(test3)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test3.Field1, r3)
	if err != nil {
		t.Fatal(err)
	}

	iterated := []bool{false, false, false}
	// iterate over structs.
	s.Iterate("test", func(k, v []byte) error {
		var media TestStruct
		if err = json.Unmarshal(v, &media); err != nil {
			t.Fatal(err)
		}

		if media.Field1 == test1.Field1 {
			if !reflect.DeepEqual(media, test1) {
				t.Fatal(err)
			}
			iterated[0] = true
		} else if media.Field1 == test2.Field1 {
			if !reflect.DeepEqual(media, test2) {
				t.Fatal(err)
			}
			iterated[1] = true
		} else if media.Field1 == test3.Field1 {
			if !reflect.DeepEqual(media, test3) {
				t.Fatal(err)
			}
			iterated[2] = true
		}
		return nil
	})

	for k, v := range iterated {
		if !v {
			t.Fatalf("failed to iterate over %d", k)
		}
	}
}

// TestService_DeleteRecords tests removing records from a database collection.
func TestService_DeleteRecords(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	// insert struct 1.
	test1 := TestStruct{Field1: "01", Field2: 1}
	r1, err := json.Marshal(test1)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test1.Field1, r1)
	if err != nil {
		t.Fatal(err)
	}

	// insert struct 2.
	test2 := TestStruct{Field1: "02", Field2: 2}
	r2, err := json.Marshal(test2)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test2.Field1, r2)
	if err != nil {
		t.Fatal(err)
	}

	// insert struct 3.
	test3 := TestStruct{Field1: "03", Field2: 3}
	r3, err := json.Marshal(test3)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Save("test", test3.Field1, r3)
	if err != nil {
		t.Fatal(err)
	}

	// delete records.
	if err := s.Delete("test", test1.Field1, test2.Field1); err != nil {
		t.Fatal(err)
	}

	iterated := []bool{false, false, false}
	// iterate over structs.
	s.Iterate("test", func(k, v []byte) error {
		var media TestStruct
		if err = json.Unmarshal(v, &media); err != nil {
			t.Fatal(err)
		}

		if media.Field1 == test1.Field1 {
			if !reflect.DeepEqual(media, test1) {
				t.Fatal(err)
			}
			iterated[0] = true
		} else if media.Field1 == test2.Field1 {
			if !reflect.DeepEqual(media, test2) {
				t.Fatal(err)
			}
			iterated[1] = true
		} else if media.Field1 == test3.Field1 {
			if !reflect.DeepEqual(media, test3) {
				t.Fatal(err)
			}
			iterated[2] = true
		}
		return nil
	})

	// check which keys got iterated over.
	if !reflect.DeepEqual(iterated, []bool{false, false, true}) {
		t.Fatalf("failed to delete keys: %v", iterated)
	}
}

// TestService_LoadRecord_NoRecord tests retrieving a database record that does not exist.
func TestService_LoadRecord_NoRecord(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	_, err := s.Load("TestCollection", "FAKE")
	if err != database.ErrRecordNotFound {
		t.Fatal(err)
	}
}

// ConfigMessage represents the message structure of the configuration messages.
type ConfigMessage struct {
	Screen struct {
		Brightness float32 `json:"brightness,omitempty"`
		Resolution string  `json:"resolution,omitempty"`
		Position   string  `json:"position,omitempty"`
		Output     string  `json:"output,omitempty"`
	} `json:"screen,omitempty"`
}

// TestStruct represents a test structure for database records.
type TestStruct struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2,omitempty"`
}
