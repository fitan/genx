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

    func GenFuncArg(name, value string) (res FuncArg) {
        res.Name = name
        res.Value = value
        return
    }

    type FuncArg struct {
        Name string
        Value string
    }
}


options {
    tokenVocab=TLexer;
}

doc: line*;

line: func
    | INSET S? CLOSE
    | NEWLINE
    ;

func
locals [
    FuncArgs: []FuncArg,
    FuncName: string
]
: ATID {$FuncName = $ATID.text} ( LPAREN (String {$FuncArgs = append($FuncArgs, GenFuncArg("",trimQuotation($String.text)))} | argument {$FuncArgs = append($FuncArgs, $argument.res)} ) (',' (String {$FuncArgs = append($FuncArgs, GenFuncArg("",trimQuotation($String.text)))} | argument {$FuncArgs = append($FuncArgs, $argument.res)}))* RPAREN ) NEWLINE
    | ATID {$FuncName = $ATID.text} LPAREN RPAREN NEWLINE
    | FieldFuncName {$FuncName = strings.TrimSpace($FieldFuncName.text)} (FIELD {$FuncArgs = append($FuncArgs, GenFuncArg("",$FIELD.text))} (FIELD {$FuncArgs = append($FuncArgs, GenFuncArg("",$FIELD.text))})*)? OLDFUNCCLOSE
    | ATID
    ;



argument returns [FuncArg res]
: ID EQ String {
        $res = GenFuncArg($ID.text, trimQuotation($String.text));
    }
  ;


//fields: ATID TEXT* ;



