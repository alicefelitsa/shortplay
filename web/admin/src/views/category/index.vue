<template>
  <div>
    <el-card>
      <div slot="header">
        <el-form :inline="true" class="queryForm query-form-inline" size="small">
          <el-form-item label="分类名">
            <el-input v-model="query.name" placeholder="请输入分类名" clearable class="queryElInput"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" icon="el-icon-search" @click="fetchData">查询</el-button>
            <el-button icon="el-icon-refresh" @click="resetQuery">重置</el-button>
          </el-form-item>
        </el-form>
        <div class="toolbar">
          <el-button type="primary" size="small" icon="el-icon-plus" @click="handleAdd">添加分类</el-button>
        </div>
      </div>

      <el-table :data="tableData" border class="tableData" v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" align="center"/>
        <el-table-column prop="name" label="分类名称" min-width="150"/>
        <el-table-column prop="sort" label="排序" width="80" align="center"/>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'info'" size="small">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
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

    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="450px" class="responsive-dialog">
      <el-form :model="formData" :rules="formRules" ref="categoryForm" label-width="80px" size="small">
        <el-form-item label="分类名" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分类名称"/>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="formData.sort" :min="0"/>
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :label="1">启用</el-radio>
            <el-radio :label="0">禁用</el-radio>
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
  name: 'Category',
  data() {
    return {
      loading: false,
      tableData: [],
      query: {name: ''},
      page: 1,
      limit: 20,
      total: 0,
      dialogVisible: false,
      dialogTitle: '添加分类',
      formData: {id: null, name: '', sort: 0, status: 1},
      formRules: {
        name: [{required: true, message: '请输入分类名称', trigger: 'blur'}],
      },
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.loading = true
      request.get('/GetCategoryList', {params: {page: this.page, limit: this.limit, name: this.query.name}}).then(res => {
        this.loading = false
        if (res.data.code === 0) {
          this.tableData = res.data.data || []
          this.total = res.data.count || 0
        }
      }).catch(() => { this.loading = false })
    },
    resetQuery() {
      this.query = {name: ''}
      this.page = 1
      this.fetchData()
    },
    handlePageChange(val) {
      this.page = val
      this.fetchData()
    },
    handleAdd() {
      this.dialogTitle = '添加分类'
      this.formData = {id: null, name: '', sort: 0, status: 1}
      this.dialogVisible = true
      this.$nextTick(() => { this.$refs.categoryForm && this.$refs.categoryForm.clearValidate() })
    },
    handleEdit(row) {
      this.dialogTitle = '编辑分类'
      this.formData = {...row}
      this.dialogVisible = true
      this.$nextTick(() => { this.$refs.categoryForm && this.$refs.categoryForm.clearValidate() })
    },
    handleSubmit() {
      this.$refs.categoryForm.validate((valid) => {
        if (!valid) return
        const api = this.formData.id ? '/SaveCategory' : '/AddCategory'
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
      this.$confirm('确定要删除该分类吗?', '提示', {type: 'warning'}).then(() => {
        request.get('/DelCategory', {params: {ids: row.id}}).then(res => {
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
