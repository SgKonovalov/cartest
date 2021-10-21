package definition

/*
ResultOk - Структура для возвращения JSON документа, в случае успешного выполнения функции.
ResultError - Структура для возвращения JSON документа, в случае ошибки бизнес-логики.
*/

type ResultOk struct {
	Result string `json:"result"`
}

type ResultError struct {
	Result string `json:"error"`
}
