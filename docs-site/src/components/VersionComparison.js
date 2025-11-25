import React, { useState } from 'react';
import clsx from 'clsx';

const VersionComparison = () => {
  const [selectedVersion, setSelectedVersion] = useState('2.0.0');

  const versions = {
    '1.0.0': {
      name: 'v1.0.0 - Basic',
      releaseDate: '2024-06-01',
      features: [
        { name: 'Execute Template', status: '✅' },
        { name: 'Basic Error Handling', status: '✅' },
        { name: 'Simple Metadata', status: '✅' },
        { name: 'Rate Limiting', status: '❌' },
        { name: 'Batch Operations', status: '❌' },
        { name: 'Webhooks', status: '❌' },
        { name: 'Analytics', status: '❌' },
        { name: 'Admin API', status: '❌' },
        { name: 'Enterprise Features', status: '❌' }
      ],
      performance: {
        throughput: '500 RPS',
        latency: '500ms',
        concurrentUsers: '1000'
      }
    },
    '1.2.1': {
      name: 'v1.2.1 - Enhanced',
      releaseDate: '2024-09-15',
      features: [
        { name: 'Execute Template', status: '✅' },
        { name: 'Enhanced Error Handling', status: '✅' },
        { name: 'Rich Metadata', status: '✅' },
        { name: 'Basic Rate Limiting', status: '✅' },
        { name: 'Batch Operations', status: '❌' },
        { name: 'Webhooks', status: '❌' },
        { name: 'Analytics', status: '❌' },
        { name: 'Admin API', status: '❌' },
        { name: 'Enterprise Features', status: '❌' }
      ],
      performance: {
        throughput: '750 RPS',
        latency: '350ms',
        concurrentUsers: '2500'
      }
    },
    '2.0.0': {
      name: 'v2.0.0 - Enterprise',
      releaseDate: '2025-01-18',
      features: [
        { name: 'Execute Template', status: '✅' },
        { name: 'Advanced Error Handling', status: '✅' },
        { name: 'Enterprise Metadata', status: '✅' },
        { name: 'Advanced Rate Limiting', status: '✅' },
        { name: 'Batch Operations', status: '✅' },
        { name: 'Webhooks', status: '✅' },
        { name: 'Analytics', status: '✅' },
        { name: 'Admin API', status: '✅' },
        { name: 'Enterprise Features', status: '✅' }
      ],
      performance: {
        throughput: '2000+ RPS',
        latency: '150ms',
        concurrentUsers: '10000+'
      }
    }
  };

  const currentVersion = versions[selectedVersion];

  return (
    <div className="version-comparison">
      <div className="version-comparison__header">
        <h3>📊 Сравнение версий Nexus Protocol</h3>
        <p>Выберите версию для сравнения возможностей и производительности</p>
      </div>

      <div className="version-comparison__selector">
        {Object.keys(versions).map(version => (
          <button
            key={version}
            className={clsx('version-comparison__version-btn', {
              active: selectedVersion === version
            })}
            onClick={() => setSelectedVersion(version)}
          >
            <span className="version-number">{version}</span>
            <span className="version-name">{versions[version].name.split(' - ')[1]}</span>
            <span className="version-date">{versions[version].releaseDate}</span>
          </button>
        ))}
      </div>

      <div className="version-comparison__content">
        <div className="version-comparison__section">
          <h4>🚀 Возможности</h4>
          <div className="version-comparison__features">
            {currentVersion.features.map((feature, index) => (
              <div key={index} className="version-comparison__feature">
                <span className="feature-status">{feature.status}</span>
                <span className="feature-name">{feature.name}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="version-comparison__section">
          <h4>⚡ Производительность</h4>
          <div className="version-comparison__metrics">
            <div className="metric-card">
              <div className="metric-value">{currentVersion.performance.throughput}</div>
              <div className="metric-label">Throughput</div>
            </div>
            <div className="metric-card">
              <div className="metric-value">{currentVersion.performance.latency}</div>
              <div className="metric-label">Avg Latency</div>
            </div>
            <div className="metric-card">
              <div className="metric-value">{currentVersion.performance.concurrentUsers}</div>
              <div className="metric-label">Concurrent Users</div>
            </div>
          </div>
        </div>

        <div className="version-comparison__section">
          <h4>🎯 Рекомендации</h4>
          <div className="version-comparison__recommendations">
            {selectedVersion === '1.0.0' && (
              <div className="recommendation-card recommendation-card--warning">
                <h5>⚠️ Legacy Version</h5>
                <p>Версия 1.0.0 больше не поддерживается. Рекомендуется миграция на v2.0.0 для получения enterprise возможностей.</p>
                <a href="/migration/step-by-step" className="cta-button">📚 Руководство миграции</a>
              </div>
            )}
            {selectedVersion === '1.2.1' && (
              <div className="recommendation-card recommendation-card--info">
                <h5>🔄 LTS Version</h5>
                <p>Версия 1.2.1 является LTS (Long Term Support) и будет поддерживаться до июня 2026. Рассмотрите миграцию на v2.0.0 для новых возможностей.</p>
                <a href="/migration/step-by-step" className="cta-button">🚀 Обновиться до v2.0.0</a>
              </div>
            )}
            {selectedVersion === '2.0.0' && (
              <div className="recommendation-card recommendation-card--success">
                <h5>✅ Latest Version</h5>
                <p>Версия 2.0.0 - это текущая enterprise-ready версия со всеми новейшими возможностями и лучшей производительностью.</p>
                <a href="/sdk/quick-start" className="cta-button">🚀 Начать работу</a>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default VersionComparison;
