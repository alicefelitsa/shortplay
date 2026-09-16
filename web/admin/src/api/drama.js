import request from './request'

// 获取短剧列表
export function getDramaList(params) {
    return request.get('/GetDramaList', { params })
}

// 获取短剧详情
export function getDramaDetail(id) {
    return request.get('/GetDramaDetail', { params: { id } })
}

// 添加短剧
export function addDrama(data) {
    return request.post('/AddDrama', data)
}

// 修改短剧
export function saveDrama(data) {
    return request.post('/SaveDrama', data)
}

// 删除短剧
export function delDrama(ids) {
    return request.get('/DelDrama', { params: { ids } })
}
