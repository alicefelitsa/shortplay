<template>
  <header class="site-header">
    <div class="header-inner">
      <router-link to="/" class="logo">短剧</router-link>
      <nav class="nav-links">
        <router-link to="/">首页</router-link>
        <router-link to="/category">分类</router-link>
        <router-link to="/search">搜索</router-link>
      </nav>
    </div>
  </header>
</template>

<script>
export default {
  name: 'SiteHeader'
}
</script>

<style scoped>
.site-header {
  background: #fff;
  border-bottom: 1px solid #eee;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  font-size: 20px;
  font-weight: 700;
  color: #1890ff;
  text-decoration: none;
}

.nav-links a {
  margin-left: 20px;
  color: #333;
  text-decoration: none;
  font-size: 14px;
}

.nav-links a.router-link-active {
  color: #1890ff;
  font-weight: 500;
}
</style>
