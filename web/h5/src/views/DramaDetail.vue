<template>
  <div class="drama-detail">
    <div class="container" v-if="drama">
      <!-- 视频播放区域（置顶） -->
      <div class="player-wrap" v-if="currentEpisode">
        <video
            controls
            class="video-player"
            :src="videoOf(currentEpisode)"
            :poster="currentEpisode.cover || drama.cover_show"
        ></video>
        <h3 class="now-playing">正在播放：{{ currentEpisode.chapter_name || ('第' + (currentEpisode.chapter_index + 1) + '集') }}</h3>
      </div>

      <div class="detail-header">
        <img :src="drama.cover_show" :alt="drama.book_name" class="detail-cover"/>
        <div class="detail-info">
          <h1>{{ drama.book_name }}</h1>
          <p class="meta">
            <span v-if="drama.ratings" class="rating">★ {{ drama.ratings }}</span>
            <span v-if="drama.language">语言：{{ drama.language }}</span>
            <span>共{{ drama.chapter_count || episodes.length }}集</span>
            <span v-if="drama.view_count_text">{{ drama.view_count_text }}播放</span>
          </p>
          <p class="author" v-if="drama.author">作者：{{ drama.author }}</p>
          <p class="desc">{{ drama.introduction }}</p>
        </div>
      </div>

      <h2 class="section-title">剧集列表</h2>
      <div class="episode-list">
        <div
            class="episode-item"
            v-for="ep in episodes"
            :key="ep.chapter_id"
            :class="{active: currentEpisode && currentEpisode.chapter_id === ep.chapter_id}"
            @click="playEpisode(ep)"
        >
          <span class="ep-sort">{{ ep.chapter_index_str || (ep.chapter_index + 1) }}</span>
          <span class="ep-title">{{ ep.chapter_name }}</span>
          <el-tag v-if="ep.is_unlock === 0" size="mini" type="warning">锁</el-tag>
        </div>
      </div>
      <div v-if="episodes.length === 0" class="empty">暂无剧集</div>

      <!-- 推荐位 -->
      <template v-if="recommends.length > 0">
        <h2 class="section-title">相关推荐</h2>
        <div class="recommend-grid">
          <div class="rec-card" v-for="r in recommends" :key="r.book_id" @click="goDetail(r.book_id)">
            <img v-lazy :data-src="r.cover_show" :alt="r.book_name" class="rec-cover"/>
            <p class="rec-title">{{ r.book_name }}</p>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import {getDramaDetail} from '@/api/drama'

export default {
  name: 'DramaDetail',
  data() {
    return {
      drama: null,
      episodes: [],
      recommends: [],
      currentEpisode: null,
    }
  },
  created() {
    this.fetchData()
  },
  watch: {
    '$route.params.id'() {
      this.fetchData()
    }
  },
  methods: {
    videoOf(ep) {
      // 优先用后端签名的播放地址（CF Worker 校验），回退旧采集源
      return ep.play_url || ep.mp4_url || ''
    },
    fetchData() {
      const bookId = this.$route.params.id
      getDramaDetail(bookId).then(res => {
        if (res.data.code === 0) {
          this.drama = res.data.data && res.data.data[0]
          this.episodes = res.data.episodes || []
          this.recommends = res.data.recommends || []
          // 自动选中第一集
          this.currentEpisode = this.episodes.length > 0 ? this.episodes[0] : null
        }
      })
    },
    playEpisode(ep) {
      this.currentEpisode = ep
      window.scrollTo({top: 0, behavior: 'smooth'})
    },
    goDetail(bookId) {
      this.$router.push(`/drama/${bookId}`)
    }
  }
}
</script>

<style scoped>
.drama-detail {
  padding: 16px 0;
}

.container {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 16px;
}

.player-wrap {
  margin-bottom: 20px;
}

.video-player {
  width: 100%;
  border-radius: 8px;
  background: #000;
  max-height: 480px;
}

.now-playing {
  font-size: 15px;
  margin: 12px 0 0;
  color: #333;
}

.detail-header {
  display: flex;
  gap: 20px;
  margin-bottom: 24px;
}

.detail-cover {
  width: 160px;
  height: 220px;
  object-fit: cover;
  border-radius: 8px;
  flex-shrink: 0;
  background: #f0f0f0;
}

.detail-info h1 {
  font-size: 22px;
  margin: 0 0 10px;
  color: #333;
}

.meta {
  font-size: 13px;
  color: #999;
  margin: 0 0 8px;
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.rating {
  color: #ff9900;
}

.author {
  font-size: 13px;
  color: #666;
  margin: 0 0 8px;
}

.desc {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.episode-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 8px;
  margin-bottom: 24px;
}

.episode-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  border: 1px solid #eee;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.episode-item:hover {
  border-color: #1890ff;
  background: #f0f7ff;
}

.episode-item.active {
  border-color: #1890ff;
  background: #e6f7ff;
}

.ep-sort {
  font-size: 13px;
  font-weight: 500;
  color: #1890ff;
  white-space: nowrap;
  min-width: 28px;
}

.ep-title {
  font-size: 13px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.recommend-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
}

.rec-card {
  cursor: pointer;
}

.rec-cover {
  width: 100%;
  aspect-ratio: 3 / 4;
  object-fit: cover;
  border-radius: 6px;
  background: #f0f0f0;
}

.rec-title {
  font-size: 13px;
  color: #333;
  margin: 6px 0 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty {
  text-align: center;
  color: #999;
  padding: 30px 0;
}

@media (max-width: 480px) {
  .detail-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  .meta {
    justify-content: center;
  }
  .detail-cover {
    width: 120px;
    height: 168px;
  }
  .episode-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
