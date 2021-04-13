
start SimElevatorServer.exe --port 15657
start SimElevatorServer.exe --port 15658

cd ..\

start main.exe -port=15657
start main.exe -port=15658
