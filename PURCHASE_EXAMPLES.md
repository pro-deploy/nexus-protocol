# 🛒 Примеры запросов с покупками

## 📋 Обзор

Nexus Protocol поддерживает обработку запросов с намерением покупки товаров и услуг. Система автоматически определяет коммерческие намерения и предоставляет возможности для выполнения покупок.

## 🎯 Типы запросов

### Query Types

Nexus Protocol определяет три типа запросов:

1. **`information_only`** - только информационный запрос (рецепты, инструкции)
2. **`with_purchases_services`** - запрос с возможностью покупки
3. **`mixed`** - смешанный тип (информация + покупки)

## 📝 Примеры запросов

### Пример 1: Поиск товара с геолокацией

**Запрос:**
```json
{
  "query": "Найди где рядом продается кокакола и купи литровую бутылку колы заберу самостоятельно",
  "language": "ru",
  "context": {
    "user_id": "user-123",
    "location": {
      "latitude": 55.7558,
      "longitude": 37.6173,
      "accuracy": 50
    },
    "locale": "ru-RU",
    "currency": "RUB",
    "region": "RU"
  }
}
```

**Обработка:**
1. AI сервис определяет домен: `commerce`
2. Определяет тип запроса: `with_purchases_services`
3. Извлекает параметры:
   - Товар: "кокакола", "литровая бутылка"
   - Локация: "рядом" (использует GPS из context)
   - Способ получения: "самовывоз"

**Анализ доменов AI:**
```
AI определяет домен "commerce" на основе ключевых слов:
- "купить", "заказать" → намерение покупки
- "кокакола", "бутылка" → товар
- "рядом", "самовывоз" → локация и способ получения
- Общая уверенность: 0.91
- Релевантность: 0.95
```

**Ответ:**
```json
{
  "execution_id": "exec-789",
  "intent_id": "intent-abc-123",
  "status": "completed",
  "query_type": "with_purchases_services",
  "sections": [
    {
      "domain_id": "commerce",
      "title": "Коммерческие предложения",
      "status": "success",
      "response_time_ms": 245,
      "results": [
        {
          "id": "product-456",
          "type": "product_purchase",
          "title": "Coca-Cola 1л бутылка",
          "description": "Найдено в 3 магазинах рядом с вами",
          "data": {
            "price": "89 ₽",
            "availability": "в наличии",
            "rating": "4.5",
            "stores": [
              {
                "name": "Пятерочка",
                "distance": "200м",
                "address": "ул. Ленина, 15",
                "pickup_available": true,
                "work_hours": "Круглосуточно",
                "phone": "+7 (495) 123-45-67"
              },
              {
                "name": "Магнит",
                "distance": "350м",
                "address": "ул. Пушкина, 8",
                "pickup_available": true,
                "work_hours": "08:00-22:00"
              }
            ]
          },
          "relevance": 0.95,
          "confidence": 0.88,
          "actions": [
            {
              "type": "reserve_product",
              "label": "Зарезервировать товар",
              "method": "POST",
              "url": "/api/v1/commerce/reserve",
              "confirm_text": "Подтвердить резервирование?"
            },
            {
              "type": "purchase",
              "label": "Купить сейчас",
              "method": "POST",
              "url": "/api/v1/commerce/purchase"
            }
          ]
        }
      ]
    }
  ],
  "domain_analysis": {
    "selected_domains": [
      {
        "domain_id": "commerce",
        "name": "Коммерция и покупки",
        "type": "commerce",
        "confidence": 0.91,
        "relevance": 0.95,
        "reason": "Высокая уверенность: найдены ключевые слова покупки и локации",
        "priority": 80,
        "capabilities": [
          {"type": "search", "description": "Поиск товаров"},
          {"type": "purchase", "description": "Оформление покупок"}
        ]
      }
    ],
    "rejected_domains": [
      {
        "domain_id": "travel",
        "confidence": 0.12,
        "reason": "Низкая релевантность: нет признаков путешествий"
      }
    ],
    "confidence": 0.91,
    "analysis_algorithm": "hybrid_keyword_semantic"
  },
  "processing_time_ms": 267
}
```

### Пример 2: Покупка с доставкой

**Запрос:**
```json
{
  "query": "купи пиццу пепперони доставь на дом",
  "language": "ru",
  "context": {
    "user_id": "user-123",
    "location": {
      "latitude": 55.7558,
      "longitude": 37.6173
    },
    "address": {
      "street": "ул. Тверская, 10",
      "apartment": "45"
    }
  }
}
```

**Ответ:**
```json
{
  "execution_id": "exec-790",
  "intent_id": "intent-def-456",
  "status": "completed",
  "query_type": "with_purchases_services",
  "sections": [
    {
      "domain_id": "commerce",
      "title": "Заказ еды",
      "status": "success",
      "response_time_ms": 189,
      "results": [
        {
          "id": "food-order-789",
          "type": "product_purchase",
          "title": "Пицца Пепперони",
          "description": "Доставка из ближайшего ресторана",
          "data": {
            "price": "599 ₽",
            "delivery_time": "30-45 минут",
            "delivery_fee": "150 ₽",
            "total": "749 ₽",
            "restaurants": [
              {
                "name": "Додо Пицца",
                "distance": "2.5 км",
                "rating": 4.8,
                "delivery_available": true,
                "work_hours": "10:00-23:00"
              }
            ]
          },
          "relevance": 0.98,
          "confidence": 0.92,
          "actions": [
            {
              "type": "purchase_with_delivery",
              "label": "Заказать с доставкой",
              "method": "POST",
              "url": "/api/v1/commerce/purchase",
              "confirm_text": "Подтвердить заказ на 749 ₽ с доставкой?"
            },
            {
              "type": "add_to_cart",
              "label": "Добавить в корзину",
              "method": "POST",
              "url": "/api/v1/commerce/cart/add"
            }
          ]
        }
      ]
    }
  ],
  "domain_analysis": {
    "selected_domains": [
      {
        "domain_id": "commerce",
        "name": "Коммерция и покупки",
        "type": "commerce",
        "confidence": 0.88,
        "relevance": 0.94,
        "reason": "Высокая уверенность: заказ еды с доставкой",
        "priority": 80
      }
    ],
    "confidence": 0.88,
    "analysis_algorithm": "hybrid_keyword_semantic"
  },
  "processing_time_ms": 201
}
```

### Пример 3: Сравнение цен

**Запрос:**
```json
{
  "query": "где дешевле купить iPhone 15 рядом со мной",
  "language": "ru",
  "context": {
    "location": {
      "latitude": 55.7558,
      "longitude": 37.6173
    }
  }
}
```

**Ответ:**
```json
{
  "execution_id": "exec-791",
  "intent_id": "intent-ghi-789",
  "status": "completed",
  "query_type": "with_purchases_services",
  "sections": [
    {
      "domain_id": "commerce",
      "title": "Сравнение цен",
      "status": "success",
      "response_time_ms": 312,
      "results": [
        {
          "id": "price-comparison-101",
          "type": "product_comparison",
          "title": "iPhone 15 128GB - сравнение цен",
          "description": "Найдено 5 предложений рядом с вами",
          "data": {
            "best_price": "89990 ₽",
            "price_range": "89990 ₽ - 95990 ₽",
            "average_price": "92790 ₽",
            "stores": [
              {
                "name": "М.Видео",
                "price": "89990 ₽",
                "distance": "500м",
                "in_stock": true,
                "rating": 4.3,
                "delivery_available": true
              },
              {
                "name": "Эльдорадо",
                "price": "91990 ₽",
                "distance": "800м",
                "in_stock": true,
                "rating": 4.1,
                "delivery_available": true
              },
              {
                "name": "DNS",
                "price": "93990 ₽",
                "distance": "1.2 км",
                "in_stock": false,
                "rating": 4.0,
                "delivery_available": true
              }
            ]
          },
          "relevance": 0.96,
          "confidence": 0.89,
          "actions": [
            {
              "type": "purchase",
              "label": "Купить в М.Видео (89990 ₽)",
              "method": "POST",
              "url": "/api/v1/commerce/purchase?store=mvideo&product=iphone15",
              "confirm_text": "Подтвердить покупку iPhone 15 за 89990 ₽?"
            },
            {
              "type": "reserve",
              "label": "Зарезервировать в М.Видео",
              "method": "POST",
              "url": "/api/v1/commerce/reserve?store=mvideo&product=iphone15"
            },
            {
              "type": "compare_prices",
              "label": "Показать все цены",
              "method": "GET",
              "url": "/api/v1/commerce/compare?product=iphone15&location=current"
            }
          ]
        }
      ]
    }
  ],
  "domain_analysis": {
    "selected_domains": [
      {
        "domain_id": "commerce",
        "name": "Коммерция и покупки",
        "type": "commerce",
        "confidence": 0.93,
        "relevance": 0.97,
        "reason": "Высокая уверенность: запрос сравнения цен на товар",
        "priority": 85
      }
    ],
    "confidence": 0.93,
    "analysis_algorithm": "hybrid_keyword_semantic"
  },
  "ranking": {
    "items": [
      {"id": "price-comparison-101", "score": 0.96, "rank": 1}
    ],
    "algorithm": "weighted_relevance_confidence"
  },
  "processing_time_ms": 334
}
```

## 🔧 Использование через SDK

### Go SDK

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    nexus "github.com/pro-deploy/nexus-protocol/sdk/go/client"
    "github.com/pro-deploy/nexus-protocol/sdk/go/types"
)

func main() {
    client := nexus.NewClient(nexus.Config{
        BaseURL: "https://api.nexus.dev",
        Token:   "your-jwt-token",
    })
    
    // Запрос с покупкой
    req := &types.ExecuteTemplateRequest{
        Query:    "Найди где рядом продается кокакола и купи литровую бутылку колы заберу самостоятельно",
        Language: "ru",
        Context: &types.UserContext{
            UserID: "user-123",
            Location: &types.UserLocation{
                Latitude:  55.7558,
                Longitude: 37.6173,
                Accuracy: 50,
            },
            Locale:   "ru-RU",
            Currency: "RUB",
            Region:   "RU",
        },
    }
    
    result, err := client.ExecuteTemplate(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }
    
    // Анализ доменов AI
    if result.DomainAnalysis != nil {
        fmt.Printf("🤖 AI анализ доменов:\n")
        fmt.Printf("   Общая уверенность: %.2f\n", result.DomainAnalysis.Confidence)
        fmt.Printf("   Алгоритм: %s\n", result.DomainAnalysis.AnalysisAlgorithm)

        fmt.Printf("   Выбранные домены:\n")
        for _, domain := range result.DomainAnalysis.SelectedDomains {
            fmt.Printf("     • %s (уверенность: %.2f, приоритет: %d)\n",
                domain.Name, domain.Confidence, domain.Priority)
        }

        if len(result.DomainAnalysis.RejectedDomains) > 0 {
            fmt.Printf("   Отклоненные домены:\n")
            for _, domain := range result.DomainAnalysis.RejectedDomains {
                fmt.Printf("     • %s (уверенность: %.2f)\n",
                    domain.Name, domain.Confidence)
            }
        }
    }

    // Проверка типа запроса
    if result.QueryType == "with_purchases_services" {
        fmt.Println("\n✅ Запрос с возможностью покупки")

        // Обработка результатов по доменам
        for _, section := range result.Sections {
            fmt.Printf("\n🏪 Домен: %s (%s)\n", section.Title, section.DomainID)
            fmt.Printf("   Статус: %s, Время: %dms\n", section.Status, section.ResponseTimeMS)

            for _, item := range section.Results {
                fmt.Printf("\n   📦 %s\n", item.Title)
                if item.Description != "" {
                    fmt.Printf("      %s\n", item.Description)
                }
                fmt.Printf("      ⭐ Релевантность: %.2f, Уверенность: %.2f\n",
                    item.Relevance, item.Confidence)

                // Показать ключевые данные
                if price, ok := item.Data["price"]; ok {
                    fmt.Printf("      💰 Цена: %s\n", price)
                }
                if time, ok := item.Data["estimated_time"]; ok {
                    fmt.Printf("      ⏱️ Время: %s\n", time)
                }

                // Действия для взаимодействия
                if len(item.Actions) > 0 {
                    fmt.Printf("      🎯 Действия:\n")
                    for _, action := range item.Actions {
                        confirm := ""
                        if action.ConfirmText != "" {
                            confirm = fmt.Sprintf(" (Подтверждение: %s)", action.ConfirmText)
                        }
                        fmt.Printf("         • %s: %s %s%s\n",
                            action.Type, action.Label, action.URL, confirm)
                    }
                }
            }
        }

        // Обработка workflow если есть
        if result.Workflow != nil {
            fmt.Printf("\n🔄 Workflow (%d шагов):\n", len(result.Workflow.Steps))
            for _, step := range result.Workflow.Steps {
                status := "⏳"
                if step.Status == "completed" {
                    status = "✅"
                } else if step.Status == "failed" {
                    status = "❌"
                }
                fmt.Printf("   %s Шаг %d: %s (%s) - %s\n",
                    status, step.Step, step.Action, step.Domain, step.Status)
            }
        }

        // Статистика обработки
        fmt.Printf("\n📊 Статистика:\n")
        fmt.Printf("   Общее время обработки: %dms\n", result.ProcessingTimeMS)
        fmt.Printf("   Домены обработано: %d\n", len(result.Sections))
        fmt.Printf("   Execution ID: %s\n", result.ExecutionID)
    }
}
```

## 🎯 Логика обработки

### 1. Определение намерения покупки

Система анализирует ключевые слова:
- **Покупка**: "купить", "заказать", "приобрести", "продается"
- **Локация**: "рядом", "близко", "недалеко", "поблизости"
- **Способ получения**: "самовывоз", "заберу", "доставка", "привези"

### 2. Извлечение параметров

- **Товар/услуга**: название, характеристики, количество
- **Локация**: использование GPS из `context.location`
- **Предпочтения**: способ получения, цена, магазин

### 3. Поиск и ранжирование

- Поиск товаров в магазинах поблизости
- Расчет расстояний от пользователя
- Сравнение цен и наличия
- Ранжирование по релевантности

### 4. Формирование ответа

- Список доступных вариантов
- Информация о магазинах
- Действия для покупки/резервирования
- Метаданные (цены, расстояния, время работы)

## 🔗 Дальнейшие действия

После получения ответа пользователь может:

1. **Резервировать товар** - через action `reserve_product`
2. **Купить товар** - через action `purchase`
3. **Получить детали** - запрос дополнительной информации
4. **Сравнить варианты** - просмотр других магазинов

### Пример 4: Комплексный многошаговый сценарий

**Запрос:**
```json
{
  "query": "закажи в макдоналдсе карточку фри, оплати, введи адрес доставки, и напоминай когда курьер выедет с заказом выпить таблетки, и через два часа выпить еще одни таблетки",
  "language": "ru",
  "context": {
    "user_id": "user-123",
    "location": {
      "latitude": 55.7558,
      "longitude": 37.6173,
      "accuracy": 50
    },
    "locale": "ru-RU",
    "currency": "RUB",
    "region": "RU"
  }
}
```

**Анализ доменов AI:**
```
AI определяет несколько доменов на основе комплексного анализа:
1. **Commerce** (confidence: 0.94) - "закажи в макдоналдсе карточку фри"
2. **Payment** (confidence: 0.89) - "оплати"
3. **Delivery** (confidence: 0.91) - "введи адрес доставки"
4. **Notifications** (confidence: 0.87) - "напоминай", "выпить таблетки"
5. **Health** (confidence: 0.82) - "таблетки"

Общая уверенность анализа: 0.91
Алгоритм: hybrid_keyword_semantic + context_analysis
```

**Ответ:**
```json
{
  "execution_id": "exec-890",
  "intent_id": "intent-multi-999",
  "status": "completed",
  "query_type": "with_purchases_services",
  "sections": [
    {
      "domain_id": "commerce",
      "title": "Заказ еды",
      "status": "success",
      "response_time_ms": 234,
      "results": [
        {
          "id": "order-mcd-123",
          "type": "food_order",
          "title": "Картошка фри (стандартная)",
          "description": "Заказ в Макдональдсе с доставкой",
          "data": {
            "restaurant": "Макдональдс",
            "item": "Картошка фри",
            "size": "стандартная",
            "price": "89 ₽",
            "delivery_available": true,
            "estimated_time": "30-45 минут",
            "restaurant_rating": 4.2
          },
          "relevance": 0.98,
          "confidence": 0.95,
          "actions": [
            {
              "type": "add_to_cart",
              "label": "Добавить в корзину",
              "method": "POST",
              "url": "/api/v1/commerce/cart/add"
            },
            {
              "type": "order_now",
              "label": "Заказать сейчас",
              "method": "POST",
              "url": "/api/v1/commerce/order",
              "confirm_text": "Подтвердить заказ картошки фри за 89 ₽?"
            }
          ]
        }
      ]
    },
    {
      "domain_id": "payment",
      "title": "Оплата заказа",
      "status": "pending",
      "response_time_ms": 145,
      "results": [
        {
          "id": "payment-456",
          "type": "payment_processing",
          "title": "Оплата заказа",
          "description": "Выберите способ оплаты",
          "data": {
            "amount": "89 ₽",
            "order_id": "order-mcd-123",
            "currency": "RUB",
            "payment_methods": ["card", "apple_pay", "google_pay", "cash"],
            "estimated_processing_time": "5 сек"
          },
          "relevance": 0.95,
          "confidence": 0.90,
          "actions": [
            {
              "type": "process_payment",
              "label": "Оплатить картой",
              "method": "POST",
              "url": "/api/v1/payment/process",
              "confirm_text": "Подтвердить оплату 89 ₽?"
            },
            {
              "type": "select_payment_method",
              "label": "Выбрать способ оплаты",
              "method": "POST",
              "url": "/api/v1/payment/method"
            }
          ]
        }
      ]
    },
    {
      "domain_id": "delivery",
      "title": "Адрес доставки",
      "status": "pending",
      "response_time_ms": 167,
      "results": [
        {
          "id": "delivery-789",
          "type": "delivery_address",
          "title": "Введите адрес доставки",
          "description": "Укажите адрес для доставки заказа из Макдональдса",
          "data": {
            "order_id": "order-mcd-123",
            "delivery_type": "courier",
            "estimated_delivery_time": "30-45 минут",
            "delivery_fee": "0 ₽",
            "free_delivery_threshold": "500 ₽",
            "current_address": "ул. Тверская, 10"
          },
          "relevance": 0.92,
          "confidence": 0.88,
          "actions": [
            {
              "type": "set_delivery_address",
              "label": "Подтвердить адрес доставки",
              "method": "POST",
              "url": "/api/v1/delivery/address",
              "confirm_text": "Доставить на ул. Тверская, 10?"
            },
            {
              "type": "edit_address",
              "label": "Изменить адрес",
              "method": "PUT",
              "url": "/api/v1/delivery/address"
            },
            {
              "type": "use_saved_address",
              "label": "Выбрать другой адрес",
              "method": "GET",
              "url": "/api/v1/delivery/addresses"
            }
          ]
        }
      ]
    },
    {
      "domain_id": "notifications",
      "title": "Напоминания и уведомления",
      "status": "success",
      "response_time_ms": 98,
      "results": [
        {
          "id": "reminder-courier-001",
          "type": "reminder",
          "title": "Напоминание о выезде курьера",
          "description": "Уведомление когда курьер выедет с заказом",
          "data": {
            "reminder_type": "courier_departure",
            "order_id": "order-mcd-123",
            "trigger": "order_status_changed",
            "status": "courier_assigned",
            "notification_methods": ["push", "sms"],
            "priority": "high"
          },
          "relevance": 0.90,
          "confidence": 0.85,
          "actions": [
            {
              "type": "create_reminder",
              "label": "Создать напоминание",
              "method": "POST",
              "url": "/api/v1/notifications/reminders/create"
            },
            {
              "type": "customize_notifications",
              "label": "Настроить уведомления",
              "method": "PUT",
              "url": "/api/v1/notifications/preferences"
            }
          ]
        },
        {
          "id": "reminder-medication-001",
          "type": "medication_reminder",
          "title": "Выпить таблетки (сейчас)",
          "description": "Напоминание о приеме лекарств",
          "data": {
            "reminder_type": "medication",
            "medication": "таблетки",
            "dosage": "1 таблетка",
            "time": "немедленно",
            "repeat": false,
            "importance": "high",
            "notification_methods": ["push", "voice"]
          },
          "relevance": 0.88,
          "confidence": 0.82,
          "actions": [
            {
              "type": "create_reminder",
              "label": "Напомнить сейчас",
              "method": "POST",
              "url": "/api/v1/notifications/reminders/create"
            },
            {
              "type": "snooze_reminder",
              "label": "Отложить на 5 минут",
              "method": "POST",
              "url": "/api/v1/notifications/reminders/snooze"
            }
          ]
        },
        {
          "id": "reminder-medication-002",
          "type": "medication_reminder",
          "title": "Выпить таблетки (через 2 часа)",
          "description": "Отложенное напоминание о приеме лекарств",
          "data": {
            "reminder_type": "medication",
            "medication": "таблетки",
            "dosage": "1 таблетка",
            "time": "через 2 часа",
            "delay_hours": 2,
            "scheduled_time": "2025-01-15T16:30:00Z",
            "repeat": false,
            "importance": "high"
          },
          "relevance": 0.88,
          "confidence": 0.82,
          "actions": [
            {
              "type": "create_reminder",
              "label": "Напомнить через 2 часа",
              "method": "POST",
              "url": "/api/v1/notifications/reminders/create"
            },
            {
              "type": "edit_schedule",
              "label": "Изменить время",
              "method": "PUT",
              "url": "/api/v1/notifications/reminders/schedule"
            }
          ]
        }
      ]
    },
    {
      "domain_id": "health",
      "title": "Здоровье и лекарства",
      "status": "success",
      "response_time_ms": 134,
      "results": [
        {
          "id": "medication-tracker-001",
          "type": "medication_tracking",
          "title": "Отслеживание приема лекарств",
          "description": "Контроль регулярного приема таблеток",
          "data": {
            "medication_name": "таблетки",
            "schedule": ["сейчас", "через 2 часа"],
            "completed": [false, false],
            "next_reminder": "через 2 часа",
            "adherence_rate": "0%",
            "side_effects_tracking": true
          },
          "relevance": 0.85,
          "confidence": 0.78,
          "actions": [
            {
              "type": "mark_taken",
              "label": "Отметить прием",
              "method": "POST",
              "url": "/api/v1/health/medication/taken"
            },
            {
              "type": "view_schedule",
              "label": "Показать расписание",
              "method": "GET",
              "url": "/api/v1/health/medication/schedule"
            },
            {
              "type": "report_side_effects",
              "label": "Сообщить о побочных эффектах",
              "method": "POST",
              "url": "/api/v1/health/side-effects/report"
            }
          ]
        }
      ]
    }
  ],
  "domain_analysis": {
    "selected_domains": [
      {
        "domain_id": "commerce",
        "name": "Коммерция и покупки",
        "type": "commerce",
        "confidence": 0.94,
        "relevance": 0.98,
        "reason": "Высокая уверенность: прямое указание заказа еды",
        "priority": 85
      },
      {
        "domain_id": "delivery",
        "name": "Доставка",
        "type": "delivery",
        "confidence": 0.91,
        "relevance": 0.93,
        "reason": "Высокая уверенность: указание адреса доставки",
        "priority": 80
      },
      {
        "domain_id": "payment",
        "name": "Платежи",
        "type": "payment",
        "confidence": 0.89,
        "relevance": 0.90,
        "reason": "Высокая уверенность: указание оплаты",
        "priority": 75
      },
      {
        "domain_id": "notifications",
        "name": "Уведомления",
        "type": "notifications",
        "confidence": 0.87,
        "relevance": 0.88,
        "reason": "Высокая уверенность: множественные напоминания",
        "priority": 70
      },
      {
        "domain_id": "health",
        "name": "Здоровье и медицина",
        "type": "health",
        "confidence": 0.82,
        "relevance": 0.85,
        "reason": "Средняя уверенность: прием лекарств",
        "priority": 60
      }
    ],
    "rejected_domains": [
      {
        "domain_id": "travel",
        "confidence": 0.08,
        "reason": "Очень низкая релевантность"
      },
      {
        "domain_id": "finance",
        "confidence": 0.05,
        "reason": "Нет финансовых аспектов"
      }
    ],
    "confidence": 0.91,
    "analysis_algorithm": "hybrid_keyword_semantic"
  },
  "workflow": {
    "steps": [
      {
        "step": 1,
        "action": "order_food",
        "domain": "commerce",
        "status": "completed",
        "result_id": "order-mcd-123"
      },
      {
        "step": 2,
        "action": "process_payment",
        "domain": "payment",
        "status": "pending",
        "depends_on": ["order-mcd-123"]
      },
      {
        "step": 3,
        "action": "set_delivery_address",
        "domain": "delivery",
        "status": "pending",
        "depends_on": ["order-mcd-123", "payment-456"]
      },
      {
        "step": 4,
        "action": "create_reminders",
        "domain": "notifications",
        "status": "success",
        "depends_on": ["delivery-789"],
        "result_ids": ["reminder-courier-001", "reminder-medication-001", "reminder-medication-002"]
      },
      {
        "step": 5,
        "action": "track_medication",
        "domain": "health",
        "status": "success",
        "result_id": "medication-tracker-001"
      }
    ]
  },
  "ranking": {
    "items": [
      {"id": "order-mcd-123", "score": 0.98, "rank": 1},
      {"id": "payment-456", "score": 0.95, "rank": 2},
      {"id": "delivery-789", "score": 0.92, "rank": 3},
      {"id": "reminder-courier-001", "score": 0.90, "rank": 4}
    ],
    "algorithm": "weighted_relevance_confidence"
  },
  "processing_time_ms": 778
}
```

**Логика обработки AI:**

1. **Интеллектуальный анализ доменов:**
   - **Commerce** (94%): Распознавание заказа еды по ключевым словам
   - **Delivery** (91%): Определение необходимости доставки
   - **Payment** (89%): Автоматическое включение оплаты для покупок
   - **Notifications** (87%): Обнаружение запросов на напоминания
   - **Health** (82%): Распознавание медицинских аспектов (таблетки)

2. **Анализ качества ответов:**
   - Проверка полноты информации (цены, адреса, время)
   - Оценка релевантности результатов
   - Проверка доступности действий
   - Генерация предложений по улучшению

3. **Динамическая маршрутизация:**
   - Автоматический выбор наиболее подходящих доменов
   - Исключение нерелевантных доменов (< 0.3 confidence)
   - Приоритизация по важности и уверенности

4. **Многошаговый workflow:**
   ```
   Шаг 1: Заказ еды (Commerce) → Выполнено
   Шаг 2: Оплата (Payment) → Ожидает (зависит от Шага 1)
   Шаг 3: Адрес доставки (Delivery) → Ожидает (зависит от Шагов 1,2)
   Шаг 4: Напоминания (Notifications) → Выполнено (параллельно)
   Шаг 5: Отслеживание лекарств (Health) → Выполнено (параллельно)
   ```

5. **Обработка напоминаний:**
   - **Курьер**: Триггер по изменению статуса заказа
   - **Медикаменты**: Немедленное + отложенное напоминание
   - **Умные уведомления**: Разные каналы (push, SMS, voice)
   - **Отслеживание приема**: Контроль adherence rate

## 📊 Метрики и аналитика

### Метрики запросов с покупками:
- **Conversion Rate** - процент завершенных покупок
- **Domain Selection Accuracy** - точность выбора доменов AI (0.0-1.0)
- **Response Quality Score** - средняя оценка качества ответов
- **Multi-domain Success Rate** - успех обработки комплексных запросов
- **Workflow Completion Rate** - процент завершения многошаговых сценариев
- **Store Selection** - выбор магазинов пользователями
- **Distance Impact** - влияние расстояния на покупки
- **Price Sensitivity** - чувствительность к ценам
- **Reminder Effectiveness** - эффективность напоминаний
- **User Satisfaction** - оценка удовлетворенности пользователей

### AI аналитика:
- **Domain Confidence Distribution** - распределение уверенности по доменам
- **Relevance Score Trends** - тренды релевантности результатов
- **Routing Decision Success** - эффективность маршрутизации запросов
- **Quality Improvement Rate** - скорость улучшения качества ответов
- **False Positive Rate** - процент ошибочно выбранных доменов
- **Processing Time per Domain** - время обработки по доменам
- **Cache Hit Rate** - эффективность кэширования
- **Error Rate by Domain** - ошибки по доменам

### Пример аналитики:
```json
{
  "period": "2025-01-15",
  "metrics": {
    "total_requests": 15420,
    "purchase_requests": 3420,
    "conversion_rate": 0.23,
    "avg_domain_confidence": 0.87,
    "avg_response_quality": 0.91,
    "workflow_completion_rate": 0.78,
    "top_domains": [
      {"domain": "commerce", "requests": 2890, "success_rate": 0.94},
      {"domain": "delivery", "requests": 1850, "success_rate": 0.89},
      {"domain": "payment", "requests": 1640, "success_rate": 0.96}
    ]
  }
}
```

## 🔄 Workflow и последовательность действий

### Многошаговые сценарии

Система поддерживает сложные сценарии с несколькими шагами:

1. **Определение зависимостей** - система автоматически определяет порядок выполнения
2. **Workflow tracking** - отслеживание прогресса по шагам
3. **Conditional execution** - выполнение шагов зависит от предыдущих
4. **Error handling** - обработка ошибок на каждом шаге
5. **Rollback support** - возможность отмены при ошибках

### Пример использования через SDK

```go
req := &types.ExecuteTemplateRequest{
    Query: "закажи в макдоналдсе карточку фри, оплати, введи адрес доставки, и напоминай когда курьер выедет с заказом выпить таблетки, и через два часа выпить еще одни таблетки",
    Language: "ru",
    Context: &types.UserContext{
        UserID: "user-123",
        Location: &types.UserLocation{
            Latitude:  55.7558,
            Longitude: 37.6173,
        },
    },
}

result, err := client.ExecuteTemplate(ctx, req)

// Обработка многошагового сценария
if result.QueryType == "with_purchases_services" {
    // Обработка workflow
    for _, section := range result.Sections {
        switch section.DomainID {
        case "commerce":
            // Обработка заказа еды
            for _, item := range section.Results {
                if item.Type == "food_order" {
                    // Выполнение заказа
                    executeAction(item.Actions[0]) // order_now
                }
            }
        case "payment":
            // Обработка оплаты
            for _, item := range section.Results {
                if item.Type == "payment_processing" {
                    // Выполнение оплаты
                    executeAction(item.Actions[0]) // process_payment
                }
            }
        case "delivery":
            // Обработка адреса доставки
            for _, item := range section.Results {
                if item.Type == "delivery_address" {
                    // Ввод адреса
                    executeAction(item.Actions[0]) // set_delivery_address
                }
            }
        case "notifications":
            // Создание напоминаний
            for _, item := range section.Results {
                if item.Type == "reminder" || item.Type == "medication_reminder" {
                    // Создание напоминания
                    executeAction(item.Actions[0]) // create_reminder
                }
            }
        }
    }
}
```

---

*Документация обновлена: v1.2.0 - Добавлен AI анализ доменов, качество ответов и динамическая маршрутизация*
