package server

type BookInfo struct {
	BookId int    // 书编号
	Name   string // 书名
}

// 书籍信息列表
type BookList struct {
	BookList []BookInfo
}

// 查询书籍信息参数
type BookInfoParams struct {
	BookId int
}

// 查询书籍列表参数
type BookListParams struct {
	Page  int
	Limit int
}

type Library int

// 获取书籍信息的函数
func (t *Library) GetBookInfo(param *BookInfoParams, reply *BookInfo) error {
	reply = findBookInfo(param.BookId)
	return nil
}

// 批量获取书籍信息的函数
func (t *Library) GetBookList(param *BookListParams, reply *BookList) error {
	reply = batchFindBookInfo(param)
	return nil
}

func findBookInfo(id int) *BookInfo {
	if id == 1 {
		return &BookInfo{1, "book1"}
	} else {
		return &BookInfo{id, "book2"}
	}
}

func batchFindBookInfo(list *BookListParams) *BookList {
	return nil
}
