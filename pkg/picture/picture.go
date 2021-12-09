package picture

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"

	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

/*
ToBase64 bitmap to image base64
*/
func ToBase64(image image.Image) (string, error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, image)
	if (err != nil) {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

/*

*/
func toGrayscale(img image.Image) {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			gray.Set(x, y, img.At(x, y));
		}
	}
}

/*
CompareGray compare gray
*/
func CompareGray(source image.Gray, target image.Gray, hashType string) (float64, error) {
	mat1, err1 := gocv.ImageGrayToMatGray(&source)
	defer mat1.Close()
	if err1 != nil {
		return 0, err1
	}
	
	mat2, err2 := gocv.ImageGrayToMatGray(&target)
	defer mat2.Close()
	if err2 != nil {
		return 0, err2
	}

	return CompareMat(mat1, mat2, hashType)
}

/*
Compare compare two images similarity
hashType: phash | average | blockmean0 | blockmean1 | colormoment | marrhildreth | radialvariance
*/
func Compare(source image.Image, target image.Image, hashType string) (float64, error) {

	mat1, err1 := gocv.ImageToMatRGBA(source)
	defer mat1.Close()
	if err1 != nil {
		return 0, err1
	}
	
	mat2, err2 := gocv.ImageToMatRGBA(target)
	defer mat2.Close()
	if err2 != nil {
		return 0, err2
	}

	return CompareMat(mat1, mat2, hashType)
	
}
/*
CompareMat compare mat
*/
func CompareMat(mat1 gocv.Mat, mat2 gocv.Mat, hashType string) (float64, error) {

	var hash contrib.ImgHashBase

	switch hashType {
		case "phash":
			hash = contrib.PHash{}
		case "average":
			hash = contrib.AverageHash{}
		case "blockmean0":
			hash = contrib.BlockMeanHash{}
		case "blockmean1":
			hash = contrib.BlockMeanHash{ Mode: contrib.BlockMeanHashMode1 }
		case "colormoment":
			hash = contrib.ColorMomentHash{}
		case "marrhildreth":
			hash = contrib.NewMarrHildrethHash()
		case "radialvariance":
			hash = contrib.NewRadialVarianceHash()
		default:
			hash = contrib.PHash{}
	}

	result1 := gocv.NewMat()
	defer result1.Close()
	result2 := gocv.NewMat()
	defer result2.Close()
	hash.Compute(mat1, &result1)
	hash.Compute(mat2, &result2)

	similar := hash.Compare(result1, result2)

	return similar, nil

}

/*
MatchValue MatchTemplate returned value
*/
type MatchValue struct {
	MinVal float32
	MaxVal float32
	MinLoc image.Point
	MaxLoc image.Point
}
/*
FindImage find image
*/
func FindImage(source, temp image.Image, matchMode gocv.TemplateMatchMode) (MatchValue, error) {
	mat1, err1 := gocv.ImageToMatRGB(source)
	defer mat1.Close()
	if err1 != nil {
		return MatchValue{}, err1
	}
	
	mat2, err2 := gocv.ImageToMatRGB(temp)
	defer mat2.Close()
	if err2 != nil {
		return MatchValue{}, err2
	}
	return FindImageMat(mat1, mat2, matchMode), nil
}
/*
FindImageMat find image
source 大图像，temp 小图像
*/
func FindImageMat(source, temp gocv.Mat, matchMode gocv.TemplateMatchMode) MatchValue {
	res := gocv.NewMat()
	defer res.Close()
	msk := gocv.NewMat()
	defer msk.Close()

	gocv.MatchTemplate(source, temp, &res, matchMode, msk)

	minVal, maxVal, minLoc, maxLoc := gocv.MinMaxLoc(res)

	return MatchValue{minVal, maxVal, minLoc, maxLoc}
}
