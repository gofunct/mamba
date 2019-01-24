//go:generate mamba walk grpc -i third_party
//go:generate mamba walk gogo -i third_party
//go:generate go fmt ./...
//go:generate go vet ./...
//go:generate go install ./...
//go:generate git add .
//go:generate git commit -m "successful go generate 🐍"
//go:generate git push origin master

package mamba
