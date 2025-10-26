package utils

import (
	"context"
	"log"
	"mime/multipart"
	"os"
	"privat-unmei/internal/customerrors"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func CreateCloudinaryUtil() *CloudinaryUtil {
	cld, err := cloudinary.NewFromParams("dk8rlicon", os.Getenv("CLOUDINARY_KEY"), os.Getenv("CLOUDINARY_SECRET"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &CloudinaryUtil{cld: cld}
}

type CloudinaryUtil struct {
	cld *cloudinary.Cloudinary
}

func (cu *CloudinaryUtil) UploadFile(ctx context.Context, file multipart.File, filename string, folder string) (*uploader.UploadResult, error) {
	overwrite := new(bool)
	*overwrite = true
	params := uploader.UploadParams{
		ResourceType: "raw",
		Folder:       folder,
		PublicID:     filename,
		Overwrite:    overwrite,
	}
	res, err := cu.cld.Upload.Upload(ctx, file, params)
	if err != nil {
		return nil, customerrors.NewError(
			"failed to upload file",
			err,
			customerrors.CommonErr,
		)
	}
	return res, nil
}
