Cab order storage
================

The purpose of this module is to provide functionality for storing and loading cab orders to/from disk.
N_FILE_DUPLICATES (from config) copies of the hall orders are stored in BACKUP_FILE_PATH (from config)
as .txt files.

Backups are formatted as `bool bool ... bool`, where bool i represents if there is a cab order at floor i.
When loading cab orders, a simple majority-wins voting system is used, should one file differ from the rest.

### Interface
* StoreCabOrders(orders [N_FLOORS][N_BUTTONS]bool)
* LoadCabOrders() [N_FLOORS]bool
