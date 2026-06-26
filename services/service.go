// 1、service 定义了“我能干什么”
package services

import "github.com/ximu-leo/s9-tss-gin/model"

type SignService interface {
	// 密钥生成，需要有入参
	Keygen(request model.KeygenRequest) (string, error)

	// 门限签名
	TransactionSign(request model.TransactionSignRequest) ([]byte, error)
}

// 在 Go 里，只有 结构体 才能实现接口，裸函数是没法实现接口的。
type Manager struct {
}

// &Manager 完整实现了 接口方法，所以可以这样写
func NewManager() (SignService, error) {
	return &Manager{}, nil
}

func (m *Manager) TransactionSign(request model.TransactionSignRequest) ([]byte, error) {
	return []byte("0x000000000000000000000000000000000000000000000000000000000000"), nil
}

func (m *Manager) Keygen(request model.KeygenRequest) (string, error) {
	return "0x00000000000000000000000000000000", nil
}
