import React, { useState, useEffect } from 'react';
import clsx from 'clsx';

const PerformanceMonitor = ({ endpoint = '/api/v1/health' }) => {
  const [metrics, setMetrics] = useState({
    responseTime: 0,
    status: 'unknown',
    uptime: 0,
    requests: 0,
    errors: 0
  });
  const [isMonitoring, setIsMonitoring] = useState(false);
  const [history, setHistory] = useState([]);
  const [interval, setInterval] = useState(5000);

  const checkPerformance = async () => {
    const startTime = Date.now();

    try {
      const response = await fetch(endpoint, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      const responseTime = Date.now() - startTime;
      const data = await response.json();

      const newMetrics = {
        responseTime,
        status: response.ok ? 'healthy' : 'error',
        uptime: data.uptime || 0,
        requests: data.total_requests || 0,
        errors: data.error_count || 0,
        timestamp: Date.now()
      };

      setMetrics(newMetrics);

      // Добавляем в историю
      setHistory(prev => {
        const newHistory = [...prev, newMetrics];
        return newHistory.slice(-20); // Храним последние 20 измерений
      });

    } catch (error) {
      const responseTime = Date.now() - startTime;
      const errorMetrics = {
        responseTime,
        status: 'error',
        uptime: metrics.uptime,
        requests: metrics.requests,
        errors: metrics.errors + 1,
        timestamp: Date.now()
      };

      setMetrics(errorMetrics);
      setHistory(prev => [...prev, errorMetrics].slice(-20));
    }
  };

  useEffect(() => {
    let intervalId;

    if (isMonitoring) {
      // Немедленная проверка
      checkPerformance();

      // Установка интервала
      intervalId = setInterval(checkPerformance, interval);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [isMonitoring, interval, endpoint]);

  const toggleMonitoring = () => {
    setIsMonitoring(!isMonitoring);
  };

  const clearHistory = () => {
    setHistory([]);
  };

  const getAverageResponseTime = () => {
    if (history.length === 0) return 0;
    const sum = history.reduce((acc, item) => acc + item.responseTime, 0);
    return Math.round(sum / history.length);
  };

  const getSuccessRate = () => {
    if (history.length === 0) return 0;
    const successCount = history.filter(item => item.status === 'healthy').length;
    return Math.round((successCount / history.length) * 100);
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'healthy': return 'var(--nexus-success)';
      case 'warning': return 'var(--nexus-warning)';
      case 'error': return 'var(--nexus-error)';
      default: return 'var(--nexus-gray-400)';
    }
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case 'healthy': return '🟢';
      case 'warning': return '🟡';
      case 'error': return '🔴';
      default: return '⚪';
    }
  };

  return (
    <div className="performance-monitor">
      <div className="performance-monitor__header">
        <div className="monitor-info">
          <h3>📊 Performance Monitor</h3>
          <p>Мониторинг производительности API в реальном времени</p>
          <code className="monitor-endpoint">{endpoint}</code>
        </div>

        <div className="monitor-controls">
          <div className="control-group">
            <label>Интервал:</label>
            <select
              value={interval}
              onChange={(e) => setInterval(Number(e.target.value))}
              className="interval-select"
            >
              <option value={1000}>1 сек</option>
              <option value={2000}>2 сек</option>
              <option value={5000}>5 сек</option>
              <option value={10000}>10 сек</option>
              <option value={30000}>30 сек</option>
            </select>
          </div>

          <div className="control-buttons">
            <button
              onClick={toggleMonitoring}
              className={clsx('monitor-btn', {
                'monitor-btn--active': isMonitoring
              })}
            >
              {isMonitoring ? '⏸️ Стоп' : '▶️ Старт'}
            </button>
            <button
              onClick={clearHistory}
              className="monitor-btn monitor-btn--secondary"
            >
              🗑️ Очистить
            </button>
          </div>
        </div>
      </div>

      <div className="performance-monitor__metrics">
        <div className="metrics-grid">
          <div className="metric-card metric-card--primary">
            <div className="metric-header">
              <span className="metric-icon">⚡</span>
              <span className="metric-title">Response Time</span>
            </div>
            <div className="metric-value">
              {metrics.responseTime}ms
            </div>
            <div className="metric-subtitle">
              Avg: {getAverageResponseTime()}ms
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-header">
              <span className="metric-icon">{getStatusIcon(metrics.status)}</span>
              <span className="metric-title">Status</span>
            </div>
            <div
              className="metric-value"
              style={{ color: getStatusColor(metrics.status) }}
            >
              {metrics.status.toUpperCase()}
            </div>
            <div className="metric-subtitle">
              Success: {getSuccessRate()}%
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-header">
              <span className="metric-icon">🔄</span>
              <span className="metric-title">Requests</span>
            </div>
            <div className="metric-value">
              {metrics.requests.toLocaleString()}
            </div>
            <div className="metric-subtitle">
              Total served
            </div>
          </div>

          <div className="metric-card">
            <div className="metric-header">
              <span className="metric-icon">⏱️</span>
              <span className="metric-title">Uptime</span>
            </div>
            <div className="metric-value">
              {Math.round(metrics.uptime / 1000 / 60)}m
            </div>
            <div className="metric-subtitle">
              Service uptime
            </div>
          </div>
        </div>
      </div>

      <div className="performance-monitor__chart">
        <h4>📈 Response Time History</h4>
        <div className="chart-container">
          <div className="chart-grid">
            {history.map((item, index) => (
              <div
                key={index}
                className={clsx('chart-bar', {
                  'chart-bar--error': item.status === 'error',
                  'chart-bar--warning': item.responseTime > 1000,
                  'chart-bar--success': item.status === 'healthy' && item.responseTime <= 1000
                })}
                style={{
                  height: `${Math.min(item.responseTime / 10, 100)}%`,
                  backgroundColor: item.status === 'error' ? 'var(--nexus-error)' :
                                   item.responseTime > 1000 ? 'var(--nexus-warning)' :
                                   'var(--nexus-success)'
                }}
                title={`${item.responseTime}ms - ${item.status}`}
              />
            ))}
          </div>
          <div className="chart-labels">
            <span>0ms</span>
            <span>500ms</span>
            <span>1000ms</span>
            <span>1500ms</span>
          </div>
        </div>

        <div className="chart-legend">
          <div className="legend-item">
            <div className="legend-color" style={{ backgroundColor: 'var(--nexus-success)' }}></div>
            <span>&lt; 1000ms</span>
          </div>
          <div className="legend-item">
            <div className="legend-color" style={{ backgroundColor: 'var(--nexus-warning)' }}></div>
            <span>&gt; 1000ms</span>
          </div>
          <div className="legend-item">
            <div className="legend-color" style={{ backgroundColor: 'var(--nexus-error)' }}></div>
            <span>Error</span>
          </div>
        </div>
      </div>

      <div className="performance-monitor__footer">
        <div className="monitor-stats">
          <span>Всего проверок: {history.length}</span>
          <span>Последняя проверка: {new Date(metrics.timestamp || Date.now()).toLocaleTimeString()}</span>
        </div>

        <div className="monitor-note">
          <small>
            💡 <strong>Примечание:</strong> Это демонстрационная версия монитора.
            В production используйте Prometheus + Grafana для полноценного мониторинга.
          </small>
        </div>
      </div>
    </div>
  );
};

export default PerformanceMonitor;
