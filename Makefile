all:
	go build -o ewfreader ./cmd/*.go

generate:
	cd parser/ && binparsegen conversion.spec.yaml > ewf_gen.go
