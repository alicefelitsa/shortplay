<template>
  <div class="content">
    <el-card shadow="always">
      <div slot="header" class="setting-header">
        <span>站点设置</span>
      </div>
      <div class="setting-form" v-loading="loading">
        <el-form ref="form" :model="form" label-width="90px">
          <el-form-item label="平台域名">
            <el-input v-model="form.domain" placeholder="如：http://127.0.0.1:8100"></el-input>
          </el-form-item>
          <el-form-item label="访问方式">
            <el-radio-group v-model="form.access_mode">
              <el-radio label="pc">PC</el-radio>
              <el-radio label="h5">H5</el-radio>
              <el-radio label="all">全部</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="展示类型">
            <el-radio-group v-model="form.display_type">
              <el-radio label="blank">空白</el-radio>
              <el-radio label="jump">跳转</el-radio>
            </el-radio-group>
            <div class="setting-tip setting-tip-block">访问方式不匹配时，跳转展示提示页，空白展示空白页</div>
          </el-form-item>
          <el-form-item label="视频密钥">
            <el-input v-model="form.video_secret_key" placeholder="与 CF Worker 的 SECRET_KEY 保持一致"></el-input>
            <div class="setting-tip setting-tip-block">用于生成视频播放地址的 HMAC 签名，改动后需同步更新 CF Worker</div>
          </el-form-item>
          <el-form-item label="过期时间">
            <el-input-number v-model="form.video_expire_minutes" :min="1" :max="10080" controls-position="right" style="width: 120px;"></el-input-number>
            <span class="setting-tip">分钟，视频播放签名的有效期</span>
          </el-form-item>
          <el-form-item class="setting-submit">
            <el-button type="primary" icon="el-icon-check" :loading="saving" @click="save">保存</el-button>
            <span class="setting-tip">域名用于前台短剧封面等资源地址的拼接</span>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script>
import {getConfigSetting, saveConfigSetting} from "@/api/setting";

export default {
  name: "Setting",
  data() {
    return {
      loading: false,
      saving: false,
      form: {
        domain: '',
        access_mode: 'all',
        display_type: 'jump',
        video_secret_key: '',
        video_expire_minutes: 1440
      }
    }
  },
  mounted() {
    this.getSetting()
  },
  methods: {
    //获取站点配置
    async getSetting() {
      this.loading = true
      try {
        let res = await getConfigSetting()
        if (res.data.code === 0 && res.data.data && res.data.data.length > 0) {
          this.form = {
            domain: res.data.data[0].domain || '',
            access_mode: res.data.data[0].access_mode || 'all',
            display_type: res.data.data[0].display_type || 'jump',
            video_secret_key: res.data.data[0].video_secret_key || '',
            video_expire_minutes: res.data.data[0].video_expire_minutes || 1440
          }
        }
      } catch (e) {
        this.$message.error(e.message);
      } finally {
        this.loading = false
      }
    },
    //保存站点配置
    async save() {
      this.saving = true
      try {
        let res = await saveConfigSetting({...this.form})
        if (res.data.code === 0) {
          this.$message.success(res.data.message)
        } else {
          this.$message.error(res.data.message)
        }
      } catch (e) {
        this.$message.error(e.message);
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style scoped>
/* 标题行：行高32px + 底部10px留白 + 字号14px */
.setting-header {
  margin-bottom: 10px;
}

.setting-header span {
  display: inline-block;
  line-height: 32px;
  font-size: 14px;
  font-weight: normal;
  color: #606266;
}

/* 表单定宽单列，避免输入框过宽拉伸 */
.setting-form {
  max-width: 560px;
  margin-top: 10px;
}

/* 访问方式/展示类型单选文字加深突出（选中态仍保留主题蓝） */
.setting-form >>> .el-radio__label {
  color: #303133;
}

.setting-form >>> .el-form-item {
  margin-bottom: 20px;
}

/* 保存按钮与说明文字同行，紧凑收尾 */
.setting-submit {
  margin-top: 8px;
  margin-bottom: 0;
}

.setting-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

/* 展示类型提示语单独一行置于单选组下方：去掉左缩进并换行 */
.setting-tip-block {
  display: block;
  margin-left: 0;
  margin-top: 2px;
  line-height: 1.6;
}
</style>
