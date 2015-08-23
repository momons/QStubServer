// ファイル管理用
package FileManager

import (
	"os"
)

// ディレクトリを作成
func CreateDir(
	dirPath string,
) error {

	// ファイル情報取得
	fileInfo, err := os.Stat(dirPath)
	if err != nil || !fileInfo.IsDir() {
		// ディレクトリでなかったら作成
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}
