package main

type Field struct {
	name          string
	typ           string
	invariant     string
	allowCreateIf string
	allowUpdateIf string
}

type Doc struct {
	path          string
	fields        []Field
	allowCreateIf string
	allowUpdateIf string
}

var issueDoc = Doc{
	path: "/issues/{doc}",
	fields: []Field{
		{
			name:          "id",
			typ:           "string",
			invariant:     "${value} == ${doc} && ${value}.size() == 10",
			allowUpdateIf: "${value} == ${value.prev}",
		},
		{
			name: "author",
			typ:  "string",
		},
		{
			name:          "created",
			typ:           "timestamp",
			allowCreateIf: "${created} =~ ${now}",
			allowUpdateIf: "${created} == ${created.prev}",
		},
		{
			name:          "modified",
			typ:           "timestamp",
			allowCreateIf: "${modified} == ${created}",
			allowUpdateIf: "${modified} =~ request.time()",
		},
		{
			name:          "title",
			typ:           "string",
			invariant:     "${title}.size() > 10 && ${title}.size() < 100",
			allowUpdateIf: "${author} == request.auth.uid",
		},
	},
	allowCreateIf: "request.auth.uid != null && request.resource.author == request.auth.uid",
	allowUpdateIf: "request.auth.uid == request.resource.data.author",
}

func main() {

}
