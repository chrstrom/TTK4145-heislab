package config

import "time"

// Parameters for the elevator itself
const N_FLOORS = 4
const N_BUTTONS = 3
const DOOR_OPEN_DURATION = 2
const TRAVEL_TIME = 2.5
const MOTOR_STOP_DETECTION_TIME = time.Millisecond * 3000

// For the cab order storage
const N_FILE_DUPLICATES = 3
const BACKUP_FILE_PATH = "orderBackup/"

// For the shared hall orders
const ORDER_REPLY_TIME = time.Millisecond * 300
const ORDER_DELEGATION_TIME = time.Millisecond * 500
const ORDER_COMPLETION_TIMEOUT = time.Second * 20
const FSM_ORDER_TIMEOUT = time.Second * 3

// For the network module
const N_MESSAGE_DUPLICATES = 10
const NETWORK_CHANNEL_BUFFER_SIZE = 10

// For the local elevator FSM
const FSM_CHANNEL_BUFFER_SIZE = 10
