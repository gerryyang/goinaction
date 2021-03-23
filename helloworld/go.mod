module example.com/m

go 1.16

require (
	example.com/greetings v0.0.0-00010101000000-000000000000
	github.com/gerryyang/goinaction/module/hello v0.0.0-20210323092231-9586d180a662 // indirect
	rsc.io/quote v1.5.2
)

replace example.com/greetings => ./greetings
