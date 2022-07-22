if [[ $OSTYPE =~ "linux" ]]
then
## Linux
GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -ldflags '-H windowsgui -w -s' -o down_tip
fi
if [[ $OSTYPE =~ "darwin" ]]
then
## Mac
GOARCH=amd64 CGO_ENABLED=1 GOOS=darwin go build -ldflags '-w -s' -o down_tip
cp -r ./build/DownTip.app ./
mv down_tip ./DownTip.app/Contents/MacOS
fi

