<template>
  <div class="content">
    <el-card shadow="always">
      <div slot="header">

        <!--搜索栏-->
        <div class="queryForm">
          <el-form :inline="true" :model="where" class="query-form-inline" size="small">
            <el-form-item label="剧集ID">
              <el-input v-model="where.book_id" placeholder="请输入 book_id" clearable class="queryElInput"
                        style="width: 220px;"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="search">查询</el-button>
              <el-button icon="el-icon-refresh" @click="reset">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>

      <!--工具栏-->
      <div class="toolbar">
        <el-button type="primary" size="small" icon="el-icon-plus" @click="add">添加</el-button>
        <el-button type="danger" size="small" icon="el-icon-delete" @click="del">删除</el-button>
        <el-button type="success" size="small" icon="el-icon-unlock" @click="batchUnlock(1)">批量解锁</el-button>
        <el-button type="warning" size="small" icon="el-icon-lock" @click="batchUnlock(0)">批量锁定</el-button>
      </div>

      <!--数据表格-->
      <el-table ref="table" class="tableData" :data="tableData" :highlight-selection-row="true"
                height="calc(100vh - 182px)"
                :border="true"
                v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center"></el-table-column>
        <el-table-column prop="id" label="ID" width="70px" align="center">
          <template v-slot="{row}">
            {{ row.id }}
          </template>
        </el-table-column>
        <el-table-column prop="book_name" label="剧名" min-width="160px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.book_name }}
          </template>
        </el-table-column>
        <el-table-column prop="book_id" label="剧集ID" width="120px" align="center">
          <template v-slot="{row}">
            {{ row.book_id }}
          </template>
        </el-table-column>
        <el-table-column prop="chapter_index" label="集序号" width="80px" align="center">
          <template v-slot="{row}">
            {{ row.chapter_index }}
          </template>
        </el-table-column>
        <el-table-column prop="chapter_name" label="分集名称" width="110px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.chapter_name }}
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(分钟)" width="110px" align="center">
          <template v-slot="{row}">
            {{ (row.duration / 60000).toFixed(1) }}
          </template>
        </el-table-column>
        <el-table-column prop="chapter_price" label="价格" width="70px" align="center">
          <template v-slot="{row}">
            {{ row.chapter_price }}
          </template>
        </el-table-column>
        <el-table-column prop="is_unlock" label="解锁" width="80px" align="center">
          <template v-slot="{row}">
            <el-tag :type="row.is_unlock === 1 ? 'success' : 'info'" size="small">
              {{ row.is_unlock === 1 ? '已解锁' : '锁定' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="150px" fixed="right">
          <template v-slot="{row}">
            <el-button size="mini" @click="play(row)">播放</el-button>
            <el-button size="mini" @click="edit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!--表格分页-->
      <div style="margin-top: 10px; text-align: center;" class="currentPage">
        <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="where.page"
            :page-sizes="pageSizes"
            :page-size="where.limit"
            background
            layout="total, sizes, prev, pager, next, jumper"
            :total="totalData">
        </el-pagination>
      </div>

    </el-card>

    <!--编辑数据-->
    <save :visible.sync="dialogVisible" :edit-data="editData" :book-id="where.book_id"
          @done="getChapterList"></save>

    <!--播放视频（ArtPlayer，原生解析 SRT，齿轮设置里切换字幕；仅允许点 X 关闭）-->
    <el-dialog :visible.sync="playVisible" width="50%" top="3vh"
               :close-on-click-modal="false" :close-on-press-escape="false" @close="closePlay">
      <!--标题单行省略，避免过长挡住右上角 X（悬停可看全文）-->
      <div slot="title" :title="playTitle"
           style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; padding-right: 30px; font-size: 18px; line-height: 24px; color: #303133;">
        {{ playTitle }}
      </div>
      <div ref="playContainer" style="width: 100%; height: 70vh; margin-bottom: 10px; margin-top: -10px"></div>
    </el-dialog>

  </div>
</template>

<script>
import {delChapter, getChapterList, getChapterPlay, setChapterUnlock} from "@/api/chapter";
import Artplayer from 'artplayer'
import save from "./save";

export default {
  name: 'Chapter',
  components: {save},
  data() {
    return {
      tableData: [],
      totalData: 0,
      pageSizes: [30, 50, 100, 200],
      multipleSelection: [],
      where: {
        book_id: '',
        page: 1,
        limit: 30,
      },
      loading: false,
      dialogVisible: false,
      editData: {},
      playVisible: false,
      playTitle: '',
      bookName: '',
    }
  },
  mounted() {
    const savedSize = Number(localStorage.getItem('adminPageSize'))
    this.where.limit = this.pageSizes.includes(savedSize) ? savedSize : this.pageSizes[0]
    // 支持从短剧列表页跳转带 book_id
    if (this.$route.query.book_id) {
      this.where.book_id = this.$route.query.book_id
    }
    this.bookName = this.$route.query.book_name || ''
    this.getChapterList()
  },
  methods: {
    //获取分集列表
    async getChapterList() {
      this.loading = true;
      setTimeout(async () => {
        try {
          let res = await getChapterList({...this.where})
          if (res.data.code === 0) {
            this.tableData = res.data.data || [];
            this.totalData = res.data.count || 0
            //剧名由后端按剧集ID查 drama_book 返回，直接输入ID查询也能显示
            this.bookName = res.data.book_name || ''
            //翻页/查询后滚动回顶部（页面 + 表格内部）
            this.$nextTick(() => {
              window.scrollTo(0, 0)
              if (this.$refs.table) this.$refs.table.bodyWrapper.scrollTop = 0
            })
          }
        } catch (e) {
          this.$message.error(e.message);
        } finally {
          this.loading = false;
        }
      }, 200)
    },
    //查询
    search() {
      this.where.page = 1
      this.getChapterList()
    },
    //重置搜索条件
    reset() {
      this.where.book_id = ''
      this.where.page = 1
      this.getChapterList()
    },
    //页数
    handleSizeChange(val) {
      localStorage.setItem('adminPageSize', val)
      this.where.limit = val
      this.getChapterList()
    },
    //页码
    handleCurrentChange(val) {
      this.where.page = val
      this.getChapterList()
    },
    //多选
    handleSelectionChange(val) {
      this.multipleSelection = []
      val.forEach(item => {
        this.multipleSelection.push(item.id);
      })
    },
    //添加（清空回显数据，避免残留上次编辑内容）
    add() {
      this.editData = {}
      this.dialogVisible = true
    },
    //获取修改数据
    edit(row) {
      this.editData = {...row}
      this.dialogVisible = true;
    },
    //播放该集视频：签名在点击播放时按需获取（列表不再预生成 play_url，避免加载后久置过期点播放失败）；
    //字幕无需签名，解析行内 subtitle JSON 数组取语言名，地址走本组 GetSubtitle 代理（ArtPlayer 原生支持 SRT，控制栏可切换）
    async play(row) {
      if (!row.video_url) {
        this.$message.warning('该集暂无视频')
        return
      }
      let playUrl = ''
      try {
        const res = await getChapterPlay(row.id)
        if (res.data.code !== 0) {
          this.$message.error(res.data.message)
          return
        }
        playUrl = res.data.data.play_url
      } catch (e) {
        this.$message.error(e.message)
        return
      }
      this.playTitle = '剧名：' + (this.bookName || '') + '，集名：' + (row.chapter_name || '')
      //字幕列表：idx 保持 subtitle 数组原下标（与 GetSubtitle 代理取路径的下标一致）；name 取文件名前缀（如 en.srt → EN）
      let subPaths = []
      try {
        const arr = JSON.parse(row.subtitle || '[]')
        subPaths = Array.isArray(arr) ? arr : []
      } catch (e) {
        subPaths = []
      }
      const subtitles = []
      subPaths.forEach((p, idx) => {
        if (!p) return
        const name = p.split('/').pop()
        const lang = name.split('.')[0] || 'sub'
        subtitles.push({url: apiUrl + '/GetSubtitle/' + row.id + '_' + idx + '.srt', name: lang.toUpperCase(), lang: lang})
      })
      //默认 en，无 en 则第一条
      const enIndex = subtitles.findIndex(s => s.lang === 'en')
      const defaultSub = subtitles.length > 0 ? subtitles[enIndex >= 0 ? enIndex : 0] : null
      this.playVisible = true
      this.$nextTick(() => {
        if (this.player) {
          this.player.destroy(false)
          this.player = null
        }
        this.player = new Artplayer({
          container: this.$refs.playContainer,
          url: playUrl,
          autoplay: true,
          volume: 0.7,
          lang: 'zh-cn',
          //短剧手机端竖屏观看：字幕从默认贴底 15px 抬高到 15%屏高，避开底部控制/手势区（控制栏显示时还会自动叠加控制栏高度）
          cssVar: {'--art-subtitle-bottom': '15%'},
          //齿轮菜单内置项，与官方演示一致（字幕偏移仅有字幕时有意义）
          playbackRate: true,
          aspectRatio: true,
          flip: true,
          subtitleOffset: subtitles.length > 0,
          setting: true,
          fullscreen: true,
          fullscreenWeb: true,
          //5.4.0 仅支持 subtitle 单个对象（ArtPlayer 自行解析 SRT），多语言切换用 setting.add 挂进齿轮
          //无字幕时不能传 subtitle 键（传 undefined 会触发 Type Error 致播放器构造失败、整部剧播不了），用条件展开按需注入
          ...(defaultSub ? {subtitle: {url: defaultSub.url, type: 'srt', name: defaultSub.name}} : {}),
        })
        //齿轮 -> 字幕菜单：切换语言 / 无字幕
        if (subtitles.length > 0) {
          this.player.setting.add({
            name: 'subtitle',
            html: '字幕',
            tooltip: defaultSub.name,
            selector: [
              {html: '无字幕', off: true},
              ...subtitles.map(s => ({html: s.name, url: s.url, default: s === defaultSub})),
            ],
            onSelect: (item) => {
              if (item.off) {
                this.player.subtitle.style({display: 'none'})
              } else {
                this.player.subtitle.style({display: ''})
                this.player.subtitle.switch(item.url, {type: 'srt', name: item.html})
              }
              return item.html
            },
          })
        }
      })
    },
    //关闭播放弹窗，销毁播放器
    closePlay() {
      if (this.player) {
        this.player.destroy(false)
        this.player = null
      }
    },
    //删除
    async del() {
      if (this.multipleSelection.length === 0) {
        this.$message.warning("请选择要删除的数据");
        return
      }
      this.$confirm('即将删除，是否继续?').then(async _ => {
        let ids = this.multipleSelection.join(',')
        try {
          let res = await delChapter(ids)
          this.$message.success(res.data.message)
          await this.getChapterList()
        } catch (e) {
          this.$message.error(e.message);
        }
      }).catch(_ => {
      });
    },
    //批量设置解锁/锁定（先勾选行，交互同删除）
    batchUnlock(val) {
      if (this.multipleSelection.length === 0) {
        this.$message.warning("请选择要操作的数据");
        return
      }
      const tip = val === 1 ? '解锁' : '锁定'
      this.$confirm('即将' + tip + ' ' + this.multipleSelection.length + ' 集，是否继续?').then(async _ => {
        const ids = this.multipleSelection.join(',')
        try {
          const res = await setChapterUnlock(ids, val)
          if (res.data.code === 0) {
            this.$message.success(res.data.message)
            await this.getChapterList()
          } else {
            this.$message.error(res.data.message)
          }
        } catch (e) {
          this.$message.error(e.message);
        }
      }).catch(_ => {
      });
    }
  }
}
</script>
