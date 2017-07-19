# screenshot
screenshot uploader for mac

1. `go get github.com/nuuls/ni` and add it to PATH
2. `go get github.com/nuuls/screenshot`
3. open terminal and run `defaults write com.apple.screencapture location /your/screenshot/path`
4. edit main.go config so it uses the dir from above
5. `go build && ./screenshot`
6. use cmd+shift+4 to take screenshots
