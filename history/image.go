package history

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	aw "github.com/deanishe/awgo"
	"gorm.io/gorm"
)

const (
	iconsTableName   = "icon_mapping"
	bitmapsTableName = "favicon_bitmaps"
)

var (
	cacheDir string
	hasCache bool
)

func init() {
	cacheDir = os.Getenv("alfred_workflow_cache")
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), "chrome_history_cache")
	}
	slog.Info("Init cache:", "dir", cacheDir)

	if PathExists(cacheDir) {
		hasCache = true
	} else {
		hasCache = os.MkdirAll(cacheDir, os.ModePerm) == nil
	}
}

func ObtainIcon(db *gorm.DB, entryDao *EntryDao) *aw.Icon {
	if !hasCache || db == nil {
		slog.Error("No cache capability!")
		return nil
	}

	var iconMapping IconMappingDao
	db.Table(iconsTableName).Where("page_url = ?", entryDao.URL).Find(&iconMapping)
	if iconMapping.ID > 0 {
		var iconBitmap IconBitmapDao
		db.Table(bitmapsTableName).Where("icon_id = ?", iconMapping.IconId).First(&iconBitmap)
		icon, err := doObtainIcon(iconBitmap)
		if err != nil {
			slog.Error("Failed to obtain icon from database:", "error", err)
			return nil
		}
		return icon
	}
	return nil
}

func doObtainIcon(dao IconBitmapDao) (*aw.Icon, error) {
	imagePath := filepath.Join(cacheDir, strconv.FormatInt(dao.IconId, 10))
	if PathExists(imagePath) {
		return &aw.Icon{
			Value: imagePath,
		}, nil
	}

	img, _, err := image.Decode(bytes.NewReader(dao.ImageData))
	if err != nil {
		return nil, fmt.Errorf("decode image failed: %w", err)
	}
	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("create image file failed: %w", err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return nil, fmt.Errorf("encode image failed: %w", err)
	}
	return &aw.Icon{
		Value: imagePath,
	}, nil
}