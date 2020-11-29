package mapper

import (
	"cloud-store.com/service/dbproxy/orm"
	"fmt"
	"reflect"
)

var funcs = map[string]interface{}{
	"OnFileUploadFinished": orm.OnFileUploadFinished,
	"GetFileMeta":          orm.GetFileMeta,
	"GetFileMetaList":      orm.GetFileMetaList,
	"UpdateFileLocation":   orm.UpdateFileLocation,

	"UserSignUp":  orm.UserSignUp,
	"UserSignIn":  orm.UserSignIn,
	"UpdateToken": orm.UpdateToken,
	"GetUserInfo": orm.GetUserInfo,
	"UserExist":   orm.UserExist,

	"OnUserFileUploadFinished": orm.OnUserFileUploadFinished,
	"QueryUserFileMetas":       orm.QueryUserFileMetas,
	"DeleteUserFile":           orm.DeleteUserFile,
	"RenameFileName":           orm.RenameFileName,
	"QueryUserFileMeta":        orm.QueryUserFileMeta,
}

//FuncCall 通过函数名调用对应的函数，用到反射
func FuncCall(funcName string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := funcs[funcName]; !ok {
		err = fmt.Errorf("Error: function:%v not exist", funcName)
		return
	}

	//通过反射可以动态调用对象的导出方法
	f := reflect.ValueOf(funcs[funcName])
	if len(params) != f.Type().NumIn() {
		err = fmt.Errorf("Error: the number of parameters is not the expected number, expect:%v, but had:%v", f.Type().NumIn(), len(params))
		return
	}

	//构造一个Value的slice，用作Call()的方法传入
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	//执行方法f，并将方法结果赋值给result
	result = f.Call(in)
	return
}
