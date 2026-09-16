import request from '@/api/request';

// 获取短剧列表
export function getDramaList(params) {
    return request.get('/GetDramaList', {params});
}

// 添加短剧
export function addDrama(data) {
    return request.post('/AddDrama', data);
}

// 修改短剧
export function saveDrama(data) {
    return request.post('/SaveDrama', data);
}

// 删除短剧
export function delDrama(ids) {
    return request.get('/DelDrama', {params: {ids}});
}

// 上传图片
export function uploadImage(formData) {
    return request.post('/UploadImage', formData, {
        headers: {'Content-Type': 'multipart/form-data'}
    });
}
