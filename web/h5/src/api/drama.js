import request from '@/api/request';

// 获取短剧列表（可传 type_id 按分类过滤）
export function getDramaList(params) {
    return request.get('/GetDramaList', {params});
}

// 获取短剧详情 + 分集（book_id）
export function getDramaDetail(bookId) {
    return request.get('/GetDramaDetail', {params: {book_id: bookId}});
}

// 获取分类列表
export function getCategory() {
    return request.get('/GetCategory');
}

// 搜索短剧
export function searchDrama(params) {
    return request.get('/Search', {params});
}

// 获取站点公开配置
export function getSiteConfig() {
    return request.get('/GetSiteConfig');
}
