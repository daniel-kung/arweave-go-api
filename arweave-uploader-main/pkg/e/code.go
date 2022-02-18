package e

// ERRCode error code of this api server
//
// 0 success
// 1 to 999, reserved
// 1000 to 1999, client error
// 2000 to 2999, server error
type ERRCode int

const (
	ERRCodeSuccess ERRCode = 0 // success
)

const (
	ERRCodeInvalidParam ERRCode = 1000 // invalid param
)

const (
	ERRCodeChainConnection ERRCode = 2000 // chain connection
)
