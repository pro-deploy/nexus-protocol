/**
 * Боковое меню с якорными ссылками на разделы одной страницы
 * Docusaurus автоматически создает якоря из заголовков
 */

module.exports = {
  mainSidebar: [
    {
      type: 'doc',
      id: 'index',
      label: '🏠 Главная',
    },
    {
      type: 'doc',
      id: 'interactive-examples',
      label: '🎮 Интерактивные примеры',
    },
    {
      type: 'doc',
      id: 'audit-results',
      label: '📋 Результаты аудита',
    },

    // Protocol Section
    {
      type: 'category',
      label: '📋 Протокол',
      collapsed: false,
      items: [
        {
          type: 'doc',
          id: 'protocol/intro',
          label: 'Введение',
        },
        {
          type: 'doc',
          id: 'protocol/message-format',
          label: 'Формат сообщений',
        },
        {
          type: 'doc',
          id: 'protocol/metadata',
          label: 'Метаданные',
        },
        {
          type: 'doc',
          id: 'protocol/error-handling',
          label: 'Обработка ошибок',
        },
        {
          type: 'doc',
          id: 'protocol/versioning',
          label: 'Версионирование',
        },
      ],
    },

    // API Section
    {
      type: 'category',
      label: '🔌 API',
      collapsed: false,
      items: [
        {
          type: 'doc',
          id: 'protocol/rest-api',
          label: 'REST API',
        },
        {
          type: 'doc',
          id: 'protocol/grpc-api',
          label: 'gRPC API',
        },
        {
          type: 'doc',
          id: 'protocol/websocket-api',
          label: 'WebSocket API',
        },
        {
          type: 'link',
          href: '/api/rest/openapi.yaml',
          label: '📄 OpenAPI 3.0 (YAML)',
        },
        {
          type: 'link',
          href: '/api/grpc/nexus.proto',
          label: '📄 Protocol Buffers',
        },
        {
          type: 'link',
          href: '/api/websocket/protocol.json',
          label: '📄 WebSocket Protocol',
        },
      ],
    },

    // SDK Section
    {
      type: 'category',
      label: '🛠️ SDK',
      collapsed: false,
      items: [
        {
          type: 'doc',
          id: 'sdk/intro',
          label: 'Введение',
        },
        {
          type: 'doc',
          id: 'sdk/installation',
          label: 'Установка',
        },
        {
          type: 'doc',
          id: 'sdk/quick-start',
          label: 'Быстрый старт',
        },
        {
          type: 'doc',
          id: 'sdk/basic-usage',
          label: 'Базовое использование',
        },
        {
          type: 'doc',
          id: 'sdk/usage-guide',
          label: 'Руководство',
        },
        {
          type: 'doc',
          id: 'sdk/advanced-guide',
          label: 'Продвинутое использование',
        },
        {
          type: 'doc',
          id: 'sdk/examples',
          label: 'Примеры',
        },
        {
          type: 'doc',
          id: 'sdk/error-handling',
          label: 'Обработка ошибок',
        },
        {
          type: 'doc',
          id: 'sdk/batch-operations',
          label: 'Batch операции',
        },
        {
          type: 'doc',
          id: 'sdk/webhooks',
          label: 'Webhooks',
        },
        {
          type: 'doc',
          id: 'sdk/analytics',
          label: 'Аналитика',
        },
        {
          type: 'doc',
          id: 'sdk/admin-api',
          label: 'Admin API',
        },
        {
          type: 'doc',
          id: 'sdk/client-api',
          label: 'Client API',
        },
        {
          type: 'doc',
          id: 'sdk/types',
          label: 'Types',
        },
      ],
    },

    // Schemas Section
    {
      type: 'category',
      label: '📋 Schemas',
      items: [
        {
          type: 'doc',
          id: 'schemas/schemas-index',
          label: 'Обзор',
        },
        {
          type: 'doc',
          id: 'schemas/validation-examples',
          label: 'Примеры валидации',
        },
        {
          type: 'link',
          href: '/schemas/message-schema.json',
          label: '📄 Message Schema (JSON)',
        },
      ],
    },

    // Additional Resources
    {
      type: 'category',
      label: '📚 Ресурсы',
      items: [
        {
          type: 'doc',
          id: 'analytics/analytics-index',
          label: '📊 Аналитика и метрики',
        },
        {
          type: 'doc',
          id: 'migration/migration-index',
          label: '🚀 Миграция',
        },
        {
          type: 'link',
          href: 'https://github.com/nexus-protocol',
          label: 'GitHub Repository',
        },
        {
          type: 'link',
          href: 'https://github.com/nexus-protocol/nexus-protocol/issues',
          label: 'Сообщить о проблеме',
        },
        {
          type: 'link',
          href: 'https://github.com/nexus-protocol/nexus-protocol/discussions',
          label: 'Обсуждения',
        },
        {
          type: 'link',
          href: 'https://nexus.dev',
          label: 'Nexus Platform',
        },
      ],
    },
  ],
};
