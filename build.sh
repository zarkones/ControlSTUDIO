rm -r export
mkdir export
mkdir export/c2
mkdir export/agents

# Build Agent.
echo "Building Linux Agents"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/agents/linux_amd64 ./agent
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/agents/linux_386 ./agent
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/agents/linux_arm ./agent
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/agents/linux_arm64 ./agent
echo "Building Windows Agents"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static' -H=windowsgui" -tags "netgo,windows,purego" -o export/agents/windows_amd64 ./agent
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w -extldflags='-static' -H=windowsgui" -tags "netgo,windows,purego" -o export/agents/windows_386 ./agent
CGO_ENABLED=0 GOOS=windows GOARCH=arm go build -ldflags="-s -w -extldflags='-static' -H=windowsgui" -tags "netgo,windows,purego" -o export/agents/windows_arm ./agent
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w -extldflags='-static' -H=windowsgui" -tags "netgo,windows,purego" -o export/agents/windows_arm64 ./agent
echo "Building Mac Agents"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,darwin,purego" -o export/agents/darwin_amd64 ./agent
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,darwin,purego" -o export/agents/darwin_arm64 ./agent
echo "Building OpenBSD Agents"
CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,openbsd,purego" -o export/agents/openbsd_amd64 ./agent
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,openbsd,purego" -o export/agents/openbsd_arm64 ./agent
CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,openbsd,purego" -o export/agents/openbsd_386 ./agent
CGO_ENABLED=0 GOOS=openbsd GOARCH=arm go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,openbsd,purego" -o export/agents/openbsd_arm ./agent
echo "Building Solaris Agents"
CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,solaris,purego" -o export/agents/solaris_amd64 ./agent

# Build C2.
echo "Building Linux C2"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/c2/linux_amd64 ./c2
CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/c2/linux_386 ./c2
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/c2/linux_arm ./c2
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,linux,purego" -o export/c2/linux_arm64 ./c2
echo "Building Windows C2"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,windows,purego" -o export/c2/windows_amd64.exe ./c2
# CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,windows,purego" -o export/c2/windows_386.exe ./c2
echo "Building MacOS C2"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -extldflags='-static'" -tags "netgo,darwin,purego" -o export/c2/darwin_amd64 ./c2

# Build User Interface.
echo "Building Linux UI"
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -tags "netgo,linux" -o export/ControlSTUDIO_linux_amd64 ./ui