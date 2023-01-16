// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repo_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "OnlineShopBackend/internal/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockItemStore is a mock of ItemStore interface.
type MockItemStore struct {
	ctrl     *gomock.Controller
	recorder *MockItemStoreMockRecorder
}

// MockItemStoreMockRecorder is the mock recorder for MockItemStore.
type MockItemStoreMockRecorder struct {
	mock *MockItemStore
}

// NewMockItemStore creates a new mock instance.
func NewMockItemStore(ctrl *gomock.Controller) *MockItemStore {
	mock := &MockItemStore{ctrl: ctrl}
	mock.recorder = &MockItemStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockItemStore) EXPECT() *MockItemStoreMockRecorder {
	return m.recorder
}

// CreateItem mocks base method.
func (m *MockItemStore) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateItem", ctx, item)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateItem indicates an expected call of CreateItem.
func (mr *MockItemStoreMockRecorder) CreateItem(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateItem", reflect.TypeOf((*MockItemStore)(nil).CreateItem), ctx, item)
}

// DeleteItem mocks base method.
func (m *MockItemStore) DeleteItem(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockItemStoreMockRecorder) DeleteItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockItemStore)(nil).DeleteItem), ctx, id)
}

// GetItem mocks base method.
func (m *MockItemStore) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", ctx, id)
	ret0, _ := ret[0].(*models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockItemStoreMockRecorder) GetItem(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockItemStore)(nil).GetItem), ctx, id)
}

// GetItemsByCategory mocks base method.
func (m *MockItemStore) GetItemsByCategory(ctx context.Context, categoryName string) (chan models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItemsByCategory", ctx, categoryName)
	ret0, _ := ret[0].(chan models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItemsByCategory indicates an expected call of GetItemsByCategory.
func (mr *MockItemStoreMockRecorder) GetItemsByCategory(ctx, categoryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItemsByCategory", reflect.TypeOf((*MockItemStore)(nil).GetItemsByCategory), ctx, categoryName)
}

// ItemsList mocks base method.
func (m *MockItemStore) ItemsList(ctx context.Context) (chan models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsList", ctx)
	ret0, _ := ret[0].(chan models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ItemsList indicates an expected call of ItemsList.
func (mr *MockItemStoreMockRecorder) ItemsList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsList", reflect.TypeOf((*MockItemStore)(nil).ItemsList), ctx)
}

// SearchLine mocks base method.
func (m *MockItemStore) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchLine", ctx, param)
	ret0, _ := ret[0].(chan models.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchLine indicates an expected call of SearchLine.
func (mr *MockItemStoreMockRecorder) SearchLine(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchLine", reflect.TypeOf((*MockItemStore)(nil).SearchLine), ctx, param)
}

// UpdateItem mocks base method.
func (m *MockItemStore) UpdateItem(ctx context.Context, item *models.Item) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", ctx, item)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockItemStoreMockRecorder) UpdateItem(ctx, item interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockItemStore)(nil).UpdateItem), ctx, item)
}

// MockCategoryStore is a mock of CategoryStore interface.
type MockCategoryStore struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryStoreMockRecorder
}

// MockCategoryStoreMockRecorder is the mock recorder for MockCategoryStore.
type MockCategoryStoreMockRecorder struct {
	mock *MockCategoryStore
}

// NewMockCategoryStore creates a new mock instance.
func NewMockCategoryStore(ctrl *gomock.Controller) *MockCategoryStore {
	mock := &MockCategoryStore{ctrl: ctrl}
	mock.recorder = &MockCategoryStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryStore) EXPECT() *MockCategoryStoreMockRecorder {
	return m.recorder
}

// CreateCategory mocks base method.
func (m *MockCategoryStore) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", ctx, category)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory.
func (mr *MockCategoryStoreMockRecorder) CreateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockCategoryStore)(nil).CreateCategory), ctx, category)
}

// DeleteCategory mocks base method.
func (m *MockCategoryStore) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory.
func (mr *MockCategoryStoreMockRecorder) DeleteCategory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockCategoryStore)(nil).DeleteCategory), ctx, id)
}

// GetCategory mocks base method.
func (m *MockCategoryStore) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", ctx, id)
	ret0, _ := ret[0].(*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory.
func (mr *MockCategoryStoreMockRecorder) GetCategory(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockCategoryStore)(nil).GetCategory), ctx, id)
}

// GetCategoryByName mocks base method.
func (m *MockCategoryStore) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryByName", ctx, name)
	ret0, _ := ret[0].(*models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryByName indicates an expected call of GetCategoryByName.
func (mr *MockCategoryStoreMockRecorder) GetCategoryByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryByName", reflect.TypeOf((*MockCategoryStore)(nil).GetCategoryByName), ctx, name)
}

// GetCategoryList mocks base method.
func (m *MockCategoryStore) GetCategoryList(ctx context.Context) (chan models.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategoryList", ctx)
	ret0, _ := ret[0].(chan models.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategoryList indicates an expected call of GetCategoryList.
func (mr *MockCategoryStoreMockRecorder) GetCategoryList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategoryList", reflect.TypeOf((*MockCategoryStore)(nil).GetCategoryList), ctx)
}

// UpdateCategory mocks base method.
func (m *MockCategoryStore) UpdateCategory(ctx context.Context, category *models.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory.
func (mr *MockCategoryStoreMockRecorder) UpdateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockCategoryStore)(nil).UpdateCategory), ctx, category)
}

// MockUserStore is a mock of UserStore interface.
type MockUserStore struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoreMockRecorder
}

// MockUserStoreMockRecorder is the mock recorder for MockUserStore.
type MockUserStoreMockRecorder struct {
	mock *MockUserStore
}

// NewMockUserStore creates a new mock instance.
func NewMockUserStore(ctrl *gomock.Controller) *MockUserStore {
	mock := &MockUserStore{ctrl: ctrl}
	mock.recorder = &MockUserStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserStore) EXPECT() *MockUserStoreMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserStore) Create(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserStoreMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserStore)(nil).Create), ctx, user)
}

// GetRightsId mocks base method.
func (m *MockUserStore) GetRightsId(ctx context.Context, name string) (models.Rights, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRightsId", ctx, name)
	ret0, _ := ret[0].(models.Rights)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRightsId indicates an expected call of GetRightsId.
func (mr *MockUserStoreMockRecorder) GetRightsId(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRightsId", reflect.TypeOf((*MockUserStore)(nil).GetRightsId), ctx, name)
}

// GetUserByEmail mocks base method.
func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserStoreMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserStore)(nil).GetUserByEmail), ctx, email)
}

// SaveSession mocks base method.
func (m *MockUserStore) SaveSession(ctx context.Context, token string, t int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSession", ctx, token, t)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSession indicates an expected call of SaveSession.
func (mr *MockUserStoreMockRecorder) SaveSession(ctx, token, t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSession", reflect.TypeOf((*MockUserStore)(nil).SaveSession), ctx, token, t)
}

// UpdateUserData mocks base method.
func (m *MockUserStore) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserData", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserData indicates an expected call of UpdateUserData.
func (mr *MockUserStoreMockRecorder) UpdateUserData(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserData", reflect.TypeOf((*MockUserStore)(nil).UpdateUserData), ctx, user)
}

// MockCartStore is a mock of CartStore interface.
type MockCartStore struct {
	ctrl     *gomock.Controller
	recorder *MockCartStoreMockRecorder
}

// MockCartStoreMockRecorder is the mock recorder for MockCartStore.
type MockCartStoreMockRecorder struct {
	mock *MockCartStore
}

// NewMockCartStore creates a new mock instance.
func NewMockCartStore(ctrl *gomock.Controller) *MockCartStore {
	mock := &MockCartStore{ctrl: ctrl}
	mock.recorder = &MockCartStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCartStore) EXPECT() *MockCartStoreMockRecorder {
	return m.recorder
}

// AddItemToCart mocks base method.
func (m *MockCartStore) AddItemToCart(ctx context.Context, cartId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItemToCart", ctx, cartId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddItemToCart indicates an expected call of AddItemToCart.
func (mr *MockCartStoreMockRecorder) AddItemToCart(ctx, cartId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItemToCart", reflect.TypeOf((*MockCartStore)(nil).AddItemToCart), ctx, cartId, itemId)
}

// Create mocks base method.
func (m *MockCartStore) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, userId)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCartStoreMockRecorder) Create(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCartStore)(nil).Create), ctx, userId)
}

// DeleteCart mocks base method.
func (m *MockCartStore) DeleteCart(ctx context.Context, cartId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCart", ctx, cartId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCart indicates an expected call of DeleteCart.
func (mr *MockCartStoreMockRecorder) DeleteCart(ctx, cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCart", reflect.TypeOf((*MockCartStore)(nil).DeleteCart), ctx, cartId)
}

// DeleteItemFromCart mocks base method.
func (m *MockCartStore) DeleteItemFromCart(ctx context.Context, cartId, itemId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItemFromCart", ctx, cartId, itemId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItemFromCart indicates an expected call of DeleteItemFromCart.
func (mr *MockCartStoreMockRecorder) DeleteItemFromCart(ctx, cartId, itemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItemFromCart", reflect.TypeOf((*MockCartStore)(nil).DeleteItemFromCart), ctx, cartId, itemId)
}

// GetCart mocks base method.
func (m *MockCartStore) GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCart", ctx, cartId)
	ret0, _ := ret[0].(*models.Cart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCart indicates an expected call of GetCart.
func (mr *MockCartStoreMockRecorder) GetCart(ctx, cartId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCart", reflect.TypeOf((*MockCartStore)(nil).GetCart), ctx, cartId)
}

// MockOrderStore is a mock of OrderStore interface.
type MockOrderStore struct {
	ctrl     *gomock.Controller
	recorder *MockOrderStoreMockRecorder
}

// MockOrderStoreMockRecorder is the mock recorder for MockOrderStore.
type MockOrderStoreMockRecorder struct {
	mock *MockOrderStore
}

// NewMockOrderStore creates a new mock instance.
func NewMockOrderStore(ctrl *gomock.Controller) *MockOrderStore {
	mock := &MockOrderStore{ctrl: ctrl}
	mock.recorder = &MockOrderStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderStore) EXPECT() *MockOrderStoreMockRecorder {
	return m.recorder
}

// ChangeAddress mocks base method.
func (m *MockOrderStore) ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAddress", ctx, order, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeAddress indicates an expected call of ChangeAddress.
func (mr *MockOrderStoreMockRecorder) ChangeAddress(ctx, order, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAddress", reflect.TypeOf((*MockOrderStore)(nil).ChangeAddress), ctx, order, address)
}

// ChangeStatus mocks base method.
func (m *MockOrderStore) ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", ctx, order, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus.
func (mr *MockOrderStoreMockRecorder) ChangeStatus(ctx, order, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockOrderStore)(nil).ChangeStatus), ctx, order, status)
}

// Create mocks base method.
func (m *MockOrderStore) Create(ctx context.Context, order *models.Order) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, order)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockOrderStoreMockRecorder) Create(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderStore)(nil).Create), ctx, order)
}

// DeleteOrder mocks base method.
func (m *MockOrderStore) DeleteOrder(ctx context.Context, order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockOrderStoreMockRecorder) DeleteOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockOrderStore)(nil).DeleteOrder), ctx, order)
}

// GetOrderByID mocks base method.
func (m *MockOrderStore) GetOrderByID(ctx context.Context, id uuid.UUID) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByID", ctx, id)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByID indicates an expected call of GetOrderByID.
func (mr *MockOrderStoreMockRecorder) GetOrderByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByID", reflect.TypeOf((*MockOrderStore)(nil).GetOrderByID), ctx, id)
}

// GetOrdersForUser mocks base method.
func (m *MockOrderStore) GetOrdersForUser(ctx context.Context, user *models.User) (chan models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersForUser", ctx, user)
	ret0, _ := ret[0].(chan models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersForUser indicates an expected call of GetOrdersForUser.
func (mr *MockOrderStoreMockRecorder) GetOrdersForUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersForUser", reflect.TypeOf((*MockOrderStore)(nil).GetOrdersForUser), ctx, user)
}
