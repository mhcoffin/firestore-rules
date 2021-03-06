// Code generated by "stringer -type=Kind"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Error-0]
	_ = x[Eof-1]
	_ = x[Word-2]
	_ = x[Dot-3]
	_ = x[IntLiteral-4]
	_ = x[FloatLiteral-5]
	_ = x[LeftBrace-6]
	_ = x[RightBrace-7]
	_ = x[LeftParen-8]
	_ = x[RightParen-9]
	_ = x[StringLiteral-10]
	_ = x[Slash-11]
	_ = x[Minus-12]
	_ = x[Plus-13]
	_ = x[Comma-14]
	_ = x[Eq-15]
	_ = x[EqEq-16]
	_ = x[SemiColon-17]
	_ = x[LeftSquareBracket-18]
	_ = x[RightSquareBracket-19]
	_ = x[Less-20]
	_ = x[LessEq-21]
	_ = x[Greater-22]
	_ = x[GreaterEq-23]
	_ = x[Colon-24]
	_ = x[QuestionMark-25]
	_ = x[And-26]
	_ = x[AndAnd-27]
	_ = x[Or-28]
	_ = x[OrOr-29]
	_ = x[NotEq-30]
	_ = x[Bang-31]
	_ = x[Star-32]
	_ = x[StarStar-33]
	_ = x[Bytes-34]
	_ = x[Percent-35]
	_ = x[Identifier-36]
	_ = x[Service-37]
	_ = x[Match-38]
	_ = x[Allow-39]
	_ = x[Create-40]
	_ = x[Update-41]
	_ = x[Delete-42]
	_ = x[Write-43]
	_ = x[Get-44]
	_ = x[List-45]
	_ = x[Read-46]
	_ = x[If-47]
	_ = x[Function-48]
	_ = x[True-49]
	_ = x[False-50]
	_ = x[In-51]
	_ = x[Is-52]
	_ = x[Return-53]
	_ = x[Let-54]
	_ = x[RulesVersion-55]
}

const _Kind_name = "ErrorEofWordDotIntLiteralFloatLiteralLeftBraceRightBraceLeftParenRightParenStringLiteralSlashMinusPlusCommaEqEqEqSemiColonLeftSquareBracketRightSquareBracketLessLessEqGreaterGreaterEqColonQuestionMarkAndAndAndOrOrOrNotEqBangStarStarStarBytesPercentIdentifierServiceMatchAllowCreateUpdateDeleteWriteGetListReadIfFunctionTrueFalseInIsReturnLetRulesVersion"

var _Kind_index = [...]uint16{0, 5, 8, 12, 15, 25, 37, 46, 56, 65, 75, 88, 93, 98, 102, 107, 109, 113, 122, 139, 157, 161, 167, 174, 183, 188, 200, 203, 209, 211, 215, 220, 224, 228, 236, 241, 248, 258, 265, 270, 275, 281, 287, 293, 298, 301, 305, 309, 311, 319, 323, 328, 330, 332, 338, 341, 353}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
