rm main
rm main.zip
GOOS=linux go build -o main main.go
zip main.zip main