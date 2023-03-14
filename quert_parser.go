package query_parser

import (
	"strings"

	"github.com/orkes-io/query_parser/stack"
	"github.com/orkes-io/query_parser/util"
)

func ConvertWhereClauseToMongoFilter(whereClause string) map[string]any {
	query := ReplaceLogicalOperators(whereClause)
	tokens := Tokenize(query)
	return ConvertEquivalentMongoQuery(tokens)
}

func Tokenize(query string) []string {
	tokens := make([]string, 0)

	temp := ""
	for i := 0; i < len(query); i++ {
		if query[i] == '(' {
			tokens = append(tokens, string(query[i]))
		} else if query[i] == ')' {
			tokens = append(tokens, temp, string(query[i]))
			temp = ""
		} else if IsLogicalOperator(string(query[i])) {
			tokens = append(tokens, temp, string(query[i]))
			temp = ""
		} else {
			temp = temp + string(query[i])
		}
	}

	if len(temp) != 0 {
		tokens = append(tokens, temp)
	}

	tokens = util.SliceMap(tokens, func(s string, _ int) string {
		return strings.TrimSpace(s)
	})

	return util.SliceFilter(tokens, func(s string, _ int) bool {
		return len(strings.TrimSpace(s)) > 0
	})
}

func ReplaceLogicalOperators(query string) string {

	query = strings.ReplaceAll(query, " AND ", " & ")
	query = strings.ReplaceAll(query, " and ", " & ")

	query = strings.ReplaceAll(query, " OR ", " | ")
	query = strings.ReplaceAll(query, " or ", " | ")

	return query
}

func IsLogicalOperator(char string) bool {
	operators := []string{"&", "|"}
	return util.SliceContains(operators, char)
}

func ConvertEquivalentMongoQuery(tokens []string) map[string]any {
	postfixTokens := ConvertInfixToPostfix(tokens)

	evaluationStack := stack.NewStack[map[string]any]()

	for _, token := range postfixTokens {
		if IsLogicalOperator(token) {
			operand2 := evaluationStack.Pop()
			operand1 := evaluationStack.Pop()
			evaluationStack.Push(map[string]any{
				ConvertEquivalentMongoOperator(token): []map[string]any{
					operand1, operand2,
				},
			})
		} else {
			evaluationStack.Push(ConvertEquivalentMongoMap(token))
		}
	}
	return evaluationStack.Pop()
}

func ConvertEquivalentMongoMap(token string) map[string]any {
	fields := strings.Fields(token)
	operand1, operator, operand2 := fields[0], fields[1], fields[2]

	mongoMap := map[string]any{
		operand1: map[string]any{
			ConvertEquivalentMongoOperator(operator): operand2,
		},
	}

	return mongoMap
}

func ConvertEquivalentMongoOperator(operator string) string {
	switch strings.TrimSpace(operator) {
	case "&":
		return "$and"
	case "|":
		return "$or"
	case "=":
		return "$eq"
	case "!=":
		return "$ne"
	case ">":
		return "$gt"
	case ">=":
		return "$ge"
	case "<":
		return "$lt"
	case "<=":
		return "$le"
	case "IN":
		return "$in"
	case "NOT IN":
		return "$ni"
	default:
		return ""
	}
}

func ConvertInfixToPostfix(infixTokens []string) []string {
	postfixTokens := make([]string, 0, len(infixTokens))

	operatorStack := stack.NewStack[string]()

	for _, token := range infixTokens {
		if IsLogicalOperator(token) {
			// pop high priority operators
			for !operatorStack.IsEmpty() && Precedence(token) <= Precedence(operatorStack.Top()) {
				postfixTokens = append(postfixTokens, operatorStack.Pop())
			}
			operatorStack.Push(token)
		} else if token == "(" {
			operatorStack.Push(token)
		} else if token == ")" {
			for operatorStack.Top() != "(" {
				postfixTokens = append(postfixTokens, operatorStack.Pop())
			}
			operatorStack.Pop()
		} else {
			postfixTokens = append(postfixTokens, token)
		}
	}

	for !operatorStack.IsEmpty() {
		postfixTokens = append(postfixTokens, operatorStack.Pop())
	}

	return postfixTokens
}

func Precedence(s string) int {
	switch s {
	case "&":
		return 9
	case "|":
		return 8
	default:
		return 0
	}
}
