package markdown

/**
 * A markdown documentation generator
 */
type Generator struct {
}

/**
 * Produce a new emitter
 */
func New() *Generator {
  return &Generator{}
}

/**
 * Generate documentation
 */
func (g *Generator) Generate(w io.Writer, c test.Case) error {
  return nil
}
