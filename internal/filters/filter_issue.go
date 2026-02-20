package filters

import "go/token"

type IssueFix struct {
    Message string    
    Pos     token.Pos 
    End     token.Pos 
    NewText string    
}

type FilterIssue struct {
    Message string
    Pos     token.Pos
    Fix     *IssueFix
}