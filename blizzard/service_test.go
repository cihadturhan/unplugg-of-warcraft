package configuration_test

import (
	"encoding/json"
	"github.com/whitesmith/brand-digital-box"
	"github.com/whitesmith/brand-digital-box/configuration"
	"reflect"
	"testing"
	"time"
)

// TestService_UpdateConfig tests updating the configuration.
func TestService_UpdateConfig(t *testing.T) {
	c := NewClient()

	// mock downloader service.
	c.DownloaderService.DownloadMediaFn = func(media box.Media) error {
		return nil
	}

	// mock player service.
	c.PlayerService.SetOptionFn = func(key, value string) error {
		if key == "geometry" && value == "1920x1080+10+10" {
			return nil
		} else if key == "key" && value == "value" {
			return nil
		}
		t.Fatalf("invalid option: %v;%v", key, value)
		return nil
	}
	c.PlayerService.SetFlagFn = func(key string, value bool) error {
		if key == "key" && value == true {
			return nil
		}
		t.Fatalf("invalid option: %v;%v", key, value)
		return nil
	}

	// mock communication module.
	var subscriptions int
	c.MQTTService.SubscribeFn = func(topic string, channel chan box.Message, onConnect func()) {
		onConnect()
		if topic != "$aws/things/nuc-aspire/shadow/get/accepted" && topic != "$aws/things/nuc-aspire/shadow/update/delta" {
			t.Fatalf("invalid topic: %s", topic)
		}
		subscriptions++
	}
	var publishes int
	c.MQTTService.PublishFn = func(topic string, message []byte) error {
		if topic != "$aws/things/nuc-aspire/shadow/get" {
			t.Fatal("invalid topic")
		}
		publishes++
		return nil
	}

	// mock screen configuration.
	var configurations int
	c.ScreenService.SetFullConfigFn = func(b float32, r string, p string, o string) {
		if b != 0.5 || r != "1920x1080" || o != "HDMI" || p != "0x0" {
			t.Fatal("invalid configuration received")
		}
		configurations++
	}

	done := make(chan bool)
	c.ScreenService.ApplyFn = func() error {
		return nil
	}

	// mock new messages.
	config := box.ConfigMessage{}
	config.Screen.Brightness = 0.5
	config.Screen.Output = "HDMI"
	config.Screen.Position = "0x0"
	config.Screen.Resolution = "1920x1080"
	config.Fallback.ID = "fallback"
	config.Fallback.URL = "URL"
	config.Fallback.Type = box.TypeVideo
	config.Player.Resolution = "1920x1080"
	config.Player.Position = "10x10"
	config.Player.Options = map[string]string{"key": "value"}
	config.Player.Flags = map[string]bool{"key": true}

	// mock database module.
	c.DatabaseService.LoadFn = func(collection string, key string) ([]byte, error) {
		if collection != configuration.Collection || key != configuration.Record {
			t.Fatal("invalid collection or record")
		}
		e, _ := json.Marshal(configuration.StorageConfig{Version: 0, Config: config})

		return e, nil
	}

	c.DatabaseService.SaveFn = func(collection string, key string, value []byte) error {
		if collection != configuration.Collection || key != configuration.Record {
			t.Fatal("invalid collection or record")
		}

		e, _ := json.Marshal(configuration.StorageConfig{Version: 2, Config: config})

		if !reflect.DeepEqual(e, value) {
			t.Fatal("invalid configuration saved")
		}

		done <- true
		return nil
	}

	get := configuration.GetMessage{}
	get.Version = 2
	get.Timestamp = time.Now().Unix()
	get.State.Desired = config

	msg := box.Message{}
	msg.Topic = "get"
	msg.Payload, _ = json.Marshal(get)

	// mock connection
	if err := c.Open(); err != nil {
		panic(err)
	}

	c.GetMessages <- msg

	if !<-done {
		t.Fatal("failed to apply settings")
	} else {
		if subscriptions != 2 {
			t.Fatal("invalid subscription count")
		}
		if publishes != 1 {
			t.Fatal("invalid publishes count")
		}
		if configurations != 1 {
			t.Fatal("invalid configuration count")
		}
	}
}
