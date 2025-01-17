package shortcode

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&shortCodeRenderer{}, 0),
	))
}

type shortCodeRenderer struct{}

func (s *shortCodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindShortCode, s.render)
}

func (s *shortCodeRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	node, ok := n.(*ShortCodeNode)
	if !ok {
		return ast.WalkContinue, nil
	}

	content := source[node.start:node.end]
	output := node.fun(Markdown(content))
	w.Write([]byte(output))

	return ast.WalkContinue, nil
}
