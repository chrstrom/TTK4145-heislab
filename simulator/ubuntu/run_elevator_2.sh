gnome-terminal --title="ELEVATOR SERVER 2" -- ./SimElevatorServer --port 15658

cd ../..

gnome-terminal --title="ELEVATOR NODE 2" -- go run main.go -port=15658
