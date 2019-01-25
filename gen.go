//go:generate mamba walk grpc
//go:generate mamba walk gogo
//go:generate go fmt ./...
//go:generate go vet ./...
//go:generate go install ./...
//go:generate git add .
//go:generate git commit -m "successful go generate 🐍"
//go:generate git push origin master

package mamba
