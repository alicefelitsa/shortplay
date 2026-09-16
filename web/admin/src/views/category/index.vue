<template>
  <div class="content">
    <el-card shadow="always">
      <div slot="header">

        <!--搜索栏-->
        <div class="queryForm">
          <el-form :inline="true" :model="where" class="query-form-inline" size="small">
            <el-form-item label="分类名">
              <el-input v-model="where.type_name" placeholder="请输入" clearable class="queryElInput"></el-input>
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
      </div>

      <!--数据表格-->
      <el-table class="tableData" :data="tableData" :highlight-selection-row="true" height="calc(100vh - 182px)"
                :border="true"
                v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center"></el-table-column>
        <el-table-column prop="id" label="ID" width="80px">
          <template v-slot="{row}">
            {{ row.id }}
          </template>
        </el-table-column>
        <el-table-column prop="type_id" label="分类ID" width="100px" align="center">
          <template v-slot="{row}">
            {{ row.type_id }}
          </template>
        </el-table-column>
        <el-table-column prop="type_name" label="分类名称" min-width="160px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.type_name }}
          </template>
        </el-table-column>
        <el-table-column prop="replace_name" label="URL别名" min-width="160px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.replace_name }}
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="100px">
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

    <!--添加/编辑分类-->
    <el-dialog :title="form.id?'修改分类':'添加分类'" :visible.sync="dialogVisible" width="400px"
               :close-on-click-modal="false">
      <el-form ref="form" :model="form" label-width="80px">
        <el-form-item label="分类ID">
          <el-input-number v-model="form.type_id" :min="0" style="width: 100%;"></el-input-number>
        </el-form-item>
        <el-form-item label="名称">
          <el-input v-model="form.type_name" placeholder="请输入分类名称"></el-input>
        </el-form-item>
        <el-form-item label="URL别名">
          <el-input v-model="form.replace_name" placeholder="分类URL别名"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible=false" size="small">取 消</el-button>
        <el-button type="primary" @click="save" size="small">确 定</el-button>
      </span>
    </el-dialog>

  </div>
</template>

<script>
import {getTypeList, addType, saveType, delType} from "@/api/category";

export default {
  name: "Category",
  data() {
    return {
      tableData: [],
      totalData: 0,
      pageSizes: [30, 50, 100, 200],
      multipleSelection: [],
      where: {
        type_name: '',
        page: 1,
        limit: 30,
      },
      loading: false,
      dialogVisible: false,
      form: {
        id: '',
        type_id: 0,
        type_name: '',
        replace_name: ''
      }
    }
  },
  mounted() {
    this.where.limit = this.pageSizes[0]
    this.getTypeList()
  },
  methods: {
    //获取分类列表
    async getTypeList() {
      this.loading = true;
      setTimeout(async () => {
        try {
          let res = await getTypeList({...this.where})
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
      this.getTypeList()
    },
    //重置搜索条件
    reset() {
      this.where.type_name = ''
      this.where.page = 1
      this.getTypeList()
    },
    //页数
    handleSizeChange(val) {
      this.where.limit = val
      this.getTypeList()
    },
    //页码
    handleCurrentChange(val) {
      this.where.page = val
      this.getTypeList()
    },
    //多选
    handleSelectionChange(val) {
      this.multipleSelection = []
      val.forEach(item => {
        this.multipleSelection.push(item.id);
      })
    },
    //添加
    add() {
      this.form = {id: '', type_id: 0, type_name: '', replace_name: ''}
      this.dialogVisible = true
    },
    //编辑
    edit(row) {
      this.form = {...row}
      this.dialogVisible = true
    },
    //添加或修改
    async save() {
      if (!this.form.type_name || !this.form.type_name.trim()) {
        this.$message.warning("请输入分类名称");
        return
      }
      try {
        const addOrSave = this.form.id ? saveType : addType;
        let res = await addOrSave({...this.form})
        if (res.data.code === 0) {
          this.$message.success(res.data.message)
          this.dialogVisible = false
          await this.getTypeList()
        } else {
          this.$message.error(res.data.message)
        }
      } catch (e) {
        this.$message.error(e.message);
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
          let res = await delType(ids)
          this.$message.success(res.data.message)
          await this.getTypeList()
        } catch (e) {
          this.$message.error(e.message);
        }
      }).catch(_ => {
      });
    }
  }
}
</script>
