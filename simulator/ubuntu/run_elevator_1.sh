gnome-terminal --title="ELEVATOR SERVER 1" -- ./SimElevatorServer --port 15657

cd ../..

gnome-terminal --title="ELEVATOR NODE 1" -- go run main.go -port=15657
