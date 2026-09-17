<template>
  <div class="content">
    <el-card shadow="always">
      <div slot="header">

        <!--搜索栏-->
        <div class="queryForm">
          <el-form :inline="true" :model="where" class="query-form-inline" size="small">
            <el-form-item label="剧名">
              <el-input v-model="where.book_name" placeholder="请输入" clearable class="queryElInput"></el-input>
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="where.status" placeholder="请选择" clearable class="queryElInput">
                <el-option label="已发布" value="PUBLISHED"></el-option>
                <el-option label="已下架" value="OFFLINE"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="分类">
              <el-select v-model="where.type_id" placeholder="请选择" clearable class="queryElInput">
                <el-option v-for="item in typeList" :key="item.type_id" :label="item.type_name"
                           :value="item.type_id"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="剧集ID">
              <el-input v-model="where.book_id" placeholder="请输入" clearable class="queryElInput"></el-input>
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
      <el-table ref="table" class="tableData" :data="tableData" :highlight-selection-row="true" height="calc(100vh - 182px)"
                :border="true"
                v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" align="center"></el-table-column>
        <el-table-column prop="id" label="ID" width="60px">
          <template v-slot="{row}">
            {{ row.id }}
          </template>
        </el-table-column>
        <el-table-column prop="book_id" label="剧集ID" width="110px" align="center">
          <template v-slot="{row}">
            {{ row.book_id }}
          </template>
        </el-table-column>
        <el-table-column prop="cover" label="封面" align="center" width="90px">
          <template v-slot="{row}">
            <el-image
                style="height: 72px; width: 54px;"
                :src="coverUrl(row)"
                fit="cover"
                :preview-src-list="[coverUrl(row)]"
            >
              <div slot="error" class="image-slot">
                <i class="el-icon-picture-outline"></i>
              </div>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="book_name" label="剧名" min-width="180px" show-overflow-tooltip>
          <template v-slot="{row}">
            {{ row.book_name }}
          </template>
        </el-table-column>
        <el-table-column prop="language" label="语言" width="80px" align="center">
          <template v-slot="{row}">
            {{ row.language }}
          </template>
        </el-table-column>
        <el-table-column prop="ratings" label="评分" width="70px" align="center">
          <template v-slot="{row}">
            {{ row.ratings }}
          </template>
        </el-table-column>
        <el-table-column prop="chapter_count" label="集数" width="70px" align="center">
          <template v-slot="{row}">
            {{ row.chapter_count }}
          </template>
        </el-table-column>
        <el-table-column prop="view_count" label="播放量" width="90px" align="center">
          <template v-slot="{row}">
            {{ row.view_count }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80px" align="center">
          <template v-slot="{row}">
            <el-switch
                v-model="row.status"
                active-value="PUBLISHED"
                inactive-value="OFFLINE"
                @change="toggleStatus(row)">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="center" width="150px" fixed="right">
          <template v-slot="{row}">
            <el-button size="mini" @click="goChapter(row)">分集</el-button>
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
    <save :visible.sync="dialogVisible" :edit-data="editData" :type-list="typeList"
          @done="getDramaList"></save>

  </div>
</template>

<script>
import {getDramaList, delDrama, saveDrama} from "@/api/drama";
import {getTypeList} from "@/api/category";
import save from "./save";

export default {
  name: 'Drama',
  components: {save},
  data() {
    return {
      tableData: [],
      totalData: 0,
      pageSizes: [30, 50, 100, 200],
      multipleSelection: [],
      where: {
        book_name: '',
        status: '',
        type_id: '',
        book_id: '',
        page: 1,
        limit: 30,
      },
      loading: false,
      dialogVisible: false,
      editData: {},
      typeList: []
    }
  },
  mounted() {
    const savedSize = Number(localStorage.getItem('adminPageSize'))
    this.where.limit = this.pageSizes.includes(savedSize) ? savedSize : this.pageSizes[0]
    this.getDramaList()
    this.getTypeList()
  },
  methods: {
    //封面地址：后端已拼好播放域名 + cover2
    coverUrl(row) {
      return row.cover_show || ''
    },
    //获取短剧列表
    async getDramaList() {
      this.loading = true;
      setTimeout(async () => {
        try {
          let res = await getDramaList({...this.where})
          if (res.data.code === 0) {
            this.tableData = res.data.data || [];
            this.totalData = res.data.count || 0
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
    //获取分类，供搜索和编辑下拉使用
    async getTypeList() {
      try {
        let res = await getTypeList({page: 1, limit: 1000})
        if (res.data.code === 0) {
          this.typeList = res.data.data || []
        }
      } catch (e) {
        this.$message.error(e.message);
      }
    },
    //查询
    search() {
      this.where.page = 1
      this.getDramaList()
    },
    //重置搜索条件
    reset() {
      this.where.book_name = ''
      this.where.status = ''
      this.where.type_id = ''
      this.where.book_id = ''
      this.where.page = 1
      this.getDramaList()
    },
    //页数
    handleSizeChange(val) {
      localStorage.setItem('adminPageSize', val)
      this.where.limit = val
      this.getDramaList()
    },
    //页码
    handleCurrentChange(val) {
      this.where.page = val
      this.getDramaList()
    },
    //多选
    handleSelectionChange(val) {
      this.multipleSelection = []
      val.forEach(item => {
        this.multipleSelection.push(item.id);
      })
    },
    //跳转分集管理
    goChapter(row) {
      this.$router.push({path: '/chapter', query: {book_id: row.book_id, book_name: row.book_name}})
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
    //切换发布状态（开关），失败回滚
    async toggleStatus(row) {
      try {
        let res = await saveDrama({id: row.id, status: row.status})
        if (res.data.code === 0) {
          this.$message.success(res.data.message)
        } else {
          this.$message.error(res.data.message)
          await this.getDramaList()
        }
      } catch (e) {
        this.$message.error(e.message)
        await this.getDramaList()
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
          let res = await delDrama(ids)
          this.$message.success(res.data.message)
          await this.getDramaList()
        } catch (e) {
          this.$message.error(e.message);
        }
      }).catch(_ => {
      });
    }
  }
}
</script>

<style>
.image-slot {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
}

.image-slot i {
  font-size: 30px;
  color: #909399;
}
</style>
