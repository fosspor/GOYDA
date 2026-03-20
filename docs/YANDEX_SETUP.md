# Yandex Cloud Integration Guide

## Подготовка к использованию Yandex LLM API

Это руководство описывает, как настроить и использовать Yandex LLM для генерации маршрутов в GOYDA.

## 1. Создание Yandex Cloud Account

### Требования
- Активная учетная запись Yandex
- Платежный метод для облака (даже для бесплатного тарифа)

### Шаги
1. Перейти на https://cloud.yandex.com
2. Создать или авторизоваться в аккаунте
3. В консоли создать новый проект/организацию
4. Перейти в **IAM & Controls** → **Service Accounts**
5. Создать новый сервис-аккаунт
6. Добавить роль **AI.GeneralUser** или **AI.Generalist**

## 2. Получение API ключей

### Способ 1: API Key (простой)

```bash
# В консоли Yandex Cloud:
# IAM & Controls → Service Accounts → выбрать аккаунт
# Вкладка "API Keys" → "Create API Key"

# Скопировать ключ и добавить в .env:
YANDEX_API_KEY=AQVNxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
YANDEX_FOLDER_ID=b1g12345678901234567890  # Видно в консоли
```

### Способ 2: IAM Token (более безопасный)

```bash
# Получить IAM token через OAuth
curl -X POST https://iam.api.cloud.yandex.net/iam/v1/tokens \
  -d "{\"yandexPassportOauthToken\":\"YA_TOKEN_HERE\"}" \
  -H "Content-Type: application/json"

# Добавить в .env:
YANDEX_IAM_TOKEN=t1.9eu...  # Token (действует 12 часов)
YANDEX_FOLDER_ID=b1g12345678901234567890
```

## 3. Настройка .env

```env
# backend/.env

# Yandex Cloud
YANDEX_FOLDER_ID=b1g12345678901234567890
YANDEX_API_KEY=AQVNxxxxxxxxxxxxxxxxxxxxxxxxx
# Или вместо API Key:
# YANDEX_IAM_TOKEN=t1.9eu...
```

## 4. Тестирование API

### Python тест

```python
# backend/test_yandex.py
import os
from ai_recommendations.services import get_yandex_service

service = get_yandex_service()

# Тестовый запрос
prompt = """Напишите короткое приветствие на русском языке."""

response = service.call_completion_api(prompt)
print(f"Response: {response}")
```

### Запуск теста

```bash
cd backend
source ../venv/bin/activate
python test_yandex.py
```

## 5. Использование в приложении

### Генерация маршрута

```python
from ai_recommendations.services import get_yandex_service

service = get_yandex_service()

# Параметры от пользователя
interests = ['wine', 'gastronomy', 'nature']
budget = 50000
duration = 3
season = 'autumn'
group_info = {
    'size': 2,
    'ages': [30, 35],
    'with_children': False,
}

# Генерировать prompt
prompt = service.generate_route_prompt(
    user_interests=interests,
    budget=budget,
    duration=duration,
    season=season,
    group_info=group_info
)

# Получить ответ от LLM
response = service.call_completion_api(prompt)

# Парсить и сохранить
import json
route_data = json.loads(response)
# ... create RouteRecommendation
```

## 6. API Documentation

### Models доступные в Yandex Cloud

```
gpt://folder_id/yandexgpt-lite       # Легкая версия (быстрая, дешевая)
gpt://folder_id/yandexgpt            # Обычная версия
gpt://folder_id/yandexgpt-pro         # Pro версия (мощная)
```

### Параметры запроса

```python
{
    'modelUri': 'gpt://folder_id/yandexgpt-lite',
    'completionOptions': {
        'stream': False,
        'temperature': 0.7,      # 0 = детерминированный, 1 = креативный
        'maxTokens': 2000        # макс токены в ответе
    },
    'messages': [
        {
            'role': 'user',
            'text': 'Your prompt here'
        }
    ]
}
```

### Example Response

```json
{
    "result": {
        "alternatives": [
            {
                "message": {
                    "role": "assistant",
                    "text": "Response text..."
                },
                "status": "ALTERNATIVE_STATUS_FINAL"
            }
        ],
        "usage": {
            "inputTextTokens": 150,
            "completionTokens": 200
        },
        "modelVersion": "19-12-2024"
    }
}
```

## 7. Pricing

### Yandex LLM Pricing (примерно)

- **yandexgpt-lite**: 1₽ за 1000 input токов, 3₽ за 1000 output токенов
- **yandexgpt**: 1.5₽ за 1000 input, 6₽ за 1000 output
- **yandexgpt-pro**: 3₽ за 1000 input, 15₽ за 1000 output

Первые месяцы могут быть дешевле/бесплатнее по промо.

## 8. Оптимизация затрат

### Советы
1. Использовать **yandexgpt-lite** для MVP
2. Кэшировать часто запрашиваемые промпты
3. Логировать все запросы в БД (есть в `AIPromptLog`)
4. Батчить запросы где возможно
5. Отключать streaming для простых задач

## 9. Troubleshooting

### "Invalid API Key"
```bash
# Проверить ключ в консоли
# Settings → API Keys → убедиться что ключ активен
```

### "Folder ID not found"
```bash
# Взять из консоли:
# Home → Project → Settings → Folder ID
```

### "Quota exceeded"
```bash
# Yandex Cloud обычно дает 500 запросов в день на бесплатном плане
# Если нужно больше, обновить тариф
```

### "Connection timeout"
```bash
# Проверить статус API:
# https://status.yandex.cloud/
# Увеличить timeout в services.py
```

## 10. Production Deployment

### Рекомендации
1. **Использовать IAM Token вместо API Key** для security
2. **Кэшировать успешные ответы** в Redis
3. **Rate limit** на генерацию маршрутов
4. **Мониторить** использование токенов
5. **Fallback** на cached версии при сбое API

### Пример Rate Limiting

```python
from django.core.cache import cache

def generate_route_with_limit(user, interests):
    cache_key = f'route_gen_{user.id}'
    
    # Лимит 5 запросов в час
    if cache.get(cache_key, 0) >= 5:
        raise RateLimitException("Too many requests")
    
    # Генерировать
    route = service.call_completion_api(prompt)
    
    # Инкрементировать счетчик
    cache.set(cache_key, cache.get(cache_key, 0) + 1, timeout=3600)
    
    return route
```

## 11. Примеры промптов для маршрутов

### Для гастротуризма
```
Напишите 3-дневный маршрут по локальным ресторанам и фермам Краснодарского края для пары, 
интересующейся органической едой и местной кухней. Бюджет 50000₽.
```

### Для винного туризма
```
Создайте винный маршрут «Краснодар по винтажам» на 4 дня для 4 человек. 
Включить 3 винодельни, винный бар и ресторан с красными винами.
```

### Для активного туризма
```
Маршрут для трекинга и природы в Краснодарском крае на выходные (2 дня) для молодой пары (25-30 лет).
Включить пешеходные маршруты, горные виды и палаточный отдых.
```

## Полезные ссылки

- [Yandex Cloud Documentation](https://cloud.yandex.com/docs)
- [Yandex LLM API Reference](https://cloud.yandex.com/docs/foundation-models/quickstart)
- [Pricing Calculator](https://cloud.yandex.com/pricing)
- [API Explorer](https://cloud.yandex.com/marketplace/products/yandexcloud/foundation-models)

---

**Note**: Убедитесь, что все API ключи хранятся в `.env` и никогда не коммитятся в Git!
