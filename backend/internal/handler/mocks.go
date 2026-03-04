package handler

import "backend/internal/types"

type MockProductStore struct {
	InsertProductErr	error
	FetchProductsErr	error
	DeleteProductErr	error
}
type MockUserStore struct{
	InsertUserErr	error
	LoginUserErr	error
}

func (m *MockProductStore) InsertProductForUser(userID int, productName string) (types.Product, error) {
	return types.Product{}, m.InsertProductErr
}
func (m *MockProductStore) FetchUserTrackedProducts(userID int) ([]types.UserProduct, error) {
	return []types.UserProduct{}, m.FetchProductsErr
}
func (m *MockProductStore) DeleteProductForUser(userID, productID int) error { return m.DeleteProductErr }

func (m *MockUserStore) InsertNewUser(username, email, password string) error { return m.InsertUserErr }
func (m *MockUserStore) LoginUser(username, password string) (string, error) { return "", m.LoginUserErr }

