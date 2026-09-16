<template>
  <div class="search-page">
    <div class="container">
      <div class="search-bar">
        <input
            v-model="keyword"
            type="text"
            placeholder="搜索短剧..."
            class="search-input"
            @keyup.enter="handleSearch"
        />
        <button class="search-btn" @click="handleSearch">搜索</button>
      </div>

      <div class="drama-grid" v-if="dramaList.length > 0">
        <div class="drama-card" v-for="item in dramaList" :key="item.book_id" @click="goDetail(item.book_id)">
          <div class="cover-wrap">
            <img v-lazy :data-src="item.cover_show" :alt="item.book_name" class="cover"/>
            <span class="episode-tag" v-if="item.chapter_count">{{ item.chapter_count }}集</span>
          </div>
          <div class="info">
            <h3 class="title">{{ item.book_name }}</h3>
            <p class="category">
              <span v-if="item.ratings" class="rating">★ {{ item.ratings }}</span>
            </p>
          </div>
        </div>
      </div>

      <div class="empty" v-else-if="searched">
        <p>未找到相关短剧</p>
      </div>
    </div>
  </div>
</template>

<script>
import {searchDrama} from '@/api/drama'

export default {
  name: 'Search',
  data() {
    return {
      keyword: '',
      dramaList: [],
      searched: false,
    }
  },
  methods: {
    handleSearch() {
      if (!this.keyword.trim()) return
      searchDrama({keyword: this.keyword}).then(res => {
        this.searched = true
        if (res.data.code === 0) {
          this.dramaList = res.data.data || []
        }
      })
    },
    goDetail(bookId) {
      this.$router.push(`/drama/${bookId}`)
    }
  }
}
</script>

<style scoped>
.search-page {
  padding: 16px 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 16px;
}

.search-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.search-input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
}

.search-input:focus {
  border-color: #1890ff;
}

.search-btn {
  padding: 10px 24px;
  background: #1890ff;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
}

.search-btn:hover {
  background: #40a9ff;
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

.category {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.rating {
  color: #ff9900;
}

.empty {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

@media (max-width: 480px) {
  .drama-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }
}
</style>
