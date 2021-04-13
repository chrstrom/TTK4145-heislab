
start SimElevatorServer.exe --port 15657
start SimElevatorServer.exe --port 15658

cd ..\

go build main.go
start main.exe -port=15657
start main.exe -port=15658
