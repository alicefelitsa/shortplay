<template>
  <div>
    <el-card>
      <div slot="header">
        <el-form :inline="true" class="queryForm query-form-inline" size="small">
          <el-form-item label="短剧ID">
            <el-input v-model="query.drama_id" placeholder="请输入短剧ID" clearable class="queryElInput"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" @click="fetchData">查询</el-button>
            <el-button icon="el-icon-refresh" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
        <div class="toolbar">
          <el-button type="primary" size="small" icon="el-icon-plus" @click="handleAdd">添加剧集</el-button>
        </div>
      </div>

      <el-table :data="tableData" border class="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" align="center"/>
        <el-table-column prop="drama_id" label="短剧ID" width="80" align="center"/>
        <el-table-column prop="title" label="剧集标题" min-width="150"/>
        <el-table-column prop="sort" label="集数" width="80" align="center"/>
        <el-table-column prop="video_url" label="视频地址" min-width="200" show-overflow-tooltip/>
        <el-table-column prop="duration" label="时长(秒)" width="90" align="center"/>
        <el-table-column prop="is_free" label="免费" width="70" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.is_free === 1 ? 'success' : 'warning'" size="small">
              {{ scope.row.is_free === 1 ? '免费' : '付费' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" size="small" style="color:#F56C6C" @click="handleDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="currentPage" style="margin-top: 15px; text-align: right;">
        <el-pagination
            @current-change="handlePageChange"
            :current-page="page"
            :page-size="limit"
            :total="total"
            layout="total, prev, pager, next"
        />
      </div>
    </el-card>

    <!-- 添加/编辑弹窗 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="550px" class="responsive-dialog">
      <el-form :model="formData" :rules="formRules" ref="episodeForm" label-width="90px" size="small">
        <el-form-item label="短剧ID" prop="drama_id">
          <el-input-number v-model="formData.drama_id" :min="1"/>
        </el-form-item>
        <el-form-item label="剧集标题" prop="title">
          <el-input v-model="formData.title" placeholder="请输入剧集标题"/>
        </el-form-item>
        <el-form-item label="集数" prop="sort">
          <el-input-number v-model="formData.sort" :min="1"/>
        </el-form-item>
        <el-form-item label="视频地址" prop="video_url">
          <el-input v-model="formData.video_url" placeholder="请输入视频地址"/>
        </el-form-item>
        <el-form-item label="时长(秒)">
          <el-input-number v-model="formData.duration" :min="0"/>
        </el-form-item>
        <el-form-item label="是否免费">
          <el-radio-group v-model="formData.is_free">
            <el-radio :label="1">免费</el-radio>
            <el-radio :label="0">付费</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button size="small" @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" size="small" @click="handleSubmit">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import request from "@/api/request";

export default {
  name: 'Episode',
  data() {
    return {
      loading: false,
      tableData: [],
      query: {drama_id: ''},
      page: 1,
      limit: 20,
      total: 0,
      dialogVisible: false,
      dialogTitle: '添加剧集',
      formData: {
        id: null,
        drama_id: 1,
        title: '',
        sort: 1,
        video_url: '',
        duration: 0,
        is_free: 1,
      },
      formRules: {
        drama_id: [{required: true, message: '请输入短剧ID', trigger: 'blur'}],
        title: [{required: true, message: '请输入剧集标题', trigger: 'blur'}],
        video_url: [{required: true, message: '请输入视频地址', trigger: 'blur'}],
      },
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      request.get('/GetEpisodeList', {params: {page: this.page, limit: this.limit, drama_id: this.query.drama_id}}).then(res => {
        this.loading = false
        if (res.data.code === 0) {
          this.tableData = res.data.data || []
          this.total = res.data.count || 0
        }
      }).catch(() => { this.loading = false })
    },
    resetQuery() {
      this.query = {drama_id: ''}
      this.page = 1
      this.fetchData()
    },
    handlePageChange(val) {
      this.page = val
      this.fetchData()
    },
    handleAdd() {
      this.dialogTitle = '添加剧集'
      this.formData = {id: null, drama_id: 1, title: '', sort: 1, video_url: '', duration: 0, is_free: 1}
      this.dialogVisible = true
      this.$nextTick(() => { this.$refs.episodeForm && this.$refs.episodeForm.clearValidate() })
    },
    handleEdit(row) {
      this.dialogTitle = '编辑剧集'
      this.formData = {...row}
      this.dialogVisible = true
      this.$nextTick(() => { this.$refs.episodeForm && this.$refs.episodeForm.clearValidate() })
    },
    handleSubmit() {
      this.$refs.episodeForm.validate((valid) => {
        if (!valid) return
        const api = this.formData.id ? '/SaveEpisode' : '/AddEpisode'
        request.post(api, this.formData).then(res => {
          if (res.data.code === 0) {
            this.$message.success(res.data.message)
            this.dialogVisible = false
            this.fetchData()
          } else {
            this.$message.error(res.data.message)
          }
        })
      })
    },
    handleDel(row) {
      this.$confirm('确定要删除该剧集吗?', '提示', {type: 'warning'}).then(() => {
        request.get('/DelEpisode', {params: {ids: row.id}}).then(res => {
          this.$message.success(res.data.message)
          this.fetchData()
        })
      }).catch(() => {})
    },
  }
}
</script>

<style scoped>
</style>
