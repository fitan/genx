lexer grammar TLexer;

@lexer::members {
    var nesting int;
}


ID: '@' [a-zA-Z_][a-zA-Z0-9_]*;
String:  ( '\'' ('\'\''|~'\'')* '\'' ) | ( '"' ('""'|~'"')* '"' );

FieldFuncName: ID  ' '+ -> pushMode(OLDFUNC);
//TEXT : ~[ <>()@,\r\n"]+;
Comma: ',';

LPAREN    : '(' {fmt.Println(nesting);nesting++;} ;

RPAREN    : ')' {fmt.Println(nesting);nesting--;} ;

IGNORE_NEWLINE
    :   '\r'? '\n' {nesting>0}? -> skip
    ;


NEWLINE
   :    '\r'? '\n'
   ;


WS: [ \t]+ -> skip ;

INSET: ~'@' {fmt.Println("ent INSIDE")} -> pushMode(INSIDE) ;

mode INSIDE;
S:   ~[\r\n]+;
CLOSE : '\r'? '\n' {fmt.Println("out INSIDE")} -> popMode;
//S: . -> more;
//CLOSE : '\r'? '\n' {fmt.Println("out INSIDE")} -> popMode;

mode OLDFUNC;
OLDFUNCCLOSE : '\r'? '\n' {fmt.Println("out Fied")} -> popMode;
FIELD: ~[ \t\r\n]+;
OLDFUNCWS: [ \t]+ ;

