# the build script, which I use (to run this script, use: sh build.sh(of course in a Unix environment))
gofmt -w *.go && go vet && golint && go build -o tanks

# run
./tanks
