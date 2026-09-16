<template>
  <div>
    <el-card>
      <div slot="header">
        <span>系统设置</span>
      </div>
      <el-form :model="formData" label-width="100px" size="small" style="max-width: 500px;">
        <el-form-item label="站点名称">
          <el-input v-model="formData.site_name" placeholder="请输入站点名称"/>
        </el-form-item>
        <el-form-item label="站点域名">
          <el-input v-model="formData.domain" placeholder="请输入站点域名"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave">保存设置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import request from "@/api/request";

export default {
  name: 'Setting',
  data() {
    return {
      formData: {
        site_name: '',
        domain: '',
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      request.get('/GetConfigSetting').then(res => {
        if (res.data.code === 0 && res.data.data && res.data.data.length > 0) {
          this.formData = {...res.data.data[0]}
        }
      })
    },
    handleSave() {
      request.post('/SaveConfigSetting', this.formData).then(res => {
        if (res.data.code === 0) {
          this.$message.success(res.data.message)
        } else {
          this.$message.error(res.data.message)
        }
      })
    },
  }
}
</script>

<style scoped>
</style>
