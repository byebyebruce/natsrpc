package generator

func (g *Generator) toolName() string {
	return "protoc-gen-" + g.name
}
func (g *Generator) pluginName() string {
	return plugins[0].Name()
}
