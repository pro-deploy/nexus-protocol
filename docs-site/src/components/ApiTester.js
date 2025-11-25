import React, { useState, useEffect } from 'react';
import clsx from 'clsx';

const ApiTester = ({ endpoint, method = 'GET', initialData = '{}' }) => {
  const [requestData, setRequestData] = useState(initialData);
  const [response, setResponse] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [activeTab, setActiveTab] = useState('request');

  const handleSubmit = async () => {
    setLoading(true);
    setError(null);

    try {
      const options = {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer demo-token',
        },
      };

      if (method !== 'GET' && method !== 'DELETE') {
        options.body = requestData;
      }

      // Mock API call for demo
      await new Promise(resolve => setTimeout(resolve, 1000));

      const mockResponse = {
        metadata: {
          request_id: '550e8400-e29b-41d4-a716-446655440000',
          protocol_version: '2.0.0',
          server_version: '2.0.0',
          timestamp: Math.floor(Date.now() / 1000),
          processing_time_ms: Math.floor(Math.random() * 500) + 100,
        },
        data: method === 'GET' ? {
          message: 'API endpoint работает корректно',
          timestamp: new Date().toISOString(),
        } : JSON.parse(requestData),
      };

      setResponse(mockResponse);
      setActiveTab('response');
    } catch (err) {
      setError(err.message);
      setActiveTab('error');
    } finally {
      setLoading(false);
    }
  };

  const methodColors = {
    GET: 'api-method--get',
    POST: 'api-method--post',
    PUT: 'api-method--put',
    DELETE: 'api-method--delete',
  };

  return (
    <div className="api-tester">
      <div className="api-tester__header">
        <div className="api-tester__endpoint">
          <span className={clsx('api-method', methodColors[method])}>
            {method}
          </span>
          <code className="api-tester__url">{endpoint}</code>
        </div>
        <button
          className="cta-button api-tester__button"
          onClick={handleSubmit}
          disabled={loading}
        >
          {loading ? '⏳ Отправка...' : '🚀 Отправить запрос'}
        </button>
      </div>

      <div className="api-tester__tabs">
        <button
          className={clsx('api-tester__tab', { active: activeTab === 'request' })}
          onClick={() => setActiveTab('request')}
        >
          📝 Запрос
        </button>
        <button
          className={clsx('api-tester__tab', { active: activeTab === 'response' })}
          onClick={() => setActiveTab('response')}
          disabled={!response}
        >
          📤 Ответ
        </button>
        <button
          className={clsx('api-tester__tab', { active: activeTab === 'error' })}
          onClick={() => setActiveTab('error')}
          disabled={!error}
        >
          ❌ Ошибка
        </button>
      </div>

      <div className="api-tester__content">
        {activeTab === 'request' && (
          <div className="api-tester__panel">
            <h4>Request Body</h4>
            <textarea
              className="api-tester__textarea"
              value={requestData}
              onChange={(e) => setRequestData(e.target.value)}
              placeholder="Введите JSON данные..."
              disabled={method === 'GET' || method === 'DELETE'}
            />
          </div>
        )}

        {activeTab === 'response' && response && (
          <div className="api-tester__panel">
            <h4>Response</h4>
            <pre className="api-tester__code">
              {JSON.stringify(response, null, 2)}
            </pre>
          </div>
        )}

        {activeTab === 'error' && error && (
          <div className="api-tester__panel api-tester__panel--error">
            <h4>Ошибка</h4>
            <div className="alert alert--danger">
              <div className="alert__content">
                <strong>Ошибка выполнения запроса:</strong> {error}
              </div>
            </div>
          </div>
        )}
      </div>

      <div className="api-tester__note">
        <small>
          💡 <strong>Примечание:</strong> Это демонстрационная версия API Tester.
          В реальном приложении запросы будут отправляться на настоящий сервер.
        </small>
      </div>
    </div>
  );
};

export default ApiTester;
