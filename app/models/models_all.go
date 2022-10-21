package models

type ApifoxModel struct {
	Code int64 `json:"code"`
	Data Pet   `json:"data"`
}

// Pet
type Pet struct {
	Category  Category `json:"category"`  // 分组
	ID        int64    `json:"id"`        // 宠物ID编号
	Name      string   `json:"name"`      // 名称
	PhotoUrls []string `json:"photoUrls"` // 照片URL
	Status    Status   `json:"status"`    // 宠物销售状态
	Tags      []Tag    `json:"tags"`      // 标签
}

// 分组
//
// Category
type Category struct {
	ID   *int64  `json:"id,omitempty"`   // 分组ID编号
	Name *string `json:"name,omitempty"` // 分组名称
}

// Tag
type Tag struct {
	ID   *int64  `json:"id,omitempty"`   // 标签ID编号
	Name *string `json:"name,omitempty"` // 标签名称
}

// 宠物销售状态
type Status string

const (
	Available Status = "available"
	Pending   Status = "pending"
	Sold      Status = "sold"
)
