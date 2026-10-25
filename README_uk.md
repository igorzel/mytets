# mytets

[English](README.md)

`mytets` — невелика утиліта командного рядка, яка виводить випадкові фрази з
творів українського письменника й драматурга [Леся Подерв'янського](https://uk.wikipedia.org/wiki/%D0%9B%D0%B5%D1%81%D1%8C_%D0%9F%D0%BE%D0%B4%D0%B5%D1%80%D0%B2%27%D1%8F%D0%BD%D1%81%D1%8C%D0%BA%D0%B8%D0%B9).
Можна отримати одну фразу, вивести список унікальних фраз або вибрати формат
JSON для використання у скриптах. Усі фрази вбудовані в програму, тож для
роботи не потрібні обліковий запис чи підключення до інтернету.

## Встановлення

Рекомендований спосіб — встановити програму зі Snap Store:

```bash
sudo snap install mytets
```

Щоб спробувати останню розробницьку версію замість стабільного випуску:

```bash
sudo snap install mytets --channel=edge
```

Також програму можна встановити за допомогою Go:

```bash
go install github.com/igorzel/mytets/cmd/mytets@latest
```

Якщо після встановлення команда `mytets` не знаходиться, переконайтеся, що
каталог `bin` Go додано до змінної `PATH`.

## Швидкий старт

Вивести одну випадкову фразу:

```bash
mytets one
```

Вивести п'ять унікальних фраз (значення за замовчуванням):

```bash
mytets list
```

Вказати кількість фраз:

```bash
mytets list --count 3
```

За замовчуванням команди виводять звичайний текст. Додайте `--output json`,
щоб отримати дані у форматі JSON для обробки скриптами:

```bash
mytets one --output json
mytets list --count 3 --output json
mytets version --output json
```

Наприклад, команда `mytets one --output json` виводить:

```json
{"message":"Ето кайф, марожено!"}
```

## Команди

| Команда | Дія |
|---------|-----|
| `mytets one` | Вивести одну випадкову фразу. |
| `mytets list` | Вивести до п'яти унікальних фраз. |
| `mytets list --count N` | Вивести до `N` унікальних фраз. |
| `mytets version` | Вивести версію програми. |
| `mytets --help` | Показати доступні команди й параметри. |

Додайте `--output json` до команд `one`, `list` або `version`, щоб отримати
результат у форматі JSON. Щоб переглянути довідку для окремої команди, виконайте
`mytets <команда> --help`.

## Збирання з вихідного коду

Щоб самостійно зібрати виконуваний файл, установіть Go та виконайте:

```bash
git clone https://github.com/igorzel/mytets.git
cd mytets
make build
./bin/mytets one
```

Щоб зібрати й установити локальний Snap-пакунок:

```bash
make snap
sudo snap install ./mytets_*.snap --dangerous
```

## Публікація в Snap Store

Ці інструкції призначені для супровідників проєкту. Зареєструйте назву Snap
один раз, а потім авторизуйтеся на кожній машині, з якої публікуватимете:

```bash
make snap-register
make snap-login
```

Опублікувати збірку в каналі `edge`:

```bash
make snap-publish SNAP_CHANNEL=edge
```

Перевірити опубліковані ревізії та канали:

```bash
make snap-status
```

Пізніше наявну ревізію можна перевести, наприклад, у стабільний канал:

```bash
make snap-release REVISION=<store-revision> SNAP_CHANNEL=stable
```

`make setup` готує локальне середовище LXD для збирання Snap і потребує
`sudo`. Команда `make snap-upload` завантажує пакунок, але не публікує його в
каналі.
