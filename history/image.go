package history

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	aw "github.com/deanishe/awgo"
	"gorm.io/gorm"
)

const (
	iconMappingTableName = "icon_mapping"
	faviconsTableName    = "favicons"
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

// ObtainIcon get the image according to the entryDao info.
func ObtainIcon(db *gorm.DB, entryDao *EntryDao) *aw.Icon {
	if !hasCache || db == nil {
		slog.Error("No cache capability!")
		return nil
	}

	var iconMapping IconMappingDao
	db.Table(iconMappingTableName).Where("page_url = ?", entryDao.URL).Find(&iconMapping)
	if iconMapping.ID > 0 {
		var faviconsDao FaviconsDao
		db.Table(faviconsTableName).Where("id = ?", iconMapping.IconId).First(&faviconsDao)
		icon, err := doObtainIcon(faviconsDao)
		if err != nil {
			slog.Error("Failed to obtain icon from database:", "error", err)
			return nil
		}
		return icon
	}
	return nil
}

func doObtainIcon(dao FaviconsDao) (*aw.Icon, error) {
	imagePath := filepath.Join(cacheDir, strconv.FormatInt(dao.ID, 10))
	if PathExists(imagePath) {
		return &aw.Icon{
			Value: imagePath,
		}, nil
	}

	response, err := http.Get(dao.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain icon from URL: %s, error: %w", dao.URL, err)
	}
	defer response.Body.Close()

	file, err := os.Create(imagePath)
	if err != nil {
		return nil, fmt.Errorf("create image file failed: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, fmt.Errorf("copy to file failed: %w", err)
	}
	return &aw.Icon{
		Value: imagePath,
	}, nil
}
