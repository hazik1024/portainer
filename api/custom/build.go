package custom

// BuildHandler 编译镜像
type BuildHandler struct {
}

// NewBuildHandler 返回新的Handler
func NewBuildHandler() *BuildHandler {
	return &BuildHandler{}
}
