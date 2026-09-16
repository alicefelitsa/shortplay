<template>
  <div class="not-found">
    <h1>404</h1>
    <p>您访问的页面不存在或已被移除</p>
    <router-link to="/" class="back-link">返回首页</router-link>
  </div>
</template>

<script>
export default {
  name: 'NotFound'
}
</script>

<style scoped>
.not-found {
  text-align: center;
  padding: 100px 16px;
}

.not-found h1 {
  font-size: 48px;
  color: #1890ff;
  margin-bottom: 10px;
}

.not-found p {
  color: #999;
  margin-bottom: 20px;
}

.back-link {
  color: #1890ff;
  text-decoration: none;
  font-size: 14px;
}

.back-link:hover {
  text-decoration: underline;
}
</style>
