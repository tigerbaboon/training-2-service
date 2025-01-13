package image

type ImageController struct {
	imageSvc *ImageService
}

func newController(imageService *ImageService) *ImageController {
	return &ImageController{
		imageSvc: imageService,
	}
}
