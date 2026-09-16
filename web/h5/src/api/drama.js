import request from '@/api/request';

// 获取短剧列表
export function getDramaList(params) {
    return request.get('/GetDramaList', {params});
}

// 获取短剧详情
export function getDramaDetail(id) {
    return request.get('/GetDramaDetail', {params: {id}});
}

// 获取分类列表
export function getCategory() {
    return request.get('/GetCategory');
}

// 搜索短剧
export function searchDrama(params) {
    return request.get('/Search', {params});
}
