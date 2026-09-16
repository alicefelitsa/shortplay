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
            <el-form-item label="MP4地址">
              <el-input v-model="form.mp4_url" type="textarea" :rows="2" placeholder="视频 MP4 地址"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="本地路径">
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
        book_id: this.bookId || '',
        chapter_id: '',
        chapter_name: '',
        chapter_index: 0,
        mp4_url: '',
        video_url: '',
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
          const addOrSave = this.form.id ? saveChapter : addChapter;
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
