import request from './request'

// 获取分类列表
export function getTypeList(params) {
    return request.get('/GetTypeList', { params })
}

// 添加分类
export function addType(data) {
    return request.post('/AddType', data)
}

// 修改分类
export function saveType(data) {
    return request.post('/SaveType', data)
}

// 删除分类
export function delType(ids) {
    return request.get('/DelType', { params: { ids } })
}
