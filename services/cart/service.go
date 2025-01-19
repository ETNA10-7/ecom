package cart

import (
	"fmt"

	"github.com/ETNA10-7/ecom/types"
)

// Why this exists the below function
func getCartItemsIDs(items []types.CartCheckoutItem) ([]int, error) {
	productIds := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

// func checkIfCartIsInStock(cartItems []types.CartCheckoutItem, products map[int]types.ProductStock) error {
// 	if len(cartItems) == 0 {
// 		return fmt.Errorf("Cart is Empty")

// 	}

// 	for _, item := range cartItems {
// 		product, ok := products[item.ProductID]
// 		if !ok {
// 			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
// 		}
// 		if int(product.Stock) <= item.Quantity {
// 			return fmt.Errorf("product %d is not available in the quantity requested", product.ProductID)
// 		}
// 	}
// 	return nil
// }

func calculateTotalPrice(cartItems []types.CartCheckoutItem, products map[int]types.Product) float64 {
	var total float64

	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}

// ACID Concept needs to be used

func (h *Handler) createOrder(products []types.Product, productsq []types.ProductStock, cartItems []types.CartCheckoutItem, userID int) (int, float64, error) {
	//For types.Products
	productsMap := make(map[int]types.Product)

	for _, product := range products {
		productsMap[product.ID] = product
	}

	stocksMap := make(map[int]types.ProductStock)

	for _, productq := range productsq {
		stocksMap[productq.ProductID] = productq
	}
	//Check if all Products are available
	// if err := checkIfCartIsInStock(cartItems, stocksMap); err != nil {
	// 	return 0, 0, err
	// }

	err := h.store.UpdateAndRestock(cartItems)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to update stock: %v", err)
	}
	//Calculate the total Price
	totalPrice := calculateTotalPrice(cartItems, productsMap)

	//reduce the quantity of products in the store
	//Better solution as multiple request and no ACID Applyy ACID
	// for _, item := range cartItems {
	// 	stock := stocksMap[item.ProductID]
	// 	stock.Stock -= item.Quantity
	// 	h.store.UpdateProductStock(stock)
	// }
	//h.store.UpdateProduct(stock, cartItems)
	orderID, err := h.orderStore.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address", // fetch address from a user addresses table
	})
	if err != nil {
		return 0, 0, err
	}

	// create order the items records
	for _, item := range cartItems {
		h.orderStore.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productsMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}
