# 


<p align="center">
  <img width="400" src="https://user-images.githubusercontent.com/8252317/83985151-9e8eaf00-a96a-11ea-9b3c-b654dc9bee2f.png">
  <h3 align="center">steam-box</h3>
  <p align="center"><img width="20" height="20" src="https://store.steampowered.com/favicon.ico"></img> Update a pinned gist to contain your Steam playtime leaderboard. </p>
  
   <p align="center">
    <a href="https://github.com/YouEclipse/steam-box/workflows/Update%20gist%20with%20Steam%20Playtime/badge.svg"><img src="https://github.com/YouEclipse/steam-box/workflows/Update%20gist%20with%20Steam%20Playtime/badge.svg" alt="Update gist with Steam Playtime"></a>
  </p>
</p>

---
[English](./README.md) | 简体中文



> 📌✨ 查看更多像这样的 Pinned Gist 项目,传送门:  https://github.com/matchai/awesome-pinned-gists



## 💻 安装

### 🎒 前置工作

1. 创建一个公开的 GitHub Gist (https://gist.github.com/)
1. 创建一个拥有 `gist` 权限的 token 并复制. (https://github.com/settings/tokens/new)
1. 创建你的 Steam  API key. (https://steamcommunity.com/dev/apikey)
1. 找到你的账号的 64 位 ID. (https://steamid.io)


### 🚀 开始安装

1. Fork 这个仓库
1. 编辑  `.github/workflows/schedule.yml` 中的[环境变量](https://github.com/YouEclipse/steam-box/actions/runs/126970182/workflow#L17-L19) :

   - **GIST_ID:** ID 是 gist url 的后缀 : `https://gist.github.com/YouEclipse/`**`9bc7025496e478f439b9cd43eba989a4`**.

1. 前往 fork 后的仓库的 **Settings > Secrets**
1. 添加以下环境变量:
   - **GH_TOKEN:** 前置工作中生成的 github token.
   - **STEAM_API_KEY:** 前置工作中创建的 steam API key. 
   - **STEAM_ID:** 你的 steam 64 位 ID. 
1. 如果你想展示某几个指定的游戏,可以把他们的 ID 设置在环境变量 **APP_ID**：
  - 如 `APP_ID=431960,730`
  - 你可以在对应游戏的 steam 商店的 url 获取到游戏 id: `https://store.steampowered.com/app/`**730**`/CounterStrike_Global_Offensive/`
  
## 🕵️ 工作原理
- 基于 [Steam API](https://partner.steamgames.com/doc/webapi)  获取游戏的游玩时间
- 基于 Github API 更新 Gist
- 使用 Github Actions 自动更新 Gist  

## 📄  开源协议
本项目使用 [Apache-2.0](./LICENSE) 协议
