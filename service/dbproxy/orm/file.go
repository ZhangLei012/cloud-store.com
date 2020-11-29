package orm

import (
	"database/sql"
	"log"

	"cloud-store.com/service/dbproxy/conn"
)

//OnFileUploadFinished 上传文件结束时更新文件元数据
func OnFileUploadFinished(fileHash, fileName string, fileSize int64, fileAddr string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare("insert ignore tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()
	result, err := statement.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		log.Printf("Error: executing statement:%v, err:%v", statement, err)
		return
	}

	if rowAffected, err := result.RowsAffected(); nil == err {
		if rowAffected <= 0 {
			log.Printf("Info: inserting affects no rows, file hash %v had been inserted", fileHash)
		}
		res.Success = true
		return
	}
	return
}

//GetFileMeta 从mysql获取文件元数据
func GetFileMeta(fileHash string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		return
	}
	tableFile := TableFile{}
	err = statement.QueryRow(fileHash).Scan(&tableFile.FileHash, &tableFile.FileName, &tableFile.FileSize, &tableFile.FileAddr)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error: select no rows")
			res.Success = true
			res.Data = nil
		} else {
			log.Printf("Error: querying row by file hash:%v, err:%v", fileHash, err)
			res.Success = false
			res.Msg = err.Error()
		}
		return
	}
	res.Success = true
	res.Data = tableFile
	return
}

//GetFileMetaList 从mysql批量获取元数据
func GetFileMetaList(limit int64) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"select file_sha1, file_name, file_size, file_addr from tbl_file where status = 1 limit ?")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	rows, err := statement.Query(limit)
	if err != nil {
		log.Printf("Error: querying statment:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}

	var tableFiles []TableFile
	for rows.Next() {
		tableFile := TableFile{}
		err = rows.Scan(&tableFile.FileHash, &tableFile.FileName, &tableFile.FileSize, &tableFile.FileAddr)
		if err != nil {
			log.Printf("Error: scaning row, err:%v", err)
			continue
		}
		tableFiles = append(tableFiles, tableFile)
	}
	res.Success = true
	res.Data = tableFiles
	return
}

//UpdateFileLocation 更新文件位置
func UpdateFileLocation(fileHash string, fileAddr string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"update tbl_file set `file_addr` = ? where `file_sha1` = ? limit 1")
	if err != nil {
		log.Printf("Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	result, err := statement.Exec(fileAddr, fileHash)
	if err != nil {
		log.Printf("Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}
	if rowAffected, err := result.RowsAffected(); nil == err {
		if rowAffected <= 0 {
			log.Printf("Error: updating file address affects no rows, file hash:%v ", fileHash)
			res.Data = "No row to update"
			return
		}
		res.Success = true
	} else {
		res.Msg = err.Error()
	}
	return
}
