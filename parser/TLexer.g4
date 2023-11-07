lexer grammar TLexer;

@lexer::members {
    var nesting int;
}


ATID: '@' [a-zA-Z_][a-zA-Z0-9_]*;


FieldFuncName: ATID  ' '+ -> pushMode(OLDFUNC);
//TEXT : ~[ <>()@,\r\n"]+;

LPAREN    : '(' {nesting++;} -> pushMode(PAREN) ;

//
//IGNORE_NEWLINE
//    :   '\r'? '\n' {nesting>0}? -> skip
//    ;
//

NEWLINE
   :    ('\r'? '\n') | EOF
   ;



WS: [ \t]+ -> skip ;

INSET: ~'@' -> pushMode(INSIDE) ;

mode PAREN;
EQ: '=';
Comma: ',';
PARENWS: [ \t]+ -> skip;
ID: [a-zA-Z_][a-zA-Z0-9_]*;
String:  ( '\'' ('\'\''|~'\'')* '\'' ) | ( '"' ('""'|~'"')* '"' );
RPAREN: ')' {nesting--;} -> popMode;


mode INSIDE;
S:   ~[\r\n]+;
CLOSE : (('\r'? '\n') | EOF) -> popMode;

mode OLDFUNC;
OLDFUNCCLOSE : (('\r'? '\n') | EOF) -> popMode;
OLDFUNCWS: [ \t]+ -> skip ;
FIELD: ( '\'' ('\'\''|~'\'')* '\'' ) | ( '"' ('""'|~'"')* '"' ) | ~[ \t\r\n]+;

