module main

go 1.19

replace github.com/stevegt/grokker => /home/stevegt/lab/grokker

require (
	github.com/stevegt/goadapt v0.4.0
	github.com/stevegt/grokker v0.0.0-20230904151314-12fff637e9c1
	github.com/yuin/goldmark v1.5.6
)

require (
	github.com/fabiustech/openai v0.4.0 // indirect
	github.com/sashabaranov/go-openai v1.9.0 // indirect
	github.com/stevegt/semver v0.0.0-20230512043732-92220054a49f // indirect
)
