<template>
  <div>
    <el-dialog :title="form.id?'修改分集':'添加分集'" :visible.sync="dialogVisible" width="50%"
               class="responsive-dialog" :close-on-click-modal="false" :before-close="handleClose" top="5vh">
      <el-form ref="form" :model="form" :rules="rules" label-width="90px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="剧集ID" prop="book_id">
              <el-input v-model="form.book_id" placeholder="所属短剧的 book_id"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分集ID" prop="chapter_id">
              <el-input v-model="form.chapter_id" placeholder="唯一分集ID" :disabled="!!form.id"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="分集名称">
              <el-input v-model="form.chapter_name" placeholder="请输入分集名称"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="集序号">
              <el-input-number v-model="form.chapter_index" :min="0" :max="9999" style="width: 100%;"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="时长(秒)">
              <el-input-number v-model="form.duration" :min="0" style="width: 100%;"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="价格">
              <el-input-number v-model="form.chapter_price" :min="0" style="width: 100%;"></el-input-number>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="字幕路径">
              <el-input v-model="subtitleText" type="textarea" :rows="4"
                        placeholder="字幕路径每行一条（如 /video/42000024160/701452748_episode_5/en.srt），留空为无字幕"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="视频路径">
              <el-input v-model="form.video_url" placeholder="本地视频路径（可选）"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="解锁">
              <el-radio-group v-model="form.is_unlock">
                <el-radio :label="1">已解锁</el-radio>
                <el-radio :label="0">锁定</el-radio>
              </el-radio-group>
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
import {addChapter, saveChapter} from "@/api/chapter";

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
    //当前列表筛选的 book_id，用于新增时预填
    bookId: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      form: this.emptyForm(),
      //字幕编辑文本（每行一条路径），保存时转 JSON 数组
      subtitleText: '',
      rules: {
        book_id: [{required: true, message: '请输入剧集ID', trigger: 'blur'}],
        chapter_id: [{required: true, message: '请输入分集ID', trigger: 'blur'}],
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
    }
  },
  watch: {
    visible(val) {
      if (val) {
        this.form = {...this.emptyForm(), ...(this.editData || {})}
        this.subtitleText = this.toSubtitleText(this.form.subtitle)
        this.$nextTick(() => {
          this.$refs.form && this.$refs.form.clearValidate()
        })
      }
    }
  },
  methods: {
    //subtitle JSON 数组 -> 换行文本供 textarea 回显
    toSubtitleText(s) {
      if (!s || s === '[]') return ''
      try {
        const arr = JSON.parse(s)
        return Array.isArray(arr) ? arr.join('\n') : String(s)
      } catch (e) {
        return String(s)
      }
    },
    emptyForm() {
      return {
        id: null,
        book_id: this.bookId || '',
        chapter_id: '',
        chapter_name: '',
        chapter_index: 0,
        video_url: '',
        subtitle: '',
        duration: 0,
        chapter_price: 0,
        is_unlock: 1,
      }
    },
    //添加或修改
    save() {
      this.$refs.form.validate(async (valid) => {
        if (!valid) return
        try {
          //只提交真实表字段（白名单拼装），列表注入的 book_name 等不会进 payload，避免更新不存在列失败
          const paths = this.subtitleText.split('\n').map(s => s.trim()).filter(s => s)
          const payload = {
            id: this.form.id,
            book_id: this.form.book_id,
            chapter_id: this.form.chapter_id,
            chapter_name: this.form.chapter_name,
            chapter_index: this.form.chapter_index,
            duration: this.form.duration,
            chapter_price: this.form.chapter_price,
            video_url: this.form.video_url,
            is_unlock: this.form.is_unlock,
            subtitle: JSON.stringify(paths),
          }
          if (!payload.id) delete payload.id
          const addOrSave = this.form.id ? saveChapter : addChapter;
          let res = await addOrSave(payload)
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
