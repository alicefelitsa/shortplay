import request from './request'

// 获取分集列表（按 book_id 过滤）
export function getChapterList(params) {
    return request.get('/GetChapterList', { params })
}

// 添加分集
export function addChapter(data) {
    return request.post('/AddChapter', data)
}

// 修改分集
export function saveChapter(data) {
    return request.post('/SaveChapter', data)
}

// 删除分集
export function delChapter(ids) {
    return request.get('/DelChapter', { params: { ids } })
}
