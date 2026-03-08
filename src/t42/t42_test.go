// Package t42 演示 github.com/google/uuid 的常见用法：
//   - 生成 UUID v4（随机）
//   - 生成 UUID v5（基于命名空间 + 名称的确定性 UUID）
//   - UUID 的字符串解析与格式验证
//   - UUID 作为结构体字段（配合 JSON 序列化）
package t42

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

// TestUUID_Generate 演示生成 UUID v4（随机 UUID）
//   - uuid.New() 每次调用返回一个全局唯一的随机 ID
//   - UUID 格式：8-4-4-4-12 共 36 个字符，例如 550e8400-e29b-41d4-a716-446655440000
func TestUUID_Generate(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	t.Logf("UUID v4 #1: %s", id1)
	t.Logf("UUID v4 #2: %s", id2)

	// 两次生成的 UUID 必然不同
	if id1 == id2 {
		t.Error("两次生成的 UUID 不应相同")
	}

	// UUID 字符串长度固定为 36（含 4 个连字符）
	if len(id1.String()) != 36 {
		t.Errorf("UUID 字符串长度期望 36，得到 %d", len(id1.String()))
	}

	// Version() 返回 UUID 版本号，v4 = 4
	if id1.Version() != 4 {
		t.Errorf("期望 UUID v4，得到 v%d", id1.Version())
	}
}

// TestUUID_Parse 演示将字符串解析为 UUID 类型，以及非法字符串的错误处理
func TestUUID_Parse(t *testing.T) {
	const raw = "550e8400-e29b-41d4-a716-446655440000"

	// 合法解析
	id, err := uuid.Parse(raw)
	if err != nil {
		t.Fatalf("解析合法 UUID 失败: %v", err)
	}
	t.Logf("解析结果: %s，版本: v%d", id, id.Version())

	// 非法字符串解析应返回错误
	_, err = uuid.Parse("not-a-uuid")
	if err == nil {
		t.Error("非法字符串应返回解析错误")
	}
	t.Logf("非法字符串解析错误（符合预期）: %v", err)
}

// TestUUID_V5 演示 UUID v5：基于命名空间 + 名称生成确定性 UUID
//   - 相同的命名空间 + 名称，每次生成的 UUID 完全一致（可复现）
//   - 不同的名称，生成的 UUID 不同
//   - 适合用于：根据业务主键生成稳定 ID、去重等场景
func TestUUID_V5(t *testing.T) {
	// uuid.NameSpaceDNS 是标准命名空间之一（还有 URL、OID、X500）
	ns := uuid.NameSpaceDNS

	id1 := uuid.NewSHA1(ns, []byte("example.com"))
	id2 := uuid.NewSHA1(ns, []byte("example.com"))
	id3 := uuid.NewSHA1(ns, []byte("other.com"))

	t.Logf("example.com → %s（v%d）", id1, id1.Version())
	t.Logf("example.com → %s（v%d）", id2, id2.Version())
	t.Logf("other.com   → %s（v%d）", id3, id3.Version())

	// 相同输入，结果确定性一致
	if id1 != id2 {
		t.Error("相同输入的 UUID v5 应完全相同")
	}

	// 不同输入，结果必然不同
	if id1 == id3 {
		t.Error("不同名称的 UUID v5 不应相同")
	}
}

// TestUUID_JSON 演示 UUID 作为结构体字段配合 JSON 序列化/反序列化
//   - uuid.UUID 实现了 json.Marshaler / json.Unmarshaler 接口
//   - JSON 中以标准字符串形式存储，例如 "id":"550e8400-..."
func TestUUID_JSON(t *testing.T) {
	type Order struct {
		ID     uuid.UUID `json:"id"`
		Item   string    `json:"item"`
		Amount int       `json:"amount"`
	}

	// 序列化
	order := Order{
		ID:     uuid.New(),
		Item:   "Go 语言实战",
		Amount: 3,
	}
	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}
	t.Logf("序列化结果: %s", data)

	// 反序列化
	var decoded Order
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}
	t.Logf("反序列化结果: ID=%s Item=%s Amount=%d", decoded.ID, decoded.Item, decoded.Amount)

	// 序列化再反序列化后 UUID 应保持不变
	if order.ID != decoded.ID {
		t.Errorf("UUID 序列化往返后不一致: %s != %s", order.ID, decoded.ID)
	}
}
