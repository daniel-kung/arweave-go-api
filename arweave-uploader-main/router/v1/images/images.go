package images

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/everFinance/goar/types"
	"github.com/everFinance/goar/utils"

	//"github.com/Dev43/arweave-go/tx"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"ccian.cc/really/arweave-api/middleware/logger"
	"ccian.cc/really/arweave-api/pkg/e"
	"ccian.cc/really/arweave-api/pkg/setting"
	"github.com/everFinance/goar"
	//"github.com/Dev43/arweave-go/wallet"
	"github.com/gin-gonic/gin"
	//"github.com/everFinance/goar/utils"
	//"github.com/everFinance/goar/types"
)

type SupportedImageSuffix string

const (
	SupportedImageSuffixJPEG SupportedImageSuffix = ".jpeg"
	SupportedImageSuffixJPG  SupportedImageSuffix = ".jpg"
	SupportedImageSuffixPNG  SupportedImageSuffix = ".png"
	SupportedImageSuffixMP4  SupportedImageSuffix = ".mp4"
	SupportedImageSuffixWEBM  SupportedImageSuffix = ".webm"
	SupportedImageSuffixGLB  SupportedImageSuffix = ".glb"
)

type MimeType string

const (
	MimeTypeJPEG MimeType = "image/jpeg"
	MimeTypePNG  MimeType = "image/png"
	MimeTypeMP4  MimeType = "video/mp4"
	MimeTypeGLB  MimeType = "model/glb"
)

var (
	gSuffix2MimeTypes = map[SupportedImageSuffix]MimeType{
		SupportedImageSuffixJPEG: MimeTypeJPEG,
		SupportedImageSuffixJPG:  MimeTypeJPEG,
		SupportedImageSuffixPNG:  MimeTypePNG,
		SupportedImageSuffixMP4:  MimeTypeMP4,
		SupportedImageSuffixWEBM: MimeTypeMP4,
		SupportedImageSuffixGLB:  MimeTypeGLB,
	}
	gSupportedMimeTypes = map[MimeType]struct{}{
		MimeTypeJPEG: {},
		MimeTypePNG:  {},
		MimeTypeMP4:  {},
		MimeTypeGLB:  {},
	}
)

// UploadImage upload image to arweave chain
//
// @Summary upload image to arweave chain
// @Id images-upload
// @Tags images
// @Accept mpfd
// @Produce json
// @Param file formData file true "图片，文件类型目前仅支持png,jpg,jpeg"
// @Success 200 {object} e.Response{data=string}
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /images [post]
func UploadImage(c *gin.Context) {
	logger := logger.GetContextLogger(c)
	// Get the file from the request.
	file, info, err := c.Request.FormFile("file")
	if err != nil {
		logger.WithError(err).Debug(`get "file" failed`)
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, "file field empty")
		return
	}
	defer file.Close()

	// get supported suffix mime type
	mimeType, err := supportedSuffix2MimeType(path.Ext(info.Filename))
	if err != nil {
		logger.WithError(err).Debug("invalid file suffix")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, "file type not supported")
		return
	}

	// check file size
	if info.Size > setting.App().Image.MaxFileSize {
		logger.WithField("size", info.Size).Debug("file size too large")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, fmt.Sprintf("file size too large, should be less equal than %d bytes", setting.App().Image.MaxFileSize))
		return
	}

	// get image data
	data := streamToByte(file)

	// upload to chain
	uri, err := uploadImageToChain(data, mimeType)
	if err != nil {
		logger.WithError(err).Error("upload file content failed")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeChainConnection, "upload file content failed")
		return
	}

	// return uri
	e.WriteResponse(c, http.StatusOK, e.ERRCodeSuccess, "", uri)
}

func supportedSuffix2MimeType(suffix string) (string, error) {
	fileSuffix := SupportedImageSuffix(strings.ToLower(suffix))
	mimeType, ok := gSuffix2MimeTypes[fileSuffix]
	if !ok {
		return "", fmt.Errorf("not supported file suffix: %s", suffix)
	}

	return string(mimeType), nil
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func uploadImageToChain(data []byte, mimeType string) (uri string, err error) {
	// create a new transactor client
	ipAddress := setting.App().Node
	//ar, err := transactor.NewTransactor(ipAddress)
	//if err != nil {
	//	err = fmt.Errorf("NewTransactor failed: %w", err)
	//	return
	//}

	// create a new wallet instance
	//w := wallet.NewWallet()
	// extract the key from the wallet instance
	keyfile := setting.App().Keyfile
	//if err = w.LoadKeyFromFile(keyfile); err != nil {
	//	err = fmt.Errorf("LoadKeyFromFile failed: %w", err)
	//	return
	//}
	w, err := goar.NewWalletFromPath(keyfile, ipAddress)
	if err != nil {
		err = fmt.Errorf("NewWalletFromPath failed: %w", err)
		return
	}

	anchor, err := w.Client.GetTransactionAnchor()
	if err != nil {
		err = fmt.Errorf("GetTransactionAnchor failed: %w", err)
		return
	}
	// upload image
	// create a transaction
	reward, err := w.Client.GetTransactionPrice(data, nil)
	if err != nil {
		err = fmt.Errorf("GetTransactionPrice failed: %w", err)
		return
	}
	tags := make([]types.Tag,1,1)
	tags[0] = types.Tag{"Content-Type", mimeType}
	tx := &types.Transaction{
		Format:   2,
		Target:   "",
		Quantity: "0",
		Tags:     utils.TagsEncode(tags),
		Data:     utils.Base64Encode(data),
		DataSize: fmt.Sprintf("%d", len(data)),
		Reward:   fmt.Sprintf("%d", reward*(100)/100),
	}

	tx.LastTx = anchor
	tx.Owner = utils.Base64Encode(w.PubKey.N.Bytes())

	if err = utils.SignTransaction(tx,  w.PrvKey); err != nil {
		return
	}

	id := tx.ID
	//txBuilder, err := ar.CreateTransaction(context.TODO(), w, "0", data, "")
	//if err != nil {
	//	err = fmt.Errorf("CreateTransaction failed: %w", err)
	//	return
	//}
	//txBuilder.AddTag("Content-Type", mimeType)
	//txBuilder.AddTag("User-Agent", fmt.Sprintf(`ArweaveAPI/%s`, "v1.0.0"))

	// sign the transaction
	//txn, err := txBuilder.Sign(w)
	//if err != nil {
	//	err = fmt.Errorf("Sign failed: %w", err)
	//	return
	//}

	// send the transaction
	//if _, err = ar.SendTransaction(context.TODO(), txn); err != nil {
	//	err = fmt.Errorf("SendTransaction failed: %w", err)
	//	return
	//}

	u, err := url.Parse(setting.App().Node)
	if err != nil {
		err = fmt.Errorf("Parse failed: %w", err)
		return
	}

	uploader, err := goar.CreateUploader(w.Client, tx, nil)
	if err != nil {
		err = fmt.Errorf("CreateUploader failed: %w", err)
		return
	}

	err = uploader.Once()
	if err != nil {
		err = fmt.Errorf("uploader.Once failed: %w", err)
		return
	}
	u.Path = path.Join(u.Path, id)

	return u.String(), nil
}

type SaveImageUriReq struct {
	URI string `json:"uri"` // 图片URI
}

// SaveImageUri give image uri and resave it to arweave chain
//
// @Summary give image uri and resave it to arweave chain
// @Id images-save-uri
// @Tags images
// @Accept json
// @Produce json
// @Param request body SaveImageUriReq true "request"
// @Success 200 {object} e.Response{data=string}
// @Failure 400 {object} e.Response
// @Failure 500 {object} e.Response
// @Router /images/uri [post]
func SaveImageUri(c *gin.Context) {
	logger := logger.GetContextLogger(c)
	// parse parameter
	var param SaveImageUriReq
	if err := c.ShouldBindJSON(&param); err != nil {
		logger.WithError(err).Debug(`get "uri" failed`)
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, "uri field empty")
		return
	}

	// download url data
	data, err := downloadURIContent(param.URI)
	if err != nil {
		logger.WithError(err).Info("download image uri failed")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, "download image uri failed")
		return
	}

	// detect mime type
	mimeType := http.DetectContentType(data)
	if _, ok := gSupportedMimeTypes[MimeType(mimeType)]; !ok {
		logger.WithField("mimeType", mimeType).Info("not supported mime type")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeInvalidParam, "not supported mime type")
		return
	}

	// upload to chain
	uri, err := uploadImageToChain(data, mimeType)
	if err != nil {
		logger.WithError(err).Error("upload file content failed")
		e.WriteError(c, http.StatusBadRequest, e.ERRCodeChainConnection, "upload file content failed")
		return
	}

	// return uri
	e.WriteResponse(c, http.StatusOK, e.ERRCodeSuccess, "", uri)
}

func downloadURIContent(uri string) ([]byte, error) {
	rsp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		return nil, errors.New("not supported redirect uri")
	}
	data := new(bytes.Buffer)
	if _, err := io.Copy(data, rsp.Body); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
