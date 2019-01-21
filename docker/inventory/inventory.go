package inventory

type dockerInventory struct {
	Dgraph bool
}

var (
	Inventory = &dockerInventory{}
)
