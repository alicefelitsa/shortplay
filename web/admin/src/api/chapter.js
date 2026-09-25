import request from './request'

// 获取分集列表（按 book_id 过滤）
export function getChapterList(params) {
    return request.get('/GetChapterList', { params })
}

// 获取单集播放地址（点击播放时按需签发视频签名，避免列表预生成后过期）
export function getChapterPlay(id) {
    return request.get('/GetChapterPlay', { params: { id } })
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

// 批量设置解锁状态（isUnlock: 1 解锁 0 锁定）
export function setChapterUnlock(ids, isUnlock) {
    return request.get('/SetChapterUnlock', { params: { ids, is_unlock: isUnlock } })
}
