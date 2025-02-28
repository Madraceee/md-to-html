Document = { Chunks } eof
Chunks = Chunk ( newLine Chunks | eof )
Chunk =  Heading | Paragraph | Code | List | Blockquote | Line
Heading = Header Paragraph
Header = ( "#" | "##" | "###" | "####" | "#####" | "######" ) { Whitespace }
Paragraph = { ( String | Format | Whitespace )} newLine
String = (Alphabets | Digit | SpecialCharacters)
Format = Bold | Italic
Bold =  Astrick Astrick Paragraph Astrick Astrick |
        Underscore Underscore Paragraph Underscore Underscore 
Italic= Astrick Paragraph  Astrick |
        Underscore Paragraph  Underscore 
Alphabets = "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I" | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R" | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z" | "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z"
Digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
SpecialCharacters = "@" | "$" | "%" | "^" | "&" | "'" | """ | "," | ";" | ":"
Whitespace = " " | "\t"
Astrick = "*"
Underscore = "_"
FormatCharacters = "#" | ">" | "{" | "}" | "[" | "]" | "-" | "_" | "(" | ")" | "\" | "`"

Code =  Backtick { ( String | FormatCharacters | Whitespace ) } Backtick |
        Backtick Backtick Backtick { Paragraph }  Backtick Backtick Backtick
List = { Tab } ( "-" | "*" | Digit "." ) Whitespace Paragraph { newLine List }
Blockquote = ">" Whitespace Paragraph { newLine Blockquote }
Line = "---" | "***"

Backtick = "`"


Line is limited to three dash or astricks followed by newline
List can be nested
Blockquote is minimal
