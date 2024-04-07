package gstring

import "regexp"

/**
Regex examples

The following table shows a couple of regular expression strings.

Regex	Meaning
.	Matches any single character.
?	Matches the preceding element once or not at all.
+	Matches the preceding element once or more times.
*	Matches the preceding element zero or more times.
^	Matches the starting position within the string.
$	Matches the ending position within the string.
|	Alternation operator.
[abc]	Matches a or b, or c.
[a-c]	Range; matches a or b, or c.
[^abc]	Negation, matches everything except a, or b, or c.
\s	Matches white space character.
\w	Matches a word character; equivalent to [a-zA-Z_0-9]
*/

func RegexMatch(s string, exp string) (bool, error) {
	return regexp.MatchString(exp, s)
}
