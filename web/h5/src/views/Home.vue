<template>
  <div class="home">
    <div class="container">
      <h2 class="section-title">热门短剧</h2>
      <div class="drama-grid">
        <div class="drama-card" v-for="item in dramaList" :key="item.book_id" @click="goDetail(item.book_id)">
          <div class="cover-wrap">
            <img v-lazy :data-src="item.cover_show" :alt="item.book_name" class="cover"/>
            <span class="episode-tag" v-if="item.chapter_count">{{ item.chapter_count }}集</span>
            <span class="free-tag" v-if="item.is_free === 1">免费</span>
          </div>
          <div class="info">
            <h3 class="title">{{ item.book_name }}</h3>
            <p class="meta">
              <span v-if="item.ratings" class="rating">★ {{ item.ratings }}</span>
              <span v-if="item.view_count_text" class="views">{{ item.view_count_text }}</span>
            </p>
          </div>
        </div>
      </div>
      <div v-if="dramaList.length === 0 && !loading" class="empty">暂无数据</div>
      <div v-if="total > dramaList.length" class="load-more">
        <button @click="loadMore" :disabled="loading">
          {{ loading ? '加载中...' : '加载更多' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import {getDramaList} from '@/api/drama'

export default {
  name: 'Home',
  data() {
    return {
      dramaList: [],
      page: 1,
      limit: 12,
      total: 0,
      loading: false,
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      getDramaList({page: this.page, limit: this.limit}).then(res => {
        this.loading = false
        if (res.data.code === 0) {
          this.dramaList = [...this.dramaList, ...(res.data.data || [])]
          this.total = res.data.count || 0
        }
      }).catch(() => { this.loading = false })
    },
    loadMore() {
      this.page++
      this.fetchData()
    },
    goDetail(bookId) {
      this.$router.push(`/drama/${bookId}`)
    }
  }
}
</script>

<style scoped>
.home {
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

.free-tag {
  position: absolute;
  top: 6px;
  left: 6px;
  background: rgba(103, 194, 58, 0.9);
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
  display: flex;
  gap: 8px;
}

.rating {
  color: #ff9900;
}

.empty {
  text-align: center;
  color: #999;
  padding: 40px 0;
}

.load-more {
  text-align: center;
  margin-top: 24px;
}

.load-more button {
  padding: 8px 24px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: #fff;
  cursor: pointer;
  font-size: 14px;
}

.load-more button:hover {
  border-color: #1890ff;
  color: #1890ff;
}

@media (max-width: 480px) {
  .drama-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }
}
</style>
