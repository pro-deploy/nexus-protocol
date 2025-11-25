import React, { useState, useEffect } from 'react';
import clsx from 'clsx';

const InteractiveTutorial = ({ tutorialId = 'basic-api' }) => {
  const [currentStep, setCurrentStep] = useState(0);
  const [completedSteps, setCompletedSteps] = useState([]);
  const [userInput, setUserInput] = useState('');
  const [feedback, setFeedback] = useState(null);

  const tutorials = {
    'basic-api': {
      title: '🚀 Основы работы с Nexus API',
      description: 'Интерактивный туториал по основным операциям API',
      steps: [
        {
          title: 'Шаг 1: Создание клиента',
          content: 'Сначала создадим клиент для работы с Nexus API',
          code: `import NexusClient from 'nexus-protocol';

const client = new NexusClient({
  baseURL: 'https://api.nexus.dev',
  token: 'your-jwt-token',
  protocolVersion: '2.0.0'
});`,
          task: 'Создайте клиент с правильными параметрами',
          validation: (code) => {
            return code.includes('NexusClient') &&
                   code.includes('baseURL') &&
                   code.includes('protocolVersion') &&
                   code.includes('2.0.0');
          },
          hint: 'Используйте конструктор NexusClient с baseURL, token и protocolVersion'
        },
        {
          title: 'Шаг 2: Простой запрос',
          content: 'Отправим простой запрос на выполнение шаблона',
          code: `const request = {
  metadata: {
    request_id: "req-" + Date.now(),
    protocol_version: "2.0.0",
    client_version: "2.0.0"
  },
  data: {
    query: "хочу борщ",
    language: "ru"
  }
};

const response = await client.executeTemplate(request);`,
          task: 'Создайте правильный запрос с metadata и data',
          validation: (code) => {
            return code.includes('request_id') &&
                   code.includes('protocol_version') &&
                   code.includes('client_version') &&
                   code.includes('executeTemplate');
          },
          hint: 'Не забудьте добавить request_id, protocol_version и client_version в metadata'
        },
        {
          title: 'Шаг 3: Обработка ответа',
          content: 'Обработаем ответ и извлечем нужную информацию',
          code: `if (response.metadata.processing_time_ms > 1000) {
  console.log('Запрос выполнен долго');
}

const executionId = response.data.execution_id;
const status = response.data.status;

console.log(\`Execution \${executionId}: \${status}\`);`,
          task: 'Извлеките execution_id и status из ответа',
          validation: (code) => {
            return code.includes('execution_id') &&
                   code.includes('status') &&
                   code.includes('console.log');
          },
          hint: 'Посмотрите структуру ответа в документации API'
        },
        {
          title: 'Шаг 4: Обработка ошибок',
          content: 'Добавим правильную обработку ошибок',
          code: `try {
  const response = await client.executeTemplate(request);
  console.log('Успех:', response.data);
} catch (error) {
  if (error.code === 'VALIDATION_FAILED') {
    console.error('Ошибка валидации:', error.message);
  } else if (error.code === 'AUTHENTICATION_FAILED') {
    console.error('Ошибка аутентификации');
  } else {
    console.error('Неизвестная ошибка:', error.message);
  }
}`,
          task: 'Добавьте обработку разных типов ошибок',
          validation: (code) => {
            return code.includes('try') &&
                   code.includes('catch') &&
                   code.includes('error.code') &&
                   code.includes('VALIDATION_FAILED');
          },
          hint: 'Используйте error.code для определения типа ошибки'
        }
      ]
    }
  };

  const tutorial = tutorials[tutorialId];
  const currentStepData = tutorial.steps[currentStep];

  const handleNext = () => {
    if (currentStep < tutorial.steps.length - 1) {
      setCompletedSteps([...completedSteps, currentStep]);
      setCurrentStep(currentStep + 1);
      setUserInput('');
      setFeedback(null);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
      setUserInput(tutorial.steps[currentStep - 1].code);
      setFeedback(null);
    }
  };

  const handleCheck = () => {
    const isValid = currentStepData.validation(userInput);
    setFeedback({
      type: isValid ? 'success' : 'error',
      message: isValid ? 'Отлично! Код правильный.' : 'Есть ошибки. Проверьте подсказку.'
    });

    if (isValid && !completedSteps.includes(currentStep)) {
      setCompletedSteps([...completedSteps, currentStep]);
    }
  };

  const handleReset = () => {
    setUserInput('');
    setFeedback(null);
  };

  useEffect(() => {
    setUserInput(currentStepData.code);
  }, [currentStep]);

  const progress = ((completedSteps.length) / tutorial.steps.length) * 100;

  return (
    <div className="interactive-tutorial">
      <div className="interactive-tutorial__header">
        <div className="tutorial-info">
          <h3>{tutorial.title}</h3>
          <p>{tutorial.description}</p>
        </div>
        <div className="tutorial-progress">
          <div className="progress-bar">
            <div
              className="progress-fill"
              style={{ width: `${progress}%` }}
            />
          </div>
          <span className="progress-text">
            {completedSteps.length} из {tutorial.steps.length} шагов
          </span>
        </div>
      </div>

      <div className="interactive-tutorial__content">
        <div className="tutorial-step">
          <div className="step-header">
            <h4>{currentStepData.title}</h4>
            <div className="step-indicator">
              <span className="step-number">{currentStep + 1}</span>
              <span className="step-total">/ {tutorial.steps.length}</span>
            </div>
          </div>

          <div className="step-content">
            <p>{currentStepData.content}</p>

            <div className="step-task">
              <strong>Задание:</strong> {currentStepData.task}
            </div>

            <div className="code-editor">
              <div className="code-editor__header">
                <span className="code-language">JavaScript</span>
                <div className="code-actions">
                  <button onClick={handleCheck} className="code-btn code-btn--check">
                    ✅ Проверить
                  </button>
                  <button onClick={handleReset} className="code-btn code-btn--reset">
                    🔄 Сбросить
                  </button>
                </div>
              </div>

              <textarea
                className="code-editor__textarea"
                value={userInput}
                onChange={(e) => setUserInput(e.target.value)}
                placeholder="Введите ваш код здесь..."
              />

              {feedback && (
                <div className={clsx('code-feedback', `code-feedback--${feedback.type}`)}>
                  {feedback.type === 'success' ? '✅' : '❌'} {feedback.message}
                  {feedback.type === 'error' && (
                    <div className="feedback-hint">
                      💡 <strong>Подсказка:</strong> {currentStepData.hint}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="interactive-tutorial__footer">
        <button
          onClick={handlePrevious}
          disabled={currentStep === 0}
          className="tutorial-btn tutorial-btn--secondary"
        >
          ← Назад
        </button>

        <div className="step-navigation">
          {tutorial.steps.map((_, index) => (
            <button
              key={index}
              onClick={() => setCurrentStep(index)}
              className={clsx('step-dot', {
                active: currentStep === index,
                completed: completedSteps.includes(index)
              })}
            />
          ))}
        </div>

        <button
          onClick={handleNext}
          disabled={!completedSteps.includes(currentStep)}
          className="tutorial-btn tutorial-btn--primary"
        >
          {currentStep === tutorial.steps.length - 1 ? '🎉 Завершить' : 'Далее →'}
        </button>
      </div>
    </div>
  );
};

export default InteractiveTutorial;
