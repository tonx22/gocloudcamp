package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

var version int32
var ErrConfigUsed = errors.New("Specified config is used")

func TestGRPCClient(t *testing.T) {

	client, err := NewGRPCClient()
	if err != nil {
		t.Fatalf("failed to create grpc client: %v", err)
	}

	r := ConfigRequest{Service: "some-service"}
	_ = json.Unmarshal([]byte("{\"key1\":\"value1\",\"key2\":\"value2\"}"), &r.Data)

	t.Run("(1) создание конфига", func(t *testing.T) {
		res, err := client.SetConfig(context.TODO(), r)
		if err != nil {
			t.Fatalf("failed to create config: %v", err)
		}
		version = res.Version
	})
	fmt.Println("текущая версия конфига:", version)

	t.Run("(2) создание новой версии конфига", func(t *testing.T) {
		_, err := client.SetConfig(context.TODO(), r)
		if err != nil {
			t.Fatalf("failed to create config: %v", err)
		}
	})

	t.Run("(3) получение текущего конфига и проверка версии", func(t *testing.T) {
		r.Version = 0
		res, err := client.GetConfig(context.TODO(), r)
		if err != nil {
			t.Fatalf("failed to get config: %v", err)
		}
		require.Equal(t, res.Version, int32(version+1), "некорректная версия текущего конфига")
	})

	t.Run("(4) попытка удалить актуальный (используемый) конфиг", func(t *testing.T) {
		r.Version = version + 1
		_, err = client.DelConfig(context.TODO(), r)
		require.Containsf(t, err.Error(), ErrConfigUsed.Error(), "Error 'ErrConfigUsed' required")
	})

	t.Run("(5) откат конфига (возврат к предыдущей версии)", func(t *testing.T) {
		r.Version = version
		r.Used = true
		res, err := client.UpdConfig(context.TODO(), r)
		if err != nil {
			t.Fatalf("failed to update config: %v", err)
		}
		require.Equal(t, res.Version, int32(version), "некорректная версия текущего конфига")
	})

	t.Run("(6) удаление ненужного конфига", func(t *testing.T) {
		r.Version = version + 1
		_, err := client.DelConfig(context.TODO(), r)
		if err != nil {
			t.Fatalf("failed to delete config: %v", err)
		}
	})

}
