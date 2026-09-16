/**
 * axios 实例
 */
import axios from 'axios';

const service = axios.create({
    baseURL: apiUrl,
    timeout: 30000,
});

/**
 * 添加响应拦截器
 */
service.interceptors.response.use(response => {
        return response;
    }, (error) => {
        if (!error.response) {
            if (error.code === 'ECONNABORTED' && error.message.includes('timeout')) {
                console.error('请求超时，请检查网络连接');
            } else if (error.message.includes('ERR_CONNECTION_REFUSED') || error.message.includes('Network Error')) {
                console.error('无法连接到服务器，请检查网络或稍后重试');
            } else {
                console.error('网络连接异常，请检查网络设置');
            }
        }
        return Promise.reject(error);
    }
);

export default service;
