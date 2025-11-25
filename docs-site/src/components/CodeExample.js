import React from 'react';
import CodeBlock from '@theme/CodeBlock';
import clsx from 'clsx';

const CodeExample = ({
  title,
  language = 'json',
  children,
  className,
  showCopyButton = true,
  ...props
}) => {
  return (
    <div className={clsx('code-example', className)}>
      {title && (
        <div className="code-example__title">
          <span className="code-example__icon">💻</span>
          {title}
        </div>
      )}
      <CodeBlock
        language={language}
        showLineNumbers={false}
        {...props}
      >
        {children}
      </CodeBlock>
    </div>
  );
};

export default CodeExample;
