# screenshot
screenshot uploader for mac

1. `go get github.com/nuuls/ni` and add it to PATH
1. `go get github.com/nuuls/screenshot`
1. open terminal and run `defaults write com.apple.screencapture location /your/screenshot/path`
1. run `killall SystemUIServer` to apply changes from above
1. edit main.go config so it uses the dir from above
1. `go build && ./screenshot`
1. use cmd+shift+4 to take screenshots
