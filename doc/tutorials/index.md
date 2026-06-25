# OverstepLab 使用教程

## 学习路径

欢迎使用 OverstepLab 越权测试靶场！本教程系列将从基础入门到高级技巧，帮助你系统性地掌握越权漏洞的发现和利用。

### 推荐学习顺序

```
00 跨平台一键启动 (one-click-start)
    ↓
01 快速入门 (quickstart)
    ↓
02 水平越权 (horizontal-idor)
    ↓
03 垂直越权 (vertical-escalation)
    ↓
04 上下文越权 (context-escalation)
    ↓
05 高级组合攻击 (advanced-combo)
    ↓
06 详细 WriteUp (writeups)
```

### 各教程说明

| 教程 | 难度 | 内容 |
|------|------|------|
| [00-跨平台一键启动](00-one-click-start.md) | ⭐ | 环境准备、一键启动脚本使用、常见问题 |
| [01-快速入门](01-quickstart.md) | ⭐ | 系统概览、首次登录、界面介绍 |
| [02-水平越权](02-horizontal-idor.md) | ⭐⭐ | 发现并利用 IDOR 漏洞 |
| [03-垂直越权](03-vertical-escalation.md) | ⭐⭐ | 低权限执行高权限操作 |
| [04-上下文越权](04-context-escalation.md) | ⭐⭐⭐ | 跨业务上下文的逻辑越权 |
| [05-高级组合](05-advanced-combo.md) | ⭐⭐⭐⭐ | 组合多个漏洞实现深度渗透 |
| [06-WriteUp](06-writeups.md) | 参考 | 每个漏洞的详细分析 |

### 推荐工具

- **浏览器开发者工具**（Network 面板）
- **curl** 或 **Postman**（API 测试）
- **Burp Suite Community**（抓包改包）

### 安全模式

系统默认运行在**漏洞模式**（所有漏洞可触发）。你可以通过右上角的开关切换到**安全模式**（所有漏洞已修复）。

### 重置环境

如果你搞乱了环境，平台管理员可以一键重置数据库。

祝学习愉快！🎯
