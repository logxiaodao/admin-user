package main

import (
	"admin/user/rpc/test"
	pb "admin/user/rpc/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "localhost:8002"
)

// 这里目前是手写的测试用例
func main() {
	// Set up a connection to the server.
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzYxODUzODIsImlhdCI6MTYzNjA5ODk4MiwiaXNTdXBlckFkbWluIjoxLCJwbGF0Zm9ybUlEIjoxLCJ1c2VySWQiOjF9.Sdui9-yCXVBQRqE8LwM4zZ_1VAIoEAlvrV1OrIYffJk"
	//token   :=  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzYxODYyOTEsImlhdCI6MTYzNjA5OTg5MSwiaXNTdXBlckFkbWluIjowLCJwbGF0Zm9ybUlEIjoxLCJ1c2VySWQiOjMxfQ.axEkB_l9O9uUZmNCSmjFvV2gj_YX3dr5HcBQlmqL1PQ"

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(test.AuthToekn{token}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserClient(conn)

	// user
	//data, err := c.Login(context.Background(), &pb.LoginRequest{
	//	Username:   "lxj2",
	//	Password:   "123456789",
	//	PlatformID: 1,
	//})

	// account
	//data, err := c.LoginOut(context.Background(), &pb.LoginOutRequest{})

	//data, err := c.UpdatePassword(context.Background(), &pb.UpdatePasswordRequest{OldPassword: "123456789", NewPassword: "12345678", ConfirmPassword: "12345678"})

	data, err := c.CheckPermission(context.Background(), &pb.CheckPermissionRequest{HttpPath: "/v1/auth/permission", HttpMethod: "POST"})

	// api
	//data, err := c.GetApi(context.Background(), &pb.GetApiRequest{CurrentPage: 1, PageSize: 10})

	//data, err := c.FindApiByIds(context.Background(), &pb.FindApiByIdsRequest{Ids: "1,2,3"})

	//data, err := c.AddApi(context.Background(), &pb.AddApiRequest{Name: "测试", HttpMethod: "POST", HttpPath: "/test"})

	//data, err := c.BatchApi(context.Background(), &pb.AddBatchApiRequest{ItemList: [] *pb.AddApiRequest{
	//	{
	//		Name: "批量测试1",
	//		HttpMethod: "POST",
	//		HttpPath: "/test1",
	//	},
	//	{
	//		Name: "批量测试2",
	//		HttpMethod: "POST",
	//		HttpPath: "/test2",
	//	},
	//}})

	//data, err := c.EditApi(context.Background(), &pb.EditApiRequest{Id: 1614, Name: "测试11", HttpMethod: "POST", HttpPath: "/test11"})

	//data, err := c.DeleteApi(context.Background(), &pb.DeleteApiRequest{IdList: []uint64{1614,1615}})

	// permission
	//data, err := c.GetPermission(context.Background(), &pb.GetPermissionRequest{CurrentPage: 1, PageSize: 10})
	//
	//data, err := c.FindPermissionByIds(context.Background(), &pb.FindPermissionByIdsRequest{Ids: "1,2,3"})
	//
	//data, err := c.AddPermission(context.Background(), &pb.AddPermissionRequest{
	//	PermissionName: "测试",
	//	ApiIdList:      []uint64{1,2,3,4},
	//})
	//
	//data, err := c.EditPermission(context.Background(), &pb.EditPermissionRequest{
	//	Id: 4,
	//	PermissionName: "测试2",
	//	ApiIdList:      []uint64{1,2,3,4,5},
	//})
	//
	//data, err := c.DeletePermission(context.Background(), &pb.DeletePermissionRequest{IdList: []uint64{4}})

	// role
	//data, err := c.GetRole(context.Background(), &pb.GetRoleRequest{CurrentPage: 1, PageSize: 10})

	//data, err := c.FindRoleByIds(context.Background(), &pb.FindRoleByIdsRequest{Ids: "1,2"})

	//data, err := c.AddRole(context.Background(), &pb.AddRoleRequest{
	//	RoleName:         "测试",
	//	PermissionIdList: []uint64{3,2},
	//})

	//data, err := c.EditRole(context.Background(), &pb.EditRoleRequest{
	//	Id: 5,
	//	RoleName:         "测试",
	//	PermissionIdList: []uint64{3},
	//})
	//
	//data, err := c.DeleteRole(context.Background(), &pb.DeleteRoleRequest{
	//	IdList: []uint64{5},
	//})

	// admin
	//data, err := c.GetAdmin(context.Background(), &pb.GetAdminRequest{CurrentPage: 1, PageSize: 10})

	//data, err := c.FindAdminByIds(context.Background(), &pb.FindAdminByIdsRequest{Ids: "1,31"})

	//data, err := c.AddAdmin(context.Background(), &pb.AddAdminRequest{
	//	Password:   "12345678",
	//	Account:    "xiaolong",
	//	NickName:   "xiaolong",
	//	Phone:      "13698814566",
	//	Email:      "1234@11.com",
	//	RoleIdList: [] uint64{1,2,3},
	//})

	//data, err := c.EditAdmin(context.Background(), &pb.EditAdminRequest{
	//	Id: 115,
	//	Account:    "xiaolong1",
	//	NickName:   "xiaolong1",
	//	Phone:      "13698814566",
	//	Email:      "1234@11.com",
	//	RoleIdList: [] uint64{1,2},
	//})

	//data, err := c.DeleteAdmin(context.Background(), &pb.DeleteAdminRequest{IdList: [] uint64{115}})

	checkError(err)
	fmt.Println(data)

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("发生错误:", err)
	}
}
