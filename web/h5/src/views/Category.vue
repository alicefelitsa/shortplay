<template>
  <div class="category-page">
    <div class="container">
      <h2 class="section-title">分类浏览</h2>
      <div class="category-list">
        <div
            class="category-item"
            :class="{active: currentType === null}"
            @click="selectType(null)"
        >全部
        </div>
        <div
            class="category-item"
            :class="{active: currentType === item.type_id}"
            v-for="item in categoryList"
            :key="item.type_id"
            @click="selectType(item.type_id)"
        >
          {{ item.type_name }}
        </div>
      </div>

      <div class="drama-grid" v-if="dramaList.length > 0">
        <div class="drama-card" v-for="d in dramaList" :key="d.book_id" @click="goDetail(d.book_id)">
          <div class="cover-wrap">
            <img v-lazy :data-src="d.cover_show" :alt="d.book_name" class="cover"/>
            <span class="episode-tag" v-if="d.chapter_count">{{ d.chapter_count }}集</span>
          </div>
          <div class="info">
            <h3 class="title">{{ d.book_name }}</h3>
            <p class="meta">
              <span v-if="d.ratings" class="rating">★ {{ d.ratings }}</span>
            </p>
          </div>
        </div>
      </div>
      <div class="empty" v-else>该分类下暂无短剧</div>
    </div>
  </div>
</template>

<script>
import {getCategory, getDramaList} from '@/api/drama'

export default {
  name: 'Category',
  data() {
    return {
      categoryList: [],
      dramaList: [],
      currentType: null,
      page: 1,
      limit: 24,
    }
  },
  created() {
    this.fetchCategories()
    this.fetchDramas()
  },
  methods: {
    fetchCategories() {
      getCategory().then(res => {
        if (res.data.code === 0) {
          // 过滤掉 type_id=0 的 all 聚合分类（已由独立的“全部”按钮代替）
          this.categoryList = (res.data.data || []).filter(t => t.type_id !== 0)
        }
      })
    },
    fetchDramas() {
      const params = {page: this.page, limit: this.limit}
      if (this.currentType !== null) {
        params.type_id = this.currentType
      }
      getDramaList(params).then(res => {
        if (res.data.code === 0) {
          this.dramaList = res.data.data || []
        }
      })
    },
    selectType(typeId) {
      this.currentType = typeId
      this.page = 1
      this.fetchDramas()
    },
    goDetail(bookId) {
      this.$router.push(`/drama/${bookId}`)
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
  margin-bottom: 24px;
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

.category-item.active {
  border-color: #1890ff;
  color: #fff;
  background: #1890ff;
}

.drama-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.drama-card {
  cursor: pointer;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: transform 0.2s;
}

.drama-card:hover {
  transform: translateY(-3px);
}

.cover-wrap {
  position: relative;
  width: 100%;
  padding-top: 140%;
  overflow: hidden;
  background: #f0f0f0;
}

.cover {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.episode-tag {
  position: absolute;
  bottom: 6px;
  right: 6px;
  background: rgba(0, 0, 0, 0.6);
  color: #fff;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 3px;
}

.info {
  padding: 8px 10px;
}

.title {
  font-size: 14px;
  font-weight: 500;
  margin: 0 0 4px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.rating {
  color: #ff9900;
}

.empty {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

@media (max-width: 480px) {
  .drama-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }
}
</style>
