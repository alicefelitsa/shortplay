<template>
  <div>
    <el-dialog :title="form.id?'修改短剧':'添加短剧'" :visible.sync="dialogVisible" width="50%"
               class="responsive-dialog" :close-on-click-modal="false" :before-close="handleClose" top="5vh">
      <el-form ref="form" :model="form" :rules="rules" label-width="90px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="剧集ID" prop="book_id">
              <el-input v-model="form.book_id" placeholder="业务主键，唯一" :disabled="!!form.id"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="剧名" prop="book_name">
              <el-input v-model="form.book_name" placeholder="请输入剧名"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="英文剧名">
              <el-input v-model="form.book_name_en" placeholder="请输入英文剧名/别名"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="URL别名">
              <el-input v-model="form.slug" placeholder="slug，用于URL"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分类">
              <el-select v-model="form.main_type_id" placeholder="请选择分类" clearable style="width: 100%;">
                <el-option v-for="item in typeList" :key="item.type_id" :label="item.type_name"
                           :value="item.type_id"></el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="语言">
              <el-input v-model="form.language" placeholder="如 英语"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="评分">
              <el-input-number v-model="form.ratings" :min="0" :max="10" :step="0.1" :precision="1"
                               style="width: 100%;"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="作者">
              <el-input v-model="form.author" placeholder="请输入作者"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row type="flex">
          <el-col :span="12">
            <el-form-item label="封面" class="album-item">
              <div class="album-wrap">
                <el-upload
                    class="album-uploader"
                    :action="uploadAction"
                    :headers="uploadHeaders"
                    name="file"
                    accept="image/*"
                    :show-file-list="false"
                    :before-upload="beforeUpload"
                    :on-success="onUploadSuccess"
                    :on-error="onUploadError">
                  <img v-if="form.cover" :src="coverUrl" class="album-preview" alt="短剧封面">
                  <div v-else class="album-placeholder">
                    <i class="el-icon-plus"></i>
                    <span>上传封面</span>
                  </div>
                </el-upload>
                <div class="album-side">
                  <div class="album-tip">只支持图片，不超过 2M；点击图片可重新上传</div>
                  <el-button v-if="form.cover" size="mini" icon="el-icon-delete" @click="form.cover = ''">删除封面
                  </el-button>
                </div>
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12" style="display: flex; flex-direction: column; justify-content: center;">
            <el-form-item label="付费">
              <el-radio-group v-model="form.is_free">
                <el-radio :label="1">免费</el-radio>
                <el-radio :label="0">付费</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio label="PUBLISHED">已发布</el-radio>
                <el-radio label="OFFLINE">已下架</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="24">
            <el-form-item label="简介">
              <el-input v-model="form.introduction" type="textarea" :rows="4" placeholder="请输入剧情简介"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="handleClose" size="small">取 消</el-button>
        <el-button type="primary" @click="save" size="small">确 定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import {addDrama, saveDrama} from "@/api/drama";

// 添加或修改
export default {
  props: {
    //弹窗是否打开
    visible: {
      type: Boolean,
      default: false
    },
    //修改回显的数据
    editData: {
      type: Object,
      default: () => ({})
    },
    //分类列表
    typeList: {
      type: Array,
      default: () => ([])
    }
  },
  data() {
    return {
      form: this.emptyForm(),
      rules: {
        book_id: [{required: true, message: '请输入剧集ID', trigger: 'blur'}],
        book_name: [{required: true, message: '请输入剧名', trigger: 'blur'}],
      }
    }
  },
  computed: {
    dialogVisible: {
      get() {
        return this.visible
      },
      set(val) {
        this.$emit('update:visible', val)
      }
    },
    //上传接口地址
    uploadAction() {
      return apiUrl + '/UploadImage'
    },
    //上传携带的鉴权头
    uploadHeaders() {
      return {Authorization: localStorage.getItem('token')}
    },
    //封面预览地址：相对路径拼上后端服务地址
    coverUrl() {
      if (!this.form.cover) {
        return ''
      }
      if (/^https?:\/\//.test(this.form.cover)) {
        return this.form.cover
      }
      return apiUrl.replace('/api/boss', '') + this.form.cover
    }
  },
  watch: {
    editData: {
      handler(data) {
        this.form = {...this.emptyForm(), ...(data || {})}
      },
      immediate: true
    },
    visible(val) {
      if (val) {
        this.form = {...this.emptyForm(), ...(this.editData || {})}
        this.$nextTick(() => {
          this.$refs.form && this.$refs.form.clearValidate()
        })
      }
    }
  },
  methods: {
    emptyForm() {
      return {
        id: null,
        book_id: '',
        book_name: '',
        book_name_en: '',
        slug: '',
        cover: '',
        ratings: 0,
        main_type_id: null,
        language: '',
        author: '',
        introduction: '',
        is_free: 1,
        status: 'PUBLISHED',
      }
    },
    //上传前校验
    beforeUpload(file) {
      const isImage = /^image\/(png|jpe?g|gif)$/.test(file.type)
      if (!isImage) {
        this.$message.error('仅支持上传 png/jpg/gif/jpeg 图片')
      }
      const isLt2M = file.size / 1024 / 1024 < 2
      if (!isLt2M) {
        this.$message.error('图片大小不能超过 2MB')
      }
      return isImage && isLt2M
    },
    //上传成功
    onUploadSuccess(res) {
      if (res.code === 0) {
        this.form.cover = res.url
        this.$message.success('封面上传成功')
      } else {
        this.$message.error(res.message || '封面上传失败')
      }
    },
    //上传失败
    onUploadError() {
      this.$message.error('封面上传失败，请检查网络连接')
    },
    //添加或修改
    save() {
      this.$refs.form.validate(async (valid) => {
        if (!valid) return
        try {
          const addOrSave = this.form.id ? saveDrama : addDrama;
          let res = await addOrSave({...this.form})
          if (res.data.code === 0) {
            this.$message.success(res.data.message)
            this.handleClose()
            this.$emit('done')
          } else {
            this.$message.error(res.data.message)
          }
        } catch (e) {
          this.$message.error(e.message);
        }
      })
    },
    //关闭编辑对话框
    handleClose() {
      this.dialogVisible = false
      this.resetForm()
    },
    //重置表单
    resetForm() {
      this.form = this.emptyForm()
      if (this.$refs.form) {
        this.$refs.form.clearValidate()
      }
    }
  }
}
</script>

<style scoped>
/* 压缩图片行的占高，减少与下方简介的空白 */
.album-item {
  margin-bottom: 10px;
}

.album-wrap {
  display: flex;
  align-items: center;
}

.album-side {
  margin-left: 12px;
}

.album-uploader >>> .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color .3s;
}

.album-uploader >>> .el-upload:hover {
  border-color: #409EFF;
}

.album-preview {
  width: 120px;
  height: 160px;
  display: block;
  object-fit: cover;
}

.album-placeholder {
  width: 120px;
  height: 160px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #8c939d;
}

.album-placeholder i {
  font-size: 28px;
}

.album-placeholder span {
  margin-top: 8px;
  font-size: 13px;
}

.album-tip {
  margin: 0 0 10px;
  font-size: 12px;
  line-height: 1.8;
  color: #909399;
}
</style>
