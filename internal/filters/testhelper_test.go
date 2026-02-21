package filters

import (
	"go/token"

	"lingo/internal/analyzer/log"
)

// makeParts — быстрое создание []LogPart для юнит-тестов.
// Каждая пара (value, isLiteral) → LogPart с фиктивными позициями.
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
