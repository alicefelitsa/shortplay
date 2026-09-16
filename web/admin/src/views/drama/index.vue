<template>
  <div>
    <el-card>
      <div slot="header">
        <!-- 查询区域 -->
        <el-form :inline="true" class="queryForm query-form-inline" size="small">
          <el-form-item label="剧名">
            <el-input v-model="query.title" placeholder="请输入剧名" clearable class="queryElInput"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" @click="fetchData">查询</el-button>
            <el-button icon="el-icon-refresh" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
        <!-- 工具栏 -->
        <div class="toolbar">
          <el-button type="primary" size="small" icon="el-icon-plus" @click="handleAdd">添加短剧</el-button>
          <el-button type="danger" size="small" icon="el-icon-delete" @click="handleBatchDel">批量删除</el-button>
        </div>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border class="tableData" @selection-change="handleSelectionChange" v-loading="loading">
        <el-table-column type="selection" width="45" align="center"/>
        <el-table-column prop="id" label="ID" width="60" align="center"/>
        <el-table-column prop="title" label="剧名" min-width="150"/>
        <el-table-column prop="cover" label="封面" width="100" align="center">
          <template slot-scope="scope">
            <el-image v-if="scope.row.cover" :src="scope.row.cover" style="width: 60px; height: 80px" fit="cover"/>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="100" align="center"/>
        <el-table-column prop="total_episodes" label="总集数" width="80" align="center"/>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'" size="small">
              {{ scope.row.status === 1 ? '上架' : '下架' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ctime" label="创建时间" width="160" align="center"/>
        <el-table-column label="操作" width="150" align="center">
          <template slot-scope="scope">
            <el-button type="text" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="text" size="small" style="color:#F56C6C" @click="handleDel(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="600px" class="responsive-dialog">
      <el-form :model="formData" :rules="formRules" ref="dramaForm" label-width="80px" size="small">
        <el-form-item label="剧名" prop="title">
          <el-input v-model="formData.title" placeholder="请输入剧名"/>
        </el-form-item>
        <el-form-item label="封面" prop="cover">
          <el-upload
              :action="uploadUrl"
              :headers="uploadHeaders"
              :on-success="handleUploadSuccess"
              :show-file-list="false"
              accept="image/*"
          >
            <el-image v-if="formData.cover" :src="formData.cover" style="width: 100px; height: 140px" fit="cover"/>
            <el-button v-else size="small" type="primary">点击上传</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-input v-model="formData.category" placeholder="请输入分类"/>
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入简介"/>
        </el-form-item>
        <el-form-item label="总集数">
          <el-input-number v-model="formData.total_episodes" :min="1" :max="999"/>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">上架</el-radio>
            <el-radio :label="0">下架</el-radio>
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
import {getDramaList, addDrama, saveDrama, delDrama} from "@/api/drama";

export default {
  name: 'Drama',
  data() {
    return {
      loading: false,
      tableData: [],
      query: {title: ''},
      page: 1,
      limit: 15,
      total: 0,
      selection: [],
      dialogVisible: false,
      dialogTitle: '添加短剧',
      formData: {
        id: null,
        title: '',
        cover: '',
        category: '',
        description: '',
        total_episodes: 1,
        status: 1,
      },
      formRules: {
        title: [{required: true, message: '请输入剧名', trigger: 'blur'}],
      },
      uploadUrl: apiUrl + '/UploadImage',
    }
  },
  computed: {
    uploadHeaders() {
      return {Authorization: localStorage.getItem('token') || ''}
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      getDramaList({page: this.page, limit: this.limit, title: this.query.title}).then(res => {
        this.loading = false
        if (res.data.code === 0) {
          this.tableData = res.data.data || []
          this.total = res.data.count || 0
        }
      }).catch(() => {
        this.loading = false
      })
    },
    resetQuery() {
      this.query = {title: ''}
      this.page = 1
      this.fetchData()
    },
    handlePageChange(val) {
      this.page = val
      this.fetchData()
    },
    handleSelectionChange(val) {
      this.selection = val
    },
    handleAdd() {
      this.dialogTitle = '添加短剧'
      this.formData = {id: null, title: '', cover: '', category: '', description: '', total_episodes: 1, status: 1}
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dramaForm && this.$refs.dramaForm.clearValidate()
      })
    },
    handleEdit(row) {
      this.dialogTitle = '编辑短剧'
      this.formData = {...row}
      this.dialogVisible = true
      this.$nextTick(() => {
        this.$refs.dramaForm && this.$refs.dramaForm.clearValidate()
      })
    },
    handleUploadSuccess(res) {
      if (res.code === 0) {
        this.formData.cover = res.url
        this.$message.success('上传成功')
      } else {
        this.$message.error(res.message || '上传失败')
      }
    },
    handleSubmit() {
      this.$refs.dramaForm.validate((valid) => {
        if (!valid) return
        const api = this.formData.id ? saveDrama : addDrama
        api(this.formData).then(res => {
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
      this.$confirm('确定要删除该短剧吗?', '提示', {type: 'warning'}).then(() => {
        delDrama(row.id).then(res => {
          this.$message.success(res.data.message)
          this.fetchData()
        })
      }).catch(() => {})
    },
    handleBatchDel() {
      if (this.selection.length === 0) {
        this.$message.warning('请选择要删除的数据')
        return
      }
      const ids = this.selection.map(item => item.id).join(',')
      this.$confirm(`确定要删除选中的 ${this.selection.length} 条数据吗?`, '提示', {type: 'warning'}).then(() => {
        delDrama(ids).then(res => {
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
