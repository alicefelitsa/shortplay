<template>
  <div class="category-page">
    <div class="container">
      <h2 class="section-title">分类浏览</h2>
      <div class="category-list">
        <div class="category-item" v-for="item in categoryList" :key="item.id">
          {{ item.name }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {getCategory} from '@/api/drama'

export default {
  name: 'Category',
  data() {
    return {
      categoryList: [],
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      getCategory().then(res => {
        if (res.data.code === 0) {
          this.categoryList = res.data.data || []
        }
      })
    }
  }
}
</script>

<style scoped>
.category-page {
  padding: 16px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
  color: #333;
}

.category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.category-item {
  padding: 8px 20px;
  border: 1px solid #eee;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.category-item:hover {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f7ff;
}
</style>
