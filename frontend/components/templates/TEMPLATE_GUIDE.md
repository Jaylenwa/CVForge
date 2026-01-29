# 模板开发规范（Templates）

本目录下的模板组件负责把 `ResumeData` 渲染成可打印的简历页面。为保证不同模板的一致性（字号档位、密度控制、日期展示、富文本渲染与分区排序），新增模板时应遵循以下规范，并优先复用 `templates/shared`。

## 必用复用层

- 排序/过滤：`getOrderedVisibleSections()`、`getOrderedItems()`（避免每个模板各写一套 `filter/sort`）。
- 间距密度：`getSpacingTokens(styles)`（基于 `styles.spacingMultiplier` 统一 compact/normal/spacious）。
- 日期区间：`formatDateRange(item, t, options)`（避免 `~` / `-` / `YYYY-MM` / `YYYY.MM` 混用）。
- 自定义字段：`parseCustomPairs()` + `normalizeCustomPairs()`（统一 JSON 解析、空值过滤与 trim）。
- 头像尺寸：`getIdPhotoClassName()`（统一证件照物理尺寸，避免 px/随意大小）。
- 富文本描述：`<RichText />`（统一 sanitize、基础排版与注入）。
- Exam 分区：`<ExamSection />`（统一 meta+scores 的解释与表格 props）。

位置：`frontend/components/templates/shared/`

## Typography 规范（字号档位）

模板内建议使用“档位”而不是随意选择大小，优先通过 Tailwind 的字号 class 表达：

- 姓名：`text-3xl` 或 `text-4xl`（模板可选其一，但同一模板内保持一致）
- 求职意向/岗位：`text-base` 或 `text-xl`
- Section 标题：`text-base` 或 `text-lg`（搭配 `font-bold`）
- 条目标题（item.title）：`text-sm` 或 `text-lg`（根据模板风格选择）
- 辅助信息（日期/地点/副标题）：`text-xs` 或 `text-sm`（搭配更浅颜色）

## 个人信息规范（Header Info）

个人信息区（姓名、求职意向/岗位、联系方式/基础字段/自定义字段）在不同模板中必须使用同一套字号档位：

- 姓名：`text-3xl`
- 个人信息（联系方式/字段）：`text-sm`
- label（如“电话/邮箱/城市”）：保持更浅色，但字号仍为 `text-sm`

## 头像规范（证件照）

模板头像必须统一为证件照物理尺寸，避免使用 `px` 或随意 `w/h`：

- 尺寸：`35mm × 49mm`（比例 5:7）
- 渲染：`object-cover`（保证不同照片比例不破版）

用法示例：

```tsx
import { getIdPhotoClassName } from './shared/templateTokens';

<img className={getIdPhotoClassName('rounded-md object-cover')} />
```

## Density 规范（间距与密度控制）

密度受 `styles.spacingMultiplier` 控制，不应在模板里硬编码多套 `space-y-*` 逻辑。

推荐用法：

```tsx
import { getSpacingTokens } from './shared/templateTokens';

const { lineHeight, contentGapClass, headerSpaceClass, listTightClass, listMediumClass } = getSpacingTokens(styles);
```

- `contentGapClass`：用于 section 之间整体间距
- `listTightClass`：用于紧凑型列表（技能/证书/自评等）
- `listMediumClass`：用于经验/项目等需要更大段落空间的列表
- `lineHeight`：用于正文段落（尤其是富文本）统一行高

## 日期展示规范

模板展示日期区间时必须使用 `formatDateRange`：

```tsx
import { formatDateRange } from './shared/templateTokens';

const range = formatDateRange(item, t, { separatorVariant: 'tilde' });
```

- `separatorVariant: 'tilde'`：输出 `start ~ end`（Classic/Slate 风格）
- `separatorVariant: 'dash'`：输出 `start - end`（Mint 风格）
- `normalizeMonthSeparator: '.'`：将 `YYYY-MM` 显示为 `YYYY.MM`

## 富文本描述规范

模板渲染 `item.description` 必须使用 `RichText`，不直接写 `dangerouslySetInnerHTML`：

```tsx
import { RichText } from './shared/RichText';

<RichText html={item.description} className="text-gray-700 mt-1" fontSize={styles.fontSize} />
```

说明：`RichText` 内部会调用 `sanitizeHtml` 并提供 `.resume-rich-content` 基础排版；颜色与额外间距由 `className` 控制。

## Exam 分区规范

Exam 分区渲染必须使用 `ExamSection`，不要在模板里重复 `meta=items[0]`、`scores=items.slice(1)` 的拆解逻辑：

```tsx
import { ExamSection } from './shared/ExamSection';

<ExamSection section={section} color={accent} t={t} />
```

## 新模板 checklist

- 使用 `getOrderedVisibleSections` 渲染分区
- 使用 `getOrderedItems` 渲染条目
- 使用 `getSpacingTokens` 把密度绑定到 `styles.spacingMultiplier`
- 使用 `formatDateRange` 输出日期
- 使用 `RichText` 输出描述
- Exam 分区使用 `ExamSection`
- 自定义字段用 `parseCustomPairs + normalizeCustomPairs`
