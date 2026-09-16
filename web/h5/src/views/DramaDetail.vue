<template>
  <div class="drama-detail">
    <div class="container" v-if="drama">
      <div class="detail-header">
        <img :src="drama.cover" :alt="drama.title" class="detail-cover"/>
        <div class="detail-info">
          <h1>{{ drama.title }}</h1>
          <p class="meta">分类：{{ drama.category }} | 共{{ drama.total_episodes }}集</p>
          <p class="desc">{{ drama.description }}</p>
        </div>
      </div>

      <h2 class="section-title">剧集列表</h2>
      <div class="episode-list">
        <div
            class="episode-item"
            v-for="ep in episodes"
            :key="ep.id"
            :class="{active: currentEpisode && currentEpisode.id === ep.id}"
            @click="playEpisode(ep)"
        >
          <span class="ep-sort">第{{ ep.sort }}集</span>
          <span class="ep-title">{{ ep.title }}</span>
          <el-tag v-if="ep.is_free === 0" size="mini" type="warning">付费</el-tag>
        </div>
      </div>

      <!-- 视频播放区域 -->
      <div class="player-wrap" v-if="currentEpisode">
        <h3>正在播放：{{ currentEpisode.title }}</h3>
        <video controls class="video-player" :src="currentEpisode.video_url"></video>
      </div>
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
      currentEpisode: null,
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      const id = this.$route.params.id
      getDramaDetail(id).then(res => {
        if (res.data.code === 0) {
          this.drama = res.data.data && res.data.data[0]
          this.episodes = res.data.episodes || []
        }
      })
    },
    playEpisode(ep) {
      this.currentEpisode = ep
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
}

.detail-info h1 {
  font-size: 22px;
  margin: 0 0 10px;
  color: #333;
}

.meta {
  font-size: 13px;
  color: #999;
  margin: 0 0 10px;
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
}

.ep-title {
  font-size: 13px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.player-wrap {
  margin-top: 16px;
}

.player-wrap h3 {
  font-size: 15px;
  margin-bottom: 12px;
  color: #333;
}

.video-player {
  width: 100%;
  border-radius: 8px;
  background: #000;
}

@media (max-width: 480px) {
  .detail-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
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
