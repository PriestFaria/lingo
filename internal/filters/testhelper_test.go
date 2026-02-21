package filters

import (
	"go/token"

	"github.com/PriestFaria/lingo/internal/analyzer/log"
)

// makeParts creates a []LogPart from alternating (value string, isLiteral bool) pairs.
// Positions are synthetic and suitable for unit tests only.
func makeParts(args ...interface{}) []log.LogPart {
	var parts []log.LogPart
	for i := 0; i+1 < len(args); i += 2 {
		value := args[i].(string)
		isLiteral := args[i+1].(bool)
		pos := token.Pos(100 + i*10)
		parts = append(parts, log.LogPart{
			Value:     value,
			IsLiteral: isLiteral,
			Pos:       pos,
			End:       pos + token.Pos(len(value)) + 2, 
		})
	}
	return parts
}

func makeCtx(parts []log.LogPart) *log.LogContext {

	return &log.LogContext{Parts: parts}

}	
