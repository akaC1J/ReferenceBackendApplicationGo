// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: api/loms/v1/loms.proto

package loms

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string  `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	User   int64   `protobuf:"varint,2,opt,name=user,proto3" json:"user,omitempty"`
	Items  []*Item `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *Order) Reset() {
	*x = Order{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Order) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Order) ProtoMessage() {}

func (x *Order) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Order.ProtoReflect.Descriptor instead.
func (*Order) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{0}
}

func (x *Order) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *Order) GetUser() int64 {
	if x != nil {
		return x.User
	}
	return 0
}

func (x *Order) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sku   uint32 `protobuf:"varint,1,opt,name=sku,proto3" json:"sku,omitempty"`
	Count uint32 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{1}
}

func (x *Item) GetSku() uint32 {
	if x != nil {
		return x.Sku
	}
	return 0
}

func (x *Item) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

type OrderInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId int64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *OrderInfoRequest) Reset() {
	*x = OrderInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderInfoRequest) ProtoMessage() {}

func (x *OrderInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderInfoRequest.ProtoReflect.Descriptor instead.
func (*OrderInfoRequest) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{2}
}

func (x *OrderInfoRequest) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type OrderInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Order *Order `protobuf:"bytes,1,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *OrderInfoResponse) Reset() {
	*x = OrderInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderInfoResponse) ProtoMessage() {}

func (x *OrderInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderInfoResponse.ProtoReflect.Descriptor instead.
func (*OrderInfoResponse) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{3}
}

func (x *OrderInfoResponse) GetOrder() *Order {
	if x != nil {
		return x.Order
	}
	return nil
}

type OrderCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Order *Order `protobuf:"bytes,1,opt,name=order,proto3" json:"order,omitempty"`
}

func (x *OrderCreateRequest) Reset() {
	*x = OrderCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderCreateRequest) ProtoMessage() {}

func (x *OrderCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderCreateRequest.ProtoReflect.Descriptor instead.
func (*OrderCreateRequest) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{4}
}

func (x *OrderCreateRequest) GetOrder() *Order {
	if x != nil {
		return x.Order
	}
	return nil
}

type OrderCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId int64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *OrderCreateResponse) Reset() {
	*x = OrderCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderCreateResponse) ProtoMessage() {}

func (x *OrderCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderCreateResponse.ProtoReflect.Descriptor instead.
func (*OrderCreateResponse) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{5}
}

func (x *OrderCreateResponse) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type OrderPayRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId int64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *OrderPayRequest) Reset() {
	*x = OrderPayRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderPayRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderPayRequest) ProtoMessage() {}

func (x *OrderPayRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderPayRequest.ProtoReflect.Descriptor instead.
func (*OrderPayRequest) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{6}
}

func (x *OrderPayRequest) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type OrderCancelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId int64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *OrderCancelRequest) Reset() {
	*x = OrderCancelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OrderCancelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderCancelRequest) ProtoMessage() {}

func (x *OrderCancelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderCancelRequest.ProtoReflect.Descriptor instead.
func (*OrderCancelRequest) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{7}
}

func (x *OrderCancelRequest) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type StocksInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sku uint32 `protobuf:"varint,1,opt,name=sku,proto3" json:"sku,omitempty"`
}

func (x *StocksInfoRequest) Reset() {
	*x = StocksInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StocksInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StocksInfoRequest) ProtoMessage() {}

func (x *StocksInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StocksInfoRequest.ProtoReflect.Descriptor instead.
func (*StocksInfoRequest) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{8}
}

func (x *StocksInfoRequest) GetSku() uint32 {
	if x != nil {
		return x.Sku
	}
	return 0
}

type StocksInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint64 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *StocksInfoResponse) Reset() {
	*x = StocksInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_loms_v1_loms_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StocksInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StocksInfoResponse) ProtoMessage() {}

func (x *StocksInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_loms_v1_loms_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StocksInfoResponse.ProtoReflect.Descriptor instead.
func (*StocksInfoResponse) Descriptor() ([]byte, []int) {
	return file_api_loms_v1_loms_proto_rawDescGZIP(), []int{9}
}

func (x *StocksInfoResponse) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_api_loms_v1_loms_proto protoreflect.FileDescriptor

var file_api_loms_v1_loms_proto_rawDesc = []byte{
	0x0a, 0x16, 0x61, 0x70, 0x69, 0x2f, 0x6c, 0x6f, 0x6d, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6c, 0x6f,
	0x6d, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x50, 0x0a, 0x05, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x05, 0x69, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x2e, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12,
	0x10, 0x0a, 0x03, 0x73, 0x6b, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x73, 0x6b,
	0x75, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x2d, 0x0a, 0x10, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x31, 0x0a, 0x11, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x32, 0x0a, 0x12, 0x4f, 0x72, 0x64,
	0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x05, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x30, 0x0a,
	0x13, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x2c, 0x0a, 0x0f, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x50, 0x61, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2f, 0x0a,
	0x12, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x25,
	0x0a, 0x11, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x6b, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x03, 0x73, 0x6b, 0x75, 0x22, 0x2a, 0x0a, 0x12, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x32, 0xa7, 0x02, 0x0a, 0x04, 0x4c, 0x6f, 0x6d, 0x73, 0x12, 0x3a, 0x0a, 0x0b, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x36, 0x0a, 0x08, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x50,
	0x61, 0x79, 0x12, 0x10, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x50, 0x61, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x3c,
	0x0a, 0x0b, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x13, 0x2e,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x09,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x11, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x37, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x12, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x24, 0x5a, 0x22, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x32, 0x35, 0x36, 0x2f, 0x63, 0x61, 0x72, 0x74, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6c, 0x6f, 0x6d, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x6c, 0x6f, 0x6d,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_loms_v1_loms_proto_rawDescOnce sync.Once
	file_api_loms_v1_loms_proto_rawDescData = file_api_loms_v1_loms_proto_rawDesc
)

func file_api_loms_v1_loms_proto_rawDescGZIP() []byte {
	file_api_loms_v1_loms_proto_rawDescOnce.Do(func() {
		file_api_loms_v1_loms_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_loms_v1_loms_proto_rawDescData)
	})
	return file_api_loms_v1_loms_proto_rawDescData
}

var file_api_loms_v1_loms_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_loms_v1_loms_proto_goTypes = []any{
	(*Order)(nil),               // 0: Order
	(*Item)(nil),                // 1: Item
	(*OrderInfoRequest)(nil),    // 2: OrderInfoRequest
	(*OrderInfoResponse)(nil),   // 3: OrderInfoResponse
	(*OrderCreateRequest)(nil),  // 4: OrderCreateRequest
	(*OrderCreateResponse)(nil), // 5: OrderCreateResponse
	(*OrderPayRequest)(nil),     // 6: OrderPayRequest
	(*OrderCancelRequest)(nil),  // 7: OrderCancelRequest
	(*StocksInfoRequest)(nil),   // 8: StocksInfoRequest
	(*StocksInfoResponse)(nil),  // 9: StocksInfoResponse
	(*emptypb.Empty)(nil),       // 10: google.protobuf.Empty
}
var file_api_loms_v1_loms_proto_depIdxs = []int32{
	1,  // 0: Order.items:type_name -> Item
	0,  // 1: OrderInfoResponse.order:type_name -> Order
	0,  // 2: OrderCreateRequest.order:type_name -> Order
	4,  // 3: Loms.OrderCreate:input_type -> OrderCreateRequest
	6,  // 4: Loms.OrderPay:input_type -> OrderPayRequest
	7,  // 5: Loms.OrderCancel:input_type -> OrderCancelRequest
	2,  // 6: Loms.OrderInfo:input_type -> OrderInfoRequest
	8,  // 7: Loms.StocksInfo:input_type -> StocksInfoRequest
	5,  // 8: Loms.OrderCreate:output_type -> OrderCreateResponse
	10, // 9: Loms.OrderPay:output_type -> google.protobuf.Empty
	10, // 10: Loms.OrderCancel:output_type -> google.protobuf.Empty
	3,  // 11: Loms.OrderInfo:output_type -> OrderInfoResponse
	9,  // 12: Loms.StocksInfo:output_type -> StocksInfoResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_api_loms_v1_loms_proto_init() }
func file_api_loms_v1_loms_proto_init() {
	if File_api_loms_v1_loms_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_loms_v1_loms_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Order); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Item); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*OrderInfoRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*OrderInfoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*OrderCreateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*OrderCreateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*OrderPayRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*OrderCancelRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*StocksInfoRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_loms_v1_loms_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*StocksInfoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_loms_v1_loms_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_loms_v1_loms_proto_goTypes,
		DependencyIndexes: file_api_loms_v1_loms_proto_depIdxs,
		MessageInfos:      file_api_loms_v1_loms_proto_msgTypes,
	}.Build()
	File_api_loms_v1_loms_proto = out.File
	file_api_loms_v1_loms_proto_rawDesc = nil
	file_api_loms_v1_loms_proto_goTypes = nil
	file_api_loms_v1_loms_proto_depIdxs = nil
}
