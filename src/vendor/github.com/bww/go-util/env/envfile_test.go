package env

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

// Test env parsing
func TestParseEnv(t *testing.T) {
  assertParseEnv(t, `KEY=VAL
FOO=BAR`, map[string]string{"KEY": "VAL", "FOO": "BAR"})
  assertParseEnv(t, `# Comment!
    KEY=VAL
FOO=BAR`, map[string]string{"KEY": "VAL", "FOO": "BAR"})
  assertParseEnv(t, `# Comment!
    KEY=VAL # Tail comment
    FOO=BAR`, map[string]string{"KEY": "VAL", "FOO": "BAR"})
  assertParseEnv(t, `# Comment!
    KEY="Quoted #Val" # Tail comment
    FOO=BAR`, map[string]string{"KEY": "Quoted #Val", "FOO": "BAR"})
}

// Assert parse env pairs
func assertParseEnv(t *testing.T, s string, e map[string]string) {
  fmt.Println(s, " -> ", e)
  c, err := parseEnv(s)
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    assert.Equal(t, e, c)
  }
}

// Test env decl parsing
func TestParseEnvDecl(t *testing.T) {
  assertParseEnvDecl(t, "KEY=VAL", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY=VAL # Blah...", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY=\"VAL\" # Blah...", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY=\"VAL\"\n# Blah...", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY=\"WHY VAL\"\n# Blah...", "KEY", "WHY VAL")
  assertParseEnvDecl(t, "KEY=\"WHY!VAL\"\n# Blah...", "KEY", "WHY!VAL")
  assertParseEnvDecl(t, "KEY='VAL' # Blah...", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY='VAL'\n# Blah...", "KEY", "VAL")
  assertParseEnvDecl(t, "KEY='WHY VAL'\n# Blah...", "KEY", "WHY VAL")
  assertParseEnvDecl(t, "KEY='WHY!VAL'\n# Blah...", "KEY", "WHY!VAL")
  
  assertParseEnvDecl(t, `KEY="#" # Blah...`, "KEY", "#")
  assertParseEnvDecl(t, `KEY="\"#\"\"" # Blah...`, "KEY", `"#""`)
  assertParseEnvDecl(t, `KEY='"#""' # Blah...`, "KEY", `"#""`)
  assertParseEnvDecl(t, `KEY="'#'" # Blah...`, "KEY", `'#'`)
  assertParseEnvDecl(t, `KEY="\\#" # Blah...`, "KEY", `\#`)
}

// Assert env decl
func assertParseEnvDecl(t *testing.T, s, ek, ev string) {
  fmt.Print(s, " -> ")
  ak, av, err := envDecl(s)
  if assert.Nil(t, err, fmt.Sprintf("%v", err)) {
    fmt.Println(ak, av)
    assert.Equal(t, ek, ak)
    assert.Equal(t, ev, av)
  }
}
