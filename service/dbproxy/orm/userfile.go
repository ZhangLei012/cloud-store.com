package orm

import (
	"database/sql"
	"log"
	"time"

	"cloud-store.com/service/dbproxy/conn"
)

//OnUserFileUploadFinished 更新用户文件表
func OnUserFileUploadFinished(userName, fileHash, fileName string, fileSize int64) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"insert ignore tbl_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`, `status`) values(?, ?, ?, ?, ?, 1)")
	if err != nil {
		log.Printf("OnUserFileUploadFinished Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	_, err = statement.Exec(userName, fileHash, fileName, fileSize, time.Now())
	if err != nil {
		log.Printf("OnUserFileUploadFinished Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}
	res.Success = true
	return
}

//QueryUserFileMeta 查询单个用户文件元信息
func QueryUserFileMeta(userName, fileHash string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"select file_sha1, file_name, file_size, create_at, update_at from tbl_file where userName = ? and file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		log.Printf("QueryUserFileMeta Error: preparing statement, err:%v", err)
		return
	}
	userFile := TableUserFile{}
	err = statement.QueryRow(userName, fileHash).Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("QueryUserFileMeta Error: select no rows")
			res.Success = true
			res.Msg = "No user file meta found"
			res.Data = nil
		} else {
			log.Printf("QueryUserFileMeta Error: querying row by file hash:%v, err:%v", fileHash, err)
			res.Success = false
			res.Msg = err.Error()
		}
		return
	}
	res.Success = true
	res.Data = userFile
	return
}

//QueryUserFileMetas 批量查询用户数据
func QueryUserFileMetas(userName string, limit int64) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name=? and status = 1 limit ?")
	if err != nil {
		log.Printf("QueryUserFileMetas Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	rows, err := statement.Query(userName, limit)
	if err != nil {
		log.Printf("QueryUserFileMetas Error: querying statment:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}

	var userFiles []TableUserFile
	for rows.Next() {
		userFile := TableUserFile{}
		err = rows.Scan(&userFile.FileHash, &userFile.FileName, &userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdated)
		if err != nil {
			log.Printf("QueryUserFileMetas Error: scaning row, err:%v", err)
			continue
		}
		userFiles = append(userFiles, userFile)
	}
	res.Success = true
	res.Data = userFiles
	return
}

//DeleteUserFile 标记删除用户文件表记录
func DeleteUserFile(userName, fileHash string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"update tbl_user_file set `status`=2 where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Printf("DeleteUserFile Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	_, err = statement.Exec(userName, fileHash)
	if err != nil {
		log.Printf("DeleteUserFile Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}
	res.Success = true
	return
}

//RenameFileName 更新用户文件表上的文件名
func RenameFileName(userName, fileHash, newFileName string) (res ExecResult) {
	res.Success = false
	statement, err := conn.DBConn().Prepare(
		"update tbl_user_file set `file_name`=? where user_name=? and file_sha1=? limit 1")
	if err != nil {
		log.Printf("RenameFileName Error: preparing statement, err:%v", err)
		res.Msg = err.Error()
		return
	}
	defer statement.Close()

	_, err = statement.Exec(newFileName, userName, fileHash)
	if err != nil {
		log.Printf("RenameFileName Error: executing statement:%v, err:%v", statement, err)
		res.Msg = err.Error()
		return
	}
	res.Success = true
	res.Data = QueryUserFileMeta(userName, fileHash).Data
	return
}
