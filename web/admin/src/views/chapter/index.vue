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
        <span v-if="where.book_id" style="margin-left: 12px; color: #909399; font-size: 13px;">
          当前剧集：{{ where.book_id }}
        </span>
      </div>

      <!--数据表格-->
      <el-table class="tableData" :data="tableData" :highlight-selection-row="true" height="calc(100vh - 182px)"
                :border="true"
                v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center"></el-table-column>
        <el-table-column prop="id" label="ID" width="70px" align="center">
          <template v-slot="{row}">
            {{ row.id }}
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
        <el-table-column prop="chapter_name" label="分集名称" min-width="160px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.chapter_name }}
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="时长(秒)" width="90px" align="center">
          <template v-slot="{row}">
            {{ row.duration }}
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
        <el-table-column label="视频" width="70px" align="center">
          <template v-slot="{row}">
            <span :style="{color: videoOf(row) ? '#67C23A' : '#F56C6C'}">
              {{ videoOf(row) ? '有' : '无' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="90px" fixed="right">
          <template v-slot="{row}">
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

  </div>
</template>

<script>
import {getChapterList, delChapter} from "@/api/chapter";
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
    }
  },
  mounted() {
    this.where.limit = this.pageSizes[0]
    // 支持从短剧列表页跳转带 book_id
    if (this.$route.query.book_id) {
      this.where.book_id = this.$route.query.book_id
    }
    this.getChapterList()
  },
  methods: {
    videoOf(row) {
      return row.mp4_url || row.video_url || ''
    },
    //获取分集列表
    async getChapterList() {
      this.loading = true;
      setTimeout(async () => {
        try {
          let res = await getChapterList({...this.where})
          if (res.data.code === 0) {
            this.tableData = res.data.data || [];
            this.totalData = res.data.count || 0
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
    }
  }
}
</script>
