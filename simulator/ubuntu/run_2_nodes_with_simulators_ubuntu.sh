
gnome-terminal -- ./SimElevatorServer --port 15657
gnome-terminal -- ./SimElevatorServer --port 15658

cd ../..

gnome-terminal -- go run main.go -port=15657
gnome-terminal -- go run main.go -port=15658
