// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: post/post.proto

package post

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Post struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	CreatorId     uint64                 `protobuf:"varint,4,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	IsPrivate     bool                   `protobuf:"varint,7,opt,name=is_private,json=isPrivate,proto3" json:"is_private,omitempty"`
	Tags          []string               `protobuf:"bytes,8,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Post) Reset() {
	*x = Post{}
	mi := &file_post_post_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{0}
}

func (x *Post) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Post) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Post) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Post) GetCreatorId() uint64 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

func (x *Post) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Post) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Post) GetIsPrivate() bool {
	if x != nil {
		return x.IsPrivate
	}
	return false
}

func (x *Post) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type CreatePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	CreatorId     uint64                 `protobuf:"varint,3,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	IsPrivate     bool                   `protobuf:"varint,4,opt,name=is_private,json=isPrivate,proto3" json:"is_private,omitempty"`
	Tags          []string               `protobuf:"bytes,5,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreatePostRequest) Reset() {
	*x = CreatePostRequest{}
	mi := &file_post_post_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreatePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePostRequest) ProtoMessage() {}

func (x *CreatePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePostRequest.ProtoReflect.Descriptor instead.
func (*CreatePostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePostRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreatePostRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreatePostRequest) GetCreatorId() uint64 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

func (x *CreatePostRequest) GetIsPrivate() bool {
	if x != nil {
		return x.IsPrivate
	}
	return false
}

func (x *CreatePostRequest) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type GetPostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	RequesterId   uint64                 `protobuf:"varint,2,opt,name=requester_id,json=requesterId,proto3" json:"requester_id,omitempty"` // For checking privacy permissions
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPostRequest) Reset() {
	*x = GetPostRequest{}
	mi := &file_post_post_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPostRequest) ProtoMessage() {}

func (x *GetPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPostRequest.ProtoReflect.Descriptor instead.
func (*GetPostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{2}
}

func (x *GetPostRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GetPostRequest) GetRequesterId() uint64 {
	if x != nil {
		return x.RequesterId
	}
	return 0
}

type UpdatePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	CreatorId     uint64                 `protobuf:"varint,4,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"` // User requesting the update
	IsPrivate     bool                   `protobuf:"varint,5,opt,name=is_private,json=isPrivate,proto3" json:"is_private,omitempty"`
	Tags          []string               `protobuf:"bytes,6,rep,name=tags,proto3" json:"tags,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdatePostRequest) Reset() {
	*x = UpdatePostRequest{}
	mi := &file_post_post_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdatePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdatePostRequest) ProtoMessage() {}

func (x *UpdatePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdatePostRequest.ProtoReflect.Descriptor instead.
func (*UpdatePostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{3}
}

func (x *UpdatePostRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdatePostRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *UpdatePostRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdatePostRequest) GetCreatorId() uint64 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

func (x *UpdatePostRequest) GetIsPrivate() bool {
	if x != nil {
		return x.IsPrivate
	}
	return false
}

func (x *UpdatePostRequest) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

type DeletePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatorId     uint64                 `protobuf:"varint,2,opt,name=creator_id,json=creatorId,proto3" json:"creator_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeletePostRequest) Reset() {
	*x = DeletePostRequest{}
	mi := &file_post_post_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeletePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePostRequest) ProtoMessage() {}

func (x *DeletePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePostRequest.ProtoReflect.Descriptor instead.
func (*DeletePostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{4}
}

func (x *DeletePostRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *DeletePostRequest) GetCreatorId() uint64 {
	if x != nil {
		return x.CreatorId
	}
	return 0
}

type DeletePostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeletePostResponse) Reset() {
	*x = DeletePostResponse{}
	mi := &file_post_post_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeletePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePostResponse) ProtoMessage() {}

func (x *DeletePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePostResponse.ProtoReflect.Descriptor instead.
func (*DeletePostResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{5}
}

func (x *DeletePostResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type ListPostsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Page          int32                  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	RequesterId   uint64                 `protobuf:"varint,3,opt,name=requester_id,json=requesterId,proto3" json:"requester_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPostsRequest) Reset() {
	*x = ListPostsRequest{}
	mi := &file_post_post_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPostsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPostsRequest) ProtoMessage() {}

func (x *ListPostsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPostsRequest.ProtoReflect.Descriptor instead.
func (*ListPostsRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{6}
}

func (x *ListPostsRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ListPostsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListPostsRequest) GetRequesterId() uint64 {
	if x != nil {
		return x.RequesterId
	}
	return 0
}

type ListPostsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Posts         []*Post                `protobuf:"bytes,1,rep,name=posts,proto3" json:"posts,omitempty"`
	TotalCount    uint64                 `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPostsResponse) Reset() {
	*x = ListPostsResponse{}
	mi := &file_post_post_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPostsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPostsResponse) ProtoMessage() {}

func (x *ListPostsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPostsResponse.ProtoReflect.Descriptor instead.
func (*ListPostsResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{7}
}

func (x *ListPostsResponse) GetPosts() []*Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

func (x *ListPostsResponse) GetTotalCount() uint64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

type PostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Post          *Post                  `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PostResponse) Reset() {
	*x = PostResponse{}
	mi := &file_post_post_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostResponse) ProtoMessage() {}

func (x *PostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostResponse.ProtoReflect.Descriptor instead.
func (*PostResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{8}
}

func (x *PostResponse) GetPost() *Post {
	if x != nil {
		return x.Post
	}
	return nil
}

type ViewPostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ViewPostRequest) Reset() {
	*x = ViewPostRequest{}
	mi := &file_post_post_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ViewPostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ViewPostRequest) ProtoMessage() {}

func (x *ViewPostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ViewPostRequest.ProtoReflect.Descriptor instead.
func (*ViewPostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{9}
}

func (x *ViewPostRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ViewPostRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ViewPostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ViewPostResponse) Reset() {
	*x = ViewPostResponse{}
	mi := &file_post_post_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ViewPostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ViewPostResponse) ProtoMessage() {}

func (x *ViewPostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ViewPostResponse.ProtoReflect.Descriptor instead.
func (*ViewPostResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{10}
}

func (x *ViewPostResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type LikePostRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LikePostRequest) Reset() {
	*x = LikePostRequest{}
	mi := &file_post_post_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LikePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikePostRequest) ProtoMessage() {}

func (x *LikePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikePostRequest.ProtoReflect.Descriptor instead.
func (*LikePostRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{11}
}

func (x *LikePostRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LikePostRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type LikePostResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	TotalLikes    uint64                 `protobuf:"varint,2,opt,name=total_likes,json=totalLikes,proto3" json:"total_likes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LikePostResponse) Reset() {
	*x = LikePostResponse{}
	mi := &file_post_post_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LikePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikePostResponse) ProtoMessage() {}

func (x *LikePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikePostResponse.ProtoReflect.Descriptor instead.
func (*LikePostResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{12}
}

func (x *LikePostResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *LikePostResponse) GetTotalLikes() uint64 {
	if x != nil {
		return x.TotalLikes
	}
	return 0
}

type Comment struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            uint64                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	PostId        uint64                 `protobuf:"varint,2,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	UserId        uint64                 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Comment) Reset() {
	*x = Comment{}
	mi := &file_post_post_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Comment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Comment) ProtoMessage() {}

func (x *Comment) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Comment.ProtoReflect.Descriptor instead.
func (*Comment) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{13}
}

func (x *Comment) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Comment) GetPostId() uint64 {
	if x != nil {
		return x.PostId
	}
	return 0
}

func (x *Comment) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Comment) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Comment) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type AddCommentRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PostId        uint64                 `protobuf:"varint,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	UserId        uint64                 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string                 `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AddCommentRequest) Reset() {
	*x = AddCommentRequest{}
	mi := &file_post_post_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AddCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddCommentRequest) ProtoMessage() {}

func (x *AddCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddCommentRequest.ProtoReflect.Descriptor instead.
func (*AddCommentRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{14}
}

func (x *AddCommentRequest) GetPostId() uint64 {
	if x != nil {
		return x.PostId
	}
	return 0
}

func (x *AddCommentRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AddCommentRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type CommentResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comment       *Comment               `protobuf:"bytes,1,opt,name=comment,proto3" json:"comment,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CommentResponse) Reset() {
	*x = CommentResponse{}
	mi := &file_post_post_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommentResponse) ProtoMessage() {}

func (x *CommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommentResponse.ProtoReflect.Descriptor instead.
func (*CommentResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{15}
}

func (x *CommentResponse) GetComment() *Comment {
	if x != nil {
		return x.Comment
	}
	return nil
}

type GetCommentsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	PostId        uint64                 `protobuf:"varint,1,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	Page          uint64                 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      uint64                 `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	RequesterId   uint64                 `protobuf:"varint,4,opt,name=requester_id,json=requesterId,proto3" json:"requester_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCommentsRequest) Reset() {
	*x = GetCommentsRequest{}
	mi := &file_post_post_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCommentsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentsRequest) ProtoMessage() {}

func (x *GetCommentsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentsRequest.ProtoReflect.Descriptor instead.
func (*GetCommentsRequest) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{16}
}

func (x *GetCommentsRequest) GetPostId() uint64 {
	if x != nil {
		return x.PostId
	}
	return 0
}

func (x *GetCommentsRequest) GetPage() uint64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetCommentsRequest) GetPageSize() uint64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetCommentsRequest) GetRequesterId() uint64 {
	if x != nil {
		return x.RequesterId
	}
	return 0
}

type GetCommentsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Comments      []*Comment             `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
	TotalCount    uint64                 `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetCommentsResponse) Reset() {
	*x = GetCommentsResponse{}
	mi := &file_post_post_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetCommentsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCommentsResponse) ProtoMessage() {}

func (x *GetCommentsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_post_post_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCommentsResponse.ProtoReflect.Descriptor instead.
func (*GetCommentsResponse) Descriptor() ([]byte, []int) {
	return file_post_post_proto_rawDescGZIP(), []int{17}
}

func (x *GetCommentsResponse) GetComments() []*Comment {
	if x != nil {
		return x.Comments
	}
	return nil
}

func (x *GetCommentsResponse) GetTotalCount() uint64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

var File_post_post_proto protoreflect.FileDescriptor

const file_post_post_proto_rawDesc = "" +
	"\n" +
	"\x0fpost/post.proto\x12\x04post\x1a\x1fgoogle/protobuf/timestamp.proto\"\x96\x02\n" +
	"\x04Post\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x1d\n" +
	"\n" +
	"creator_id\x18\x04 \x01(\x04R\tcreatorId\x129\n" +
	"\n" +
	"created_at\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt\x12\x1d\n" +
	"\n" +
	"is_private\x18\a \x01(\bR\tisPrivate\x12\x12\n" +
	"\x04tags\x18\b \x03(\tR\x04tags\"\x9d\x01\n" +
	"\x11CreatePostRequest\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12\x1d\n" +
	"\n" +
	"creator_id\x18\x03 \x01(\x04R\tcreatorId\x12\x1d\n" +
	"\n" +
	"is_private\x18\x04 \x01(\bR\tisPrivate\x12\x12\n" +
	"\x04tags\x18\x05 \x03(\tR\x04tags\"C\n" +
	"\x0eGetPostRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12!\n" +
	"\frequester_id\x18\x02 \x01(\x04R\vrequesterId\"\xad\x01\n" +
	"\x11UpdatePostRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x14\n" +
	"\x05title\x18\x02 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x1d\n" +
	"\n" +
	"creator_id\x18\x04 \x01(\x04R\tcreatorId\x12\x1d\n" +
	"\n" +
	"is_private\x18\x05 \x01(\bR\tisPrivate\x12\x12\n" +
	"\x04tags\x18\x06 \x03(\tR\x04tags\"B\n" +
	"\x11DeletePostRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x1d\n" +
	"\n" +
	"creator_id\x18\x02 \x01(\x04R\tcreatorId\".\n" +
	"\x12DeletePostResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\"f\n" +
	"\x10ListPostsRequest\x12\x12\n" +
	"\x04page\x18\x01 \x01(\x05R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x02 \x01(\x05R\bpageSize\x12!\n" +
	"\frequester_id\x18\x03 \x01(\x04R\vrequesterId\"V\n" +
	"\x11ListPostsResponse\x12 \n" +
	"\x05posts\x18\x01 \x03(\v2\n" +
	".post.PostR\x05posts\x12\x1f\n" +
	"\vtotal_count\x18\x02 \x01(\x04R\n" +
	"totalCount\".\n" +
	"\fPostResponse\x12\x1e\n" +
	"\x04post\x18\x01 \x01(\v2\n" +
	".post.PostR\x04post\":\n" +
	"\x0fViewPostRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x04R\x06userId\",\n" +
	"\x10ViewPostResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\":\n" +
	"\x0fLikePostRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x04R\x06userId\"M\n" +
	"\x10LikePostResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x1f\n" +
	"\vtotal_likes\x18\x02 \x01(\x04R\n" +
	"totalLikes\"\xa0\x01\n" +
	"\aComment\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x04R\x02id\x12\x17\n" +
	"\apost_id\x18\x02 \x01(\x04R\x06postId\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\x04R\x06userId\x12\x18\n" +
	"\acontent\x18\x04 \x01(\tR\acontent\x129\n" +
	"\n" +
	"created_at\x18\x05 \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\"_\n" +
	"\x11AddCommentRequest\x12\x17\n" +
	"\apost_id\x18\x01 \x01(\x04R\x06postId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x04R\x06userId\x12\x18\n" +
	"\acontent\x18\x03 \x01(\tR\acontent\":\n" +
	"\x0fCommentResponse\x12'\n" +
	"\acomment\x18\x01 \x01(\v2\r.post.CommentR\acomment\"\x81\x01\n" +
	"\x12GetCommentsRequest\x12\x17\n" +
	"\apost_id\x18\x01 \x01(\x04R\x06postId\x12\x12\n" +
	"\x04page\x18\x02 \x01(\x04R\x04page\x12\x1b\n" +
	"\tpage_size\x18\x03 \x01(\x04R\bpageSize\x12!\n" +
	"\frequester_id\x18\x04 \x01(\x04R\vrequesterId\"a\n" +
	"\x13GetCommentsResponse\x12)\n" +
	"\bcomments\x18\x01 \x03(\v2\r.post.CommentR\bcomments\x12\x1f\n" +
	"\vtotal_count\x18\x02 \x01(\x04R\n" +
	"totalCount2\xaf\x04\n" +
	"\vPostService\x129\n" +
	"\n" +
	"CreatePost\x12\x17.post.CreatePostRequest\x1a\x12.post.PostResponse\x123\n" +
	"\aGetPost\x12\x14.post.GetPostRequest\x1a\x12.post.PostResponse\x129\n" +
	"\n" +
	"UpdatePost\x12\x17.post.UpdatePostRequest\x1a\x12.post.PostResponse\x12?\n" +
	"\n" +
	"DeletePost\x12\x17.post.DeletePostRequest\x1a\x18.post.DeletePostResponse\x12<\n" +
	"\tListPosts\x12\x16.post.ListPostsRequest\x1a\x17.post.ListPostsResponse\x129\n" +
	"\bViewPost\x12\x15.post.ViewPostRequest\x1a\x16.post.ViewPostResponse\x129\n" +
	"\bLikePost\x12\x15.post.LikePostRequest\x1a\x16.post.LikePostResponse\x12<\n" +
	"\n" +
	"AddComment\x12\x17.post.AddCommentRequest\x1a\x15.post.CommentResponse\x12B\n" +
	"\vGetComments\x12\x18.post.GetCommentsRequest\x1a\x19.post.GetCommentsResponseB\"Z social-network/common/proto/postb\x06proto3"

var (
	file_post_post_proto_rawDescOnce sync.Once
	file_post_post_proto_rawDescData []byte
)

func file_post_post_proto_rawDescGZIP() []byte {
	file_post_post_proto_rawDescOnce.Do(func() {
		file_post_post_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_post_post_proto_rawDesc), len(file_post_post_proto_rawDesc)))
	})
	return file_post_post_proto_rawDescData
}

var file_post_post_proto_msgTypes = make([]protoimpl.MessageInfo, 18)
var file_post_post_proto_goTypes = []any{
	(*Post)(nil),                  // 0: post.Post
	(*CreatePostRequest)(nil),     // 1: post.CreatePostRequest
	(*GetPostRequest)(nil),        // 2: post.GetPostRequest
	(*UpdatePostRequest)(nil),     // 3: post.UpdatePostRequest
	(*DeletePostRequest)(nil),     // 4: post.DeletePostRequest
	(*DeletePostResponse)(nil),    // 5: post.DeletePostResponse
	(*ListPostsRequest)(nil),      // 6: post.ListPostsRequest
	(*ListPostsResponse)(nil),     // 7: post.ListPostsResponse
	(*PostResponse)(nil),          // 8: post.PostResponse
	(*ViewPostRequest)(nil),       // 9: post.ViewPostRequest
	(*ViewPostResponse)(nil),      // 10: post.ViewPostResponse
	(*LikePostRequest)(nil),       // 11: post.LikePostRequest
	(*LikePostResponse)(nil),      // 12: post.LikePostResponse
	(*Comment)(nil),               // 13: post.Comment
	(*AddCommentRequest)(nil),     // 14: post.AddCommentRequest
	(*CommentResponse)(nil),       // 15: post.CommentResponse
	(*GetCommentsRequest)(nil),    // 16: post.GetCommentsRequest
	(*GetCommentsResponse)(nil),   // 17: post.GetCommentsResponse
	(*timestamppb.Timestamp)(nil), // 18: google.protobuf.Timestamp
}
var file_post_post_proto_depIdxs = []int32{
	18, // 0: post.Post.created_at:type_name -> google.protobuf.Timestamp
	18, // 1: post.Post.updated_at:type_name -> google.protobuf.Timestamp
	0,  // 2: post.ListPostsResponse.posts:type_name -> post.Post
	0,  // 3: post.PostResponse.post:type_name -> post.Post
	18, // 4: post.Comment.created_at:type_name -> google.protobuf.Timestamp
	13, // 5: post.CommentResponse.comment:type_name -> post.Comment
	13, // 6: post.GetCommentsResponse.comments:type_name -> post.Comment
	1,  // 7: post.PostService.CreatePost:input_type -> post.CreatePostRequest
	2,  // 8: post.PostService.GetPost:input_type -> post.GetPostRequest
	3,  // 9: post.PostService.UpdatePost:input_type -> post.UpdatePostRequest
	4,  // 10: post.PostService.DeletePost:input_type -> post.DeletePostRequest
	6,  // 11: post.PostService.ListPosts:input_type -> post.ListPostsRequest
	9,  // 12: post.PostService.ViewPost:input_type -> post.ViewPostRequest
	11, // 13: post.PostService.LikePost:input_type -> post.LikePostRequest
	14, // 14: post.PostService.AddComment:input_type -> post.AddCommentRequest
	16, // 15: post.PostService.GetComments:input_type -> post.GetCommentsRequest
	8,  // 16: post.PostService.CreatePost:output_type -> post.PostResponse
	8,  // 17: post.PostService.GetPost:output_type -> post.PostResponse
	8,  // 18: post.PostService.UpdatePost:output_type -> post.PostResponse
	5,  // 19: post.PostService.DeletePost:output_type -> post.DeletePostResponse
	7,  // 20: post.PostService.ListPosts:output_type -> post.ListPostsResponse
	10, // 21: post.PostService.ViewPost:output_type -> post.ViewPostResponse
	12, // 22: post.PostService.LikePost:output_type -> post.LikePostResponse
	15, // 23: post.PostService.AddComment:output_type -> post.CommentResponse
	17, // 24: post.PostService.GetComments:output_type -> post.GetCommentsResponse
	16, // [16:25] is the sub-list for method output_type
	7,  // [7:16] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_post_post_proto_init() }
func file_post_post_proto_init() {
	if File_post_post_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_post_post_proto_rawDesc), len(file_post_post_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   18,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_post_post_proto_goTypes,
		DependencyIndexes: file_post_post_proto_depIdxs,
		MessageInfos:      file_post_post_proto_msgTypes,
	}.Build()
	File_post_post_proto = out.File
	file_post_post_proto_goTypes = nil
	file_post_post_proto_depIdxs = nil
}
