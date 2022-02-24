package model

type ResultVo struct {
	Code int
	Data interface{}
}

const ResultSuccess int = 200

func Success(data interface{}) ResultVo {
	var resultVo = ResultVo{}
	resultVo.Code = ResultSuccess
	resultVo.Data = data

	return resultVo
}
