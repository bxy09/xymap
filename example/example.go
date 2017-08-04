package example

//go:generate genny -pkg example -in=$GOPATH/src/github.com/bxy09/xymap/xymap.go -out=gen-xymap.go gen "KeyType=string,int ValueType=int"
