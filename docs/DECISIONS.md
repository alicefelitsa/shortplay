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

## 8. 分集管理与播放弹窗（chapter）

- **默认列表过滤**：`GetChapterList` 不传 `book_id` 时只显示 `flag=1` 短剧关联的集（`where exists (select 1 from drama_book b where b.book_id = drama_chapter.book_id and b.flag = 1)`），与 H5 展示同口径；按剧集ID查询时不过滤，便于管理下架/测试剧。
- **剧名注入**：`book_name` 由后端返回——工具栏/弹窗用的单值按查询ID（未传取首行）查 `drama_book`；列表列按当页去重 book_id 一次 `in (?)` 批量查注入每行。
- **播放弹窗**：ArtPlayer 挂载，`top="3vh"`、容器 `height: 70vh`；仅允许点 X 关闭（`:close-on-click-modal="false"` + `:close-on-press-escape="false"`）；标题 `slot="title"` 单行省略 + `padding-right: 30px` 防挡 X，格式「剧名：xxx，集名：xxx」。
- **字幕存储**：`drama_chapter.subtitle` 为 JSON 数组（如 `["/video/xxx/en.srt"]`），`[]`/空 = 无字幕；编辑弹窗（save.vue）用多行文本编辑，**每行一条路径**，保存时转 JSON 数组；字段标签「字幕路径」「视频路径」。
- **批量操作**：工具栏「批量解锁/批量锁定」，交互同删除（勾选→确认→调接口→刷新）；接口 `GET /api/boss/SetChapterUnlock?ids=1,2&is_unlock=0|1`，后端严格校验 is_unlock 仅 0/1 后拼接 SQL，顺带更新 `updated_at`。批量接口统一 GET + 逗号分隔 ids。
- **表单提交防御**：编辑保存只提交真实表字段，剔除列表注入的 `play_url / subtitle_urls / book_name / cover_show`，否则 `Updates` 报「操作失败」。

## 9. 字幕加载与播放（ArtPlayer + 同源代理）

- **播放器选型：ArtPlayer（npm `artplayer`，实测 5.4.0）**，原生解析 SRT。不要用原生 `<video>+<track>`——track 只支持 WebVTT，SRT 必失败（曾走过 SRT→VTT 转换方案，已废弃）。
- **5.4.0 没有 `subtitles` 数组选项**（传了会被静默忽略）：默认字幕用 `subtitle: {url, type:'srt', name}` 单对象；多语言切换用 `player.setting.add({name, html, tooltip, selector, onSelect})` 挂进齿轮菜单，onSelect 里 `player.subtitle.switch(url, {type:'srt'})`，「无字幕」用 `player.subtitle.style({display:'none'})`（恢复传 `display:''`）。
- 齿轮菜单按启用项动态生成：`setting: true` + `playbackRate / aspectRatio / flip / subtitleOffset`；全不启用时菜单为空、点击像无反应；中文界面 `lang: 'zh-cn'`（内置于核心包）。
- **字幕位置**：默认 `--art-subtitle-bottom: 15px` 贴底，短剧手机端观看用 `cssVar: {'--art-subtitle-bottom': '15%'}` 抬高；控制栏显示时 ArtPlayer 自动叠加控制栏高度。
- **字幕代理**：`GET /api/web/GetSubtitle/{chapterId}_{index}.srt`（挂 **web 组免鉴权**，播放器请求带不了登录 token）。服务端查 `drama_chapter.subtitle` 取对应路径后 `http.Get(domain + "/file" + path)` **原样返回**（`text/plain`，不做格式转换）。保留代理仅为规避浏览器跨域 fetch 限制；`.srt` 后缀供播放器识别类型。
- **字幕无需签名**：CF Worker 仅校验视频签名，字幕可匿名直连（`domain + /file + path`）；只有视频才走 `GenerateSignedVideoURL`。
- **字幕地址必须用全局 `apiWebUrl` 绝对地址**（`public/index.html` 定义，同 `apiUrl` 惯例）：dev 下相对路径会打到前端 dev server（8080 Express）404；跨域由后端全局 CORS 中间件放行。

## 10. 后台列表交互补充规范

- **翻页回顶**：列表成功回调 `$nextTick` 内双重回顶——`window.scrollTo(0,0)` + `$refs.table.bodyWrapper.scrollTop = 0`；三个列表页（drama/chapter/category）均已加 `ref="table"`。
- **状态开关**：drama 列表状态列用 `el-switch`（`active-value="PUBLISHED"` / `inactive-value="OFFLINE"`），切换即调 `saveDrama({id, status})`，失败后重新拉列表回滚。
- **操作列宽度**：两个 `mini` 按钮横排需 **150px**（140px 会换行堆叠）。
- 短剧列表查询支持分类下拉（`type_id`，精确匹配 `main_type_id`）+ 剧集ID（`book_id`，like 模糊）。

## 11. 踩坑经验

- **Element 编辑表单误加 `size="small"`**：会使输入框偏小、与 shop 不一致；编辑表单应无 size 用默认。
- **Chrome DevTools「问题」面板 label for 提示**：是无障碍/自动填充最佳实践提醒，非 JS 错误；Element UI 通病，内部后台可忽略。
- **参数名与时间单位不一致**：曾出现参数叫 `expireHours` 却乘 `time.Minute` 的 bug；命名与单位必须一致（现统一分钟）。
- **JOIN 查询漏目标表过滤**：recommends 曾漏 `b.flag = 1`；多表 JOIN 时每个表的过滤条件都要检查。
- **图片 URL 拼接需考虑中间层**：直连 B2 与经 CF Worker 的路径不同（Worker 需 `/file` 前缀），拼接前确认链路。
- **编辑回显依赖列表 select 字段**：`GetDramaList` 漏 select `introduction` 导致编辑弹窗简介为空；回显字段必须在列表查询中选出。
- **el-col 浮动高度独立**：要让同行列等高/居中，所在 `el-row` 必须 `type="flex"`。
- **GORM Scan 数值列类型断言静默失败**：`drama_chapter.book_id` 是数值类型，Scan 进 `map[string]interface{}` 后为 int64，`.(string)` 断言全部落空且不报错（曾致剧名列全空）；取 map 值统一 `fmt.Sprintf("%v", ...)` 归一化 + nil 保护。
- **Windows 显示缩放影响列宽观感**：约 180% 缩放下 140 CSS px 渲染成 ~250 物理像素，截图看似列宽未生效，实际配置有效，勿反复改。

## 12. 部署注意

- 前端 `apiUrl` / `apiWebUrl` 默认 `http://127.0.0.1:8100`，部署时在 `public/index.html` 切换生产域名（两者同步改，文件内已留同源写法注释行）。
- admin 本地端口 8082（`localhost:8082/admin`）。
- 改动 `video_secret_key` 后需同步 CF Worker 密钥，否则播放签名校验失败。
