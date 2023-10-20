parser grammar TParser;
@header {
    import "strings"
}

@members {
    func trimQuotation(s string) string {
        if strings.HasPrefix(s, "\"") {
	        return strings.Trim(s, "\"")
        }
        if strings.HasPrefix(s, "'") {
	        return strings.Trim(s, "'")
        }
        return s
    }
}


options {
    tokenVocab=TLexer;
}

doc: line+;

line: func
    | INSET S CLOSE
    | NEWLINE
    ;

func
locals [
    FuncArgs: []string,
    FuncName: string
]
: ID {$FuncName = $ID.text} ( LPAREN String {$FuncArgs = append($FuncArgs, trimQuotation($String.text))} (',' String {$FuncArgs = append($FuncArgs, trimQuotation($String.text))})* RPAREN ) NEWLINE
    | ID {$FuncName = $ID.text} LPAREN RPAREN NEWLINE
    | FieldFuncName {$FuncName = strings.TrimSpace($FieldFuncName.text)} (FIELD {$FuncArgs = append($FuncArgs, $FIELD.text)} (OLDFUNCWS FIELD {$FuncArgs = append($FuncArgs, $FIELD.text)})*)? OLDFUNCWS? OLDFUNCCLOSE
    ;
//




//fields: ID TEXT* ;



