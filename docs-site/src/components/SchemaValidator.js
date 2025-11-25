import React, { useState } from 'react';
import clsx from 'clsx';

const SchemaValidator = ({ schemaUrl = '/schemas/message-schema.json' }) => {
  const [jsonInput, setJsonInput] = useState(`{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}`);
  const [validationResult, setValidationResult] = useState(null);
  const [loading, setLoading] = useState(false);

  const validateJSON = async () => {
    setLoading(true);
    setValidationResult(null);

    try {
      // Parse JSON
      const data = JSON.parse(jsonInput);

      // Mock validation for demo (in real app, this would call a validation service)
      await new Promise(resolve => setTimeout(resolve, 500));

      // Simple validation checks
      const errors = [];

      if (!data.metadata) {
        errors.push({
          keyword: 'required',
          dataPath: '',
          message: 'should have required property \'metadata\''
        });
      } else {
        if (!data.metadata.request_id) {
          errors.push({
            keyword: 'required',
            dataPath: '.metadata',
            message: 'should have required property \'request_id\''
          });
        }

        if (!data.metadata.protocol_version) {
          errors.push({
            keyword: 'required',
            dataPath: '.metadata',
            message: 'should have required property \'protocol_version\''
          });
        }

        if (!data.metadata.client_version) {
          errors.push({
            keyword: 'required',
            dataPath: '.metadata',
            message: 'should have required property \'client_version\''
          });
        }
      }

      if (errors.length > 0) {
        setValidationResult({
          valid: false,
          errors: errors
        });
      } else {
        setValidationResult({
          valid: true,
          message: 'JSON is valid according to Nexus Protocol schema'
        });
      }
    } catch (error) {
      setValidationResult({
        valid: false,
        errors: [{
          keyword: 'parse',
          message: `JSON parsing error: ${error.message}`
        }]
      });
    } finally {
      setLoading(false);
    }
  };

  const clearInput = () => {
    setJsonInput(`{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}`);
    setValidationResult(null);
  };

  const loadExample = (type) => {
    const examples = {
      valid: `{
  "metadata": {
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "protocol_version": "2.0.0",
    "client_version": "2.0.0",
    "client_id": "web-app",
    "client_type": "web",
    "timestamp": 1640995200
  },
  "data": {
    "query": "хочу борщ",
    "language": "ru",
    "context": {
      "user_id": "user-123",
      "location": {
        "latitude": 55.7558,
        "longitude": 37.6173,
        "accuracy": 50
      },
      "locale": "ru-RU",
      "currency": "RUB"
    }
  }
}`,
      invalid: `{
  "data": {
    "query": "хочу борщ"
  }
}`,
      error: `{
  "metadata": {
    "request_id": "invalid-uuid",
    "protocol_version": "latest",
    "client_version": "2.0.0"
  },
  "data": {
    "query": "хочу борщ"
  }
}`
    };

    setJsonInput(examples[type]);
    setValidationResult(null);
  };

  return (
    <div className="schema-validator">
      <div className="schema-validator__header">
        <h3>🔍 JSON Schema Validator</h3>
        <p>Проверьте корректность JSON сообщений согласно Nexus Protocol</p>
      </div>

      <div className="schema-validator__content">
        <div className="schema-validator__input-group">
          <label className="schema-validator__label">
            JSON для валидации:
          </label>
          <textarea
            className="schema-validator__textarea"
            value={jsonInput}
            onChange={(e) => setJsonInput(e.target.value)}
            placeholder="Введите JSON объект..."
          />
        </div>

        <div className="schema-validator__actions">
          <button
            className="schema-validator__button schema-validator__button--validate"
            onClick={validateJSON}
            disabled={loading}
          >
            {loading ? '⏳ Валидация...' : '✅ Проверить'}
          </button>
          <button
            className="schema-validator__button schema-validator__button--clear"
            onClick={clearInput}
          >
            🗑️ Очистить
          </button>
          <div style={{ marginLeft: 'auto', display: 'flex', gap: '8px' }}>
            <button
              className="schema-validator__button schema-validator__button--clear"
              onClick={() => loadExample('valid')}
            >
              ✅ Пример валидный
            </button>
            <button
              className="schema-validator__button schema-validator__button--clear"
              onClick={() => loadExample('invalid')}
            >
              ❌ Пример невалидный
            </button>
            <button
              className="schema-validator__button schema-validator__button--clear"
              onClick={() => loadExample('error')}
            >
              ⚠️ С ошибками
            </button>
          </div>
        </div>

        {validationResult && (
          <div className="schema-validator__input-group">
            <label className="schema-validator__label">
              Результат валидации:
            </label>
            <div className={clsx('schema-validator__result', {
              'schema-validator__result--valid': validationResult.valid,
              'schema-validator__result--invalid': !validationResult.valid
            })}>
              {validationResult.valid ? (
                <div>
                  <strong>✅ JSON валиден!</strong>
                  <br />
                  {validationResult.message}
                </div>
              ) : (
                <div>
                  <strong>❌ Найдены ошибки валидации:</strong>
                  <br />
                  {validationResult.errors.map((error, index) => (
                    <div key={index} style={{ marginTop: '8px' }}>
                      <code>{error.dataPath || '/'}</code>: {error.message}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        <div style={{ marginTop: '16px', padding: '16px', background: 'var(--nexus-gray-50)', borderRadius: '8px', fontSize: '0.875rem' }}>
          <strong>💡 Подсказка:</strong> Используйте этот инструмент для проверки JSON сообщений перед отправкой в API.
          Валидация происходит согласно <a href="/schemas/message-schema.json" target="_blank">схеме Nexus Protocol</a>.
        </div>
      </div>
    </div>
  );
};

export default SchemaValidator;
