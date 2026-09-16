import request from './request'

// 获取站点配置
export function getConfigSetting() {
    return request.get('/GetConfigSetting')
}

// 保存站点配置
export function saveConfigSetting(data) {
    return request.post('/SaveConfigSetting', data)
}
