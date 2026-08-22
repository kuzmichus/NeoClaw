<div align="center">
<img src="../../assets/logo.webp" alt="NeoClaw" width="512">

<h1>NeoClaw: Суперэффективный ИИ-ассистент на Go</h1>

<h3>Железо за $10 · 10 МБ ОЗУ · Загрузка за мс · Погнали, NeoClaw!</h3>
  <p>
    <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Arch-x86__64%2C%20ARM64%2C%20MIPS%2C%20RISC--V%2C%20LoongArch-blue" alt="Hardware">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://picoclaw.io"><img src="https://img.shields.io/badge/Website-picoclaw.io-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://docs.picoclaw.io/"><img src="https://img.shields.io/badge/Docs-Official-007acc?style=flat&logo=read-the-docs&logoColor=white" alt="Docs"></a>
    <a href="https://deepwiki.com/sipeed/picoclaw"><img src="https://img.shields.io/badge/Wiki-DeepWiki-FFA500?style=flat&logo=wikipedia&logoColor=white" alt="Wiki"></a>
    <br>
    <a href="https://x.com/SipeedIO"><img src="https://img.shields.io/badge/X_(Twitter)-SipeedIO-black?style=flat&logo=x&logoColor=white" alt="Twitter"></a>
    <a href="../../assets/wechat.png"><img src="https://img.shields.io/badge/WeChat-Group-41d56b?style=flat&logo=wechat&logoColor=white"></a>
    <a href="https://discord.gg/V4sAZ9XWpN"><img src="https://img.shields.io/badge/Discord-Community-4c60eb?style=flat&logo=discord&logoColor=white" alt="Discord"></a>
  </p>

[中文](README.zh.md) | [日本語](README.ja.md) | [한국어](README.ko.md) | [Português](README.pt-br.md) | [Tiếng Việt](README.vi.md) | [Français](README.fr.md) | [Italiano](README.it.md) | [Bahasa Indonesia](README.id.md) | [Malay](README.ms.md) | **Русский** | [English](../../README.md)

</div>

---

> **NeoClaw** — open-source проект, являющийся форком [PicoClaw](https://github.com/sipeed/picoclaw), полностью написанный на **Go**.

**NeoClaw** — сверхлёгкий персональный ИИ-ассистент, вдохновлённый [NanoBot](https://github.com/HKUDS/nanobot). Он был переписан с нуля на **Go** в процессе «самозагрузки» — сам ИИ-агент управлял миграцией архитектуры и оптимизацией кода.

**Работает на железе за $10 с <10 МБ ОЗУ** — это на 99% меньше памяти, чем у OpenClaw, и на 98% дешевле Mac mini!

<table align="center">
<tr align="center">
<td align="center" valign="top">
<p align="center">
<img src="../../assets/picoclaw_mem.gif" width="360" height="240">
</p>
</td>
<td align="center" valign="top">
<p align="center">
<img src="../../assets/licheervnano.png" width="400" height="240">
</p>
</td>
</tr>
</table>

> [!CAUTION]
> **Уведомление о безопасности**
>
> * **НИКАКОЙ КРИПТЫ:** NeoClaw **не** выпускал никаких официальных токенов или криптовалют. Всё, что появляется на `pump.fun` или других торговых площадках, — **мошенничество**.
> * **ОФИЦИАЛЬНЫЙ ДОМЕН:** **Единственный** официальный сайт — **[picoclaw.io](https://picoclaw.io)**, корпоративный сайт — **[sipeed.com](https://sipeed.com)**
> * **ОСТОРОЖНО:** Множество доменов `.ai/.org/.com/.net/...` зарегистрированы третьими лицами. Не доверяйте им.
> * **ПРИМЕЧАНИЕ:** NeoClaw находится на ранней стадии активной разработки. Возможны неисправленные проблемы безопасности. Не развертывайте в продакшн до версии v1.0.
> * **ПРИМЕЧАНИЕ:** Недавно в NeoClaw было слито много PR. Свежие сборки могут потреблять 10–20 МБ ОЗУ. Оптимизация ресурсов запланирована после стабилизации функциональности.

## 📢 Новости

2026-05-11 🛒 **LicheeRV-Claw на AliExpress!** Теперь LicheeRV-Claw можно приобрести на [AliExpress](https://www.aliexpress.com/item/1005006519668532.html) — стало проще попробовать NeoClaw на компактном железе RISC-V.

<p align="center">
  <a href="https://www.aliexpress.com/item/1005006519668532.html">
    <img src="../../assets/licheerv-claw.jpg" alt="LicheeRV-Claw on AliExpress" width="520">
  </a>
</p>

2026-05-28 🚀 **Вышел релиз v0.2.9!** Управление MCP-серверами в Web UI, настраиваемый веб-поиск на базе Sogou, анимация отклика инструментов в каналах, значения по умолчанию `pretty_print` и `disable_escape_html`, а также множество исправлений ошибок в провайдерах и каналах.

2026-05-14 🚀 **Вышел релиз v0.2.8!** CLI-команды MCP (`show`, `add`, `list`, `remove`, `test`, `edit`), пустой объект вместо null для параметров MCP-инструментов и исправления сборки.

2026-05-07 🚀 **Вышел релиз v0.2.7!** Настраиваемый веб-поиск на базе Sogou, анимация отклика инструментов в каналах, правки линтера.

2026-04-23 🚀 **Вышел релиз v0.2.6!** Хуки с действием respond и подробной документацией, поддержка изоляции, исправление баннера помощи.

2026-04-11 🚀 **Вышел релиз v0.2.5!** Zoneinfo из переменных окружения TZ/ZONEINFO, выравнивание CommonMark-рендеринга в Matrix, чтение `read_file` по строкам.

2026-03-31 📱 **Поддержка Android!** NeoClaw теперь работает на Android! Скачайте APK на [picoclaw.io](https://picoclaw.io/download)

2026-03-25 🚀 **Вышел релиз v0.2.4!** Переработка архитектуры агента (SubTurn, Hooks, Steering, EventBus), интеграция WeChat/WeCom, усиление безопасности (.security.yml, фильтрация чувствительных данных), новые провайдеры (AWS Bedrock, Azure, Xiaomi MiMo) и 35 исправлений ошибок. У NeoClaw уже **26K звёзд**!

2026-03-17 🚀 **Вышел релиз v0.2.3!** UI в системном трее (Windows и Linux), запрос статуса суб-агента (`spawn_status`), экспериментальная горячая перезагрузка шлюза, security-gating для Cron и 2 исправления безопасности. У NeoClaw уже **25K звёзд**!

2026-03-09 🎉 **v0.2.1 — крупнейшее обновление!** Поддержка протокола MCP, 4 новых канала (Matrix/IRC/WeCom/Discord Proxy), 3 новых провайдера (Kimi/Minimax/Avian), конвейер зрения (vision pipeline), JSONL-хранилище памяти, маршрутизация моделей.

2026-02-28 📦 Вышел **v0.2.0** с поддержкой Docker Compose и лаунчера Web UI.

<details>
<summary>Более ранние новости...</summary>

2026-02-26 🎉 NeoClaw набирает **20K звёзд** всего за 17 дней! Запущены авто-оркестрация каналов и интерфейсы возможностей.

2026-02-16 🎉 NeoClaw преодолевает отметку 12K звёзд за одну неделю! Официально запущены роли мейнтейнеров сообщества и [дорожная карта](../../ROADMAP.md).

2026-02-13 🎉 NeoClaw преодолевает 5000 звёзд за 4 дня! Ведётся работа над дорожной картой проекта и группами разработчиков.

2026-02-09 🎉 **Релиз NeoClaw!** Создан за 1 день, чтобы принести ИИ-агентов на железо за $10 с <10 МБ ОЗУ. Погнали, NeoClaw!

</details>

## ✨ Возможности

🪶 **Сверхлёгкость**: ядро занимает <10 МБ памяти — на 99% меньше OpenClaw.*

💰 **Минимальная стоимость**: достаточно эффективно для железа за $10 — на 98% дешевле Mac mini.

⚡️ **Молниеносная загрузка**: старт в 400 раз быстрее. Загружается менее чем за 1 с даже на одноядерном процессоре 0,6 ГГц.

🌍 **Настоящая переносимость**: один бинарник для архитектур RISC-V, ARM, MIPS и x86. Один бинарник — работает везде!

🤖 **Создан самим ИИ**: нативная реализация на чистом Go — 95% кода ядра сгенерировано агентом и доработано через ревью человеком.

🔌 **Поддержка MCP**: нативная интеграция с [Model Context Protocol](https://modelcontextprotocol.io/) — подключайте любые MCP-серверы для расширения возможностей агента.

👁️ **Конвейер зрения**: отправляйте изображения и файлы прямо агенту — автоматическое base64-кодирование для мультимодальных LLM.

🧠 **Умная маршрутизация**: маршрутизация моделей по правилам — простые запросы уходят лёгким моделям, экономя средства на API.

_*В свежих сборках может расходоваться 10–20 МБ из-за быстрого притока PR. Планируется оптимизация ресурсов. Сравнение скорости загрузки основано на бенчмарках одноядерных систем 0,8 ГГц (см. таблицу ниже)._

<div align="center">

|                                | OpenClaw      | NanoBot                  | **NeoClaw**                            |
| ------------------------------ | ------------- | ------------------------ | -------------------------------------- |
| **Язык**                       | TypeScript    | Python                   | **Go**                                 |
| **ОЗУ**                        | >1 ГБ         | >100 МБ                  | **< 10 МБ***                           |
| **Время загрузки**</br>(ядро 0,8 ГГц) | >500 с        | >30 с                    | **<1 с**                               |
| **Стоимость**                  | Mac Mini $599 | Большинство Linux-плат ~$50 | **Любая Linux-плата**</br>**от $10** |

<img src="../../assets/compare.jpg" alt="NeoClaw" width="512">

</div>

> **[Список совместимого оборудования](../guides/hardware-compatibility.md)** — все протестированные платы: от RISC-V за $5 до Raspberry Pi и Android-смартфонов. Вашей платы нет в списке? Пришлите PR!

<p align="center">
<img src="../../assets/hardware-banner.jpg" alt="NeoClaw Hardware Compatibility" width="100%">
</p>

## 🦾 Демонстрация

### 🛠️ Стандартные рабочие сценарии ассистента

<table align="center">
<tr align="center">
<th><p align="center">Режим фулстек-инженера</p></th>
<th><p align="center">Журналирование и планирование</p></th>
<th><p align="center">Веб-поиск и обучение</p></th>
</tr>
<tr>
<td align="center"><p align="center"><img src="../../assets/picoclaw_code.gif" width="240" height="180"></p></td>
<td align="center"><p align="center"><img src="../../assets/picoclaw_memory.gif" width="240" height="180"></p></td>
<td align="center"><p align="center"><img src="../../assets/picoclaw_search.gif" width="240" height="180"></p></td>
</tr>
<tr>
<td align="center">Разработка · Деплой · Масштабирование</td>
<td align="center">Расписание · Автоматизация · Память</td>
<td align="center">Поиск · Инсайты · Тренды</td>
</tr>
</table>

### 🐜 Инновационное развёртывание с минимальными ресурсами

NeoClaw можно развернуть практически на любом Linux-устройстве!

- За $9,9 [LicheeRV-Nano](https://www.aliexpress.com/item/1005006519668532.html) в версии E(Ethernet) или W(WiFi6) — минимальный домашний ассистент
- $30–50 [NanoKVM](https://www.aliexpress.com/item/1005007369816019.html) или $100 [NanoKVM-Pro](https://www.aliexpress.com/item/1005010048471263.html) — автоматизация серверных операций
- $50 [MaixCAM](https://www.aliexpress.com/item/1005008053333693.html) или $100 [MaixCAM2](https://www.kickstarter.com/projects/zepan/maixcam2-build-your-next-gen-4k-ai-camera) — умное видеонаблюдение

<https://private-user-images.githubusercontent.com/83055338/547056448-e7b031ff-d6f5-4468-bcca-5726b6fecb5c.mp4>

🌟 Впереди ещё больше примеров развёртывания!

## 📦 Установка

### Скачать с picoclaw.io (рекомендуется)

Зайдите на **[picoclaw.io](https://picoclaw.io)** — официальный сайт автоматически определит вашу платформу и предложит скачивание в один клик. Не нужно вручную выбирать архитектуру.

### Скачать готовый бинарник

Либо скачайте бинарник для своей платформы со страницы [GitHub Releases](https://github.com/sipeed/picoclaw/releases).

### Сборка из исходников (для разработки)

Предварительные требования:

- Go 1.25+
- Node.js 22+ и pnpm 10.33.0+ для сборки Web UI / лаунчера

```bash
git clone https://github.com/sipeed/picoclaw.git

cd picoclaw
make deps

# Установка зависимостей фронтенда
(cd web/frontend && pnpm install --frozen-lockfile)

# Сборка основного бинарника для текущей платформы
make build

# Сборка лаунчера Web UI (нужен для режима WebUI)
make build-launcher

# Сборка основных бинарников для всех платформ из Makefile
make build-all

# Сборка для Raspberry Pi Zero 2 W
# 32-бит: make build-linux-arm
# 64-бит: make build-linux-arm64
make build-pi-zero

# Сборка и установка
make install
```

**Raspberry Pi Zero 2 W:** используйте бинарник, соответствующий вашей ОС: 32-битная Raspberry Pi OS -> `make build-linux-arm`; 64-битная -> `make build-linux-arm64`. Либо выполните `make build-pi-zero` для сборки обеих версий.

## 🚀 Быстрый старт

### 🌐 Лаунчер WebUI (рекомендуется для десктопа)

Лаунчер WebUI предоставляет браузерный интерфейс для настройки и общения. Это самый простой способ начать — знание командной строки не требуется.

**Вариант 1: двойной клик (десктоп)**

После скачивания с [picoclaw.io](https://picoclaw.io) дважды щёлкните `picoclaw-launcher` (или `picoclaw-launcher.exe` в Windows). Браузер откроется автоматически по адресу `http://localhost:18800`.

**Вариант 2: командная строка**

```bash
picoclaw-launcher
# Откройте http://localhost:18800 в браузере
```

> [!TIP]
> **Удалённый доступ / Docker / VM:** добавьте флаг `-public`, чтобы слушать на всех интерфейсах:
> ```bash
> picoclaw-launcher -public
> ```

<p align="center">
<img src="../../assets/launcher-webui.jpg" alt="WebUI Launcher" width="600">
</p>

**Первые шаги:**

Откройте WebUI, затем: **1)** настройте провайдера (добавьте API-ключ LLM) -> **2)** настройте канал (например, Telegram) -> **3)** запустите шлюз -> **4)** общайтесь!

Подробная документация по WebUI — на [docs.picoclaw.io](https://docs.picoclaw.io).

<details>
<summary><b>Docker (альтернатива)</b></summary>

```bash
# 1. Склонируйте репозиторий
git clone https://github.com/sipeed/picoclaw.git
cd picoclaw

# 2. Первый запуск — автоматически создаёт docker/data/config.json и завершается
#    (срабатывает только если отсутствуют и config.json, и workspace/)
docker compose -f docker/docker-compose.yml --profile launcher up
# Контейнер выводит "First-run setup complete." и останавливается.

# 3. Задайте свои API-ключи
vim docker/data/config.json

# 4. Запуск
docker compose -f docker/docker-compose.yml --profile launcher up -d
# Откройте http://localhost:18800
```

> **Пользователям Docker / VM:** шлюз по умолчанию слушает `127.0.0.1`. Задайте `PICOCLAW_GATEWAY_HOST=0.0.0.0` или используйте флаг `-public`, чтобы сделать его доступным с хоста.

```bash
# Просмотр логов
docker compose -f docker/docker-compose.yml logs -f

# Остановка
docker compose -f docker/docker-compose.yml --profile launcher down

# Обновление
docker compose -f docker/docker-compose.yml pull
docker compose -f docker/docker-compose.yml --profile launcher up -d
```

</details>

<details>
<summary><b>macOS — предупреждение системы безопасности при первом запуске</b></summary>

macOS может заблокировать `picoclaw-launcher` при первом запуске, так как он скачан из интернета и не нотарифицирован через Mac App Store.

**Шаг 1:** дважды щёлкните `picoclaw-launcher`. Вы увидите предупреждение безопасности:

<p align="center">
<img src="../../assets/macos-gatekeeper-warning.jpg" alt="Предупреждение macOS Gatekeeper" width="400">
</p>

> *"picoclaw-launcher" Not Opened — Apple could not verify "picoclaw-launcher" is free of malware that may harm your Mac or compromise your privacy.*

**Шаг 2:** откройте **Системные настройки** → **Конфиденциальность и безопасность** → прокрутите до раздела **Безопасность** → нажмите «Открыть всё равно» → подтвердите кнопкой «Открыть всё равно» в диалоге.

<p align="center">
<img src="../../assets/macos-gatekeeper-allow.jpg" alt="macOS Privacy & Security — Open Anyway" width="600">
</p>

После этого одноразового шага `picoclaw-launcher` будет открываться как обычно при последующих запусках.

</details>

<a id="-run-on-old-android-phones"></a>
### 📱 Android

Дайте своему десятилетнему смартфону вторую жизнь! Превратите его в умного ИИ-ассистента с NeoClaw.

**Вариант 1: установка APK**

Предпросмотр:

<table>
  <tr>
    <td><img src="../../assets/fui_main_page.jpg" width="200"></td>
    <td><img src="../../assets/fui_web_page.jpg" width="200"></td>
    <td><img src="../../assets/fui_log_page.jpg" width="200"></td>
    <td><img src="../../assets/fui_setting_page.jpg" width="200"></td>
  </tr>
</table>

Скачайте APK с [picoclaw.io](https://picoclaw.io/download/) и установите напрямую. Termux не нужен!

**Вариант 2: Termux**

Полный чек-лист настройки через командную строку см. в [руководстве по Android Termux](../guides/android-termux.md).

<details>
<summary><b>Терминальный лаунчер (для сред с ограниченными ресурсами)</b></summary>

1. Установите [Termux](https://github.com/termux/termux-app) (скачайте с [GitHub Releases](https://github.com/termux/termux-app/releases) или найдите в F-Droid / Google Play)
2. Выполните следующие команды:

```bash
# Скачайте последний релиз
wget https://github.com/sipeed/picoclaw/releases/latest/download/picoclaw_Linux_arm64.tar.gz
tar xzf picoclaw_Linux_arm64.tar.gz
pkg install proot
termux-chroot ./picoclaw onboard   # chroot обеспечивает стандартную структуру файловой системы Linux
```

Затем следуйте разделу «Терминальный лаунчер» ниже, чтобы завершить настройку.

<img src="../../assets/termux.jpg" alt="NeoClaw on Termux" width="512">

Для минимальных окружений, где доступен только основной бинарник `picoclaw` (без UI лаунчера), всё можно настроить через командную строку и JSON-файл конфигурации.

**1. Инициализация**

```bash
picoclaw onboard
```

Команда создаёт `~/.picoclaw/config.json` и каталог workspace.

**2. Настройка** (`~/.picoclaw/config.json`)

```json
{
  "agents": {
    "defaults": {
      "model_name": "gpt-5.4"
    }
  },
  "model_list": [
    {
      "model_name": "gpt-5.4",
      "model": "openai/gpt-5.4"
      // api_key теперь загружается из .security.yml
    }
  ]
}
```

> Полный шаблон конфигурации со всеми доступными опциями см. в `config/config.example.json` в репозитории.
>
> Обратите внимание: формат config.example.json — версия 0, содержит чувствительные ключи и будет автоматически мигрирован на версию 1+; после этого config.json хранит только нечувствительные данные, а чувствительные ключи сохраняются в `.security.yml`. Если нужно изменить ключи вручную, подробности см. в `docs/security/security_configuration.md`.


**3. Общение**

```bash
# Одноразовый вопрос
picoclaw agent -m "What is 2+2?"

# Интерактивный режим
picoclaw agent

# Запуск шлюза для интеграции с чат-приложениями
picoclaw gateway
```

</details>

## 🔌 Провайдеры (LLM)

NeoClaw поддерживает 30+ LLM-провайдеров через конфигурацию `model_list`. Используйте формат `protocol/model`:

| Провайдер | Протокол | API-ключ | Примечания |
|----------|----------|---------|-------|
| [OpenAI](https://platform.openai.com/api-keys) | `openai/` | Требуется | GPT-5.4, GPT-4o, o3 и др. |
| [Anthropic](https://console.anthropic.com/settings/keys) | `anthropic/` | Требуется | Claude Opus 4.6, Sonnet 4.6 и др. |
| [Google Gemini](https://aistudio.google.com/apikey) | `gemini/` | Требуется | Gemini 3 Flash, 2.5 Pro и др. |
| [OpenRouter](https://openrouter.ai/keys) | `openrouter/` | Требуется | 200+ моделей, единый API |
| [Zhipu (GLM)](https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys) | `zhipu/` | Требуется | GLM-4.7, GLM-5 и др. |
| [DeepSeek](https://platform.deepseek.com/api_keys) | `deepseek/` | Требуется | DeepSeek-V3, DeepSeek-R1 |
| [Volcengine](https://console.volcengine.com) | `volcengine/` | Требуется | Doubao, модели Ark |
| [Qwen](https://dashscope.console.aliyun.com/apiKey) | `qwen/` | Требуется | Qwen3, Qwen-Max и др. |
| [Groq](https://console.groq.com/keys) | `groq/` | Требуется | Быстрый инференс (Llama, Mixtral) |
| [Moonshot (Kimi)](https://platform.moonshot.cn/console/api-keys) | `moonshot/` | Требуется | Модели Kimi |
| [Minimax](https://platform.minimaxi.com/user-center/basic-information/interface-key) | `minimax/` | Требуется | Модели MiniMax |
| [Mistral](https://console.mistral.ai/api-keys) | `mistral/` | Требуется | Mistral Large, Codestral |
| [NVIDIA NIM](https://build.nvidia.com/) | `nvidia/` | Требуется | Модели, размещённые NVIDIA |
| [Cerebras](https://cloud.cerebras.ai/) | `cerebras/` | Требуется | Быстрый инференс |
| [NEAR AI Cloud](https://near.ai/) | `nearai/` | Требуется | TEE-инференс, совместимость с OpenAI |
| [Novita AI](https://novita.ai/) | `novita/` | Требуется | Различные открытые модели |
| [Xiaomi MiMo](https://platform.xiaomimimo.com/) | `mimo/` | Требуется | Модели MiMo |
| [Ollama](https://ollama.com/) | `ollama/` | Не требуется | Локальные модели, self-hosted |
| [vLLM](https://docs.vllm.ai/) | `vllm/` | Не требуется | Локальное развёртывание, совместимость с OpenAI |
| [LiteLLM](https://docs.litellm.ai/) | `litellm/` | Зависит от случая | Прокси для 100+ провайдеров |
| [Azure OpenAI](https://portal.azure.com/) | `azure/` | API-ключ или Entra ID** | Корпоративное развёртывание Azure |
| [GitHub Copilot](https://github.com/features/copilot) | `github-copilot/` | OAuth | Вход по коду устройства |
| [Antigravity](https://console.cloud.google.com/) | `antigravity/` | OAuth | Google Cloud AI |
| [AWS Bedrock](https://console.aws.amazon.com/bedrock)* | `bedrock/` | Учётные данные AWS | Claude, Llama, Mistral на AWS |

> \* Для AWS Bedrock нужен build-тег: `go build -tags bedrock`. Задайте в `api_base` имя региона (например, `us-east-1`) для автоматического определения эндпоинта во всех партициях AWS (aws, aws-cn, aws-us-gov). При использовании полного URL эндпоинта также необходимо настроить `AWS_REGION` через переменную окружения или конфиг/профиль AWS.
>
> \*\* Azure OpenAI использует `api_key`, если он задан. Если `api_key` не указан, провайдер переключается на Microsoft Entra ID через `DefaultAzureCredential` (переменные окружения, workload identity, managed identity, Azure CLI и т.д.). Для пути Entra ID нужен build-тег: `go build -tags azidentity`.

<details>
<summary><b>Локальное развёртывание (Ollama, vLLM и др.)</b></summary>

**Ollama:**
```json
{
  "model_list": [
    {
      "model_name": "local-llama",
      "model": "ollama/llama3.1:8b",
      "api_base": "http://localhost:11434/v1"
    }
  ]
}
```

**vLLM:**
```json
{
  "model_list": [
    {
      "model_name": "local-vllm",
      "model": "vllm/your-model",
      "api_base": "http://localhost:8000/v1"
    }
  ]
}
```

Полные детали настройки провайдеров см. в [Providers & Models](../guides/providers.md).

</details>

## 💬 Каналы (чат-приложения)

Общайтесь с NeoClaw через 19+ мессенджеров:

| Канал | Настройка | Протокол | Док. |
|---------|----------|----------|------|
| **Telegram** | Просто (токен бота) | Long polling | [Гайд](../channels/telegram/README.md) |
| **Discord** | Просто (токен бота + intents) | WebSocket | [Гайд](../channels/discord/README.md) |
| **WhatsApp** | Просто (QR-скан или URL моста) | Нативный / Bridge | [Гайд](../guides/chat-apps.md#whatsapp) |
| **Weixin** | Просто (нативный QR-скан) | iLink API | [Гайд](../guides/chat-apps.md#weixin) |
| **QQ** | Просто (AppID + AppSecret) | WebSocket | [Гайд](../channels/qq/README.md) |
| **Slack** | Просто (токен бота + приложения) | Socket Mode | [Гайд](../channels/slack/README.md) |
| **Matrix** | Средне (homeserver + токен) | Sync API | [Гайд](../channels/matrix/README.md) |
| **Delta Chat** | Просто (скрипт аккаунта или email/пароль) | JSON-RPC (email/E2EE) | [Гайд](../channels/deltachat/README.md) |
| **DingTalk** | Средне (client credentials) | Stream | [Гайд](../channels/dingtalk/README.md) |
| **Feishu / Lark** | Средне (App ID + Secret) | WebSocket/SDK | [Гайд](../channels/feishu/README.md) |
| **LINE** | Средне (учётные данные + webhook) | Webhook | [Гайд](../channels/line/README.md) |
| **WeCom** | Просто (QR-вход или вручную) | WebSocket | [Гайд](../channels/wecom/README.md) |
| **VK** | Просто (токен группы) | Long Poll | [Гайд](../channels/vk/README.md) |
| **IRC** | Средне (сервер + ник) | Протокол IRC | [Гайд](../guides/chat-apps.md#irc) |
| **OneBot** | Средне (URL WebSocket) | OneBot v11 | [Гайд](../channels/onebot/README.md) |
| **MQTT** | Просто (брокер + agent_id) | MQTT pub/sub | [Гайд](../channels/mqtt/README.md) |
| **MaixCam** | Просто (включить) | TCP-сокет | [Гайд](../channels/maixcam/README.md) |
| **Pico** | Просто (включить) | Нативный протокол | Встроенный |
| **Pico Client** | Просто (URL WebSocket) | WebSocket | Встроенный |

> Все webhook-каналы используют общий HTTP-сервер шлюза (`gateway.host`:`gateway.port`, по умолчанию `127.0.0.1:18790`). Feishu работает в режиме WebSocket/SDK и не использует общий HTTP-сервер.

> Детализация логов управляется параметром `gateway.log_level` (по умолчанию `warn`). Поддерживаемые значения: `debug`, `info`, `warn`, `error`, `fatal`. Также задаётся через `PICOCLAW_LOG_LEVEL`. Подробнее см. [Конфигурация](../guides/configuration.md#gateway-log-level).

Подробные инструкции по настройке каналов см. в [Настройке чат-приложений](../guides/chat-apps.md).

## 🔧 Инструменты

### 🔍 Веб-поиск

NeoClaw может искать в интернете, чтобы предоставлять актуальную информацию. Настройка в `tools.web`:

| Поисковая система | API-ключ | Бесплатный лимит | Ссылка |
|--------------|---------|-----------|------|
| DuckDuckGo | Не требуется | Безлимит | Встроенный fallback |
| [Gemini Google Search](https://aistudio.google.com/apikey) | Требуется | Зависит от случая | Gemini с grounding через Google Search |
| [Baidu Search](https://cloud.baidu.com/doc/qianfan-api/s/Wmbq4z7e5) | Требуется | 1500/мес (дневная квота) | AI-поиск, оптимизирован для Китая |
| [Tavily](https://tavily.com) | Требуется | 1000 запросов/мес | Оптимизирован для ИИ-агентов |
| [Brave Search](https://brave.com/search/api) | Требуется | 2000 запросов/мес | Быстрый и приватный |
| [Kagi Search](https://help.kagi.com/kagi/api/search.html) | Требуется | Платный/лимит зависит от настройки API | Премиальные результаты поиска |
| [Perplexity](https://www.perplexity.ai) | Требуется | Платный | AI-поиск |
| [SearXNG](https://github.com/searxng/searxng) | Не требуется | Self-hosted | Бесплатная метапоисковая система |
| [GLM Search](https://open.bigmodel.cn/) | Требуется | Зависит от случая | Веб-поиск Zhipu |

### ⚙️ Другие инструменты

NeoClaw включает встроенные инструменты для работы с файлами, выполнения кода, планирования и многого другого. Подробнее см. [Конфигурацию инструментов](../reference/tools_configuration.md).

## 🎯 Навыки (Skills)

Навыки — это модульные возможности, расширяющие вашего агента. Они загружаются из файлов `SKILL.md` в вашем workspace.

**Установка навыков из ClawHub:**

```bash
picoclaw skills search "web scraping"
picoclaw skills install <skill-name>
```

**Настройка реестров навыков**:

Добавьте в свой `config.json`:
```json
{
  "tools": {
    "skills": {
      "registries": {
        "clawhub": {
          "auth_token": "your-clawhub-token"
        },
        "github": {
          "base_url": "https://github.com",
          "auth_token": "your-github-token",
          "proxy": ""
        }
      }
    }
  }
}
```

`tools.skills.github.*` устарел. Используйте вместо него `tools.skills.registries.github.*`.

Подробнее см. [Конфигурация инструментов - Skills](../reference/tools_configuration.md#skills-tool).

## 🔗 MCP (Model Context Protocol)

NeoClaw нативно поддерживает [MCP](https://modelcontextprotocol.io/) — подключайте любые MCP-серверы, чтобы расширить возможности агента внешними инструментами и источниками данных.

```json
{
  "tools": {
    "mcp": {
      "enabled": true,
      "servers": {
        "filesystem": {
          "enabled": true,
          "command": "npx",
          "args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"]
        }
      }
    }
  }
}
```

Типовыми настройками MCP можно управлять прямо из CLI, не редактируя JSON вручную:

```bash
picoclaw mcp add filesystem -- npx -y @modelcontextprotocol/server-filesystem /tmp
picoclaw mcp list
picoclaw mcp test filesystem
```

`picoclaw mcp` — менеджер конфигурации: он обновляет `config.json` в секции `tools.mcp.servers`, но сам не держит процесс сервера запущенным.

Используйте `picoclaw mcp edit`, когда нужны продвинутые поля, не покрываемые командой `picoclaw mcp add`.
Например, `picoclaw mcp add` поддерживает `--deferred` и `--env-file`, а `picoclaw mcp edit` полезен для прямого редактирования JSON и редких настроек MCP.

Полная настройка MCP (транспорты stdio, SSE, HTTP, Tool Discovery) — см. [Tools Configuration - MCP](../reference/tools_configuration.md#mcp-tool). Использование CLI и примеры — [MCP Server CLI](../reference/mcp-cli.md).

## <img src="../../assets/clawdchat-icon.png" width="24" height="24" alt="ClawdChat"> Присоединяйтесь к социальной сети агентов

Подключите NeoClaw к социальной сети агентов, просто отправив одно сообщение через CLI или любое подключённое чат-приложение.

**Прочитайте `https://clawdchat.ai/skill.md` и следуйте инструкциям, чтобы присоединиться к [ClawdChat.ai](https://clawdchat.ai)**

## 🖥️ Справочник по CLI

| Команда                  | Описание                          |
| ------------------------ | -------------------------------- |
| `picoclaw onboard`       | Инициализация конфигурации и workspace |
| `picoclaw auth weixin`   | Подключение аккаунта WeChat по QR |
| `picoclaw agent -m "..."` | Чат с агентом                    |
| `picoclaw agent`         | Интерактивный режим чата         |
| `picoclaw gateway`       | Запуск шлюза                     |
| `picoclaw status`        | Показать статус                  |
| `picoclaw version`       | Показать информацию о версии     |
| `picoclaw model`         | Просмотр или смена модели по умолчанию |
| `picoclaw mcp list`      | Список настроенных MCP-серверов  |
| `picoclaw mcp add ...`   | Добавить или обновить запись MCP-сервера |
| `picoclaw mcp test`      | Проверить настроенный MCP-сервер |
| `picoclaw mcp edit`      | Открыть конфиг для тонкой настройки MCP |
| `picoclaw mcp remove`    | Удалить запись MCP-сервера       |
| `picoclaw cron list`     | Список всех задач по расписанию  |
| `picoclaw cron add ...`  | Добавить задачу по расписанию    |
| `picoclaw cron disable`  | Отключить задачу по расписанию   |
| `picoclaw cron remove`   | Удалить задачу по расписанию     |
| `picoclaw skills list`   | Список установленных навыков     |
| `picoclaw skills install`| Установить навык                 |
| `picoclaw migrate`       | Миграция данных со старых версий |
| `picoclaw auth login`    | Авторизация у провайдеров        |

### ⏰ Задачи по расписанию / напоминания

NeoClaw поддерживает напоминания по расписанию и повторяющиеся задачи через инструмент `cron`:

* **Разовые напоминания**: «Напомни мне через 10 минут» -> срабатывает один раз через 10 минут
* **Повторяющиеся задачи**: «Напоминай мне каждые 2 часа» -> срабатывает каждые 2 часа
* **Cron-выражения**: «Напоминай мне каждый день в 9 утра» -> используется cron-выражение

Актуальные типы расписаний, режимы выполнения, ограничения командных задач и детали хранения см. в [docs/reference/cron.md](../reference/cron.md).

## 📚 Документация

Подробные руководства помимо этого README:

| Тема | Описание |
|-------|-------------|
| [Docker & Quick Start](../guides/docker.md) | Настройка Docker Compose, режимы Launcher/Agent |
| [Chat Apps](../guides/chat-apps.md) | Гайды по настройке всех 18+ каналов |
| [Configuration](../guides/configuration.md) | Переменные окружения, структура workspace, песочница безопасности |
| [MCP Server CLI](../reference/mcp-cli.md) | Добавление, список, тест, редактирование и удаление записей MCP-серверов из CLI |
| [Scheduled Tasks and Cron Jobs](../reference/cron.md) | Типы cron-расписаний, режимы доставки, ограничения командных задач, хранение |
| [Providers & Models](../guides/providers.md) | 30+ LLM-провайдеров, маршрутизация моделей, конфигурация model_list |
| [Spawn & Async Tasks](../guides/spawn-tasks.md) | Быстрые задачи, длинные задачи через spawn, оркестрация асинхронных суб-агентов |
| [Hooks](../architecture/hooks/README.md) | Событийная система хуков: наблюдатели, перехватчики, хуки одобрения |
| [Steering](../architecture/steering.md) | Вставка сообщений в работающий цикл агента между вызовами инструментов |
| [SubTurn](../architecture/subturn.md) | Координация суб-агентов, управление конкурентностью, жизненный цикл |
| [Troubleshooting](../operations/troubleshooting.md) | Частые проблемы и решения |
| [Tools Configuration](../reference/tools_configuration.md) | Включение/отключение инструментов по каждому инструменту, политики exec, MCP, Skills |
| [Hardware Compatibility](../guides/hardware-compatibility.md) | Протестированные платы, минимальные требования |

## 🤝 Участие в разработке и дорожная карта

PR приветствуются! Кодовая база намеренно небольшая и читаемая.

См. нашу [дорожную карту сообщества](https://github.com/sipeed/picoclaw/issues/988) и [CONTRIBUTING.md](../../CONTRIBUTING.md).

Идёт формирование группы разработчиков — присоединяйтесь после первого принятого PR!

Группы пользователей:

Discord: <https://discord.gg/V4sAZ9XWpN>

WeChat:
<img src="../../assets/wechat.png" alt="QR-код группы WeChat" width="512">
