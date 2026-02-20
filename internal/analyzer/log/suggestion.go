package log

import "go/token"

type Suggestion struct {
	Message string
	Replace string
	Start token.Pos
	End token.Pos
}