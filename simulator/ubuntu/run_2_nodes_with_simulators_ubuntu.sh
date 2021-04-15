
gnome-terminal --title="ELEVATOR SERVER 1" -- ./SimElevatorServer --port 15657
gnome-terminal --title="ELEVATOR SERVER 2" -- ./SimElevatorServer --port 15658

cd ../..

gnome-terminal --title="ELEVATOR NODE 1" -- go run main.go -port=15657
gnome-terminal --title="ELEVATOR NODE 2" -- go run main.go -port=15658
