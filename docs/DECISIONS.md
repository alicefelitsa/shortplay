# Shortplay 项目关键决策与开发规范

> 本文档汇总架构、数据、接口、前端 UI 的关键决策与踩坑经验。
> 换设备 / 换人接手时以本文为准；与代码冲突时以代码为准并回改本文。

## 1. 技术栈与架构

- 后端：Go + Gin + GORM + Viper + Redis，分层 `config / controller / middleware / route / tools`，入口 `web.go`（无 main.go）。
- 前端 admin：Vue 2 + Element UI（管理后台）；前端 h5：Vue 2（用户端）。
- 数据库：MySQL `drama` 库。
- 响应约定：shortplay 统一 `res.data.code === 0` 为成功。

## 2. 图片 / 封面拼接（CF Worker 回源 B2）

- `drama_book.cover2` 存 B2 相对路径，形如 `/video/42000024646/cover.jpg`。
- 展示地址统一拼接：`cover_show = config.domain + "/file" + cover2`。
  - 例：`https://v.zly163.org/file/video/42000024646/cover.jpg`。
  - `/file` 是 CF Worker 回源 B2 的路由前缀；图片无签名可直连。
- 后端 `webController.withCoverShow(rows, domain)` 统一注入 `cover_show` 字段；前端（admin + h5）一律只读 `cover_show`，不自行拼 cover2。
- `cover_show` 仅展示用、**不写库**：`AddDrama / SaveDrama` 里 `delete(data, "cover_show")` 防御。

## 3. flag 展示过滤

- `drama_book.flag`：`1`=展示，`0`=隐藏（另有少量 `2`）。
- 所有对外/后台列表口均强制 `flag = 1`：
  - H5：列表、搜索、详情、推荐位（recommends 的 JOIN 需带 `b.flag = 1`）。
  - 后台：`GetDramaList`。
- `drama_chapter.m3u8_flag` 全为 0，业务上不使用，勿混淆。

## 4. 视频播放签名（HMAC-SHA256）

- 播放地址：`domain + "/file" + video_url`（`video_url` 为 B2 相对路径；`mp4_url` 是失效旧源，勿用）。
- 签名算法（`tools/video_sign.go` → `GenerateSignedVideoURL`）：
  - `message = videoPath + ":" + expires`
  - `sig = base64.URLEncoding( HMAC-SHA256(secretKey, message) )`
  - 最终：`{domain}{videoPath}?expires={expires}&signature={sig}`
- 过期时间按**分钟**计（`video_expire_minutes`），默认 1440（=24h）；后台可配，上限 10080。
- `GetDramaDetail` 的 episodes 注入 `play_url`（已签名）；H5 `videoOf(ep)` 优先用 `play_url`。
- **CF Worker 侧密钥必须与 config 表 `video_secret_key` 同步**，改动其一需改另一。

## 5. config 配置表

列：`id, domain, access_mode, display_type, video_secret_key, video_expire_minutes`。

- `domain`：播放/图片域名，如 `https://v.zly163.org`。
- `video_secret_key`：视频签名密钥（如 `felitsa`）。
- `video_expire_minutes`：签名有效期（分钟）。
- 后台「设置」页可编辑；读取用 `playDomain / videoSecret / videoExpireMinutes` helper（取第一行，空值回退默认）。

## 6. 后台 UI 规范（对齐 shop，保持统一）

- 列表页：`div.content > el-card shadow="always"`；header 内 `div.queryForm` 包 `el-form.query-form-inline` 且 `size="small"`；header 外 `.toolbar` 按钮 `size="small"`；`el-table` **默认尺寸**（不带 size）+ `border` + `highlight-selection-row` + `height="calc(100vh - 182px)"`；操作列按钮 `size="mini"`；居中完整分页。
- 编辑弹窗：独立 `save.vue` 子组件；`el-dialog` 百分比宽（50%）；`el-form` **不带 size**（默认输入框高度）；footer 按钮 `size="small"`；上传图片的删除按钮 `size="mini"`。
- 设置页：单列 `max-width: 560px`；数字输入框显式定宽（如 120px）避免过宽。
- `label-width`：4 字标签约 90px、2 字约 80px。
- 任何新增/修改功能必须沿用以上规范，不得自行引入冲突尺寸。

## 7. 编辑弹窗布局要点（drama/save.vue）

- 行 1–4：剧集ID/剧名、英文剧名/URL别名、分类/语言、评分/作者（两列）。
- 行 5：`<el-row type="flex">` 内 封面（左 span12）+ 付费/状态（右 span12）。
  - 右列 `style="display:flex; flex-direction:column; justify-content:center;"` 垂直居中。
  - **必须用 `el-row type="flex"`**：el-col 默认浮动、兄弟列高度独立，不拉齐则居中无效。
  - 不能给含 10 个半宽列的大 row 加 flex（nowrap 会挤成一行），故封面行独立拆出。
- 行 6：简介（全宽 textarea）。
- 封面提示文字：`只支持图片，不超过 2M；点击图片可重新上传`，`.album-tip { line-height: 1.8 }`。
- 添加按钮必须清空 `editData`（`add()` 方法），防止"点编辑后再点添加仍显示编辑数据"。

## 8. 踩坑经验

- **Element 编辑表单误加 `size="small"`**：会使输入框偏小、与 shop 不一致；编辑表单应无 size 用默认。
- **Chrome DevTools「问题」面板 label for 提示**：是无障碍/自动填充最佳实践提醒，非 JS 错误；Element UI 通病，内部后台可忽略。
- **参数名与时间单位不一致**：曾出现参数叫 `expireHours` 却乘 `time.Minute` 的 bug；命名与单位必须一致（现统一分钟）。
- **JOIN 查询漏目标表过滤**：recommends 曾漏 `b.flag = 1`；多表 JOIN 时每个表的过滤条件都要检查。
- **图片 URL 拼接需考虑中间层**：直连 B2 与经 CF Worker 的路径不同（Worker 需 `/file` 前缀），拼接前确认链路。
- **编辑回显依赖列表 select 字段**：`GetDramaList` 漏 select `introduction` 导致编辑弹窗简介为空；回显字段必须在列表查询中选出。
- **el-col 浮动高度独立**：要让同行列等高/居中，所在 `el-row` 必须 `type="flex"`。

## 9. 部署注意

- 前端 `apiUrl` 默认 `http://127.0.0.1:8100`，部署时在 `public/index.html` 切换生产域名。
- admin 本地端口 8082（`localhost:8082/admin`）。
- 改动 `video_secret_key` 后需同步 CF Worker 密钥，否则播放签名校验失败。
